package util

import (
	"bufio"
	"bytes"
	"fmt"
	"os/exec"
	"strings"

	"github.com/spf13/cobra"
)

func GitAdd(_ *cobra.Command, files ...string) error {
	args := append([]string{"add"}, files...)
	cmd := exec.Command("git", args...)
	err := cmd.Run()
	if err != nil {
		err = fmt.Errorf("`%s`: %s", strings.Join(cmd.Args, " "), err.Error())
		return err
	}
	return nil
}

func GitIsClean(_ *cobra.Command) (bool, error) {
	cmd := exec.Command("git", "status", "--porcelain")
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		err = fmt.Errorf("`%s`: %s", strings.Join(cmd.Args, " "), err.Error())
		return false, err
	}

	var line string
	scanner := bufio.NewScanner(bytes.NewReader(out.Bytes()))
	for scanner.Scan() {
		// if any line starts with ?? it ain't clean
		line = scanner.Text()
		if line[:2] == "??" {
			return false, nil
		}
	}
	if err = scanner.Err(); err != nil {
		return false, err
	}

	return true, nil
}
