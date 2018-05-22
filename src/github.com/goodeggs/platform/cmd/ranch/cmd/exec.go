package cmd

import (
	"os"
	"strings"

	"golang.org/x/crypto/ssh/terminal"

	"github.com/goodeggs/platform/cmd/ranch/util"
	"github.com/spf13/cobra"
)

var execCmd = &cobra.Command{
	Use:   "exec <pid> <command>",
	Short: "Execute a command in an existing process",
	RunE: func(cmd *cobra.Command, args []string) (err error) {
		appName, err := util.AppName(cmd)
		if err != nil {
			return err
		}

		pid := args[0]
		command := strings.Join(args[1:], " ")

		exitCode, err := exec(appName, pid, command)
		if err != nil {
			return err
		}

		os.Exit(exitCode)
		return nil
	},
}

func exec(appName, pid, command string) (int, error) {
	fd := os.Stdin.Fd()

	var w, h int

	if terminal.IsTerminal(int(fd)) {
		stdinState, err := terminal.GetState(int(fd))

		if err != nil {
			return -1, err
		}

		defer terminal.Restore(int(fd), stdinState)

		w, h, err = terminal.GetSize(int(fd))
		if err != nil {
			return -1, err
		}
	}

	return util.ConvoxExec(appName, pid, command, w, h, os.Stdin, os.Stdout)
}

func init() {
	RootCmd.AddCommand(execCmd)
}
