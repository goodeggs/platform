package cmd

import (
	"fmt"

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

		app, err := util.AppName(cmd)

		if err != nil {
			util.Error(err)
			return
		}

		ps, err := util.Convox().GetProcesses(app, false)

		if err != nil {
			util.Error(err)
			return
		}

		for _, p := range ps {
			fmt.Println(p.Id, p.Name, p.Release, util.HumanizeTime(p.Started), p.Command)
		}

	},
}
