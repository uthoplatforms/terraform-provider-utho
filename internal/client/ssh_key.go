package client

import (
	"encoding/json"
	"fmt"
)

type SSHKeyCreateRequest struct {
	Name   string `json:"name"`
	SSHKey string `json:"sshkey"`
}

type SSHKeyInstance struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	Value     string `json:"value"`
	CreatedAt string `json:"created_at"`
}

type SSHKeyListResponse struct {
	Status string           `json:"status"`
	Key    []SSHKeyInstance `json:"key"`
}

func (c *Client) CreateSSHKey(req *SSHKeyCreateRequest) error {
	respBytes, err := c.Post("/key/import", req)
	if err != nil {
		return fmt.Errorf("failed to create SSH key: %w", err)
	}
	var result map[string]interface{}
	if err := json.Unmarshal(respBytes, &result); err != nil {
		return fmt.Errorf("failed to parse create SSH key response: %w", err)
	}
	if result["status"] != "success" {
		return fmt.Errorf("create SSH key failed: %s", result["message"])
	}
	return nil
}

func (c *Client) GetSSHKey(keyID string) (*SSHKeyInstance, error) {
	respBytes, err := c.Get("/key")
	if err != nil {
		return nil, fmt.Errorf("failed to list SSH keys: %w", err)
	}
	var listResp SSHKeyListResponse
	if err := json.Unmarshal(respBytes, &listResp); err != nil {
		return nil, fmt.Errorf("failed to parse SSH key list: %w", err)
	}
	for _, key := range listResp.Key {
		if key.ID == keyID {
			return &key, nil
		}
	}
	return nil, nil
}

func (c *Client) DeleteSSHKey(keyID string) error {
	respBytes, err := c.Delete(fmt.Sprintf("/key/%s/delete", keyID))
	if err != nil {
		return fmt.Errorf("failed to delete SSH key: %w", err)
	}
	var result map[string]string
	if err := json.Unmarshal(respBytes, &result); err != nil {
		return fmt.Errorf("failed to parse delete SSH key response: %w", err)
	}
	if result["status"] != "success" {
		return fmt.Errorf("delete SSH key failed: %s", result["message"])
	}
	return nil
}
