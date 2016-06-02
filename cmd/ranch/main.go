package main

import (
	"fmt"
	"os"

	"github.com/goodeggs/platform/cmd/ranch/cmd"
)

var VERSION = "dev"

func init() {
	cmd.Version = VERSION
}

func main() {
	if err := cmd.RootCmd.Execute(); err != nil {
		if err.Error() == "pflag: help requested" {
			cmd.RootCmd.Usage()
		} else {
			fmt.Printf("Error: %s\n", err.Error())
			os.Exit(1)
		}
	}
}
