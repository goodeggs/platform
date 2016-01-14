package util

import (
	"bufio"
	"bytes"
	"fmt"
	"os/exec"
	"strings"
)

func git(appDir string, args ...string) (string, error) {
	git := exec.Command("git", args...)
	git.Dir = appDir
	out, err := git.CombinedOutput()

	if err != nil {
		err = fmt.Errorf("`%s`: %s", strings.Join(git.Args, " "), err.Error())
		return "", err
	}

	return string(out), nil
}

func GitTag(appDir string, tag string, message string) error {
	_, err := git(appDir, "tag", "-am", message, tag)
	return err
}

func GitCommit(appDir string, message string) error {
	_, err := git(appDir, "commit", "-m", message)
	return err
}

func GitAdd(appDir string, files ...string) error {
	args := append([]string{"add"}, files...)
	_, err := git(appDir, args...)
	return err
}

func GitCurrentSha(appDir string) (string, error) {
	out, err := git(appDir, "rev-parse", "--short", "HEAD")
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(out), nil
}

func GitIsClean(appDir string) (bool, error) {
	git := exec.Command("git", "status", "--porcelain")
	git.Dir = appDir
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
		// if any line doesn't start with ?? it ain't clean
		line = scanner.Text()
		if line[:2] != "??" {
			return false, nil
		}
	}
	if err = scanner.Err(); err != nil {
		return false, err
	}

	return true, nil
}
