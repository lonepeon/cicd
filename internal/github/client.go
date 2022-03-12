package github

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type Client struct {
	HTTPClient *http.Client
	URL        string

	username string
	token    string
}

func NewClient(username, token string) *Client {
	return &Client{
		URL:        "https://api.github.com",
		HTTPClient: http.DefaultClient,
		username:   username,
		token:      token,
	}
}

func (c *Client) Post(path string, payload interface{}) ([]byte, int, error) {
	reqBody, err := json.Marshal(payload)
	if err != nil {
		return nil, -99, fmt.Errorf("can't encode body: %v (body=%#+v)", err, reqBody)
	}

	req, err := http.NewRequest(http.MethodPost, c.URL+path, bytes.NewReader(reqBody))
	if err != nil {
		return nil, -99, fmt.Errorf("can't initialize post request: %v", err)
	}
	req.Header.Add("Accept", "application/vnd.github.v3+json")
	req.SetBasicAuth(c.username, c.token)

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, -99, fmt.Errorf("can't post request: %v", err)
	}

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, -99, fmt.Errorf("can't read post response: %v", err)
	}

	return respBody, resp.StatusCode, nil
}
