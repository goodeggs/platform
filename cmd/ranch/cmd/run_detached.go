package cmd

import (
	"strings"

	"github.com/goodeggs/platform/cmd/ranch/Godeps/_workspace/src/github.com/spf13/cobra"
	"github.com/goodeggs/platform/cmd/ranch/util"
)

var runDetachedCmd = &cobra.Command{
	Use:   "run:detached",
	Short: "Run a detached one-off command",
	RunE: func(cmd *cobra.Command, args []string) (err error) {
		appName, err := util.AppName(cmd)
		if err != nil {
			return err
		}

		process := "web"
		command := strings.Join(args, " ")

		if err = util.ConvoxRunDetached(appName, process, command); err != nil {
			return err
		}

		return nil
	},
}

func init() {
	RootCmd.AddCommand(runDetachedCmd)
}
