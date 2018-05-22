package util

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path"
	"regexp"
	"strings"

	"github.com/fatih/color"
	"github.com/fsouza/go-dockerclient"
	"github.com/heroku/docker-registry-client/registry"
	"github.com/spf13/viper"
)

type jsonMessage struct {
	Status   string `json:"status,omitempty"`
	Progress string `json:"progress,omitempty"`
	Error    string `json:"error,omitempty"`
	Stream   string `json:"stream,omitempty"`
}

func dockerClient() (*docker.Client, error) {
	return docker.NewClientFromEnv()
}

func dockerRegistryUrl(pathname string) string {
	u, _ := url.Parse(viper.GetString("docker.registry.url"))
	u.Path = path.Join(u.Path, "v1", pathname)
	return u.String()
}

func dockerRegistryClient() (*registry.Registry, error) {
	u := viper.GetString("docker.registry.url")
	username := viper.GetString("docker.registry.username")
	password := viper.GetString("docker.registry.password")

	url := strings.TrimSuffix(u, "/")
	transport := registry.WrapTransport(http.DefaultTransport, url, username, password)
	r := &registry.Registry{
		URL: url,
		Client: &http.Client{
			Transport: transport,
		},
		Logf: registry.Quiet,
	}

	if err := r.Ping(); err != nil {
		return nil, err
	}

	return r, nil
}

func DockerResolveImageName(imageName string) (string, error) {
	host := viper.GetString("docker.registry.url")

	registryUrl, err := url.Parse(host)

	if err != nil {
		return "", err
	}

	hostname := registryUrl.Host

	org := viper.GetString("docker.registry.org")

	if org == "" {
		return strings.Join([]string{hostname, imageName}, "/"), nil
	}

	return strings.Join([]string{hostname, org, imageName}, "/"), nil
}

func registryAuth() docker.AuthConfiguration {
	serverAddress := path.Join(viper.GetString("docker.registry.url"), "v1")
	return docker.AuthConfiguration{
		Email:         viper.GetString("docker.registry.email"),
		Username:      viper.GetString("docker.registry.username"),
		Password:      viper.GetString("docker.registry.password"),
		ServerAddress: serverAddress,
	}
}

func DockerImageExists(imageNameWithTag string) (bool, error) {
	parts := strings.Split(imageNameWithTag, ":")
	imageName, tag := parts[0], parts[1]

	if org := viper.GetString("docker.registry.org"); org != "" {
		imageName = strings.Join([]string{org, imageName}, "/")
	}

	client, err := dockerRegistryClient()
	if err != nil {
		return false, err
	}

	manifest, err := client.Manifest(imageName, tag)
	if err != nil {
		if strings.Contains(err.Error(), "status=404") {
			return false, nil
		}
		return false, err
	} else if manifest == nil {
		return false, nil
	}

	return true, nil
}

func DockerPull(imageNameWithTag string) error {
	parts := strings.Split(imageNameWithTag, ":")
	imageName, tag := parts[0], parts[1]

	client, err := dockerClient()

	if err != nil {
		return err
	}

	absoluteImageName, err := DockerResolveImageName(imageName)

	if err != nil {
		return err
	}

	silencer := newDockerSilencer()

	opts := docker.PullImageOptions{
		Repository:    absoluteImageName,
		Tag:           tag,
		OutputStream:  silencer.Writer(),
		RawJSONStream: true,
	}

	fmt.Printf("🐮  Pulling docker image `%s:%s'... ", absoluteImageName, tag)

	if err = client.PullImage(opts, registryAuth()); err != nil {
		fmt.Println("✘")
		return err
	}

	if err = silencer.Finalize(); err != nil {
		fmt.Println("✘")
		return err
	}

	fmt.Println("✔")
	return nil
}

type dockerSilencer struct {
	r   *io.PipeReader
	w   *io.PipeWriter
	res chan error
}

