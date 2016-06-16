package cmd

import (
	"fmt"

	"github.com/goodeggs/platform/cmd/ranch/Godeps/_workspace/src/github.com/spf13/cobra"
	"github.com/goodeggs/platform/cmd/ranch/util"
)

var envCmd = &cobra.Command{
	Use:   "env",
	Short: "Print the application environment",
	RunE: func(cmd *cobra.Command, args []string) (err error) {
		config, err := util.LoadAppConfig(cmd)
		if err != nil {
			return err
		}

		appName, err := util.AppName(cmd)
		if err != nil {
			return err
		}

		if config.EnvId == "" {
			return fmt.Errorf("your config does not contain an env_id")
		}

		plaintext, err := util.RanchGetSecret(appName, config.EnvId)
		if err != nil {
			return err
		}

		env, err := util.ParseEnv(plaintext)
		if err != nil {
			return err
		}

		for key, value := range env {
			fmt.Printf("%s=%s\n", key, value)
		}

		return nil
	},
}

func init() {
	RootCmd.AddCommand(envCmd)
}
