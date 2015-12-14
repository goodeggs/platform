package cmd

import (
	"os"

	"github.com/goodeggs/platform/cmd/ranch/util"
	"github.com/olekukonko/tablewriter"
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

		table := tablewriter.NewWriter(os.Stdout)
		table.SetBorder(false)
		table.SetColumnSeparator("")
		table.SetCenterSeparator("")
		table.SetAutoWrapText(false)
		table.Append([]string{"ID", "NAME", "RELEASE", "STARTED", "COMMAND"})

		for _, p := range ps {
			table.Append([]string{p.Id, p.Name, p.Release, util.HumanizeTime(p.Started), p.Command})
		}

		table.Render()
	},
}
