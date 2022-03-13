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
	APIURL     string
	UploadURL  string

	username string
	token    string
}

func NewClient(username, token string) *Client {
	return &Client{
		APIURL:     "https://api.github.com",
		HTTPClient: http.DefaultClient,
		username:   username,
		token:      token,
	}
}
func (c *Client) PostUpload(path string, contentType string, data []byte) ([]byte, int, error) {
	req, err := c.buildRequest(c.UploadURL+path, data)
	if err != nil {
		return nil, -99, err
	}
	req.Header.Add("Content-Type", contentType)

	return c.executeRequest(req)
}

func (c *Client) PostAPI(path string, payload interface{}) ([]byte, int, error) {
	reqBody, err := json.Marshal(payload)
	if err != nil {
		return nil, -99, fmt.Errorf("can't encode body: %v (body=%#+v)", err, reqBody)
	}

	req, err := c.buildRequest(c.APIURL+path, reqBody)
	if err != nil {
		return nil, -99, err
	}

	return c.executeRequest(req)
}

func (c *Client) buildRequest(url string, reqBody []byte) (*http.Request, error) {
	req, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(reqBody))
	if err != nil {
		return nil, fmt.Errorf("can't initialize post request: %v", err)
	}
	req.Header.Add("Accept", "application/vnd.github.v3+json")
	req.SetBasicAuth(c.username, c.token)

	return req, nil
}

func (c *Client) executeRequest(req *http.Request) ([]byte, int, error) {
	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, -99, fmt.Errorf("can't execute request: %v", err)
	}

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, -99, fmt.Errorf("can't read response: %v", err)
	}

	return respBody, resp.StatusCode, nil
}
