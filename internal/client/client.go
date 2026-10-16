package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

const BaseURL = "https://api.utho.com/v2"

type Client struct {
	apiKey     string
	httpClient *http.Client
	baseURL    string
}

func NewClient(apiKey string) *Client {
	return &Client{
		apiKey:  apiKey,
		baseURL: BaseURL,
		httpClient: &http.Client{
			Timeout: 600 * time.Second,
		},
	}
}

func NewClientWithURL(apiKey string, baseURL string) *Client {
	url := BaseURL
	if baseURL != "" {
		url = baseURL
	}
	return &Client{
		apiKey:  apiKey,
		baseURL: url,
		httpClient: &http.Client{
			Timeout: 600 * time.Second,
		},
	}
}

func (c *Client) doRequest(method, endpoint string, body interface{}) ([]byte, error) {
	url := fmt.Sprintf("%s%s", c.baseURL, endpoint)

	var reqBody io.Reader
	if body != nil {
		jsonBytes, err := json.Marshal(body)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal request body: %w", err)
		}
		reqBody = bytes.NewBuffer(jsonBytes)
	}

	req, err := http.NewRequest(method, url, reqBody)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.apiKey))

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	respBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	// Trim any stray characters before JSON (API bug workaround)
	trimmed := respBytes
	for i, b := range respBytes {
		if b == '{' || b == '[' {
			trimmed = respBytes[i:]
			break
		}
	}
	respBytes = trimmed

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, fmt.Errorf("API error (status %d): %s", resp.StatusCode, string(respBytes))
	}

	return respBytes, nil
}

func (c *Client) Get(endpoint string) ([]byte, error) {
	return c.doRequest(http.MethodGet, endpoint, nil)
}

func (c *Client) Post(endpoint string, body interface{}) ([]byte, error) {
	return c.doRequest(http.MethodPost, endpoint, body)
}

func (c *Client) Put(endpoint string, body interface{}) ([]byte, error) {
	return c.doRequest(http.MethodPut, endpoint, body)
}

func (c *Client) Delete(endpoint string) ([]byte, error) {
	return c.doRequest(http.MethodDelete, endpoint, nil)
}
