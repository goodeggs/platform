package util

import (
	"fmt"
	"os"
	"strings"

	"github.com/fsouza/go-dockerclient"
	"github.com/spf13/viper"
)

func dockerClient() (*docker.Client, error) {
	if os.Getenv("DOCKER_HOST") != "" {
		return docker.NewClientFromEnv()
	} else {
		return docker.NewClient("unix:///var/run/docker.sock")
	}
}

func registry() string {
	return fmt.Sprintf("%s:5000", viper.GetString("convox.host"))
}

func registryAuth() docker.AuthConfiguration {
	return docker.AuthConfiguration{
		Username:      "convox",
		Password:      viper.GetString("convox.password"),
		ServerAddress: registry(),
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
		Name:         fmt.Sprintf("https://%s/%s", registry(), name),
		Tag:          tag,
		Registry:     registry(), // deprecated see https://github.com/fsouza/go-dockerclient/issues/377
		OutputStream: os.Stdout,
	}

	err = client.PushImage(opts, registryAuth())

	if err != nil {
		return err
	}

	return nil
}

func DockerImageName(appName string, sha string) string {
	return fmt.Sprintf("%s/%s:%s", "goodeggs", appName, sha)
}

func DockerBuild(appDir string, imageName string) error {
	client, err := dockerClient()

	if err != nil {
		return err
	}

	auth, err := docker.NewAuthConfigurationsFromDockerCfg()

	if err != nil {
		return err
	}

	opts := docker.BuildImageOptions{
		Name:         imageName,
		OutputStream: os.Stdout,
		ContextDir:   appDir,
		Pull:         true,
		AuthConfigs:  *auth,
	}

	return client.BuildImage(opts)
}
