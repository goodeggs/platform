package cmd

import (
	"fmt"

	"github.com/goodeggs/platform/cmd/ranch/util"
	"github.com/spf13/cobra"
)

var panicRollbackCmd = &cobra.Command{
	Use:   "panic:rollback <release>",
	Short: "Revert to a previous release",
	RunE: func(cmd *cobra.Command, args []string) (err error) {
		appName, err := util.AppName(cmd)
		if err != nil {
			return err
		}

		if len(args) == 0 {
			cmd.Usage()
			return fmt.Errorf("usage")
		}

		appVersion := args[0]

		if err = util.ConvoxPromote(appName, appVersion); err != nil {
			return err
		}

		fmt.Println("WARNING: this rollback only affects your code and environment -- scale may be inconsistent.")

		return nil
	},
}

func init() {
	RootCmd.AddCommand(panicRollbackCmd)
}