func newDockerSilencer() *dockerSilencer {
	r, w := io.Pipe()
	res := make(chan error)
	h := dockerSilencer{r, w, res}

	go func() {
		dec := json.NewDecoder(r)
		for {
			var m jsonMessage
			if err := dec.Decode(&m); err == io.EOF {
				h.res <- nil
				break
			} else if err != nil {
				h.res <- err
				break
			}

			if Verbose && m.Stream != "" {
				fmt.Print(m.Stream)
			} else if Verbose && m.Progress != "" {
				fmt.Printf("%s %s\r", m.Status, m.Progress)
			} else if m.Error != "" {
				h.res <- errors.New(m.Error)
				break
			}

			if Verbose && m.Status != "" {
				fmt.Println(m.Status)
			}
		}
	}()

	return &h
}

func (h *dockerSilencer) Writer() *io.PipeWriter {
	return h.w
}

func (h *dockerSilencer) Finalize() error {
	h.w.Close()
	return <-h.res
}

func DockerPush(imageNameWithTag string) error {
	parts := strings.Split(imageNameWithTag, ":")
	imageName, tag := parts[0], parts[1]

	client, err := dockerClient()

	if err != nil {
		return err
	}

	absoluteImageName, err := DockerResolveImageName(imageName)

	if err != nil {
		return err
	}

	silencer := newDockerSilencer()

	opts := docker.PushImageOptions{
		Name:          absoluteImageName,
		Tag:           tag,
		OutputStream:  silencer.Writer(),
		RawJSONStream: true,
	}

	fmt.Printf("🐮  Pushing docker image %s:%s... ", absoluteImageName, tag)

	if err = client.PushImage(opts, registryAuth()); err != nil {
		fmt.Println("✘")
		return err
	}

	if err = silencer.Finalize(); err != nil {
		fmt.Println("✘")
		return err
	}

	fmt.Println("✔")
	return nil
}

func parseRepositoryTag(imageName string) (string, string) {
	parts := strings.SplitN(imageName, ":", 2)
	return parts[0], parts[1]
}

func DockerTag(imageName, tag string) error {
	client, err := dockerClient()

	if err != nil {
		return err
	}

	absoluteImageName, err := DockerResolveImageName(imageName)
	if err != nil {
		return err
	}

	repo, _ := parseRepositoryTag(absoluteImageName)

	opts := docker.TagImageOptions{
		Repo: repo,
		Tag:  tag,
	}

	fmt.Printf("🐮  Tagging docker image %s as %s... ", absoluteImageName, tag)

	if err = client.TagImage(absoluteImageName, opts); err != nil {
		fmt.Println("✘")
		return err
	}

	fmt.Println("✔")
	return nil
}

func DockerBuild(appDir string, imageName string, buildEnv map[string]string) error {
	client, err := dockerClient()

	if err != nil {
		return err
	}

	absoluteImageName, err := DockerResolveImageName(imageName)

	if err != nil {
		return err
	}

	jsonEnvStr, err := json.Marshal(buildEnv)

	if err != nil {
		return err
	}

	r, w := io.Pipe()

	go func() {
		stripColors := regexp.MustCompile("\\x1b\\[[0-9;]*[mG]")
		scanner := bufio.NewScanner(r)
		for scanner.Scan() {
			text := scanner.Text()
			fmt.Printf("%s %s\n", color.WhiteString("docker build |"), stripColors.ReplaceAllString(text, ""))
		}
	}()

	buildArgs := make([]docker.BuildArg, 1)
	buildArgs[0] = docker.BuildArg{Name: "RANCH_BUILD_ENV", Value: string(jsonEnvStr)}

	opts := docker.BuildImageOptions{
		Name:         absoluteImageName,
		OutputStream: w,
		ContextDir:   appDir,
		Pull:         true,
		BuildArgs:    buildArgs,
	}

	auth, err := docker.NewAuthConfigurationsFromDockerCfg()

	if err != nil && !os.IsNotExist(err) {
		return err
	}

	if auth != nil {
		opts.AuthConfigs = *auth
	}

	return client.BuildImage(opts)
}
