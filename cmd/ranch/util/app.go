package util

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/goodeggs/platform/cmd/ranch/Godeps/_workspace/src/github.com/spf13/cobra"
)

func AppConfigPath(cmd *cobra.Command) (string, error) {
	configFile, err := cmd.Flags().GetString("filename")
	if err != nil {
		return "", err
	}

	// use specified config file
	if configFile != "" {
		return filepath.EvalSymlinks(configFile)
	}

	// No filename was specified, scan for .ranch.*.yaml files
	if files, err := filepath.Glob(".ranch.*.yaml"); err != nil || len(files) > 0 {
		return "", fmt.Errorf("cannot infer config file path, use -f")
	}

	// fallback to default .ranch.yaml file
	if _, err := os.Stat(".ranch.yaml"); err == nil {
		return filepath.EvalSymlinks(".ranch.yaml")
	}

	return "", nil
}

func LoadAppConfig(cmd *cobra.Command) (*RanchConfig, error) {

	filename, err := AppConfigPath(cmd)
	if err != nil {
		return nil, err
	}

	// no ranchfiles
	if filename == "" {
		return nil, nil
	}

	return LoadRanchConfig(filename)
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
	app, err := cmd.Flags().GetString("app")
	if err != nil {
		return "", err
	}

	// use specified app name
	if app != "" {
		return app, nil
	}

	// fall back to name in config file
	config, err := LoadAppConfig(cmd)
	if err != nil {
		return "", err
	}
	if config != nil && config.AppName != "" {
		return config.AppName, nil
	}

	return "", fmt.Errorf("unable to infer app name, use -a")
}
