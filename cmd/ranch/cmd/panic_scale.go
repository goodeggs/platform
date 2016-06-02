package cmd

import (
	"fmt"

	"github.com/goodeggs/platform/cmd/ranch/Godeps/_workspace/src/github.com/spf13/cobra"
	"github.com/goodeggs/platform/cmd/ranch/util"
)

var memory int
var count int

var panicScaleCmd = &cobra.Command{
	Use:   "panic:scale <process>",
	Short: "Adjust the scale of a process",
	RunE: func(cmd *cobra.Command, args []string) (err error) {
		appName, err := util.AppName(cmd)
		if err != nil {
			return err
		}

		if len(args) != 1 {
			cmd.Usage()
			return fmt.Errorf("usage")
		}

		process := args[0]

		if err = util.ConvoxScaleProcess(appName, process, count, memory); err != nil {
			return err
		}

		return nil
	},
}

func init() {
	panicScaleCmd.Flags().IntVar(&count, "count", -1, "Instance count")
	panicScaleCmd.Flags().IntVar(&memory, "memory", -1, "Memory in MB")
	RootCmd.AddCommand(panicScaleCmd)
}
