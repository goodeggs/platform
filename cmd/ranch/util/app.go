package util

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"

	"github.com/goodeggs/platform/cmd/ranch/Godeps/_workspace/src/github.com/ghodss/yaml"
	"github.com/goodeggs/platform/cmd/ranch/Godeps/_workspace/src/github.com/spf13/cobra"
)

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

	var config RanchConfig
	yaml.Unmarshal(src, &config)

	if err != nil {
		return nil, err
	}

	for name, proc := range config.Processes {
		if proc.Instances > 0 && proc.Count == 0 {
			fmt.Printf("deprecated: rename `instances` to `count` in your .ranch.yaml for app '%s'\n", name)
			proc.Count = proc.Instances
		}
	}

	return &config, nil
}

func AppDir(_ *cobra.Command) (string, error) {
	wd, err := os.Getwd()

	if err != nil {
		return "", err
	}

	return filepath.EvalSymlinks(wd)
}

func AppVersion(cmd *cobra.Command) (string, error) {
	appDir, err := AppDir(cmd)

	if err != nil {
		return "", err
	}

	return GitCurrentSha(appDir)
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
