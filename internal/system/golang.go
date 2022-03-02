package system

import (
	"fmt"
	"os/exec"
	"regexp"
	"strings"
)

var (
	goVersionRegex = regexp.MustCompile(`go(([\d]+.?)+)`)
)

func GoVersion() (string, error) {
	cmd := exec.Command("go", "version")

	var stdout strings.Builder
	var stderr strings.Builder

	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("can't get go version output: %v (output=%s)", err, stderr.String())
	}

	submatches := goVersionRegex.FindStringSubmatch(stdout.String())

	if len(submatches) < 2 {
		return "", fmt.Errorf("unexpected command output (output=%s)", stdout.String())
	}

	return strings.TrimSpace(submatches[1]), nil
}
