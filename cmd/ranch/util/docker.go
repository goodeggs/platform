package util

import (
	"fmt"
	"os"

	"github.com/fsouza/go-dockerclient"
)

func dockerClient() (*docker.Client, error) {
	if os.Getenv("DOCKER_HOST") != "" {
		return docker.NewClientFromEnv()
	} else {
		return docker.NewClient("unix:///var/run/docker.sock")
	}
}

func DockerBuild(appDir string, appName string, appVersion int) (string, error) {
	client, err := dockerClient()

	if err != nil {
		return "", err
	}

	imageName := fmt.Sprintf("%s/%s:v%d", "goodeggs", appName, appVersion)

	opts := docker.BuildImageOptions{
		Name:         imageName,
		OutputStream: os.Stdout,
		ContextDir:   appDir,
		Pull:         true,
	}

	err = client.BuildImage(opts)

	if err != nil {
		return "", err
	}

	return imageName, nil
}
