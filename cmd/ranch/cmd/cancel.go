package cmd

import (
	"fmt"
	"time"

	"github.com/goodeggs/platform/cmd/ranch/util"
	"github.com/spf13/cobra"
)

var cancelCmd = &cobra.Command{
	Use:   "cancel <app name>",
	Short: "Cancel a change to an app",
	RunE: func(cmd *cobra.Command, args []string) (err error) {
		if len(args) != 1 {
			cmd.Usage()
			return fmt.Errorf("expected 1 argument (app name)")
		}

		appName := args[0]

		if err = util.ConvoxCancel(appName); err != nil {
			return err
		}

		fmt.Println("waiting 10s for app update cancel")
		time.Sleep(10 * time.Second) // wait for app cancel to happen

		return util.ConvoxWaitForStatus(appName, "running")
	},
}

func init() {
	RootCmd.AddCommand(cancelCmd)
}
