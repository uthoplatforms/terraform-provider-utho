package helper

import (
	"bytes"
	"context"
	"net/http"
)

const (
	YES = "yes"
	NO  = "no"
)

// NewUthoRequest send a request with auth token and set common http headers
func NewUthoRequest(ctx context.Context, method, url string, body []byte, token string) (*http.Response, error) {
	bodyReader := bytes.NewReader(body)
	req, err := http.NewRequest(method, url, bodyReader)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Accept-Encoding", "application/json")
	client := &http.Client{}

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	return resp, nil
}
