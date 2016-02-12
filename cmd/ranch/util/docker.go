package util

import (
	"fmt"
	"os"
	"strings"

	"github.com/goodeggs/platform/cmd/ranch/Godeps/_workspace/src/github.com/fsouza/go-dockerclient"
	"github.com/goodeggs/platform/cmd/ranch/Godeps/_workspace/src/github.com/spf13/viper"
)

func dockerClient() (*docker.Client, error) {
	return docker.NewClientFromEnv()
}

func DockerRegistry() string {
	return fmt.Sprintf("%s:5000", viper.GetString("convox.host"))
}

func registryAuth() docker.AuthConfiguration {
	return docker.AuthConfiguration{
		Email:         "user@convox.io",
		Username:      "convox",
		Password:      viper.GetString("convox.password"),
		ServerAddress: DockerRegistry(),
	}
}

func DockerPush(imageName string) error {
	parts := strings.Split(imageName, ":")
	name, tag := parts[0], parts[1]

	client, err := dockerClient()

	if err != nil {
		return err
	}

	opts := docker.PushImageOptions{
		Name:         fmt.Sprintf("%s/%s", DockerRegistry(), name),
		Tag:          tag,
		Registry:     DockerRegistry(), // deprecated see https://github.com/fsouza/go-dockerclient/issues/377
		OutputStream: os.Stdout,
	}

	err = client.PushImage(opts, registryAuth())

	if err != nil {
		return err
	}

	return nil
}

func DockerImageName(appName string, sha string) string {
	return fmt.Sprintf("%s:%s", appName, sha)
}

func DockerBuild(appDir string, imageName string) error {
	client, err := dockerClient()

	if err != nil {
		return err
	}

	opts := docker.BuildImageOptions{
		Name:         fmt.Sprintf("%s/%s", DockerRegistry(), imageName),
		OutputStream: os.Stdout,
		ContextDir:   appDir,
		Pull:         true,
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
