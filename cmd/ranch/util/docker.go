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

func DockerPush(imageName string) error {
	fmt.Println("TODO: docker push $imageName")
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

	opts := docker.BuildImageOptions{
		Name:         imageName,
		OutputStream: os.Stdout,
		ContextDir:   appDir,
		Pull:         true,
	}

	return client.BuildImage(opts)
}
