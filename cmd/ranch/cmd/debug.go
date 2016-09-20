package cmd

import (
	"fmt"

	"github.com/goodeggs/platform/cmd/ranch/Godeps/_workspace/src/github.com/spf13/cobra"
	"github.com/goodeggs/platform/cmd/ranch/util"
)

var debugCmd = &cobra.Command{
	Use:   "debug",
	Short: "Debug",
	RunE: func(cmd *cobra.Command, args []string) (err error) {
		config, err := util.LoadAppConfig(cmd)
		if err != nil {
			return err
		}

		fmt.Print("\n\n------------ docker-compose.yml ------------\n\n")

		content, err := util.GenerateDockerCompose("FIXME", config)
		if err != nil {
			return err
		}

		fmt.Print(string(content))

		return nil
	},
}

func init() {
	RootCmd.AddCommand(debugCmd)
}
