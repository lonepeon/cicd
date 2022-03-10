package system

import (
	"fmt"
	"os/exec"
	"regexp"
	"strings"

	"github.com/lonepeon/cicd/internal"
)

var (
	goVersionRegex   = regexp.MustCompile(`go(([\d]+.?)+)`)
	rustVersionRegex = regexp.MustCompile(`rustc (([\d]+.?)+)`)
)

func GetVersion(lang internal.Language) (string, error) {
	switch lang {
	case internal.Go:
		return goVersion()
	case internal.Rust:
		return rustVersion()
	}

	panic(fmt.Sprintf("missing switch clause for language '%v'", lang))
}

func goVersion() (string, error) {
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

func rustVersion() (string, error) {
	cmd := exec.Command("rustc", "--version")

	var stdout strings.Builder
	var stderr strings.Builder

	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("can't get rust version output: %v (output=%s)", err, stderr.String())
	}

	submatches := rustVersionRegex.FindStringSubmatch(stdout.String())

	if len(submatches) < 2 {
		return "", fmt.Errorf("unexpected command output (output=%s)", stdout.String())
	}

	return strings.TrimSpace(submatches[1]), nil
}
