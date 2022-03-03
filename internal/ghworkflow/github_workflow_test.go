package ghworkflow_test

import (
	"strings"
	"testing"

	"github.com/lonepeon/cicd/internal"
	"github.com/lonepeon/cicd/internal/ghworkflow"
	"github.com/lonepeon/golib/testutils"
)

func TestIntegration(t *testing.T) {
	if testing.Short() {
		t.SkipNow()
	}

	t.Run("TestParseGoSuccess", testParseGoSuccess)
	t.Run("TestParseFileNotFound", testParseFileNotFound)
}

func TestGetVersionForGo(t *testing.T) {
	type Testcase struct {
		Name             string
		Object           string
		Language         internal.Language
		ExpectedVersions []string
	}

	check := func(tc Testcase) {
		t.Run(tc.Name, func(t *testing.T) {
			r := strings.NewReader(tc.Object)
			workflow, err := ghworkflow.ParseFromReader(r)
			testutils.RequireNoError(t, err, "unexpected error while parsing workflow")

			versions := workflow.GetVersions(tc.Language)
			testutils.RequireEqualInt(t, len(tc.ExpectedVersions), len(versions), "unexpected number of versions: %v", versions)

			for i := range tc.ExpectedVersions {
				testutils.AssertEqualString(t, tc.ExpectedVersions[i], versions[i], "unexpected versions")
			}
		})
	}

	check(Testcase{
		Name:             "withNoDeclaration",
		Language:         internal.Go,
		ExpectedVersions: nil,
		Object: `
jobs:
  format-tests:
    runs-on: "ubuntu-latest"
    steps:
      - name: "checkout code"
        uses: actions/checkout@v2
`,
	})

	check(Testcase{
		Name:             "withOneDeclaration",
		Language:         internal.Go,
		ExpectedVersions: []string{"1.17.7"},
		Object: `
jobs:
  format-tests:
    runs-on: "ubuntu-latest"
    steps:
      - name: "checkout code"
        uses: actions/checkout@v2
      - name: "setup go version"
        uses: actions/setup-go@v2
        with:
          go-version: "1.17.7"
`,
	})

	check(Testcase{
		Name:             "withMultipleSameDeclarations",
		Language:         internal.Go,
		ExpectedVersions: []string{"1.17.7", "1.17.7"},
		Object: `
jobs:
  format-tests:
    runs-on: "ubuntu-latest"
    steps:
      - name: "checkout code"
        uses: actions/checkout@v2
      - name: "setup go version"
        uses: actions/setup-go@v2
        with:
          go-version: "1.17.7"

  integration-tests:
    runs-on: "ubuntu-latest"
    steps:
      - name: "checkout code"
        uses: actions/checkout@v2
      - name: "setup go version"
        uses: actions/setup-go@v2
        with:
          go-version: "1.17.7"
`,
	})

	check(Testcase{
		Name:             "withMultipleDifferentDeclarations",
		Language:         internal.Go,
		ExpectedVersions: []string{"1.17.7", "1.16"},
		Object: `
jobs:
  format-tests:
    runs-on: "ubuntu-latest"
    steps:
      - name: "checkout code"
        uses: actions/checkout@v2
      - name: "setup go version"
        uses: actions/setup-go@v2
        with:
          go-version: "1.17.7"

  integration-tests:
    runs-on: "ubuntu-latest"
    steps:
      - name: "checkout code"
        uses: actions/checkout@v2
      - name: "setup go version"
        uses: actions/setup-go@v2
        with:
          go-version: "1.16"
`,
	})

}

func testParseGoSuccess(t *testing.T) {
	workflow, err := ghworkflow.Parse("testdata/workflow-go.yaml")
	testutils.RequireNoError(t, err, "expected file to be parsed")

	versions := workflow.GetVersions(internal.Go)
	testutils.RequireEqualInt(t, 2, len(versions), "unexpected number of matches in %v", versions)
	testutils.AssertEqualString(t, "1.17.7", versions[0], "expected to find correct go version")
	testutils.AssertEqualString(t, "1.17.7", versions[1], "expected to find correct go version")
}

func testParseFileNotFound(t *testing.T) {
	_, err := ghworkflow.Parse("not-a-file")
	testutils.RequireHasError(t, err, "expected file to not be parsed")
	testutils.AssertContainsString(t, "can't open", err.Error(), "expected error message to contain reason")
	testutils.AssertContainsString(t, "not-a-file", err.Error(), "expected error message to contain filename")
}
