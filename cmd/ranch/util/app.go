package util

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"

	"github.com/ghodss/yaml"
	"github.com/spf13/cobra"
)

type RanchConfig struct {
	Name    string `json:"name"`
	Version int    `json:"version"`
}

func AppConfigPath(cmd *cobra.Command) (string, error) {
	appDir, err := AppDir(cmd)
	if err != nil {
		return "", err
	}
	return path.Join(appDir, ".ranch.yaml"), nil
}

func LoadAppConfig(cmd *cobra.Command) (*RanchConfig, error) {
	filename, err := AppConfigPath(cmd)

	if err != nil {
		return nil, err
	}

	src, err := ioutil.ReadFile(filename)

	if err != nil {
		if os.IsNotExist(err) {
			return nil, fmt.Errorf(".ranch.yaml does not exist -- try `ranch init`")
		}
		return nil, err
	}

	if err != nil {
		return nil, err
	}

	var config RanchConfig
	yaml.Unmarshal(src, &config)

	if err != nil {
		return nil, err
	}

	return &config, nil
}

func AppDir(_ *cobra.Command) (string, error) {
	return os.Getwd()
}

func AppVersion(cmd *cobra.Command) (int, error) {
	config, err := LoadAppConfig(cmd)
	if err != nil {
		return -1, err
	}
	return config.Version, nil
}

func AppName(cmd *cobra.Command) (string, error) {
	// use flag
	if app := cmd.Flag("app").Value.String(); app != "" {
		return app, nil
	}

	// fall back to config
	if config, err := LoadAppConfig(cmd); err == nil {
		return config.Name, nil
	}

	// fall back to directory name
	if appDir, err := AppDir(cmd); err != nil {
		return "", err
	} else {
		return path.Base(appDir), nil
	}
}
