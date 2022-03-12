package github

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type ReleaseID int

func CreateRelease(client *Client, repository string, releaseName string, commitish string) (ReleaseID, error) {
	type Response struct {
		ID int `json:"id"`
	}

	body, httpCode, err := client.Post(fmt.Sprintf("/repos/%s/releases", repository), map[string]string{
		"target_commitish": commitish,
		"tag_name":         releaseName,
		"name":             releaseName,
	})

	if err != nil {
		return ReleaseID(-99), fmt.Errorf("can't create the release: %v", err)
	}

	if httpCode != http.StatusCreated {
		return ReleaseID(-99), fmt.Errorf("unexpected create release status code (expected=201, actual=%d)", httpCode)
	}

	var response Response
	if err := json.Unmarshal(body, &response); err != nil {
		return ReleaseID(-99), fmt.Errorf("can't decode create release body: %v (body=%v)", err, string(body))
	}

	return ReleaseID(response.ID), nil
}
