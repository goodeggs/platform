package cmd

import (
	"fmt"
	"os"

	"github.com/goodeggs/platform/cmd/ranch/Godeps/_workspace/src/github.com/spf13/cobra"
	"github.com/goodeggs/platform/cmd/ranch/util"
)

var panicRollbackCmd = &cobra.Command{
	Use:   "panic:rollback <release>",
	Short: "Revert to a previous release.",
	Run: func(cmd *cobra.Command, args []string) {
		appName, err := util.AppName(cmd)
		util.Check(err)

		if len(args) == 0 {
			cmd.Usage()
			os.Exit(1)
		}

		appVersion := args[0]

		err = util.ConvoxPromote(appName, appVersion)
		util.Check(err)

		err = util.ConvoxWaitForStatus(appName, "running")
		util.Check(err)

		fmt.Println("WARNING: this rollback only affects your code and environment -- scale may be inconsistent.")
	},
}

func init() {
	RootCmd.AddCommand(panicRollbackCmd)
}
