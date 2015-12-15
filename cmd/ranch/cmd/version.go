package cmd

import (
	"bufio"
	"bytes"
	"fmt"
	"io/ioutil"
	"regexp"
	"strings"

	"github.com/goodeggs/platform/cmd/ranch/util"
	"github.com/spf13/cobra"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Prints the application version",
	Run: func(cmd *cobra.Command, args []string) {
		version, err := util.AppVersion(cmd)
		util.Check(err)

		fmt.Printf("v%d\n", version)
	},
}

var versionBumpCmd = &cobra.Command{
	Use:   "bump",
	Short: "Increments the application version",
	Run: func(cmd *cobra.Command, args []string) {

		fmt.Println("TODO: ensure git is clean")

		oldVersion, err := util.AppVersion(cmd)
		util.Check(err)

		newVersion := oldVersion + 1

		appConfigPath, err := util.AppConfigPath(cmd)
		util.Check(err)

		input, err := ioutil.ReadFile(appConfigPath)
		util.Check(err)

		// fuzzy match ^version: 1$
		re, err := regexp.Compile(fmt.Sprintf("^\\s*version\\s*:\\s*%d\\b", oldVersion))
		util.Check(err)

		scanner := bufio.NewScanner(bytes.NewReader(input))
		var buf bytes.Buffer
		var line string
		for scanner.Scan() {
			line = scanner.Text()
			if re.Match([]byte(line)) {
				line = strings.Replace(line, fmt.Sprintf("%d", oldVersion), fmt.Sprintf("%d", newVersion), 1)
			}
			buf.WriteString(line)
			buf.WriteString("\n")
		}
		util.Check(scanner.Err())

		err = ioutil.WriteFile(appConfigPath, buf.Bytes(), 0644)
		util.Check(err)

		fmt.Println("TODO: git add .ranch.yaml")
		fmt.Println("TODO: git commit -m 'v2'")
		fmt.Println("TODO: git tag -am 'v2'")

		fmt.Printf("v%d\n", newVersion)
	},
}

func init() {
	versionCmd.AddCommand(versionBumpCmd)
	RootCmd.AddCommand(versionCmd)
}
