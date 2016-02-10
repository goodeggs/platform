package cmd

import (
	"strings"

	"github.com/goodeggs/platform/cmd/ranch/Godeps/_workspace/src/github.com/spf13/cobra"
	"github.com/goodeggs/platform/cmd/ranch/util"
)

var runDetachedCmd = &cobra.Command{
	Use:   "run:detached",
	Short: "Run a detached one-off command",
	Run: func(cmd *cobra.Command, args []string) {
		appName, err := util.AppName(cmd)
		util.Check(err)

		process := "web"
		command := strings.Join(args, " ")

		err = util.ConvoxRunDetached(appName, process, command)
		util.Check(err)
	},
}

func init() {
	RootCmd.AddCommand(runDetachedCmd)
}
