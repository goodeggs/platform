package cmd

import (
	"fmt"
	"io/ioutil"
	"os"
	"regexp"
	"strings"

	"github.com/goodeggs/platform/cmd/ranch/Godeps/_workspace/src/github.com/spf13/cobra"
	"github.com/goodeggs/platform/cmd/ranch/util"
)

var envUnsetCmd = &cobra.Command{
	Use:   "env:unset",
	Short: "Un-set environment variables",
	Run: func(cmd *cobra.Command, args []string) {
		appDir, err := util.AppDir(cmd)
		util.Check(err)

		clean, err := util.GitIsClean(appDir)
		util.Check(err)

		if !clean {
			util.Die("git working directory not clean.")
		}

		keysToDelete, err := readKeysFromEnv(args)
		util.Check(err)

		keysToDeleteMap := make(map[string]int, len(keysToDelete))

		for _, key := range keysToDelete {
			keysToDeleteMap[key] = 1
		}

		config, err := util.LoadAppConfig(cmd)
		util.Check(err)

		appName, err := util.AppName(cmd)
		util.Check(err)

		oldEnv, err := util.EnvGet(appName, config.EnvId)
		util.Check(err)

		newEnv := make(map[string]string)
		var deletedKeys []string

		for k, v := range oldEnv {
			_, ok := keysToDeleteMap[k]
			if !ok {
				newEnv[k] = v
			} else {
				deletedKeys = append(deletedKeys, k)
			}
		}

		if len(deletedKeys) == 0 {
			util.Die("key(s) not found... nothing to do.")
		}

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

		message := fmt.Sprintf("unset env %s", strings.Join(deletedKeys, ","))

		err = util.GitCommit(appDir, message)
		util.Check(err)

		sha, err := util.GitCurrentSha(appDir)
		util.Check(err)

		fmt.Printf("[%s] %s\n", sha, message)
	},
}

func readKeysFromEnv(args []string) (keys []string, err error) {

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
		data += fmt.Sprintf(" %s", value)
	}

	keys = regexp.MustCompile(`\s+`).Split(data, -1)

	return keys, nil
}

func init() {
	RootCmd.AddCommand(envUnsetCmd)
}
