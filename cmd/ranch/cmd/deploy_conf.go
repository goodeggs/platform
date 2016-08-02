package cmd

import (
	"fmt"
	"strings"

	"github.com/goodeggs/platform/cmd/ranch/Godeps/_workspace/src/github.com/spf13/cobra"
	"github.com/goodeggs/platform/cmd/ranch/util"
)

var deployConfCmd = &cobra.Command{
	Use:   "deploy:conf",
	Short: "Deploy only configuration changes, not code changes",
	RunE: func(cmd *cobra.Command, args []string) (err error) {
		config, err := util.LoadAppConfig(cmd)
		if err != nil {
			return err
		}

		if errors := util.RanchValidateConfig(config); len(errors) > 0 {
			for _, err := range errors {
				fmt.Println(err.Error())
			}
			return fmt.Errorf(".ranch.yaml did not validate")
		}

		appConfigPath, err := util.AppConfigPath(cmd)
		if err != nil {
			return err
		}

		appDir, err := util.AppDir(cmd)
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

		confSha, err := util.GitCurrentSha(appDir)
		if err != nil {
			return err
		}

		currentRelease, err := util.ConvoxCurrentVersion(config.AppName)
		if err != nil {
			return err
		}

		parts := strings.SplitN(currentRelease, "-", 2)
		codeSha := parts[0]

		return util.RanchDeploy(appDir, config, codeSha, confSha)
	},
}

func init() {
	RootCmd.AddCommand(deployConfCmd)
}
