package client

import (
	"encoding/json"
	"fmt"
)

// ── Request structs ───────────────────────────────────────────────────────

type AlertContactCreateRequest struct {
	Name         string `json:"name"`
	Email        string `json:"email"`
	MobileNumber string `json:"mobilenumber"`
	Status       string `json:"status"`
}

type AlertCreateRequest struct {
	Name     string `json:"name"`
	RefType  string `json:"ref_type"`
	Type     string `json:"type"`
	Compare  string `json:"compare"`
	Value    string `json:"value"`
	For      string `json:"for"`
	Contacts string `json:"contacts"`
	Status   string `json:"status"`
	RefIDs   string `json:"ref_ids"`
}

// ── Response structs ──────────────────────────────────────────────────────

type AlertContact struct {
	ID           string `json:"id"`
	Name         string `json:"name"`
	Email        string `json:"email"`
	Slack        string `json:"slack"`
	Status       string `json:"status"`
	MobileNumber string `json:"mobilenumber"`
}

type AlertContactListResponse struct {
	Contacts []AlertContact `json:"contacts"`
}

type AlertInstance struct {
	ID      string `json:"id"`
	Name    string `json:"name"`
	RefType string `json:"ref_type"`
	Type    string `json:"type"`
	Compare string `json:"compare"`
	Value   string `json:"value"`
	For     string `json:"for"`
	Status  string `json:"status"`
}

type AlertListResponse struct {
	Alerts []AlertInstance `json:"alerts"`
}

// ── Contact methods ───────────────────────────────────────────────────────

func (c *Client) CreateAlertContact(req *AlertContactCreateRequest) (string, error) {
	respBytes, err := c.Post("/alert/contact/add", req)
	if err != nil {
		return "", fmt.Errorf("failed to create alert contact: %w", err)
	}
	var result map[string]interface{}
	if err := json.Unmarshal(respBytes, &result); err != nil {
		return "", fmt.Errorf("failed to parse create alert contact response: %w", err)
	}
	if result["status"] != "success" {
		return "", fmt.Errorf("create alert contact failed: %s", result["message"])
	}
	switch v := result["id"].(type) {
	case float64:
		return fmt.Sprintf("%.0f", v), nil
	default:
		return fmt.Sprintf("%v", v), nil
	}
}

func (c *Client) UpdateAlertContact(contactID string, req *AlertContactCreateRequest) error {
	respBytes, err := c.Post(fmt.Sprintf("/alert/contact/%s/update", contactID), req)
	if err != nil {
		return fmt.Errorf("failed to update alert contact: %w", err)
	}
	var result map[string]string
	if err := json.Unmarshal(respBytes, &result); err != nil {
		return fmt.Errorf("failed to parse update alert contact response: %w", err)
	}
	if result["status"] != "success" {
		return fmt.Errorf("update alert contact failed: %s", result["message"])
	}
	return nil
}

func (c *Client) ListAlertContacts() ([]AlertContact, error) {
	respBytes, err := c.Get("/alert/contact/list")
	if err != nil {
		return nil, fmt.Errorf("failed to list alert contacts: %w", err)
	}
	var resp AlertContactListResponse
	if err := json.Unmarshal(respBytes, &resp); err != nil {
		return nil, fmt.Errorf("failed to parse list alert contacts response: %w", err)
	}
	return resp.Contacts, nil
}

func (c *Client) GetAlertContact(contactID string) (*AlertContact, error) {
	contacts, err := c.ListAlertContacts()
	if err != nil {
		return nil, err
	}
	for _, ct := range contacts {
		if ct.ID == contactID {
			return &ct, nil
		}
	}
	return nil, nil
}

func (c *Client) DeleteAlertContact(contactID string) error {
	respBytes, err := c.Delete(fmt.Sprintf("/alert/contact/%s/delete", contactID))
	if err != nil {
		return fmt.Errorf("failed to delete alert contact: %w", err)
	}
	var result map[string]string
	if err := json.Unmarshal(respBytes, &result); err != nil {
		return fmt.Errorf("failed to parse delete alert contact response: %w", err)
	}
	if result["status"] != "success" {
		return fmt.Errorf("delete alert contact failed: %s", result["message"])
	}
	return nil
}

// ── Alert methods ─────────────────────────────────────────────────────────

func (c *Client) CreateAlert(req *AlertCreateRequest) (string, error) {
	respBytes, err := c.Post("/alert", req)
	if err != nil {
		return "", fmt.Errorf("failed to create alert: %w", err)
	}
	var result map[string]interface{}
	if err := json.Unmarshal(respBytes, &result); err != nil {
		return "", fmt.Errorf("failed to parse create alert response: %w", err)
	}
	if result["status"] != "success" {
		return "", fmt.Errorf("create alert failed: %s", result["message"])
	}
	switch v := result["id"].(type) {
	case float64:
		return fmt.Sprintf("%.0f", v), nil
	default:
		return fmt.Sprintf("%v", v), nil
	}
}

func (c *Client) UpdateAlert(alertID string, req *AlertCreateRequest) error {
	respBytes, err := c.Post(fmt.Sprintf("/alert/%s/update", alertID), req)
	if err != nil {
		return fmt.Errorf("failed to update alert: %w", err)
	}
	var result map[string]string
	if err := json.Unmarshal(respBytes, &result); err != nil {
		return fmt.Errorf("failed to parse update alert response: %w", err)
	}
	if result["status"] != "success" {
		return fmt.Errorf("update alert failed: %s", result["message"])
	}
	return nil
}

func (c *Client) ListAlerts() ([]AlertInstance, error) {
	respBytes, err := c.Get("/alert")
	if err != nil {
		return nil, fmt.Errorf("failed to list alerts: %w", err)
	}
	var resp AlertListResponse
	if err := json.Unmarshal(respBytes, &resp); err != nil {
		return nil, fmt.Errorf("failed to parse list alerts response: %w", err)
	}
	return resp.Alerts, nil
}

func (c *Client) GetAlert(alertID string) (*AlertInstance, error) {
	alerts, err := c.ListAlerts()
	if err != nil {
		return nil, err
	}
	for _, a := range alerts {
		if a.ID == alertID {
			return &a, nil
		}
	}
	return nil, nil
}

func (c *Client) DeleteAlert(alertID string) error {
	respBytes, err := c.Delete(fmt.Sprintf("/alert/%s/delete", alertID))
	if err != nil {
		return fmt.Errorf("failed to delete alert: %w", err)
	}
	var result map[string]string
	if err := json.Unmarshal(respBytes, &result); err != nil {
		return fmt.Errorf("failed to parse delete alert response: %w", err)
	}
	if result["status"] != "success" {
		return fmt.Errorf("delete alert failed: %s", result["message"])
	}
	return nil
}
