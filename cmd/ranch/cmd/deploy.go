package cmd

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"strings"
	"text/template"
	"time"

	"github.com/goodeggs/platform/cmd/ranch/Godeps/_workspace/src/github.com/spf13/cobra"
	"github.com/goodeggs/platform/cmd/ranch/util"
)

var Build bool

var cedarDockerfileTemplate = template.Must(template.New("cedar-dockerfile").Parse(`# generated by ranch
FROM goodeggs/cedar:4e94dfd
MAINTAINER Good Eggs <open-source@goodeggs.com>

# Build-time Environment
ARG RANCH_BUILD_ENV

ENTRYPOINT ["/usr/bin/profile"]

USER app
ENV HOME /app
WORKDIR /app

COPY . /build

RUN sudo mkdir -p /cache && \
  sudo chown -R app /buildkit /build /cache && \
  /usr/bin/build /build /cache && \
  sudo rm -rf /app && \
  sudo mv /build /app
`))

var nodejsDockerfileTemplate = template.Must(template.New("nodejs-dockerfile").Parse(`# generated by ranch
FROM goodeggs/ranch-baseimage-nodejs:3c4036d
MAINTAINER Good Eggs <open-source@goodeggs.com>
`))

type composeTemplateVars struct {
	ImageName   string
	Environment map[string]string
	Config      *util.RanchConfig
}

var dockerComposeTemplate = template.Must(template.New("docker-compose").Parse(`# generated by ranch
{{ range $name, $process := .Config.Processes }}
{{ $name }}:
  image: {{ $.ImageName }}
  command: /start {{ $name }}
  volumes:
    - /var/run/docker.sock:/var/run/docker.sock
  {{ if eq $name "web" }}
  labels:
    - convox.port.443.protocol=https
  ports:
    - 443:3000
  {{ end }}
  environment:
{{ if eq $name "web" }}{{printf "    - PORT=3000\n" }}{{ end }}
{{ range $k, $v := $.Environment }}{{ printf "    - %s=%s\n" $k $v }}{{ end }}
{{ end }}

run:
  image: {{ $.ImageName }}
  command: sh -c 'while true; do echo this process should not be running; sleep 300; done'
  volumes:
    - /var/run/docker.sock:/var/run/docker.sock
  labels:
{{ range $k, $v := $.Config.Cron }}{{ printf "    - convox.cron.%s=%s\n" $k $v }}{{ end }}
  environment:
{{ range $k, $v := $.Environment }}{{ printf "    - %s=%s\n" $k $v }}{{ end }}
`))

var procfileTemplate = template.Must(template.New("procfile").Parse(`# generated by ranch
{{ range $Name, $Process := .Config.Processes }}
{{ $Name }}: {{ $Process.Command }}
{{ end }}
`))

var deployCmd = &cobra.Command{
	Use:   "deploy",
	Short: "Deploy the application",
	RunE: func(cmd *cobra.Command, args []string) (err error) {

		appDir, err := util.AppDir(cmd)
		if err != nil {
			return err
		}

		appName, err := util.AppName(cmd)
		if err != nil {
			return err
		}

		config, err := util.LoadAppConfig(cmd)
		if err != nil {
			return err
		}

		if errors := util.RanchValidateConfig(config); len(errors) > 0 {
			for _, err := range errors {
				fmt.Println(err.Error())
			}
			return fmt.Errorf(".ranch.yaml did not validate")
		}

		appVersion, err := util.AppVersion(cmd)
		if err != nil {
			return err
		}

		imageNameWithTag := strings.Join([]string{config.ImageName, appVersion}, ":")

		exists, err := util.DockerImageExists(imageNameWithTag)
		if err != nil {
			return err
		} else if exists {
			fmt.Printf("%s docker image already exists in registry, skipping build.\n", imageNameWithTag)
		} else {
			if err = dockerBuildAndPush(appDir, imageNameWithTag, config); err != nil {
				return err
			}
		}

		// both code and conf are being deployed
		releaseId := strings.Join([]string{appVersion, appVersion}, "-")

		exists, err = util.RanchReleaseExists(config.AppName, releaseId)
		if err != nil {
			return err
		} else if exists {
			currentRelease, err := util.ConvoxCurrentVersion(config.AppName)
			if err != nil {
				return err
			}

			if currentRelease != releaseId {
				fmt.Printf("promoting existing release %s\n", releaseId)
				if err = util.ConvoxPromote(config.AppName, releaseId); err != nil {
					return err
				}

				time.Sleep(10 * time.Second) // wait for promote to apply

				if err = util.ConvoxWaitForStatus(config.AppName, "running"); err != nil {
					return err
				}
			} else {
				fmt.Printf("existing release %s is currently live, skipping promote.\n", releaseId)
			}
		} else {
			buildDir, err := ioutil.TempDir("", "ranch")
			if err != nil {
				return err
			}

			fmt.Println("using build directory", buildDir)

			var env map[string]string
			if config.EnvId != "" {
				plaintext, err := util.RanchGetSecret(appName, config.EnvId)
				if err != nil {
					return err
				}

				env, err = util.ParseEnv(plaintext)
				if err != nil {
					return err
				}
			}

			if err = generateDockerCompose(imageNameWithTag, config, env, buildDir); err != nil {
				return err
			}

			if err = convoxDeploy(appName, releaseId, buildDir); err != nil {
				return err
			}
		}

		if err = util.ConvoxScale(appName, config); err != nil {
			return err
		}

		return nil
	},
}

