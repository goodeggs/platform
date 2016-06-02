package cmd

import (
	"fmt"

	"github.com/goodeggs/platform/cmd/ranch/Godeps/_workspace/src/github.com/spf13/cobra"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Prints the ranch CLI version",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println(Version)
	},
}

func init() {
	RootCmd.AddCommand(versionCmd)
}
