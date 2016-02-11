package cmd

import (
	"os"

	"github.com/goodeggs/platform/cmd/ranch/Godeps/_workspace/src/github.com/olekukonko/tablewriter"
	"github.com/goodeggs/platform/cmd/ranch/Godeps/_workspace/src/github.com/spf13/cobra"
	"github.com/goodeggs/platform/cmd/ranch/util"
)

var releasesCmd = &cobra.Command{
	Use:   "releases",
	Short: "List releases",
	Run: func(cmd *cobra.Command, args []string) {
		appName, err := util.AppName(cmd)
		util.Check(err)

		releases, err := util.ConvoxReleases(appName)
		util.Check(err)

		table := tablewriter.NewWriter(os.Stdout)
		table.SetBorder(false)
		table.SetColumnSeparator("")
		table.SetCenterSeparator("")
		table.SetAutoWrapText(false)
		table.Append([]string{"ID", "CREATED", "STATUS"})

		for _, r := range releases {
			table.Append([]string{r.Id, util.HumanizeTime(r.Created), r.Status})
		}

		table.Render()
	},
}

func init() {
	RootCmd.AddCommand(releasesCmd)
}
