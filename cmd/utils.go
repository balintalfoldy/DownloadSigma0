package cmd

import (
	"fmt"
	"io"
	"net/http"
	"time"
)

func PostRequest(url string, apiKey string, body io.Reader) (*http.Response, error) {

	client := &http.Client{
		Timeout: 30 * time.Second,
	}

	req, err := http.NewRequest("POST", url, body)
	if err != nil {
		return nil, fmt.Errorf("client: could not create request: %s", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("ApiKey", apiKey)

	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("client: error making http request: %s", err)
	}

	return resp, nil

}
