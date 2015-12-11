package cmd

import (
	"github.com/goodeggs/platform/cmd/ranch/util"
	"github.com/spf13/cobra"
)

func init() {
	RootCmd.AddCommand(psCmd)
}

var psCmd = &cobra.Command{
	Use:   "ps",
	Short: "List an app's processes",
	Run: func(cmd *cobra.Command, args []string) {
		args = append([]string{"ps"}, args...)
		util.Convox(args...)
	},
}
