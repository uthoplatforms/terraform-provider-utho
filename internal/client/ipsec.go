package client

import (
	"encoding/json"
	"fmt"
	"time"
)

// ── Request structs ───────────────────────────────────────────────────────

type IPSecCreateRequest struct {
	Name         string `json:"name"`
	DCSlug       string `json:"dcslug"`
	VPC          string `json:"vpc"`
	CPUModel     string `json:"cpumodel"`
	BillingCycle string `json:"billingcycle"`
}

type IPSecPairRequest struct {
	IPSecID     string `json:"ipsecid"`
	PeerIPSecID string `json:"peer_ipsecid"`
	LocalSubnet string `json:"local_subnet"`
	PeerSubnet  string `json:"peer_subnet"`
	Name        string `json:"name"`
}

type IPSecConnectionUpdateRequest struct {
	IPSecID          string `json:"ipsecid"`
	ID               string `json:"id"`
	Name             string `json:"name"`
	RemoteIP         string `json:"remote_ip"`
	RemoteLocalIP    string `json:"remote_local_ip"`
	LocalIP          string `json:"local_ip"`
	PSK              string `json:"psk"`
	Phase1Encryption string `json:"phase1-encryption"`
	Phase2Encryption string `json:"phase2-encryption"`
	Phase1Integrity  string `json:"phase1-integrity"`
	Phase2Integrity  string `json:"phase2-integrity"`
	Phase1DHGroup    string `json:"phase1-dh-group"`
	Phase2DHGroup    string `json:"phase2-dh-group"`
	IKEVersion       string `json:"ike_version"`
	Phase1Lifetime   string `json:"phase1_lifetime"`
	Phase2Lifetime   string `json:"phase2_lifetime"`
	RekeyMargin      string `json:"rekey_margin"`
	RekeyFuzz        string `json:"rekey_fuzz"`
	ReplayWindow     string `json:"replay_window"`
	DPDTimeout       string `json:"dpd_timeout"`
	DPDAction        string `json:"dpd_action"`
	StartupAction    string `json:"startup_action"`
}

// ── Response structs ──────────────────────────────────────────────────────

type IPSecInstance struct {
	ID           string      `json:"id"`
	Name         string      `json:"name"`
	DCSlug       string      `json:"dcslug"`
	VPC          interface{} `json:"vpc"`
	PSK          string      `json:"psk"`
	Status       string      `json:"status"`
	BillingCycle string      `json:"billingcycle"`
	CreatedAt    string      `json:"created_at"`
	CloudID      string      `json:"cloudid"`
}

type IPSecConnection struct {
	ID               string      `json:"id"`
	IPSecID          string      `json:"ipsecid"`
	Name             string      `json:"name"`
	RemoteIP         string      `json:"remote_ip"`
	RemoteLocalIP    string      `json:"remote_local_ip"`
	LocalIP          string      `json:"local_ip"`
	PSK              string      `json:"psk"`
	Phase1Encryption string      `json:"phase1-encryption"`
	Phase2Encryption string      `json:"phase2-encryption"`
	Phase1Integrity  string      `json:"phase1-integrity"`
	Phase2Integrity  string      `json:"phase2-integrity"`
	Phase1DHGroup    interface{} `json:"phase1-dh-group"`
	Phase2DHGroup    interface{} `json:"phase2-dh-group"`
	IKEVersion       string      `json:"ike_version"`
	Phase1Lifetime   string      `json:"phase1_lifetime"`
	Phase2Lifetime   string      `json:"phase2_lifetime"`
	RekeyMargin      string      `json:"rekey_margin"`
	RekeyFuzz        string      `json:"rekey_fuzz"`
	ReplayWindow     string      `json:"replay_window"`
	DPDTimeout       string      `json:"dpd_timeout"`
	DPDAction        string      `json:"dpd_action"`
	StartupAction    string      `json:"startup_action"`
	Status           string      `json:"status"`
	CreatedAt        string      `json:"created_at"`
}

// ── IPSec Tunnel methods ──────────────────────────────────────────────────

func (c *Client) CreateIPSec(req *IPSecCreateRequest) (string, error) {
	respBytes, err := c.Post("/ipsec?action=create", req)
	if err != nil {
		return "", fmt.Errorf("failed to create IPSec tunnel: %w", err)
	}
	var result map[string]interface{}
	if err := json.Unmarshal(respBytes, &result); err != nil {
		return "", fmt.Errorf("failed to parse create IPSec response: %w", err)
	}
	if result["status"] != "success" {
		return "", fmt.Errorf("create IPSec failed: %s", result["message"])
	}
	switch v := result["id"].(type) {
	case float64:
		return fmt.Sprintf("%.0f", v), nil
	default:
		return fmt.Sprintf("%v", v), nil
	}
}

