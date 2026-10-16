package client

import (
	"encoding/json"
	"fmt"
	"strings"
)

// ── Request structs ───────────────────────────────────────────────────────

type FirewallCreateRequest struct {
	Name string `json:"name"`
}

type FirewallRuleRequest struct {
	Type        string `json:"type"`
	Service     string `json:"service"`
	Protocol    string `json:"protocol"`
	Port        string `json:"port"`
	PortRange   string `json:"port_range"`
	Addresses   string `json:"addresses"`
	SourceRange string `json:"source_range"`
}

// ── Response structs ──────────────────────────────────────────────────────

type FirewallRule struct {
	ID         string `json:"id"`
	FirewallID string `json:"firewallid"`
	Type       string `json:"type"`
	Service    string `json:"service"`
	Protocol   string `json:"protocol"`
	Port       string `json:"port"`
	Addresses  string `json:"addresses"`
}

type FirewallInstance struct {
	ID           string         `json:"id"`
	Name         string         `json:"name"`
	CreatedAt    string         `json:"created_at"`
	RuleCount    string         `json:"rulecount"`
	ServersCount string         `json:"serverscount"`
	Rules        []FirewallRule `json:"rules"`
}

type FirewallListResponse struct {
	Firewalls []FirewallInstance `json:"firewalls"`
}

// ── helpers ───────────────────────────────────────────────────────────────

// parseFirewallResponse parses a standard Utho API response.
// Many firewall endpoints return an empty body on success — treat that as success.
func parseFirewallResponse(respBytes []byte) (map[string]interface{}, error) {
	trimmed := strings.TrimSpace(string(respBytes))
	if trimmed == "" || trimmed == "null" {
		return map[string]interface{}{"status": "success"}, nil
	}
	var result map[string]interface{}
	if err := json.Unmarshal(respBytes, &result); err != nil {
		// Non-JSON but non-empty — treat as success since HTTP call succeeded
		return map[string]interface{}{"status": "success"}, nil
	}
	return result, nil
}

// ── API methods ───────────────────────────────────────────────────────────

func (c *Client) CreateFirewall(name string) (string, error) {
	respBytes, err := c.Post("/firewall/create", &FirewallCreateRequest{Name: name})
	if err != nil {
		return "", fmt.Errorf("failed to create firewall: %w", err)
	}
	result, err := parseFirewallResponse(respBytes)
	if err != nil {
		return "", fmt.Errorf("failed to parse create firewall response: %w", err)
	}
	if result["status"] != nil && result["status"] != "success" {
		return "", fmt.Errorf("create firewall failed: %s", result["message"])
	}

	// Try to get ID from response first
	if id, ok := result["id"].(string); ok && id != "" {
		return id, nil
	}

	// API returns empty body — look up the firewall by name
	listResp, err2 := c.Get("/firewall")
	if err2 != nil {
		return "", fmt.Errorf("failed to look up created firewall: %w", err2)
	}
	var fwList FirewallListResponse
	if err2 := json.Unmarshal(listResp, &fwList); err2 != nil {
		return "", fmt.Errorf("failed to parse firewall list: %w", err2)
	}
	for _, fw := range fwList.Firewalls {
		if fw.Name == name {
			return fw.ID, nil
		}
	}
	return "", fmt.Errorf("firewall created but could not find ID for name: %s", name)
}

func (c *Client) GetFirewall(firewallID string) (*FirewallInstance, error) {
	respBytes, err := c.Get(fmt.Sprintf("/firewall/%s", firewallID))
	if err != nil {
		return nil, fmt.Errorf("failed to get firewall: %w", err)
	}
	var resp FirewallListResponse
	if err := json.Unmarshal(respBytes, &resp); err != nil {
		return nil, fmt.Errorf("failed to parse get firewall response: %w", err)
	}
	if len(resp.Firewalls) == 0 {
		return nil, nil
	}
	return &resp.Firewalls[0], nil
}

func (c *Client) DeleteFirewall(firewallID string) error {
	respBytes, err := c.Delete(fmt.Sprintf("/firewall/%s/destroy", firewallID))
	if err != nil {
		return fmt.Errorf("failed to delete firewall: %w", err)
	}
	result, err := parseFirewallResponse(respBytes)
	if err != nil {
		return fmt.Errorf("failed to parse delete firewall response: %w", err)
	}
	if result["status"] != "success" {
		return fmt.Errorf("delete firewall failed: %s", result["message"])
	}
	return nil
}

func (c *Client) AddFirewallRule(firewallID string, req *FirewallRuleRequest) (string, error) {
	endpoint := fmt.Sprintf("/firewall/%s/rule/add", firewallID)
	respBytes, err := c.Post(endpoint, req)
	if err != nil {
		return "", fmt.Errorf("failed to add firewall rule: %w", err)
	}
	result, err := parseFirewallResponse(respBytes)
	if err != nil {
		return "unknown", nil
	}
	if result["status"] != nil && result["status"] != "success" {
		return "", fmt.Errorf("add firewall rule failed: %s", result["message"])
	}
	if result["id"] != nil {
		return fmt.Sprintf("%v", result["id"]), nil
	}
	return "unknown", nil
}

func (c *Client) DeleteFirewallRule(firewallID string, ruleID string) error {
	endpoint := fmt.Sprintf("/firewall/%s/rule/%s/delete", firewallID, ruleID)
	respBytes, err := c.Delete(endpoint)
	if err != nil {
		// 404 means rule already deleted — treat as success
		if strings.Contains(err.Error(), "404") || strings.Contains(err.Error(), "Not Found") {
			return nil
		}
		return fmt.Errorf("failed to delete firewall rule: %w", err)
	}
	result, err := parseFirewallResponse(respBytes)
	if err != nil {
		return nil
	}
	if result["status"] != nil && result["status"] != "success" {
		return fmt.Errorf("delete firewall rule failed: %s", result["message"])
	}
	return nil
}

func (c *Client) AttachFirewallServer(firewallID string, cloudID string) error {
	endpoint := fmt.Sprintf("/firewall/%s/server/add", firewallID)
	payload := map[string]string{"cloudid": cloudID}
	respBytes, err := c.Post(endpoint, payload)
	if err != nil {
		return fmt.Errorf("failed to attach firewall to server: %w", err)
	}
	result, err := parseFirewallResponse(respBytes)
	if err != nil {
		return fmt.Errorf("failed to parse attach firewall response: %w", err)
	}
	if result["status"] != "success" {
		return fmt.Errorf("attach firewall failed: %s", result["message"])
	}
	return nil
}

func (c *Client) DetachFirewallServer(firewallID string, cloudID string) error {
	endpoint := fmt.Sprintf("/firewall/%s/server/%s/delete", firewallID, cloudID)
	respBytes, err := c.Delete(endpoint)
	if err != nil {
		return fmt.Errorf("failed to detach firewall from server: %w", err)
	}
	result, err := parseFirewallResponse(respBytes)
	if err != nil {
		return fmt.Errorf("failed to parse detach firewall response: %w", err)
	}
	if result["status"] != "success" {
		return fmt.Errorf("detach firewall failed: %s", result["message"])
	}
	return nil
}
