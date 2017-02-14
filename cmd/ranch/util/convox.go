package util

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path"
	"strconv"
	"time"

	"github.com/convox/rack/client"
	"github.com/jhoonb/archivex"
	"github.com/spf13/viper"
)

func convoxClient() (*client.Client, error) {
	if !viper.IsSet("convox.host") || !viper.IsSet("convox.password") {
		return nil, fmt.Errorf("must set 'convox.host' and 'convox.password' in $HOME/.ranch.yaml")
	}

	host := viper.GetString("convox.host")
	password := viper.GetString("convox.password")
	version := "20151211151200"
	return client.New(host, password, version), nil
}

// ConvoxReleases returns a list of releases with the Convox ReleaseId mapped
// to a Ranch ReleaseId.
func ConvoxReleases(appName string) ([]RanchRelease, error) {
	client, err := convoxClient()

	if err != nil {
		return nil, err
	}

	app, err := client.GetApp(appName)

	if err != nil {
		return nil, err
	}

	convoxReleases, err := client.GetReleases(appName)

	if err != nil {
		return nil, err
	}

	shaMap, err := buildShaMap(appName)

	if err != nil {
		return nil, err
	}

	var releases []RanchRelease

	for _, convoxRelease := range convoxReleases {
		status := ""

		if app.Release == convoxRelease.Id {
			status = "active"
		}

		appVersion, ok := shaMap[convoxRelease.Id]

		if !ok {
			continue
		}

		release := RanchRelease{
			Id:      appVersion,
			App:     appName,
			Created: convoxRelease.Created,
			Status:  status,
		}

		releases = append(releases, release)
	}

	return releases, nil
}

// ConvoxRunDetached starts a detached run of a given app, process, and command.
func ConvoxRunDetached(appName, process, command string) error {
	client, err := convoxClient()

	if err != nil {
		return err
	}

	return client.RunProcessDetached(appName, process, command)
}

// ConvoxRunAttached starts an attached run of a given app, process, and command.
func ConvoxRunAttached(appName, process, command string, input io.Reader, output io.WriteCloser) (int, error) {
	client, err := convoxClient()

	if err != nil {
		return -1, err
	}

	return client.RunProcessAttached(appName, process, command, input, output)
}

// ConvoxExec runs a command inside the given Convox pid, using `docker exec`.
func ConvoxExec(appName, pid, command string, input io.Reader, output io.WriteCloser) (int, error) {
	client, err := convoxClient()

	if err != nil {
		return -1, err
	}

	return client.ExecProcessAttached(appName, pid, command, input, output)
}

// ConvoxLogs tails (and follows) the logs for a given application.
func ConvoxLogs(appName string, output io.WriteCloser) error {
	client, err := convoxClient()

	if err != nil {
		return err
	}

	return client.StreamAppLogs(appName, output)
}

// ConvoxGetFormation returns the formation of the given app translated into a RanchFormation.
func ConvoxGetFormation(appName string) (formation RanchFormation, err error) {

	formation = make(RanchFormation)

	client, err := convoxClient()

	if err != nil {
		return nil, err
	}

	convoxFormation, err := client.ListFormation(appName)

	if err != nil {
		return nil, err
	}

	for _, convoxFormationEntry := range convoxFormation {
		formation[convoxFormationEntry.Name] = RanchFormationEntry{
			Count:    convoxFormationEntry.Count,
			Memory:   convoxFormationEntry.Memory,
			Balancer: convoxFormationEntry.Balancer,
		}
	}

	return formation, nil
}

// ConvoxScaleProcess applies the given scale (count, memory) to the given process.
func ConvoxScaleProcess(appName, processName string, count, memory int) (err error) {
	client, err := convoxClient()

	if err != nil {
		return err
	}

	message := fmt.Sprintf("scaling %s to", processName)
	if count > -1 {
		message += fmt.Sprintf(" count=%d", count)
	}
	if memory > -1 {
		message += fmt.Sprintf(" memory=%d", memory)
	}
	fmt.Println(message)

	strCount := ""
	if count != -1 {
		strCount = strconv.Itoa(count)
	}

	strMemory := ""
	if memory != -1 {
		strMemory = strconv.Itoa(memory)
	}

	if err = client.SetFormation(appName, processName, strCount, strMemory); err != nil {
		return err
	}

	return nil
}

