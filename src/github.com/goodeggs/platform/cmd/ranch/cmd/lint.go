package cmd

import (
	"fmt"

	"github.com/goodeggs/platform/cmd/ranch/util"
	"github.com/spf13/cobra"
)

var lintCmd = &cobra.Command{
	Use:   "lint",
	Short: "Lint a ranch config",
	RunE: func(cmd *cobra.Command, args []string) (err error) {
		if _, err := util.LoadAppConfig(cmd); err != nil {
			return err
		}

		fmt.Println("valid")
		return nil
	},
}

func init() {
	RootCmd.AddCommand(lintCmd)
}
