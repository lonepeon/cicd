package asdf_test

import (
	"strings"
	"testing"

	"github.com/lonepeon/golib/testutils"

	"github.com/lonepeon/cicd/internal"
	"github.com/lonepeon/cicd/internal/asdf"
)

func TestIntegration(t *testing.T) {
	if testing.Short() {
		t.SkipNow()
	}

	t.Run("TestParseSuccess", testParseSuccess)
	t.Run("TestParseFileNotFound", testParseFileNotFound)
}

func TestParseFromReaderErrors(t *testing.T) {
	type TestCase struct {
		Name                     string
		ASDFContent              string
		Error                    error
		AdditionalStringPresence []string
	}

	check := func(tc TestCase) {
		t.Run(tc.Name, func(t *testing.T) {
			content := strings.NewReader(tc.ASDFContent)
			_, err := asdf.ParseFromReader(content)
			testutils.RequireHasError(t, err, "expecting an error")
			testutils.AssertErrorIs(t, tc.Error, err, "unexpected error")

			for _, additionalString := range tc.AdditionalStringPresence {
				testutils.AssertContainsString(t, additionalString, err.Error(), "partial string not found in error")
			}
		})
	}

	check(TestCase{
		Name:                     "WithSameLanguageMultipleTime",
		ASDFContent:              "golang 1.16\nnodejs 16.13.2\ngolang 1.17.7",
		Error:                    asdf.ErrMultipleDeclaration,
		AdditionalStringPresence: []string{"golang", "3"},
	})
}

func TestGetVersion(t *testing.T) {
	type TestCase struct {
		Name               string
		ASDFContent        string
		Language           internal.Language
		ExpectedVersion    string
		ExpectedHasVersion bool
	}

	check := func(tc TestCase) {
		t.Run(tc.Name, func(t *testing.T) {
			content := strings.NewReader(tc.ASDFContent)
			versions, err := asdf.ParseFromReader(content)
			testutils.RequireNoError(t, err, "unexpected error while parsing content")
			version, found := versions.GetVersion(tc.Language)

			testutils.AssertEqualBool(t, tc.ExpectedHasVersion, found, "unexpected found value")
			testutils.AssertEqualString(t, tc.ExpectedVersion, version, "unexpected version value")
		})
	}

	check(TestCase{
		Name:               "WithOneLineForGo",
		ASDFContent:        "golang 1.17.7",
		Language:           internal.Go,
		ExpectedVersion:    "1.17.7",
		ExpectedHasVersion: true,
	})

	check(TestCase{
		Name:               "WithSeveralLinesButOneForRust",
		ASDFContent:        "ruby 3.1.1\ngolang 1.17.7\nrust 1.59.0",
		Language:           internal.Rust,
		ExpectedVersion:    "1.59.0",
		ExpectedHasVersion: true,
	})

	check(TestCase{
		Name:               "WithSeveralLinesButOneForGo",
		ASDFContent:        "ruby 3.1.1\ngolang 1.17.7\nnodejs 16.13.2",
		Language:           internal.Go,
		ExpectedVersion:    "1.17.7",
		ExpectedHasVersion: true,
	})

	check(TestCase{
		Name:               "WithGoPresentInComment",
		ASDFContent:        "#golang 1.16\ngolang 1.17.7\nnodejs 16.13.2",
		Language:           internal.Go,
		ExpectedVersion:    "1.17.7",
		ExpectedHasVersion: true,
	})

	check(TestCase{
		Name:               "WithGoFollowedByComment",
		ASDFContent:        "ruby 3.1.1\ngolang 1.17.7 # a comment\nnodejs 16.13.2",
		Language:           internal.Go,
		ExpectedVersion:    "1.17.7",
		ExpectedHasVersion: true,
	})
}

func testParseSuccess(t *testing.T) {
	versions, err := asdf.Parse("testdata/tool-versions")
	testutils.RequireNoError(t, err, "expected file to be parsed")

	golang, found := versions.GetVersion(internal.Go)
	testutils.AssertEqualBool(t, true, found, "expected to find go version")
	testutils.AssertEqualString(t, "1.17.7", golang, "expected to find correct go version")
}

func testParseFileNotFound(t *testing.T) {
	_, err := asdf.Parse("not-a-file")
	testutils.RequireHasError(t, err, "expected file to not be parsed")
	testutils.AssertContainsString(t, "can't open", err.Error(), "expected error message to contain reason")
	testutils.AssertContainsString(t, "not-a-file", err.Error(), "expected error message to contain filename")
}
