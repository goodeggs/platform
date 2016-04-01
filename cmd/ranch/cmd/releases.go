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
	RunE: func(cmd *cobra.Command, args []string) (err error) {
		appName, err := util.AppName(cmd)
		if err != nil {
			return err
		}

		releases, err := util.ConvoxReleases(appName)
		if err != nil {
			return err
		}

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

		return nil
	},
}

func init() {
	RootCmd.AddCommand(releasesCmd)
}
