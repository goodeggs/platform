package cmd

import (
	"fmt"
	"os"
	"sort"
	"strings"

	"github.com/goodeggs/platform/cmd/ranch/Godeps/_workspace/src/github.com/olekukonko/tablewriter"
	"github.com/goodeggs/platform/cmd/ranch/Godeps/_workspace/src/github.com/spf13/cobra"
	"github.com/goodeggs/platform/cmd/ranch/util"
)

var infoCmd = &cobra.Command{
	Use:   "info",
	Short: "Displays information about an application",
	RunE: func(cmd *cobra.Command, args []string) (err error) {
		appName, err := util.AppName(cmd)
		if err != nil {
			return err
		}

		app, err := util.RanchGetApp(appName)
		if err != nil {
			return err
		}

		formation, err := util.RanchGetFormation(appName)
		if err != nil {
			return err
		}

		ps := []string{}
		endpoints := []string{}

		for name, f := range formation {
			// skip the hidden run process
			if name == "run" {
				continue
			}

			ps = append(ps, name)

			for _, port := range f.Ports {
				endpoints = append(endpoints, fmt.Sprintf("%s:%d (%s)", f.Balancer, port, name))
			}
		}

		sort.Strings(ps)

		table := tablewriter.NewWriter(os.Stdout)
		table.SetBorder(false)
		table.SetColumnSeparator("")
		table.SetCenterSeparator("")
		table.SetAutoWrapText(false)
		table.Append([]string{"Name", appName})
		table.Append([]string{"Status", app.Status})
		table.Append([]string{"Release", app.Release})
		table.Append([]string{"Processes", strings.Join(ps, " ")})
		table.Append([]string{"Endpoints", strings.Join(endpoints, "\n")})
		table.Render()

		return nil
	},
}

func init() {
	RootCmd.AddCommand(infoCmd)
}
