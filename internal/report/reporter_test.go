package report_test

import (
	"strings"
	"testing"

	"github.com/lonepeon/cicd/internal/report"
	"github.com/lonepeon/golib/testutils"
)

func TestReporterSuccess(t *testing.T) {
	var stdout strings.Builder
	var stderr strings.Builder

	var r report.Reporter
	r.Stdout = &stdout
	r.Stderr = &stderr

	r.Success("this is a success")

	testutils.AssertContainsString(t, "[OK]", stdout.String(), "expected proper message")
	testutils.AssertContainsString(t, "this is a success", stdout.String(), "expected proper message")
	testutils.AssertEqualString(t, "", stderr.String(), "expected empty STDERR")
}

func TestReporterFailure(t *testing.T) {
	var stdout strings.Builder
	var stderr strings.Builder

	var r report.Reporter
	r.Stdout = &stdout
	r.Stderr = &stderr

	r.Fail("this is a failure")

	testutils.AssertEqualString(t, "", stdout.String(), "expected empty STDOUT")
	testutils.AssertContainsString(t, "[ERR]", stderr.String(), "expected proper message")
	testutils.AssertContainsString(t, "this is a failure", stderr.String(), "expected proper message")
}

func TestReporterHasFailed(t *testing.T) {
	var combined strings.Builder

	var r report.Reporter
	r.Stdout = &combined
	r.Stderr = &combined

	testutils.AssertEqualBool(t, false, r.HasFailed(), "expected reporter to be successful")
	r.Success("this is a success")
	testutils.AssertEqualBool(t, false, r.HasFailed(), "expected reporter to still be successful")
	r.Fail("this is a failure")
	testutils.AssertEqualBool(t, true, r.HasFailed(), "expected reporter to still be a failure")
	r.Success("this is another success")
	testutils.AssertEqualBool(t, true, r.HasFailed(), "expected reporter to still be a failure")
}
