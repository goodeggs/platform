package cmd

import (
	"fmt"

	"github.com/goodeggs/platform/cmd/ranch/util"
	"github.com/spf13/cobra"
)

var deployCmd = &cobra.Command{
	Use:   "deploy",
	Short: "Deploy the application",
	RunE: func(cmd *cobra.Command, args []string) (err error) {

		appDir, err := util.AppDir(cmd)
		if err != nil {
			return err
		}

		config, err := util.LoadAppConfig(cmd)
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

		codeSha, err := util.GitCurrentSha(appDir)
		if err != nil {
			return err
		}
		confSha := codeSha // same revision

		return util.RanchDeploy(appDir, config, codeSha, confSha)
	},
}

func init() {
	RootCmd.AddCommand(deployCmd)
}
