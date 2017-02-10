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
	configFile, err := cmd.Flags().GetString("filename")
	
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return "", err
	} else {
		return _appConfigPath(configFile)
	}
}

func _appConfigPath(configFile string) (string, error) {
	// use specified config file
	if configFile != "" {
		fmt.Printf("using config file %s\n", configFile)
		return filepath.EvalSymlinks(configFile)
	}
	
	// No filename was specified, scan for .ranch.*.yaml files
	files, err := filepath.Glob(".ranch.*.yaml")

	if err != nil {
		// scanning directory failed
		return "", fmt.Errorf("failed to scan directory for .ranch.*.yaml files")
	} else if len(files) >= 2 {
		// too many .ranch.*.yaml files exist, we don't want to be ambiguous!
		return "", fmt.Errorf("Multiple .ranch.*.yaml files exist, specify -f")
	} else {
		// fallback to default .ranch.yaml file
		return filepath.EvalSymlinks(".ranch.yaml")
	}
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
			config.Processes[name] = proc // write it back to the map
		}
	}

	if config.ImageName == "" {
		config.ImageName = config.AppName
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
	app := cmd.Flag("app").Value.String()
	return _appName(app, cmd)
}

func _appName(app string, cmd *cobra.Command) (string, error) {
	// use specified app name
	if app != "" {
		fmt.Printf("using app name %s\n", app)
		return app, nil
	}

	// No app name was specified, scan for .ranch.*.yaml files
	files, err := filepath.Glob(".ranch.*.yaml")

	if err != nil {
		// scanning directory failed
		return "", fmt.Errorf("failed to scan directory for .ranch.*.yaml files")
	} else if len(files) >= 2 {
		// too many .ranch.*.yaml files exist, we don't want to be ambiguous!
		return "", fmt.Errorf("Multiple .ranch.*.yaml files exist, specify -a")
	}

	// fall back to config from .ranch.yaml
	config, err := LoadAppConfig(cmd)
	if err == nil {
		return config.AppName, nil
	}

	// fall back to directory name
	if appDir, err := AppDir(cmd); err != nil {
		return "", err
	} else {
		return path.Base(appDir), nil
	}
}
