package github_test

import (
	"strconv"
	"strings"
	"testing"

	"github.com/lonepeon/cicd/internal"
	"github.com/lonepeon/cicd/internal/github"
	"github.com/lonepeon/cicd/internal/github/githubtest"
	"github.com/lonepeon/golib/testutils"
)

func TestIntegration(t *testing.T) {
	if testing.Short() {
		t.SkipNow()
	}

	t.Run("TestParseRustSuccess", testParseRustSuccess)
	t.Run("TestParseGoSuccess", testParseGoSuccess)
	t.Run("TestParseFileNotFound", testParseFileNotFound)
}

//nolint:funlen
func TestGetVersion(t *testing.T) {
	type Testcase struct {
		Name             string
		Object           string
		Language         internal.Language
		ExpectedVersions []github.Entry
	}

	check := func(tc Testcase) {
		t.Run(tc.Name, func(t *testing.T) {
			r := strings.NewReader(tc.Object)
			workflow, err := github.ParseFromReader(r)
			testutils.RequireNoError(t, err, "unexpected error while parsing workflow")

			versions := workflow.GetVersions(tc.Language)
			testutils.RequireEqualInt(t, len(tc.ExpectedVersions), len(versions), "unexpected number of versions: %v", versions)

			for i := range tc.ExpectedVersions {
				testutils.AssertEqualString(t, tc.ExpectedVersions[i].Version, versions[i].Version, "unexpected version")
				testutils.AssertEqualInt(t, tc.ExpectedVersions[i].Line, versions[i].Line, "unexpected line")
			}
		})
	}

	check(Testcase{
		Name:             "withNoGoDeclaration",
		Language:         internal.Go,
		ExpectedVersions: nil,
		Object:           githubtest.NewWorkflowFile(t).Build(),
	})

	check(Testcase{
		Name:             "withNoRustDeclaration",
		Language:         internal.Rust,
		ExpectedVersions: nil,
		Object:           githubtest.NewWorkflowFile(t).Build(),
	})

	check(Testcase{
		Name:             "withOneGoDeclaration",
		Language:         internal.Go,
		ExpectedVersions: []github.Entry{{Version: "1.17.7", Line: 18}},
		Object:           githubtest.NewWorkflowFile(t).WithActionSetupGoV2("1.17.7").Build(),
	})

	check(Testcase{
		Name:             "withOneRustDeclaration",
		Language:         internal.Rust,
		ExpectedVersions: []github.Entry{{Version: "1.59.0", Line: 18}},
		Object:           githubtest.NewWorkflowFile(t).WithActionSetupRust("1.59.0").Build(),
	})

	check(Testcase{
		Name:     "withMultipleSameGoDeclarations",
		Language: internal.Go,
		ExpectedVersions: []github.Entry{
			{Version: "1.17.7", Line: 18},
			{Version: "1.17.7", Line: 31},
		},
		Object: githubtest.NewWorkflowFile(t).WithActionSetupGoV2("1.17.7", "1.17.7").Build(),
	})

	check(Testcase{
		Name:     "withMultipleSameRustDeclarations",
		Language: internal.Rust,
		ExpectedVersions: []github.Entry{
			{Version: "1.59.0", Line: 18},
			{Version: "1.59.0", Line: 31},
		},
		Object: githubtest.NewWorkflowFile(t).WithActionSetupRust("1.59.0", "1.59.0").Build(),
	})

	check(Testcase{
		Name:     "withMultipleDifferentGoDeclarations",
		Language: internal.Go,
		ExpectedVersions: []github.Entry{
			{Version: "1.17.7", Line: 18},
			{Version: "1.16", Line: 31},
		},
		Object: githubtest.NewWorkflowFile(t).WithActionSetupGoV2("1.17.7", "1.16").Build(),
	})

	check(Testcase{
		Name:     "withMultipleDifferentRustDeclarations",
		Language: internal.Rust,
		ExpectedVersions: []github.Entry{
			{Version: "1.59.0", Line: 18},
			{Version: "1.58.0", Line: 31},
		},
		Object: githubtest.NewWorkflowFile(t).WithActionSetupRust("1.59.0", "1.58.0").Build(),
	})
}

func testParseGoSuccess(t *testing.T) {
	workflow, err := github.Parse("testdata/workflow-go.yaml")
	testutils.RequireNoError(t, err, "expected file to be parsed")

	versions := workflow.GetVersions(internal.Go)
	testutils.RequireEqualInt(t, 2, len(versions), "unexpected number of matches in %v", versions)
	testutils.AssertEqualString(t, "1.17.7", versions[0].Version, "expected to find correct go version")
	testutils.AssertEqualString(t, "1.17.7", versions[1].Version, "expected to find correct go version")
	testutils.AssertEqualString(t, "17", strconv.Itoa(versions[0].Line), "expected to find correct line")
	testutils.AssertEqualString(t, "31", strconv.Itoa(versions[1].Line), "expected to find correct line")
}

func testParseRustSuccess(t *testing.T) {
	workflow, err := github.Parse("testdata/workflow-rust.yaml")
	testutils.RequireNoError(t, err, "expected file to be parsed")

	versions := workflow.GetVersions(internal.Rust)
	testutils.RequireEqualInt(t, 2, len(versions), "unexpected number of matches in %v", versions)
	testutils.AssertEqualString(t, "1.59.0", versions[0].Version, "expected to find correct rust version")
	testutils.AssertEqualString(t, "1.59.0", versions[1].Version, "expected to find correct rust version")
	testutils.AssertEqualString(t, "17", strconv.Itoa(versions[0].Line), "expected to find correct line")
	testutils.AssertEqualString(t, "31", strconv.Itoa(versions[1].Line), "expected to find correct line")
}

func testParseFileNotFound(t *testing.T) {
	_, err := github.Parse("not-a-file")
	testutils.RequireHasError(t, err, "expected file to not be parsed")
	testutils.AssertContainsString(t, "can't open", err.Error(), "expected error message to contain reason")
	testutils.AssertContainsString(t, "not-a-file", err.Error(), "expected error message to contain filename")
}
