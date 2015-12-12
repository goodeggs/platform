package util

import (
	"os"
	"path"

	"github.com/spf13/cobra"
)

func AppDir(_ *cobra.Command) (string, error) {
	return os.Getwd()
}

func AppVersion(_ *cobra.Command) (string, error) {
	return "v1", nil
}

func AppName(cmd *cobra.Command) (string, error) {
	// use flag
	if app := cmd.Flag("app").Value.String(); app != "" {
		return app, nil
	}

	// fall back to directory name
	if appDir, err := AppDir(cmd); err != nil {
		return "", err
	} else {
		return path.Base(appDir), nil
	}
}
