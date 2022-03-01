package asdf_test

import (
	"strings"
	"testing"

	"github.com/lonepeon/golib/testutils"

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
		Language           string
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
		Name:               "WithOneLineForTheLanguage",
		ASDFContent:        "golang 1.17.7",
		Language:           "golang",
		ExpectedVersion:    "1.17.7",
		ExpectedHasVersion: true,
	})

	check(TestCase{
		Name:               "WithSeveralLinesButOneForTheLanguage",
		ASDFContent:        "ruby 3.1.1\ngolang 1.17.7\nnodejs 16.13.2",
		Language:           "golang",
		ExpectedVersion:    "1.17.7",
		ExpectedHasVersion: true,
	})

	check(TestCase{
		Name:               "WithLanguagePresentInComment",
		ASDFContent:        "#golang 1.16\ngolang 1.17.7\nnodejs 16.13.2",
		Language:           "golang",
		ExpectedVersion:    "1.17.7",
		ExpectedHasVersion: true,
	})

	check(TestCase{
		Name:               "WithLanguageFollowedByComment",
		ASDFContent:        "ruby 3.1.1\ngolang 1.17.7 # a comment\nnodejs 16.13.2",
		Language:           "golang",
		ExpectedVersion:    "1.17.7",
		ExpectedHasVersion: true,
	})
}

func testParseSuccess(t *testing.T) {
	versions, err := asdf.Parse("testdata/tool-versions")
	testutils.RequireNoError(t, err, "expected file to be parsed")

	golang, found := versions.GetVersion(asdf.Go)
	testutils.AssertEqualBool(t, true, found, "expected to find go version")
	testutils.AssertEqualString(t, "1.17.7", golang, "expected to find correct go version")

	ruby, found := versions.GetVersion(asdf.Ruby)
	testutils.AssertEqualBool(t, true, found, "expected to find ruby version")
	testutils.AssertEqualString(t, "3.1.1", ruby, "expected to find correct ruby version")

	nodejs, found := versions.GetVersion(asdf.NodeJS)
	testutils.AssertEqualBool(t, true, found, "expected to find nodejs version")
	testutils.AssertEqualString(t, "16.13.2", nodejs, "expected to find correct nodejs version")
}

func testParseFileNotFound(t *testing.T) {
	_, err := asdf.Parse("not-a-file")
	testutils.RequireHasError(t, err, "expected file to not be parsed")
	testutils.AssertContainsString(t, "can't open", err.Error(), "expected error message to contain reason")
	testutils.AssertContainsString(t, "not-a-file", err.Error(), "expected error message to contain filename")
}
