package cmd

import (
	"os"

	"github.com/goodeggs/platform/cmd/ranch/Godeps/_workspace/src/github.com/spf13/cobra"
	"github.com/goodeggs/platform/cmd/ranch/util"
)

var logsCmd = &cobra.Command{
	Use:   "logs",
	Short: "Tail the application logs",
	Run: func(cmd *cobra.Command, args []string) {
		appName, err := util.AppName(cmd)
		util.Check(err)

		err = util.ConvoxLogs(appName, os.Stdout)
		util.Check(err)
	},
}

func init() {
	RootCmd.AddCommand(logsCmd)
}
