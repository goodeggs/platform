package util

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/url"
	"os"
	"path"
	"strings"

	"github.com/goodeggs/platform/cmd/ranch/Godeps/_workspace/src/github.com/fsouza/go-dockerclient"
	"github.com/goodeggs/platform/cmd/ranch/Godeps/_workspace/src/github.com/heroku/docker-registry-client/registry"
	"github.com/goodeggs/platform/cmd/ranch/Godeps/_workspace/src/github.com/spf13/viper"
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

	return registry.New(u, username, password)
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

	silencer := NewDockerSilencer()

	opts := docker.PullImageOptions{
		Repository:    absoluteImageName,
		Tag:           tag,
		OutputStream:  silencer.Writer(),
		RawJSONStream: true,
	}

	fmt.Printf("Pulling docker image %q:", absoluteImageName)

	if err = client.PullImage(opts, registryAuth()); err != nil {
		fmt.Println("ERROR")
		return err
	}

	if err = silencer.Finalize(); err != nil {
		fmt.Println("ERROR")
		return err
	}

	fmt.Println("DONE")
	return nil
}

type DockerSilencer struct {
	r   *io.PipeReader
	w   *io.PipeWriter
	res chan error
}

func NewDockerSilencer() *DockerSilencer {
	r, w := io.Pipe()
	res := make(chan error)
	h := DockerSilencer{r, w, res}

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
			} else if m.Error != "" {
				h.res <- errors.New(m.Error)
				break
			}
		}
	}()

	return &h
}

func (h *DockerSilencer) Writer() *io.PipeWriter {
	return h.w
}

func (h *DockerSilencer) Finalize() error {
	h.w.Close()
	close(h.res)
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

	silencer := NewDockerSilencer()

	opts := docker.PushImageOptions{
		Name:          absoluteImageName,
		Tag:           tag,
		OutputStream:  silencer.Writer(),
		RawJSONStream: true,
	}

	fmt.Printf("Pushing docker image %s: ", absoluteImageName)

	if err = client.PushImage(opts, registryAuth()); err != nil {
		fmt.Println("ERROR")
		return err
	}

	if err = silencer.Finalize(); err != nil {
		fmt.Println("ERROR")
		return err
	}

	fmt.Println("DONE")
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

	return client.TagImage(absoluteImageName, opts)
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

	buildArgs := make([]docker.BuildArg, 1)
	buildArgs[0] = docker.BuildArg{Name: "RANCH_BUILD_ENV", Value: string(jsonEnvStr)}

	opts := docker.BuildImageOptions{
		Name:         absoluteImageName,
		OutputStream: os.Stdout,
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
