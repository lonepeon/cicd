package github_test

import (
	"testing"

	"github.com/lonepeon/cicd/internal/github"
	"github.com/lonepeon/cicd/internal/github/githubtest"
	"github.com/lonepeon/golib/testutils"
)

func TestCreateRelease(t *testing.T) {
	server := githubtest.NewReleaseServer(t, "<username>", "<token>")
	server.ExpectedRepository = "lonepeon/something"
	server.ExpectedReleaseName = "20210312084512"
	server.ExpectedCommmitish = "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa"
	server.ReturnedReleaseID = 42
	client := server.StartMockServer()

	releaseID, err := github.CreateRelease(
		client,
		"lonepeon/something",
		"20210312084512",
		"aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa",
	)

	testutils.RequireNoError(t, err, "can't create release")
	githubtest.AssertEqualReleaseID(t, github.ReleaseID(42), releaseID, "unexpected release ID")
}