// ConvoxScale iterates over the given RanchConfig applying the correct scale to each named process.
func ConvoxScale(appName string, config *RanchConfig) (err error) {

	fmt.Println("🐮  Scaling your app:")

	existingFormation, err := ConvoxGetFormation(appName)

	if err != nil {
		return err
	}

	// scale down hidden 'run' process, which is used by ranch run and cron.
	if existingEntry, ok := existingFormation["run"]; ok {
		fmt.Printf(" - run count=0 memory=2048 ")
		if existingEntry.Count != 0 || existingEntry.Memory != 2048 {
			if err = ConvoxScaleProcess(appName, "run", 0, 2048); err != nil {
				fmt.Println("✘")
				return err
			}

			time.Sleep(5 * time.Second) // wait for scale to apply

			if err = ConvoxWaitForStatusWithMessage(appName, "running", ""); err != nil {
				// graphics handled by parent
				return err
			}
		} else {
			fmt.Println("✔")
		}
	}

	for processName, processConfig := range config.Processes {
		fmt.Printf(" - %s count=%d memory=%d ", processName, processConfig.Count, processConfig.Memory)
		if existingEntry, ok := existingFormation[processName]; ok {
			if existingEntry.Count == processConfig.Count && existingEntry.Memory == processConfig.Memory {
				fmt.Println("✔")
				continue
			}
		}

		if err = ConvoxScaleProcess(appName, processName, processConfig.Count, processConfig.Memory); err != nil {
			fmt.Println("✘")
			return err
		}

		time.Sleep(5 * time.Second) // wait for scale to apply

		if err = ConvoxWaitForStatusWithMessage(appName, "running", ""); err != nil {
			// graphics handled by parent
			return err
		}
		fmt.Println("✔")
	}

	return nil
}

// ConvoxCurrentVersion returns the currently active Convox release mapped to Ranch release.
func ConvoxCurrentVersion(appName string) (string, error) {
	client, err := convoxClient()
	if err != nil {
		return "", err
	}

	app, err := client.GetApp(appName)
	if err != nil {
		return "", err
	}

	shaMap, err := buildShaMap(appName)
	if err != nil {
		return "", err
	}

	sha, exists := shaMap[app.Release]
	if !exists {
		return "", fmt.Errorf("current running an unknown convox release %s", app.Release)
	}

	return sha, nil
}

// ConvoxPromote promotes a given Ranch release by mapping it back to a Convox release.
func ConvoxPromote(appName string, ranchReleaseID string) error {
	convoxReleaseID, err := getConvoxRelease(appName, ranchReleaseID)
	if err != nil {
		return err
	}

	client, err := convoxClient()
	if err != nil {
		return err
	}

	fmt.Printf("🐮  Promoting release %s... ", ranchReleaseID)

	if _, err = client.PromoteRelease(appName, convoxReleaseID); err != nil {
		fmt.Println("✘")
		return err
	}

	if err = ConvoxWaitForStatusWithMessage(appName, "updating", ""); err != nil {
		// graphics handled in WaitForStatus
		return err
	}

	return ConvoxWaitForStatusWithMessage(appName, "running", fmt.Sprintf("🐮  Waiting for release %s to roll out... ", ranchReleaseID))
}

// ConvoxDeploy creates a new Convox release given an app and build directory.
func ConvoxDeploy(appName string, buildDir string) (string, error) {
	client, err := convoxClient()

	if err != nil {
		return "", err
	}

	app, err := client.GetApp(appName)

	if err != nil {
		return "", err
	}

	switch app.Status {
	case "creating":
		return "", fmt.Errorf("app is still creating: %s", appName)
	case "running", "updating":
	default:
		return "", fmt.Errorf("unable to build app: %s", appName)
	}

	tar, err := createTarball(buildDir)

	if err != nil {
		return "", err
	}

	cache := true
	config := "docker-compose.yml"

	fmt.Print("🐮  Uploading Convox build... ")

	build, err := client.CreateBuildSource(appName, tar, cache, config)
	if err != nil {
		fmt.Println("✘")
		return "", err
	}

	fmt.Println("✔")

	return finishBuild(client, appName, build)
}

// ConvoxPsStop stops a Convox process.
func ConvoxPsStop(appName string, pid string) error {
	client, err := convoxClient()

	if err != nil {
		return err
	}

	_, err = client.StopProcess(appName, pid)
	return err
}

func createTarball(buildDir string) ([]byte, error) {
	tmpDir, err := ioutil.TempDir("", "ranch")
	defer os.RemoveAll(tmpDir)

	if err != nil {
		return nil, err
	}

	tgzfile := path.Join(tmpDir, "build.tar.gz")

	tar := new(archivex.TarFile)

	err = tar.Create(tgzfile)
	if err != nil {
		return nil, err
	}

	err = tar.AddAll(buildDir, false)
	if err != nil {
		return nil, err
	}

	err = tar.Close()
	if err != nil {
		return nil, err
	}

	return ioutil.ReadFile(tgzfile)
}

