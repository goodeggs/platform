package util

import (
	"bufio"
	"bytes"
	"fmt"
	"os/exec"
	"strings"

	"github.com/spf13/cobra"
)

func git(c *cobra.Command, args ...string) ([]byte, error) {
	git := exec.Command("git", args...)
	out, err := git.CombinedOutput()

	if err != nil {
		err = fmt.Errorf("`%s`: %s", strings.Join(git.Args, " "), err.Error())
		return nil, err
	}

	return out, nil
}

func GitTag(c *cobra.Command, tag string, message string) error {
	_, err := git(c, "tag", "-am", message, tag)
	return err
}

func GitCommit(c *cobra.Command, message string) error {
	_, err := git(c, "commit", "-m", message)
	return err
}

func GitAdd(c *cobra.Command, files ...string) error {
	args := append([]string{"add"}, files...)
	_, err := git(c, args...)
	return err
}

func GitIsClean(_ *cobra.Command) (bool, error) {
	git := exec.Command("git", "status", "--porcelain")
	var out bytes.Buffer
	git.Stdout = &out
	err := git.Run()
	if err != nil {
		err = fmt.Errorf("`%s`: %s", strings.Join(git.Args, " "), err.Error())
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
