package cmd

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"text/template"

	"github.com/goodeggs/platform/cmd/ranch/util"
	"github.com/spf13/cobra"
)

const yamlTemplate string = `
name: {{.AppName}}

processes:
  web:
    command: node server.js
    count: 2
    memory: 256
`

type yamlTemplateVars struct {
	AppName string
}

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize the .ranch.yaml file",
	RunE: func(cmd *cobra.Command, args []string) (err error) {
		appName, err := util.AppName(cmd)
		if err != nil {
			return err
		}

		appDir, err := util.AppDir(cmd)
		if err != nil {
			return err
		}

		appYaml := path.Join(appDir, ".ranch.yaml")
		if _, err := os.Stat(appYaml); !os.IsNotExist(err) {
			return fmt.Errorf(".ranch.yaml already exists!")
		}

		tmpl, err := template.New(".ranch.yaml").Parse(yamlTemplate)
		if err != nil {
			return err
		}

		vars := yamlTemplateVars{appName}
		var buf bytes.Buffer
		if err = tmpl.Execute(&buf, vars); err != nil {
			return err
		}

		if err = ioutil.WriteFile(appYaml, buf.Bytes(), 0644); err != nil {
			return err
		}

		fmt.Println("generated .ranch.yaml -- check it now!")

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

		return nil
	},
}

func init() {
	RootCmd.AddCommand(initCmd)
}
