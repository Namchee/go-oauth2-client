package services

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type GithubUserRequest struct {
	Name string `json:"name"`
}

func GetUsername(accessToken string) (string, error) {
	httpClient := http.Client{
		Timeout: 10 * time.Second,
	}

	req, err := http.NewRequest(http.MethodGet, "https://api.github.com/user", nil)

	if err != nil {
		return "", err
	}

	req.Header.Set("Authorization", fmt.Sprintf("token %s", accessToken))

	res, err := httpClient.Do(req)

	if err != nil {
		return "", err
	}

	defer res.Body.Close()

	var githubResponse GithubUserRequest

	if err := json.NewDecoder(res.Body).Decode(&githubResponse); err != nil {
		return "", err
	}

	return githubResponse.Name, nil
}
