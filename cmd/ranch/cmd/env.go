package cmd

import (
	"fmt"

	"github.com/goodeggs/platform/cmd/ranch/Godeps/_workspace/src/github.com/spf13/cobra"
	"github.com/goodeggs/platform/cmd/ranch/util"
)

var envCmd = &cobra.Command{
	Use:   "env",
	Short: "Print the application environment",
	Run: func(cmd *cobra.Command, args []string) {
		config, err := util.LoadAppConfig(cmd)
		util.Check(err)

		appName, err := util.AppName(cmd)
		util.Check(err)

		if config.EnvId == "" {
			util.Die("your config does not contain an env_id")
		}

		plaintext, err := util.EcruGetSecret(appName, config.EnvId)
		util.Check(err)

		env, err := util.ParseEnv(plaintext)
		util.Check(err)

		for key, value := range env {
			fmt.Printf("%s=%s\n", key, value)
		}
	},
}

func init() {
	RootCmd.AddCommand(envCmd)
}
