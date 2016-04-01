package cmd

import (
	"os"

	"github.com/goodeggs/platform/cmd/ranch/Godeps/_workspace/src/github.com/spf13/cobra"
	"github.com/goodeggs/platform/cmd/ranch/util"
)

var logsCmd = &cobra.Command{
	Use:   "logs",
	Short: "Tail the application logs",
	RunE: func(cmd *cobra.Command, args []string) (err error) {
		appName, err := util.AppName(cmd)
		if err != nil {
			return err
		}

		if err = util.ConvoxLogs(appName, os.Stdout); err != nil {
			return err
		}

		return nil
	},
}

func init() {
	RootCmd.AddCommand(logsCmd)
}
