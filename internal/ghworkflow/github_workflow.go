package ghworkflow

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

type Entry struct {
	Version string
	Line    int
}

type VersionMatcherFunc func(string) (string, bool)

type GitHubWorkflow struct {
	content string
}

func Parse(path string) (GitHubWorkflow, error) {
	f, err := os.Open(path)
	if err != nil {
		return GitHubWorkflow{}, fmt.Errorf("can't open github workflow: %v", err)
	}
	defer f.Close()

	return ParseFromReader(f)
}

func ParseFromReader(r io.Reader) (GitHubWorkflow, error) {
	var content strings.Builder
	if _, err := io.Copy(&content, r); err != nil {
		return GitHubWorkflow{}, fmt.Errorf("")
	}

	return GitHubWorkflow{content: content.String()}, nil
}

func (w GitHubWorkflow) GetVersions(lang internal.Language) []Entry {
	var versions []Entry
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

			versions = append(versions, Entry{Version: version, Line: lineNumber})
		}
	}

	return versions
}

func (w GitHubWorkflow) versionMatchers(lang internal.Language) []VersionMatcherFunc {
	switch lang {
	case internal.Go:
		return []VersionMatcherFunc{w.actionSetupGoV2Matcher}
	case internal.Rust:
		return []VersionMatcherFunc{w.actionSetupRustMatcher}
	}

	panic(fmt.Sprintf("language '%v' is not supported", lang))
}

func (w GitHubWorkflow) actionSetupGoV2Matcher(line string) (string, bool) {
	matches := actionSetupGoV2.FindStringSubmatch(line)
	if len(matches) < 2 {
		return "", false
	}

	return matches[1], true
}

func (w GitHubWorkflow) actionSetupRustMatcher(line string) (string, bool) {
	matches := actionSetupRust.FindStringSubmatch(line)
	if len(matches) < 2 {
		return "", false
	}

	return matches[1], true
}
