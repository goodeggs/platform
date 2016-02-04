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
	Use:   "env:set",
	Short: "Set environment variables",
	Run: func(cmd *cobra.Command, args []string) {
		appDir, err := util.AppDir(cmd)
		util.Check(err)

		clean, err := util.GitIsClean(appDir)
		util.Check(err)

		if !clean {
			util.Die("git working directory not clean.")
		}

		newEnv, err := readEnvChanges(args)
		util.Check(err)

		var updatedKeys []string
		for k, _ := range newEnv {
			updatedKeys = append(updatedKeys, k)
		}

		config, err := util.LoadAppConfig(cmd)
		util.Check(err)

		appName, err := util.AppName(cmd)
		util.Check(err)

		oldEnv, err := getExistingEnv(appName, config.EnvId)
		util.Check(err)

		err = mergo.Merge(&newEnv, oldEnv)
		util.Check(err)

		data := ""
		for k, v := range newEnv {
			data += fmt.Sprintf("%s=%s\n", k, v)
		}

		envId, err := util.EcruCreateSecret(appName, data)
		util.Check(err)

		err = util.RanchUpdateEnvId(appDir, envId)
		util.Check(err)

		err = util.GitAdd(appDir, ".ranch.yaml")
		util.Check(err)

		err = util.GitCommit(appDir, fmt.Sprintf("set env %s", strings.Join(updatedKeys, ",")))
		util.Check(err)
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

func getExistingEnv(appName, envId string) (env map[string]string, err error) {

	if envId == "" {
		return env, nil
	}

	plaintext, err := util.EcruGetSecret(appName, envId)
	if err != nil {
		return nil, err
	}

	env, err = util.ParseEnv(plaintext)
	if err != nil {
		return nil, err
	}

	return env, nil
}

func init() {
	RootCmd.AddCommand(envSetCmd)
}
