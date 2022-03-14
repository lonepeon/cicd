package github

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type ReleaseID int

type Asset struct {
	Name        string
	ContentType string
	Content     []byte
}

func CreateRelease(client *Client, repository string, releaseName string, commitish string) (ReleaseID, error) {
	type Response struct {
		ID int `json:"id"`
	}

	body, httpCode, err := client.PostAPI(fmt.Sprintf("/repos/%s/releases", repository), map[string]string{
		"target_commitish": commitish,
		"tag_name":         releaseName,
		"name":             releaseName,
	})

	if err != nil {
		return ReleaseID(-99), fmt.Errorf("can't create the release: %v", err)
	}

	if httpCode != http.StatusCreated {
		return ReleaseID(-99), fmt.Errorf("unexpected create release status code (expected=201, actual=%d, body=%s)", httpCode, string(body))
	}

	var response Response
	if err := json.Unmarshal(body, &response); err != nil {
		return ReleaseID(-99), fmt.Errorf("can't decode create release body: %v (body=%v)", err, string(body))
	}

	return ReleaseID(response.ID), nil
}

func UploadAsset(client *Client, repository string, releaseID ReleaseID, asset Asset) error {
	body, httpCode, err := client.PostUpload(fmt.Sprintf("/repos/%s/releases/%d/assets?name=%s", repository, releaseID, asset.Name), asset.ContentType, asset.Content)
	if err != nil {
		return fmt.Errorf("can't upload asset to release: %v", err)
	}

	if httpCode != http.StatusCreated {
		return fmt.Errorf("unexpected upload asset status code (expected=201, actual=%d, body=%s)", httpCode, string(body))
	}

	return nil
}
