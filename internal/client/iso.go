package client

import (
	"encoding/json"
	"fmt"
)

// ── Request structs ───────────────────────────────────────────────────────

type ISOAddRequest struct {
	DCSlug string `json:"dcslug"`
	Name   string `json:"name"`
	URL    string `json:"url"`
}

// ── Response structs ──────────────────────────────────────────────────────

type ISO struct {
	Name       string  `json:"name"`
	File       string  `json:"file"`
	Size       float64 `json:"size"`
	AddedAt    string  `json:"added_at"`
	Download   string  `json:"download"`
	DC         string  `json:"dc"`
	DCLocation struct {
		DCCC     string `json:"dccc"`
		Location string `json:"location"`
	} `json:"dclocation"`
}

// ── Client methods ────────────────────────────────────────────────────────

func (c *Client) AddISO(req *ISOAddRequest) (string, error) {
	respBytes, err := c.Post("/iso/add", req)
	if err != nil {
		return "", fmt.Errorf("failed to add ISO: %w", err)
	}
	var result map[string]interface{}
	if err := json.Unmarshal(respBytes, &result); err != nil {
		return "", fmt.Errorf("failed to parse add ISO response: %w", err)
	}
	if result["status"] != "success" {
		return "", fmt.Errorf("add ISO failed: %s", result["message"])
	}
	switch v := result["id"].(type) {
	case float64:
		return fmt.Sprintf("%.0f", v), nil
	default:
		return fmt.Sprintf("%v", v), nil
	}
}

func (c *Client) ListISOs() ([]ISO, error) {
	respBytes, err := c.Get("/iso")
	if err != nil {
		return nil, fmt.Errorf("failed to list ISOs: %w", err)
	}
	var result struct {
		ISOs []ISO `json:"isos"`
	}
	if err := json.Unmarshal(respBytes, &result); err != nil {
		return nil, fmt.Errorf("failed to parse list ISOs response: %w", err)
	}
	return result.ISOs, nil
}

func (c *Client) GetISOByFile(file string) (*ISO, error) {
	isos, err := c.ListISOs()
	if err != nil {
		return nil, err
	}
	for _, iso := range isos {
		if iso.File == file {
			return &iso, nil
		}
	}
	return nil, nil
}

func (c *Client) DeleteISO(file string) error {
	respBytes, err := c.Delete(fmt.Sprintf("/iso/%s/delete", file))
	if err != nil {
		return fmt.Errorf("failed to delete ISO: %w", err)
	}
	var result map[string]interface{}
	if err := json.Unmarshal(respBytes, &result); err != nil {
		return fmt.Errorf("failed to parse delete ISO response: %w", err)
	}
	if result["status"] != "success" {
		return fmt.Errorf("delete ISO failed: %s", result["message"])
	}
	return nil
}
