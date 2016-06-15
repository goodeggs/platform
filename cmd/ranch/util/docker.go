package util

import (
	"encoding/json"
	"fmt"
	"net/url"
	"os"
	"path"
	"strings"

	"github.com/goodeggs/platform/cmd/ranch/Godeps/_workspace/src/github.com/fsouza/go-dockerclient"
	"github.com/goodeggs/platform/cmd/ranch/Godeps/_workspace/src/github.com/parnurzeal/gorequest"
	"github.com/goodeggs/platform/cmd/ranch/Godeps/_workspace/src/github.com/spf13/viper"
)

func dockerClient() (*docker.Client, error) {
	return docker.NewClientFromEnv()
}

func dockerRegistryUrl(pathname string) string {
	u, _ := url.Parse(viper.GetString("docker.registry.url"))
	u.Path = path.Join(u.Path, "v2", pathname)
	return u.String()
}

func dockerRegistryClient() (*gorequest.SuperAgent, error) {
	username := viper.GetString("docker.registry.username")
	password := viper.GetString("docker.registry.password")

	request := jsonClient().
		SetBasicAuth(username, password)

	return request, nil
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
	pathname := fmt.Sprintf("/%s/manifests/%s", imageName, tag)

	client, err := dockerRegistryClient()
	if err != nil {
		return false, err
	}

	resp, _, errs := client.Get(dockerRegistryUrl(pathname)).End()

	if len(errs) > 0 {
		return false, errs[0]
	} else if resp.StatusCode == 404 {
		return false, nil
	} else if resp.StatusCode == 200 {
		return true, nil
	}

	return false, fmt.Errorf("unexpected HTTP status code from docker registry: GET %s %d", pathname, resp.StatusCode)
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

	opts := docker.PushImageOptions{
		Name:         absoluteImageName,
		Tag:          tag,
		OutputStream: os.Stdout,
	}

	err = client.PushImage(opts, registryAuth())

	if err != nil {
		return err
	}

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
