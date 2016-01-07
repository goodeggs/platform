package util

import (
	"fmt"
	"os"
	"strings"

	"github.com/fsouza/go-dockerclient"
)

func dockerClient() (*docker.Client, error) {
	if os.Getenv("DOCKER_HOST") != "" {
		return docker.NewClientFromEnv()
	} else {
		return docker.NewClient("unix:///var/run/docker.sock")
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
		Name:         name,
		Tag:          tag,
		Registry:     "https://index.docker.io/v1/",
		OutputStream: os.Stdout,
	}

	auths, err := docker.NewAuthConfigurationsFromDockerCfg()

	if err != nil {
		return err
	}

	err = client.PushImage(opts, auths.Configs["https://index.docker.io/v1/"])

	if err != nil {
		return err
	}

	return nil
}

func DockerImageName(appName string, appVersion int) string {
	return fmt.Sprintf("%s/%s:v%d", "goodeggs", appName, appVersion)
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
