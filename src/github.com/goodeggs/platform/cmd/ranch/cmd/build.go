package cmd

import (
	"fmt"

	"github.com/goodeggs/platform/cmd/ranch/util"
	"github.com/spf13/cobra"
)

var buildCmd = &cobra.Command{
	Use:   "build",
	Short: "Build the application",
	RunE: func(cmd *cobra.Command, args []string) (err error) {

		config, err := util.LoadAppConfig(cmd)
		if err != nil {
			return err
		}

		appDir, err := util.AppDir(cmd)
		if err != nil {
			return err
		}

		appConfigPath, err := util.AppConfigPath(cmd)
		if err != nil {
			return err
		}

		isClean, err := util.GitFileIsClean(appDir, appConfigPath)
		if err != nil {
			return err
		}
		if !isClean {
			return fmt.Errorf("your ranch config file %s must be committed before deploying.", appConfigPath)
		}

		appSha, err := util.GitCurrentSha(appDir)
		if err != nil {
			return err
		}

		return util.DockerBuildAndPush(appDir, config.ImageName, appSha, config)
	},
}

func init() {
	RootCmd.AddCommand(buildCmd)
}
