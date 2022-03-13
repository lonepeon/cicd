package githubtest

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/lonepeon/cicd/internal/github"
	"github.com/lonepeon/golib/testutils"
)

func AssertEqualReleaseID(t *testing.T, expected github.ReleaseID, actual github.ReleaseID, pattern string, vars ...interface{}) {
	t.Helper()

	testutils.AssertEqualInt(t, int(expected), int(actual), pattern, vars...)
}

type UploadServer struct {
	t        *testing.T
	username string
	token    string

	ExpectedRepository  string
	ExpectedReleaseID   int
	ExpectedContentType string
	ExpectedContent     []byte
}

func (s *UploadServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	testutils.AssertEqualString(s.t, fmt.Sprintf("/repos/%s/releases/%d/assets", s.ExpectedRepository, s.ExpectedReleaseID), r.URL.Path, "unexpected request path")

	username, token, found := r.BasicAuth()
	testutils.AssertEqualBool(s.t, true, found, "missing password in basic authentication")
	testutils.AssertEqualString(s.t, s.username, username, "unexpected authenticated username")
	testutils.AssertEqualString(s.t, s.token, token, "unexpected authenticated token")

	rawBody, err := io.ReadAll(r.Body)
	testutils.AssertNoError(s.t, err, "can't read request body")

	testutils.AssertEqualBytes(s.t, s.ExpectedContent, rawBody, "unexpected asset")
	testutils.AssertEqualString(s.t, s.ExpectedContentType, r.Header.Get("Content-Type"), "unexpected request content type")

	w.WriteHeader(http.StatusCreated)

	fmt.Fprintf(w, `
 {
  "url": "https://api.github.com/repos/%[1]s/releases/assets/1",
  "browser_download_url": "https://github.com/%[1]s/releases/download/v1.0.0/example.zip",
  "id": 1,
  "node_id": "MDEyOlJlbGVhc2VBc3NldDE=",
  "name": "example.zip",
  "label": "short description",
  "state": "uploaded",
  "content_type": "application/zip",
  "size": 1024,
  "download_count": 42,
  "created_at": "2013-02-27T19:35:32Z",
  "updated_at": "2013-02-27T19:35:32Z",
  "uploader": {
    "login": "octocat",
    "id": 1,
  }
}`, s.ExpectedRepository)
}

func (s *UploadServer) StartMockServer() *github.Client {
	mock := httptest.NewServer(s)
	s.t.Cleanup(mock.Close)

	client := github.NewClient(s.username, s.token)
	client.HTTPClient = mock.Client()
	client.UploadURL = mock.URL

	return client
}

func NewUploadServer(t *testing.T, username, token string) *UploadServer {
	return &UploadServer{t: t, username: username, token: token}
}

type ReleaseServer struct {
	t        *testing.T
	username string
	token    string

	ExpectedRepository  string
	ExpectedReleaseName string
	ExpectedCommmitish  string

	ReturnedReleaseID int
}

func (s *ReleaseServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	testutils.AssertEqualString(s.t, fmt.Sprintf("/repos/%s/releases", s.ExpectedRepository), r.URL.Path, "unexpected request path")

	username, token, found := r.BasicAuth()
	testutils.AssertEqualBool(s.t, true, found, "missing password in basic authentication")
	testutils.AssertEqualString(s.t, s.username, username, "unexpected authenticated username")
	testutils.AssertEqualString(s.t, s.token, token, "unexpected authenticated token")

	rawBody, err := io.ReadAll(r.Body)
	testutils.AssertNoError(s.t, err, "can't read request body")

	var reqBody map[string]string
	err = json.Unmarshal(rawBody, &reqBody)
	testutils.AssertNoError(s.t, err, "can't decode request body: %v", rawBody)

	testutils.AssertEqualString(s.t, s.ExpectedReleaseName, reqBody["name"], "unexpected release name")
	testutils.AssertEqualString(s.t, s.ExpectedCommmitish, reqBody["target_commitish"], "unexpected targetted commit/branch")

	w.WriteHeader(http.StatusCreated)
	fmt.Fprintf(w, `
{
  "url": "https://api.github.com/repos/%[1]s/releases/%[2]d",
  "html_url": "https://github.com/%[1]s/releases/%[3]s",
  "assets_url": "https://api.github.com/repos/%[1]s/releases/%[2]d/assets",
  "upload_url": "https://uploads.github.com/repos/%[1]s/releases/%[2]d/assets{?name,label}",
  "tarball_url": "https://api.github.com/repos/%[1]s/tarball/%[3]s",
  "zipball_url": "https://api.github.com/repos/%[1]s/zipball/%[3]s",
  "discussion_url": "https://github.com/%[1]s/discussions/90",
  "id": %[2]d,
  "node_id": "MDc6UmVsZWFzZTE=",
  "tag_name": "%[3]s",
  "target_commitish": "%[4]s",
  "name": "%[3]s",
  "body": "Description of the release",
  "draft": false,
  "prerelease": false,
  "created_at": "2013-02-27T19:35:32Z",
  "published_at": "2013-02-27T19:35:32Z",
  "author": {
    "login": "octocat",
    "id": 1
  }
}
`, s.ExpectedRepository, s.ReturnedReleaseID, s.ExpectedReleaseName, s.ExpectedCommmitish)
}

func (s *ReleaseServer) StartMockServer() *github.Client {
	mock := httptest.NewServer(s)
	s.t.Cleanup(mock.Close)

	client := github.NewClient(s.username, s.token)
	client.HTTPClient = mock.Client()
	client.APIURL = mock.URL

	return client
}

func NewReleaseServer(t *testing.T, username, token string) *ReleaseServer {
	return &ReleaseServer{t: t, username: username, token: token}
}
