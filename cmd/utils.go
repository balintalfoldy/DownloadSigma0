package cmd

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/schollz/progressbar/v3"
)

type Client struct {
	BaseURL    string
	apiKey     string
	HTTPClient *http.Client
}

func NewClient(baseURL string, apiKey string, timeOutSec int) *Client {
	return &Client{
		BaseURL: baseURL,
		apiKey:  apiKey,
		HTTPClient: &http.Client{
			Timeout: time.Duration(timeOutSec) * time.Second,
		},
	}
}

func GetProgressbar(size int) *progressbar.ProgressBar {
	bar := progressbar.NewOptions(size,
		progressbar.OptionEnableColorCodes(true),
		progressbar.OptionShowBytes(true),
		progressbar.OptionSetWidth(60),
		progressbar.OptionSetDescription("Downloading file..."),
		progressbar.OptionSetTheme(progressbar.Theme{
			Saucer:        "[green]|[reset]",
			SaucerHead:    "[yellow]>[reset]",
			SaucerPadding: " ",
			BarStart:      "[",
			BarEnd:        "]",
		}))
	return bar
}

func (client *Client) DownloadFile(outPath string, url string, size int64, bar *progressbar.ProgressBar) error {

	req, err := http.NewRequest("GET", fmt.Sprintf("%s/%s", client.BaseURL, url), nil)
	if err != nil {
		return fmt.Errorf("could not create request: %s", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/octet-stream")
	req.Header.Set("ApiKey", client.apiKey)

	resp, err := client.HTTPClient.Do(req)
	if err != nil {
		return fmt.Errorf("error making http request: %s", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("bad status: %s", resp.Status)
	}

	f, err := os.Create(outPath)
	if err != nil {
		return fmt.Errorf("error creating file: %s", err)
	}
	defer f.Close()

	if _, err := io.Copy(io.MultiWriter(f, bar), resp.Body); err != nil {
		return fmt.Errorf("error while downloading: %v", err)
	}

	return nil
}

func (client *Client) PostRequest(url string, body io.Reader, v interface{}) error {

	req, err := http.NewRequest("POST", fmt.Sprintf("%s/%s", client.BaseURL, url), body)
	if err != nil {
		return fmt.Errorf("could not create request: %s", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	req.Header.Set("ApiKey", client.apiKey)

	res, err := client.HTTPClient.Do(req)
	if err != nil {
		return fmt.Errorf("error making http request: %s", err)
	}

	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		var errRes errorResponse
		if err := json.NewDecoder(res.Body).Decode(&errRes); err != nil {
			return fmt.Errorf("error decoding error response with status code %d: %s", res.StatusCode, err)
		}
		return fmt.Errorf("error making http request with status code %d: %s", errRes.Code, errRes.ErrorMessage)
	}

	if err := json.NewDecoder(res.Body).Decode(&v); err != nil {
		return fmt.Errorf("error decoding error response with status code %d: %s", res.StatusCode, err)
	}
	return nil

}

type errorResponse struct {
	Code         int    `json:"code"`
	ErrorMessage string `json:"ErrorMessage"`
}

type successResponse struct {
	Code int    `json:"code"`
	Data string `json:"ErrorMessage"`
}
