package util

import (
	"os"
	"strings"

	"github.com/spf13/cobra"
)

const TEST_RANCHY = ".ranch.yaml"
const TEST_FILE_1 = ".ranch.test1.yaml"

func removeTestFiles() {
	_ = os.Remove(TEST_RANCHY)
	_ = os.Remove(TEST_FILE_1)
}

func mockCmd(flags string) cobra.Command {
	cmd := cobra.Command{
		Use: "junk",
	}
	cmd.Flags().StringP("filename", "f", "", "config filename (defaults to .ranch.yaml)")
	cmd.Flags().StringP("app", "a", "", "app name")
	cmd.SetArgs(strings.Split(flags, " "))
	cmd.SilenceErrors = true
	cmd.SilenceUsage = true
	cmd.Execute()
	return cmd
}
