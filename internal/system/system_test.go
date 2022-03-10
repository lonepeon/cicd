package system_test

import (
	"os"
	"testing"

	"github.com/lonepeon/cicd/internal"
	"github.com/lonepeon/cicd/internal/system"
	"github.com/lonepeon/golib/testutils"
)

func TestIntegration(t *testing.T) {
	if testing.Short() {
		t.SkipNow()
	}

	t.Run("TestGoVersionSuccess", testGoVersionSuccess)
	t.Run("TestGoVersionFail", testGoVersionFail)
	t.Run("TestRustVersionSuccess", testRustVersionSuccess)
	t.Run("TestRustVersionFail", testRustVersionFail)
}

func testGoVersionSuccess(t *testing.T) {
	t.Setenv("GO_TEST_VERSION", "1.17.7")
	t.Setenv("PATH", "./testdata/gobinary/:"+os.Getenv("PATH"))

	version, err := system.GetVersion(internal.Go)
	testutils.RequireNoError(t, err, "expected go version to be found")

	testutils.AssertEqualString(t, "1.17.7", version, "expected to find correct go version")
}

func testGoVersionFail(t *testing.T) {
	t.Setenv("PATH", "./testdata/failing-gobinary/:"+os.Getenv("PATH"))

	_, err := system.GetVersion(internal.Go)
	testutils.RequireHasError(t, err, "expected to not get go version")

	testutils.AssertContainsString(t, "can't get go version", err.Error(), "expected to find correct error message")
	testutils.AssertContainsString(t, "failing go binary", err.Error(), "expected to find correct error message")
}

func testRustVersionSuccess(t *testing.T) {
	t.Setenv("RUST_TEST_VERSION", "1.59.0")
	t.Setenv("PATH", "./testdata/rustbinary/:"+os.Getenv("PATH"))

	version, err := system.GetVersion(internal.Rust)
	testutils.RequireNoError(t, err, "expected rust version to be found")

	testutils.AssertEqualString(t, "1.59.0", version, "expected to find correct rust version")
}

func testRustVersionFail(t *testing.T) {
	t.Setenv("PATH", "./testdata/failing-rustbinary/:"+os.Getenv("PATH"))

	_, err := system.GetVersion(internal.Rust)
	testutils.RequireHasError(t, err, "expected to not get rust version")

	testutils.AssertContainsString(t, "can't get rust version", err.Error(), "expected to find correct error message")
	testutils.AssertContainsString(t, "failing rust binary", err.Error(), "expected to find correct error message")
}
