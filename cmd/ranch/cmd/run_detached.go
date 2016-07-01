package cmd

import (
	"fmt"
	"strings"

	"github.com/goodeggs/platform/cmd/ranch/Godeps/_workspace/src/github.com/spf13/cobra"
	"github.com/goodeggs/platform/cmd/ranch/util"
)

var runDetachedCmd = &cobra.Command{
	Use:   "run:detached <command>",
	Short: "Run a detached one-off command",
	RunE: func(cmd *cobra.Command, args []string) (err error) {
		if len(args) == 0 {
			cmd.Usage()
			return fmt.Errorf("must specify command")
		}

		appName, err := util.AppName(cmd)
		if err != nil {
			return err
		}

		process := "run"
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
