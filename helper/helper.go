package helper

import (
	"bytes"
	"context"
	"crypto/rand"
	"encoding/base64"
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

// GeneratePassword generates a random password of the given length.
func GeneratePassword(length int) (string, error) {
	byteLen := (length*6 + 7) / 8
	b := make([]byte, byteLen)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}
	password := base64.RawURLEncoding.EncodeToString(b)
	if len(password) > length {
		password = password[:length]
	}

	return password, nil
}
