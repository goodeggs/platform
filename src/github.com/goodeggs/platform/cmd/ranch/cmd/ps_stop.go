package cmd

import (
	"fmt"

	"github.com/goodeggs/platform/cmd/ranch/util"
	"github.com/spf13/cobra"
)

var psStopCmd = &cobra.Command{
	Use:   "ps:stop <pid>",
	Short: "Stop a running process",
	RunE: func(cmd *cobra.Command, args []string) (err error) {
		appName, err := util.AppName(cmd)
		if err != nil {
			return err
		}

		if len(args) != 1 {
			cmd.Usage()
			return fmt.Errorf("usage")
		}

		pid := args[0]

		return util.ConvoxPsStop(appName, pid)
	},
}

func init() {
	RootCmd.AddCommand(psStopCmd)
}
