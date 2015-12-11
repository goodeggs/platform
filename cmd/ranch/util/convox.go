package util

import (
	"os"
	"os/exec"
)

func Convox(args ...string) error {
	cmd := exec.Command("convox", args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}
