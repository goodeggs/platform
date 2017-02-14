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
		config, err := util.LoadAppConfig(cmd)
		if err != nil {
			return err
		}

		errs := util.RanchValidateConfig(config)

		if len(errs) == 0 {
			fmt.Println("valid")
			return nil
		}

		for _, err := range errs {
			fmt.Println(err.Error())
		}
		return fmt.Errorf("ranch config had errors")
	},
}

func init() {
	RootCmd.AddCommand(lintCmd)
}
