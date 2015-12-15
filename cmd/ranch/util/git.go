package util

import (
	"bufio"
	"bytes"
	"os/exec"

	"github.com/spf13/cobra"
)

func GitIsClean(_ *cobra.Command) (bool, error) {
	cmd := exec.Command("git", "status", "--porcelain")
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
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