func finishBuild(client *client.Client, appName string, build *client.Build) (string, error) {

	if build.Id == "" {
		return "", fmt.Errorf("unable to fetch build id")
	}

	if Verbose {
		reader, writer := io.Pipe()
		go io.Copy(os.Stdout, reader)
		if err := client.StreamBuildLogs(appName, build.Id, writer); err != nil {
			return "", err
		}
	}

	return waitForBuild(client, appName, build.Id)
}

func waitForBuild(client *client.Client, appName, buildID string) (string, error) {
	timeout := time.After(30 * time.Minute)
	tick := time.Tick(10 * time.Second)

	fmt.Print("🐮  Waiting for Convox build to finish... ")

	for {
		select {
		case <-tick:
			build, err := client.GetBuild(appName, buildID)
			if err != nil {
				fmt.Println("✘")
				return "", err
			}

			switch build.Status {
			case "complete":
				fmt.Println("✔")
				return build.Release, nil
			case "error", "failed":
				fmt.Println("✘")
				return "", fmt.Errorf("%s build failed: %s", appName, build.Status)
			default:
				if Verbose {
					fmt.Print(".")
				}
			}
		case <-timeout:
			fmt.Println("✘")
			return "", fmt.Errorf("%s build failed: timeout", appName)
		}
	}
}

// ConvoxWaitForStatus waits for a Convox app to have a particular status.
func ConvoxWaitForStatus(appName, status string) error {
	return ConvoxWaitForStatusWithMessage(appName, status, fmt.Sprintf("🐮  Waiting for '%s' status... ", status))
}

// ConvoxWaitForStatusWithMessage waits for a Convox app to have a particular status and displays the given message.
func ConvoxWaitForStatusWithMessage(appName, status string, message string) error {
	timeout := time.After(30 * time.Minute)
	tick := time.Tick(10 * time.Second)

	failed := false

	client, err := convoxClient()
	if err != nil {
		return err
	}

	fmt.Print(message)

	for {
		select {
		case <-tick:
			app, err := client.GetApp(appName)
			if err != nil {
				fmt.Println("✘")
				return err
			}

			switch app.Status {
			case status:
				fmt.Println("✔")
				if failed {
					return fmt.Errorf("Your deploy was not healthy and was rolled back to the previous version.  Consult your app's logs or ask #delivery-eng for help.")
				}
				return nil
			case "rollback":
				if !failed {
					failed = true
					fmt.Print("✘ ROLLBACK\nWaiting for rollback... ")
				}
			default:
				if Verbose {
					fmt.Print(".")
				}
			}

		case <-timeout:
			fmt.Println("✘")
			return fmt.Errorf("The rollout took longer than 30 minutes so we gave up.")
		}
	}

}

// ConvoxPs returns an array of RanchProcess objects based on the currently running state of the app.
func ConvoxPs(appName string) ([]RanchProcess, error) {
	client, err := convoxClient()

	if err != nil {
		return nil, err
	}

	convoxPs, err := client.GetProcesses(appName, false) // false == no process stats

	if err != nil {
		return nil, err
	}

	shaMap, err := buildShaMap(appName)

	if err != nil {
		return nil, err
	}

	var ps []RanchProcess

	for _, v := range convoxPs {
		p := RanchProcess(v)

		sha, ok := shaMap[p.Release]

		if !ok {
			p.Release = "convox:" + p.Release
		} else {
			p.Release = sha
		}

		ps = append(ps, p)
	}

	return ps, nil
}

func buildShaMap(appName string) (map[string]string, error) {
	ranchReleases, err := RanchReleases(appName)

	if err != nil {
		return nil, err
	}

	shaMap := make(map[string]string)

	for _, ranchRelease := range ranchReleases {
		shaMap[ranchRelease.ConvoxRelease] = ranchRelease.Id
	}

	return shaMap, nil
}

func getConvoxRelease(appName, ranchReleaseID string) (convoxReleaseID string, err error) {
	ranchReleases, err := RanchReleases(appName)

	if err != nil {
		return "", err
	}

	for _, ranchRelease := range ranchReleases {
		if ranchRelease.Id == ranchReleaseID {
			return ranchRelease.ConvoxRelease, nil
		}
	}

	return "", fmt.Errorf("could not map release %s to a Convox release", ranchReleaseID)
}
