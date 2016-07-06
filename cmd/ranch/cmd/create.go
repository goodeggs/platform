package cmd

import (
	"fmt"
	"time"

	"github.com/goodeggs/platform/cmd/ranch/Godeps/_workspace/src/github.com/spf13/cobra"
	"github.com/goodeggs/platform/cmd/ranch/util"
)

var createCmd = &cobra.Command{
	Use:   "create <app name>",
	Short: "Create a new application",
	RunE: func(cmd *cobra.Command, args []string) (err error) {
		if len(args) != 1 {
			cmd.Usage()
			return fmt.Errorf("expected 1 argument (app name)")
		}

		appName := args[0]

		if err = util.RanchCreateApp(appName); err != nil {
			return err
		}

		fmt.Println("waiting 10s for app create")
		time.Sleep(10 * time.Second) // wait for app create to happen

		return util.ConvoxWaitForStatus(appName, "running")
	},
}

func init() {
	RootCmd.AddCommand(createCmd)
}
