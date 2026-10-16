package client

import (
	"encoding/json"
	"fmt"
)

// ── Request structs ───────────────────────────────────────────────────────

type APITokenCreateRequest struct {
	Name  string `json:"name"`
	Write string `json:"write"`
}

// ── Response structs ──────────────────────────────────────────────────────

type APITokenInstance struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	Write     string `json:"write"`
	CreatedAt string `json:"created_at"`
}

type APITokenListResponse struct {
	Status string             `json:"status"`
	API    []APITokenInstance `json:"api"`
}

// ── API methods ───────────────────────────────────────────────────────────

func (c *Client) CreateAPIToken(req *APITokenCreateRequest) (string, string, error) {
	respBytes, err := c.Post("/api/generate", req)
	if err != nil {
		return "", "", fmt.Errorf("failed to create API token: %w", err)
	}
	var result map[string]string
	if err := json.Unmarshal(respBytes, &result); err != nil {
		return "", "", fmt.Errorf("failed to parse create API token response: %w", err)
	}
	if result["status"] != "success" {
		return "", "", fmt.Errorf("create API token failed: %s", result["message"])
	}
	return result["apikey"], result["message"], nil
}

func (c *Client) ListAPITokens() ([]APITokenInstance, error) {
	respBytes, err := c.Get("/api")
	if err != nil {
		return nil, fmt.Errorf("failed to list API tokens: %w", err)
	}
	var resp APITokenListResponse
	if err := json.Unmarshal(respBytes, &resp); err != nil {
		return nil, fmt.Errorf("failed to parse list API tokens response: %w", err)
	}
	return resp.API, nil
}

func (c *Client) GetAPIToken(tokenID string) (*APITokenInstance, error) {
	tokens, err := c.ListAPITokens()
	if err != nil {
		return nil, err
	}
	for _, t := range tokens {
		if t.ID == tokenID {
			return &t, nil
		}
	}
	return nil, nil
}

func (c *Client) DeleteAPIToken(tokenID string) error {
	respBytes, err := c.Delete(fmt.Sprintf("/api/%s/delete", tokenID))
	if err != nil {
		return fmt.Errorf("failed to delete API token: %w", err)
	}
	var result map[string]string
	if err := json.Unmarshal(respBytes, &result); err != nil {
		return fmt.Errorf("failed to parse delete API token response: %w", err)
	}
	if result["status"] != "success" {
		return fmt.Errorf("delete API token failed: %s", result["message"])
	}
	return nil
}
