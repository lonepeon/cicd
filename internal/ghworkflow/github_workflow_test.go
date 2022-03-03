package ghworkflow_test

import (
	"strconv"
	"strings"
	"testing"

	"github.com/lonepeon/cicd/internal"
	"github.com/lonepeon/cicd/internal/ghworkflow"
	"github.com/lonepeon/cicd/internal/ghworkflow/ghworkflowtest"
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
		ExpectedVersions []ghworkflow.Entry
	}

	check := func(tc Testcase) {
		t.Run(tc.Name, func(t *testing.T) {
			r := strings.NewReader(tc.Object)
			workflow, err := ghworkflow.ParseFromReader(r)
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
		Name:             "withNoDeclaration",
		Language:         internal.Go,
		ExpectedVersions: nil,
		Object:           ghworkflowtest.NewWorkflowFile(t).Build(),
	})

	check(Testcase{
		Name:             "withOneDeclaration",
		Language:         internal.Go,
		ExpectedVersions: []ghworkflow.Entry{{Version: "1.17.7", Line: 18}},
		Object:           ghworkflowtest.NewWorkflowFile(t).WithActionSetupGoV2("1.17.7").Build(),
	})

	check(Testcase{
		Name:     "withMultipleSameDeclarations",
		Language: internal.Go,
		ExpectedVersions: []ghworkflow.Entry{
			{Version: "1.17.7", Line: 18},
			{Version: "1.17.7", Line: 31},
		},
		Object: ghworkflowtest.NewWorkflowFile(t).WithActionSetupGoV2("1.17.7", "1.17.7").Build(),
	})

	check(Testcase{
		Name:     "withMultipleDifferentDeclarations",
		Language: internal.Go,
		ExpectedVersions: []ghworkflow.Entry{
			{Version: "1.17.7", Line: 18},
			{Version: "1.16", Line: 31},
		},
		Object: ghworkflowtest.NewWorkflowFile(t).WithActionSetupGoV2("1.17.7", "1.16").Build(),
	})

}

func testParseGoSuccess(t *testing.T) {
	workflow, err := ghworkflow.Parse("testdata/workflow-go.yaml")
	testutils.RequireNoError(t, err, "expected file to be parsed")

	versions := workflow.GetVersions(internal.Go)
	testutils.RequireEqualInt(t, 2, len(versions), "unexpected number of matches in %v", versions)
	testutils.AssertEqualString(t, "1.17.7", versions[0].Version, "expected to find correct go version")
	testutils.AssertEqualString(t, "1.17.7", versions[1].Version, "expected to find correct go version")
	testutils.AssertEqualString(t, "17", strconv.Itoa(versions[0].Line), "expected to find correct line")
	testutils.AssertEqualString(t, "31", strconv.Itoa(versions[1].Line), "expected to find correct line")
}

func testParseFileNotFound(t *testing.T) {
	_, err := ghworkflow.Parse("not-a-file")
	testutils.RequireHasError(t, err, "expected file to not be parsed")
	testutils.AssertContainsString(t, "can't open", err.Error(), "expected error message to contain reason")
	testutils.AssertContainsString(t, "not-a-file", err.Error(), "expected error message to contain filename")
}