func (c *Client) GetIPSec(ipsecID string) (*IPSecInstance, error) {
	respBytes, err := c.Get(fmt.Sprintf("/ipsec?action=info&ipsecid=%s", ipsecID))
	if err != nil {
		return nil, fmt.Errorf("failed to get IPSec tunnel: %w", err)
	}
	var result struct {
		Status string        `json:"status"`
		Data   IPSecInstance `json:"data"`
	}
	if err := json.Unmarshal(respBytes, &result); err != nil {
		return nil, fmt.Errorf("failed to parse get IPSec response: %w", err)
	}
	if result.Data.ID == "" {
		return nil, nil
	}
	return &result.Data, nil
}

func (c *Client) DeleteIPSec(ipsecID string) error {
	respBytes, err := c.Delete(fmt.Sprintf("/ipsec?action=delete&ipsecid=%s", ipsecID))
	if err != nil {
		return fmt.Errorf("failed to delete IPSec tunnel: %w", err)
	}
	var result map[string]string
	if err := json.Unmarshal(respBytes, &result); err != nil {
		return fmt.Errorf("failed to parse delete IPSec response: %w", err)
	}
	if result["status"] != "success" {
		return fmt.Errorf("delete IPSec failed: %s", result["message"])
	}
	return nil
}

// ── IPSec Pair (Connection) methods ──────────────────────────────────────

func (c *Client) WaitForIPSecReady(id string) error {
	for attempt := 0; attempt < 60; attempt++ {
		if attempt > 0 {
			time.Sleep(10 * time.Second)
		}
		ipsec, err := c.GetIPSec(id)
		if err != nil || ipsec == nil {
			continue
		}
		if ipsec.Status == "active" {
			return nil
		}
		if ipsec.Status == "failed" || ipsec.Status == "error" {
			return fmt.Errorf("IPSec tunnel entered failed state")
		}
	}
	return fmt.Errorf("STILL_PROVISIONING")
}

func (c *Client) CreateIPSecPair(req *IPSecPairRequest) (string, error) {
	respBytes, err := c.Post("/ipsec?action=pair", req)
	if err != nil {
		return "", fmt.Errorf("failed to create IPSec pair: %w", err)
	}
	var result map[string]interface{}
	if err := json.Unmarshal(respBytes, &result); err != nil {
		return "", fmt.Errorf("failed to parse create IPSec pair response: %w", err)
	}
	if result["status"] != "success" {
		return "", fmt.Errorf("create IPSec pair failed: %s", result["message"])
	}
	// Extract connection_id from side_a
	if sideA, ok := result["side_a"].(map[string]interface{}); ok {
		switch v := sideA["connection_id"].(type) {
		case float64:
			return fmt.Sprintf("%.0f", v), nil
		default:
			return fmt.Sprintf("%v", v), nil
		}
	}
	return "", fmt.Errorf("could not find connection_id in response")
}

func (c *Client) UpdateIPSecConnection(req *IPSecConnectionUpdateRequest) error {
	respBytes, err := c.Put("/ipsec?action=connection", req)
	if err != nil {
		return fmt.Errorf("failed to update IPSec connection: %w", err)
	}
	var result map[string]string
	if err := json.Unmarshal(respBytes, &result); err != nil {
		return fmt.Errorf("failed to parse update IPSec connection response: %w", err)
	}
	if result["status"] != "success" {
		return fmt.Errorf("update IPSec connection failed: %s", result["message"])
	}
	return nil
}

func (c *Client) ListIPSecConnections(ipsecID string) ([]IPSecConnection, error) {
	respBytes, err := c.Get(fmt.Sprintf("/ipsec?action=connection&ipsecid=%s", ipsecID))
	if err != nil {
		return nil, fmt.Errorf("failed to list IPSec connections: %w", err)
	}
	var result struct {
		Status string            `json:"status"`
		Data   []IPSecConnection `json:"data"`
	}
	if err := json.Unmarshal(respBytes, &result); err != nil {
		return nil, fmt.Errorf("failed to parse list IPSec connections response: %w", err)
	}
	return result.Data, nil
}

func (c *Client) DeleteIPSecConnection(ipsecID string, connectionID string) error {
	endpoint := fmt.Sprintf("/ipsec?action=connection&ipsecid=%s&id=%s", ipsecID, connectionID)
	respBytes, err := c.Delete(endpoint)
	if err != nil {
		return fmt.Errorf("failed to delete IPSec connection: %w", err)
	}
	var result map[string]string
	if err := json.Unmarshal(respBytes, &result); err != nil {
		return fmt.Errorf("failed to parse delete IPSec connection response: %w", err)
	}
	if result["status"] != "success" {
		return fmt.Errorf("delete IPSec connection failed: %s", result["message"])
	}
	return nil
}
