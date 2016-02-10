package cmd

import (
	"github.com/goodeggs/platform/cmd/ranch/Godeps/_workspace/src/golang.org/x/crypto/ssh/terminal"
	"os"
	"strings"

	"github.com/goodeggs/platform/cmd/ranch/Godeps/_workspace/src/github.com/spf13/cobra"
	"github.com/goodeggs/platform/cmd/ranch/util"
)

var runCmd = &cobra.Command{
	Use:   "run",
	Short: "Run a one-off command",
	Run: func(cmd *cobra.Command, args []string) {
		appName, err := util.AppName(cmd)
		util.Check(err)

		process := "web"
		command := strings.Join(args, " ")

		exitCode, err := runAttached(appName, process, command)
		util.Check(err)

		os.Exit(exitCode)
	},
}

func runAttached(appName, process, command string) (int, error) {
	fd := os.Stdin.Fd()

	if terminal.IsTerminal(int(fd)) {
		stdinState, err := terminal.GetState(int(fd))

		if err != nil {
			return -1, err
		}

		defer terminal.Restore(int(fd), stdinState)
	}

	return util.ConvoxRunAttached(appName, process, command, os.Stdin, os.Stdout)
}

func init() {
	RootCmd.AddCommand(runCmd)
}
