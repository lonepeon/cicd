package report_test

import (
	"strings"
	"testing"

	"github.com/lonepeon/cicd/internal/report"
	"github.com/lonepeon/golib/testutils"
)

func TestVersionReporterSuccess(t *testing.T) {
	var stdout strings.Builder
	var stderr strings.Builder

	var reporter report.Reporter
	reporter.Stdout = &stdout
	reporter.Stderr = &stderr

	vReporter := report.NewVersionReporter("1.42.0")
	vReporter.Reporter = &reporter

	vReporter.AssertVersion(report.VersionEntry{
		Version: "1.42.0",
		Kind:    "system",
	})

	testutils.AssertContainsString(t, "[OK]", stdout.String(), "expected proper message")
	testutils.AssertContainsString(t, "system", stdout.String(), "expected proper message")
	testutils.AssertContainsString(t, "1.42.0", stdout.String(), "expected proper message")
	testutils.AssertEqualString(t, "", stderr.String(), "expected empty STDERR")

	vReporter.AssertVersion(report.VersionEntry{
		Version:  "1.42.0",
		Kind:     "GitHub Workflow",
		Location: ".github/workflows/test.yaml:12",
	})

	testutils.AssertContainsString(t, "[OK]", stdout.String(), "expected proper message")
	testutils.AssertContainsString(t, "GitHub", stdout.String(), "expected proper message")
	testutils.AssertContainsString(t, "1.42.0", stdout.String(), "expected proper message")
	testutils.AssertContainsString(t, "test.yaml:12", stdout.String(), "expected proper message")
	testutils.AssertEqualString(t, "", stderr.String(), "expected empty STDERR")
}

func TestVersionReporterFailure(t *testing.T) {
	var stdout strings.Builder
	var stderr strings.Builder

	var reporter report.Reporter
	reporter.Stdout = &stdout
	reporter.Stderr = &stderr

	vReporter := report.NewVersionReporter("1.42.0")
	vReporter.Reporter = &reporter

	vReporter.AssertVersion(report.VersionEntry{
		Version: "1.12.0",
		Kind:    "system",
	})

	testutils.AssertEqualString(t, "", stdout.String(), "expected empty STDOUT")
	testutils.AssertContainsString(t, "[ERR]", stderr.String(), "expected proper message")
	testutils.AssertContainsString(t, "system", stderr.String(), "expected proper message")
	testutils.AssertContainsString(t, "1.12.0", stderr.String(), "expected proper message")

	vReporter.AssertVersion(report.VersionEntry{
		Version:  "1.18.0",
		Kind:     "GitHub Workflow",
		Location: ".github/workflows/test.yaml:12",
	})

	testutils.AssertEqualString(t, "", stdout.String(), "expected empty STDOUT")
	testutils.AssertContainsString(t, "[ERR]", stderr.String(), "expected proper message")
	testutils.AssertContainsString(t, "GitHub", stderr.String(), "expected proper message")
	testutils.AssertContainsString(t, "1.18.0", stderr.String(), "expected proper message")
	testutils.AssertContainsString(t, "test.yaml:12", stderr.String(), "expected proper message")
}

func TestVersionReporterHasFailed(t *testing.T) {
	var combined strings.Builder

	var reporter report.Reporter
	reporter.Stdout = &combined
	reporter.Stderr = &combined

	vReporter := report.NewVersionReporter("1.42.0")
	vReporter.Reporter = &reporter

	testutils.AssertEqualBool(t, false, vReporter.HasFailed(), "exepected report to be a success")
	vReporter.AssertVersion(report.VersionEntry{Version: "1.42.0", Kind: "system"})
	testutils.AssertEqualBool(t, false, vReporter.HasFailed(), "exepected report to still be a success")
	vReporter.AssertVersion(report.VersionEntry{Version: "1.0.0", Kind: "system"})
	testutils.AssertEqualBool(t, true, vReporter.HasFailed(), "exepected report to be a failure")
	vReporter.AssertVersion(report.VersionEntry{Version: "1.42.0", Kind: "system"})
	testutils.AssertEqualBool(t, true, vReporter.HasFailed(), "exepected report to still be a failure")
}
