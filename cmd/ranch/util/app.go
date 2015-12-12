package util

import (
	"path"
	"path/filepath"

	"github.com/spf13/cobra"
)

func AppName(cmd *cobra.Command) (string, error) {

	// use flag
	if app := cmd.Flag("app").Value.String(); app != "" {
		return app, nil
	}

	// fall back to directory name
	if abs, err := filepath.Abs("."); err != nil {
		return "", err
	} else {
		return path.Base(abs), nil
	}
}
