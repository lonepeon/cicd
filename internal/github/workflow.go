package github

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"regexp"
	"strings"

	"github.com/lonepeon/cicd/internal"
)

var (
	actionSetupGoV2 = regexp.MustCompile(`  go-version:\s+"?((\d+\.?)+)"?`)
	actionSetupRust = regexp.MustCompile(`  toolchain:\s+"?((\d+\.?)+)"?`)
)

type WorkflowMatch struct {
	Version string
	Line    int
}

type VersionMatcherFunc func(string) (string, bool)

type Workflow struct {
	content string
}

func ParseWorkflow(path string) (Workflow, error) {
	f, err := os.Open(path)
	if err != nil {
		return Workflow{}, fmt.Errorf("can't open github workflow: %v", err)
	}
	defer f.Close()

	return ParseWorkflowFromReader(f)
}

func ParseWorkflowFromReader(r io.Reader) (Workflow, error) {
	var content strings.Builder
	if _, err := io.Copy(&content, r); err != nil {
		return Workflow{}, fmt.Errorf("")
	}

	return Workflow{content: content.String()}, nil
}

func (w Workflow) GetVersions(lang internal.Language) []WorkflowMatch {
	var versions []WorkflowMatch
	var lineNumber int
	scanner := bufio.NewScanner(strings.NewReader(w.content))
	matchers := w.versionMatchers(lang)

	for scanner.Scan() {
		lineNumber++

		for _, matcher := range matchers {
			version, ok := matcher(scanner.Text())
			if !ok {
				continue
			}

			versions = append(versions, WorkflowMatch{Version: version, Line: lineNumber})
		}
	}

	return versions
}

func (w Workflow) versionMatchers(lang internal.Language) []VersionMatcherFunc {
	switch lang {
	case internal.Go:
		return []VersionMatcherFunc{w.actionSetupGoV2Matcher}
	case internal.Rust:
		return []VersionMatcherFunc{w.actionSetupRustMatcher}
	}

	panic(fmt.Sprintf("language '%v' is not supported", lang))
}

func (w Workflow) actionSetupGoV2Matcher(line string) (string, bool) {
	matches := actionSetupGoV2.FindStringSubmatch(line)
	if len(matches) < 2 {
		return "", false
	}

	return matches[1], true
}

func (w Workflow) actionSetupRustMatcher(line string) (string, bool) {
	matches := actionSetupRust.FindStringSubmatch(line)
	if len(matches) < 2 {
		return "", false
	}

	return matches[1], true
}
