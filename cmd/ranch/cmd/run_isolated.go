package cmd

import (
	"fmt"
	"strings"

	"github.com/goodeggs/platform/cmd/ranch/util"
	"github.com/spf13/cobra"
)

var instanceType string
var debug bool

var runIsolatedCmd = &cobra.Command{
	Use:   "run:isolated -t <instance type> <command>",
	Short: "Run a one-off command on an isolated EC2 box",
	RunE: func(cmd *cobra.Command, args []string) (err error) {
		if instanceType == "" {
			cmd.Usage()
			return fmt.Errorf("must specify instance type")
		}

		if len(args) == 0 {
			cmd.Usage()
			return fmt.Errorf("must specify command")
		}

		appName, err := util.AppName(cmd)
		if err != nil {
			return err
		}

		command := strings.Join(args, " ")

		instanceId, err := util.RanchRunIsolated(appName, instanceType, debug, command)
		if err != nil {
			return err
		}

		fmt.Println(instanceId)
		return nil
	},
}

func init() {
	runIsolatedCmd.Flags().StringVarP(&instanceType, "instance-type", "t", "", "EC2 Instance Type")
	runIsolatedCmd.Flags().BoolVarP(&debug, "debug", "d", false, "Debug (leave instance running)")
	RootCmd.AddCommand(runIsolatedCmd)
}
