package client

import (
	"encoding/json"
	"fmt"
)

// ── Request structs ───────────────────────────────────────────────────────

type DNSZoneCreateRequest struct {
	Domain string `json:"domain"`
}

type DNSRecordRequest struct {
	Type     string `json:"type"`
	Hostname string `json:"hostname"`
	Value    string `json:"value"`
	TTL      string `json:"ttl"`
}

// ── Response structs ──────────────────────────────────────────────────────

type DNSRecord struct {
	ID       string `json:"id"`
	Hostname string `json:"hostname"`
	Type     string `json:"type"`
	Value    string `json:"value"`
	TTL      string `json:"ttl"`
}

type DNSZone struct {
	Domain      string      `json:"domain"`
	NSPoint     string      `json:"nspoint"`
	CreatedAt   string      `json:"created_at"`
	RecordCount string      `json:"dnsrecord_count"`
	Records     []DNSRecord `json:"records"`
}

type DNSZoneListResponse struct {
	Domains []DNSZone `json:"domains"`
}

// ── DNS Zone methods ──────────────────────────────────────────────────────

func (c *Client) CreateDNSZone(domain string) error {
	respBytes, err := c.Post("/dns/adddomain", &DNSZoneCreateRequest{Domain: domain})
	if err != nil {
		return fmt.Errorf("failed to create DNS zone: %w", err)
	}
	var result map[string]string
	if err := json.Unmarshal(respBytes, &result); err != nil {
		return fmt.Errorf("failed to parse create DNS zone response: %w", err)
	}
	if result["status"] != "success" {
		return fmt.Errorf("create DNS zone failed: %s", result["message"])
	}
	return nil
}

func (c *Client) GetDNSZone(domain string) (*DNSZone, error) {
	respBytes, err := c.Get(fmt.Sprintf("/dns/%s", domain))
	if err != nil {
		return nil, fmt.Errorf("failed to get DNS zone: %w", err)
	}
	var resp DNSZoneListResponse
	if err := json.Unmarshal(respBytes, &resp); err != nil {
		return nil, fmt.Errorf("failed to parse get DNS zone response: %w", err)
	}
	if len(resp.Domains) == 0 {
		return nil, nil
	}
	return &resp.Domains[0], nil
}

func (c *Client) DeleteDNSZone(domain string) error {
	respBytes, err := c.Delete(fmt.Sprintf("/dns/%s/delete", domain))
	if err != nil {
		return fmt.Errorf("failed to delete DNS zone: %w", err)
	}
	var result map[string]string
	if err := json.Unmarshal(respBytes, &result); err != nil {
		return fmt.Errorf("failed to parse delete DNS zone response: %w", err)
	}
	if result["status"] != "success" {
		return fmt.Errorf("delete DNS zone failed: %s", result["message"])
	}
	return nil
}

// ── DNS Record methods ────────────────────────────────────────────────────

func (c *Client) CreateDNSRecord(domain string, req *DNSRecordRequest) (string, error) {
	respBytes, err := c.Post(fmt.Sprintf("/dns/%s/record/add", domain), req)
	if err != nil {
		return "", fmt.Errorf("failed to create DNS record: %w", err)
	}
	var result map[string]string
	if err := json.Unmarshal(respBytes, &result); err != nil {
		return "", fmt.Errorf("failed to parse create DNS record response: %w", err)
	}
	if result["status"] != "success" {
		return "", fmt.Errorf("create DNS record failed: %s", result["message"])
	}
	return result["id"], nil
}

func (c *Client) UpdateDNSRecord(domain string, recordID string, req *DNSRecordRequest) error {
	respBytes, err := c.Post(fmt.Sprintf("/dns/%s/record/%s/update", domain, recordID), req)
	if err != nil {
		return fmt.Errorf("failed to update DNS record: %w", err)
	}
	var result map[string]string
	if err := json.Unmarshal(respBytes, &result); err != nil {
		return fmt.Errorf("failed to parse update DNS record response: %w", err)
	}
	if result["status"] != "success" {
		return fmt.Errorf("update DNS record failed: %s", result["message"])
	}
	return nil
}

func (c *Client) DeleteDNSRecord(domain string, recordID string) error {
	respBytes, err := c.Delete(fmt.Sprintf("/dns/%s/record/%s/delete", domain, recordID))
	if err != nil {
		return fmt.Errorf("failed to delete DNS record: %w", err)
	}
	var result map[string]string
	if err := json.Unmarshal(respBytes, &result); err != nil {
		return fmt.Errorf("failed to parse delete DNS record response: %w", err)
	}
	if result["status"] != "success" {
		return fmt.Errorf("delete DNS record failed: %s", result["message"])
	}
	return nil
}
