package cmd

import (
	"fmt"
	"io/ioutil"
	"os"
	"regexp"
	"strings"

	"github.com/goodeggs/platform/cmd/ranch/util"
	"github.com/spf13/cobra"
)

var envUnsetCmd = &cobra.Command{
	Use:   "env:unset KEY KEY2",
	Short: "Un-set environment variables",
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

		keysToDelete, err := readKeysFromEnv(args)
		if err != nil {
			return err
		}

		keysToDeleteMap := make(map[string]int, len(keysToDelete))

		for _, key := range keysToDelete {
			keysToDeleteMap[key] = 1
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
			return fmt.Errorf("env:unset is deprecated")
		}

		appName, err := util.AppName(cmd)
		if err != nil {
			return err
		}

		oldEnv, err := util.EnvGet(appName, config.EnvId)
		if err != nil {
			return err
		}

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
			return fmt.Errorf("key(s) not found... nothing to do.")
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

		message := fmt.Sprintf("unset env %s", strings.Join(deletedKeys, ","))

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
