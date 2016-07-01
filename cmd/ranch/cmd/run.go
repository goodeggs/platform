package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/goodeggs/platform/cmd/ranch/Godeps/_workspace/src/golang.org/x/crypto/ssh/terminal"

	"github.com/goodeggs/platform/cmd/ranch/Godeps/_workspace/src/github.com/spf13/cobra"
	"github.com/goodeggs/platform/cmd/ranch/util"
)

var runCmd = &cobra.Command{
	Use:   "run <command>",
	Short: "Run a one-off command",
	RunE: func(cmd *cobra.Command, args []string) (err error) {
		if len(args) == 0 {
			cmd.Usage()
			return fmt.Errorf("must specify command")
		}

		appName, err := util.AppName(cmd)
		if err != nil {
			return err
		}

		process := "run"
		command := strings.Join(args, " ")

		exitCode, err := runAttached(appName, process, command)
		if err != nil {
			return err
		}

		os.Exit(exitCode)
		return nil
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
