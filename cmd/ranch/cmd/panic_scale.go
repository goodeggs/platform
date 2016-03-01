package cmd

import (
	"os"

	"github.com/goodeggs/platform/cmd/ranch/Godeps/_workspace/src/github.com/spf13/cobra"
	"github.com/goodeggs/platform/cmd/ranch/util"
)

var memory int
var count int

var panicScaleCmd = &cobra.Command{
	Use:   "panic:scale <process>",
	Short: "Adjust the scale of a process.",
	Run: func(cmd *cobra.Command, args []string) {
		appName, err := util.AppName(cmd)
		util.Check(err)

		if len(args) != 1 {
			cmd.Usage()
			os.Exit(1)
		}

		process := args[0]

		err = util.ConvoxScaleProcess(appName, process, count, memory)
		util.Check(err)
	},
}

func init() {
	panicScaleCmd.Flags().IntVar(&count, "count", -1, "Instance count")
	panicScaleCmd.Flags().IntVar(&memory, "memory", -1, "Memory in MB")
	RootCmd.AddCommand(panicScaleCmd)
}
