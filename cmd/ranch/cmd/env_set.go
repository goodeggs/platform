package cmd

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/goodeggs/platform/cmd/ranch/util"
	"github.com/imdario/mergo"
	"github.com/spf13/cobra"
)

var envSetCmd = &cobra.Command{
	Use:   "env:set KEY=VALUE KEY2='VALUE 2'",
	Short: "Set environment variables",
	RunE: func(cmd *cobra.Command, args []string) (err error) {
		appDir, err := util.AppDir(cmd)
		if err != nil {
			return err
		}

		clean, err := util.GitIsClean(appDir)
		if err != nil {
			return err
		}

		if !clean {
			return fmt.Errorf("git working directory not clean.")
		}

		newEnv, err := readEnvChanges(args)
		if err != nil {
			return err
		}

		var updatedKeys []string
		for k, _ := range newEnv {
			updatedKeys = append(updatedKeys, k)
		}

		configPath, err := util.AppConfigPath(cmd)
		if err != nil {
			return err
		}

		config, err := util.LoadAppConfig(cmd)
		if err != nil {
			return err
		}

		if len(config.Env) > 0 {
			return fmt.Errorf("env:set is deprecated")
		}

		appName, err := util.AppName(cmd)
		if err != nil {
			return err
		}

		oldEnv, err := util.EnvGet(appName, config.EnvId)
		if err != nil {
			return err
		}

		if err = mergo.Merge(&newEnv, oldEnv); err != nil {
			return err
		}

		data := ""
		for k, v := range newEnv {
			data += fmt.Sprintf("%s=%s\n", k, v)
		}

		envId, err := util.RanchCreateSecret(appName, data)
		if err != nil {
			return err
		}

		if err = util.RanchUpdateEnvId(configPath, envId); err != nil {
			return err
		}

		if err = util.GitAdd(appDir, configPath); err != nil {
			return err
		}

		message := fmt.Sprintf("set env %s\n\n[ci skip]", strings.Join(updatedKeys, ","))

		if err = util.GitCommit(appDir, message); err != nil {
			return err
		}

		sha, err := util.GitCurrentSha(appDir)
		if err != nil {
			return err
		}

		fmt.Printf("[%s] %s\n", sha, message)
		fmt.Println("NOTE: you must deploy to apply this change, or you can use `ranch deploy:conf -f " + configPath + "` to apply it to the active release.")

		return nil
	},
}

func readEnvChanges(args []string) (env map[string]string, err error) {

	data := ""

	stat, err := os.Stdin.Stat()
	if err != nil {
		return nil, err
	}

	if (stat.Mode() & os.ModeCharDevice) == 0 {
		in, err := ioutil.ReadAll(os.Stdin)
		if err != nil {
			return nil, err
		}

		data += string(in)
	}

	for _, value := range args {
		data += fmt.Sprintf("%s\n", value)
	}

	env, err = util.ParseEnv(data)
	if err != nil {
		return nil, err
	}

	return env, nil
}

func init() {
	RootCmd.AddCommand(envSetCmd)
}
