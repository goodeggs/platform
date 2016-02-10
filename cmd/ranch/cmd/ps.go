package cmd

import (
	"os"

	"github.com/goodeggs/platform/cmd/ranch/Godeps/_workspace/src/github.com/olekukonko/tablewriter"
	"github.com/goodeggs/platform/cmd/ranch/Godeps/_workspace/src/github.com/spf13/cobra"
	"github.com/goodeggs/platform/cmd/ranch/util"
)

func init() {
	RootCmd.AddCommand(psCmd)
}

var psCmd = &cobra.Command{
	Use:   "ps",
	Short: "List an app's processes",
	Run: func(cmd *cobra.Command, args []string) {

		appName, err := util.AppName(cmd)
		util.Check(err)

		ps, err := util.ConvoxPs(appName)
		util.Check(err)

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
