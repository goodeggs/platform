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
version: 1

processes:
  web:
    command: node server.js
    instances: 2
    memory: 256
    domains:
      - {{.AppName}}
`

type yamlTemplateVars struct {
	AppName string
}

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize the .ranch.yaml file",
	Run: func(cmd *cobra.Command, args []string) {
		appName, err := util.AppName(cmd)
		util.Check(err)

		appDir, err := util.AppDir(cmd)
		util.Check(err)

		appYaml := path.Join(appDir, ".ranch.yaml")
		if _, err := os.Stat(appYaml); !os.IsNotExist(err) {
			util.Die(".ranch.yaml already exists!")
		}

		tmpl, err := template.New(".ranch.yaml").Parse(yamlTemplate)
		util.Check(err)

		vars := yamlTemplateVars{appName}
		var buf bytes.Buffer
		err = tmpl.Execute(&buf, vars)
		util.Check(err)

		err = ioutil.WriteFile(appYaml, buf.Bytes(), 0644)
		util.Check(err)

		fmt.Println("generated .ranch.yaml")
	},
}

func init() {
	RootCmd.AddCommand(initCmd)
}
