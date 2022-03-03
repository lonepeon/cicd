package ghworkflow

import (
	"fmt"
	"io"
	"os"
	"regexp"
	"strings"

	"github.com/lonepeon/cicd/internal"
)

var (
	actionSetupGoV2 = regexp.MustCompile(`  go-version:\s+"?((\d+\.?)+)"?`)
)

type VersionsMatcherFunc func() []string

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

func (w GitHubWorkflow) GetVersions(lang internal.Language) []string {
	matchers := w.versionMatchers(lang)

	var versions []string
	for _, matcher := range matchers {
		versions = append(versions, matcher()...)
	}

	return versions
}

func (w GitHubWorkflow) versionMatchers(lang internal.Language) []VersionsMatcherFunc {
	switch lang {
	case internal.Go:
		return []VersionsMatcherFunc{w.actionSetupGoV2Matcher}
	}

	panic(fmt.Sprintf("language '%v' is not supported", lang))
}

func (w GitHubWorkflow) actionSetupGoV2Matcher() []string {
	matches := actionSetupGoV2.FindAllStringSubmatch(w.content, -1)

	versions := make([]string, len(matches))
	for i := range matches {
		versions[i] = matches[i][1]
	}

	return versions
}