func init() {
	deployCmd.Flags().BoolVar(&Build, "build", true, "Build and push the Docker image")
	RootCmd.AddCommand(deployCmd)
}

func convoxDeploy(appName, releaseId, buildDir string) error {
	convoxReleaseId, err := util.ConvoxDeploy(appName, buildDir)

	if err != nil {
		return err
	}

	if err = util.RanchCreateRelease(appName, releaseId, convoxReleaseId); err != nil {
		return err
	}

	fmt.Printf("promoting release %s\n", releaseId)

	if err = util.ConvoxPromote(appName, releaseId); err != nil {
		return err
	}

	time.Sleep(10 * time.Second) // wait for promote to apply

	if err = util.ConvoxWaitForStatus(appName, "running"); err != nil {
		return err
	}

	return nil
}

func generateDockerCompose(imageName string, config *util.RanchConfig, env map[string]string, buildDir string) error {
	var out bytes.Buffer

	absoluteImageName, err := util.DockerResolveImageName(imageName)

	if err != nil {
		return err
	}

	err = dockerComposeTemplate.Execute(&out, composeTemplateVars{
		ImageName:   absoluteImageName,
		Environment: env,
		Config:      config,
	})

	if err != nil {
		return err
	}

	dockerCompose := path.Join(buildDir, "docker-compose.yml")
	err = ioutil.WriteFile(dockerCompose, out.Bytes(), 0644)

	if err != nil {
		return err
	}

	return nil
}

func dockerBuildAndPush(appDir, imageName string, config *util.RanchConfig) (err error) {

	env, err := util.EnvGet(config.AppName, config.EnvId)

	if err != nil {
		return err
	}

	dockerfile := path.Join(appDir, "Dockerfile")

	if _, err := os.Stat(dockerfile); os.IsNotExist(err) {

		var template *template.Template

		if _, err := os.Stat(path.Join(appDir, ".buildpacks")); os.IsNotExist(err) {
			template = nodejsDockerfileTemplate
		} else {
			template = cedarDockerfileTemplate
		}

		var out bytes.Buffer
		err = template.Execute(&out, struct{}{})

		if err != nil {
			return err
		}

		err = ioutil.WriteFile(dockerfile, out.Bytes(), 0644)

		if err != nil {
			return err
		}

		defer os.Remove(dockerfile) // cleanup
	} else {
		fmt.Println("WARNING: using existing Dockerfile")
	}

	procfile := path.Join(appDir, "Procfile")

	if _, err := os.Stat(procfile); os.IsNotExist(err) {
		var out bytes.Buffer
		err = procfileTemplate.Execute(&out, composeTemplateVars{
			ImageName: imageName,
			Config:    config,
		})

		if err != nil {
			return err
		}

		err = ioutil.WriteFile(procfile, out.Bytes(), 0644)

		if err != nil {
			return err
		}

		defer os.Remove(procfile) // cleanup
	} else {
		fmt.Println("WARNING: using existing Procfile")
	}

	err = util.DockerBuild(appDir, imageName, env)

	if err != nil {
		return err
	}

	err = util.DockerPush(imageName)

	if err != nil {
		return err
	}

	return nil
}
