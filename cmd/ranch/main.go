package main

import (
	"fmt"

	"github.com/goodeggs/platform/cmd/ranch/Godeps/_workspace/src/github.com/spf13/cobra"
	"github.com/goodeggs/platform/cmd/ranch/cmd"
)

var VERSION = "dev"

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Prints the ranch CLI version",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println(VERSION)
	},
}

func main() {
	cmd.RootCmd.AddCommand(versionCmd)
	cmd.Execute()
}
