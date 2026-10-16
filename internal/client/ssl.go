package client

import (
	"encoding/json"
	"fmt"
)

// ── Request structs ───────────────────────────────────────────────────────

type SSLCertificateCreateRequest struct {
	Name             string `json:"name"`
	CertificateKey   string `json:"certificate_key"`
	PrivateKey       string `json:"private_key"`
	CertificateChain string `json:"certificateChain"`
	Type             string `json:"type"`
}

// ── Response structs ──────────────────────────────────────────────────────

type SSLCertificateInstance struct {
	ID            string `json:"id"`
	Name          string `json:"name"`
	Mode          string `json:"mode"`
	Type          string `json:"type"`
	Status        string `json:"status"`
	PrimaryDomain string `json:"primary_domain"`
	IsWildcard    string `json:"is_wildcard"`
	AutoRenew     string `json:"auto_renew"`
	CreatedAt     string `json:"created_at"`
	ExpireAt      string `json:"expire_at"`
	Issuer        string `json:"issuer"`
}

type SSLCertificateListResponse struct {
	Status       string                   `json:"status"`
	Certificates []SSLCertificateInstance `json:"certificates"`
}

// ── API methods ───────────────────────────────────────────────────────────

func (c *Client) CreateSSLCertificate(req *SSLCertificateCreateRequest) (string, error) {
	respBytes, err := c.Post("/certificates", req)
	if err != nil {
		return "", fmt.Errorf("failed to create SSL certificate: %w", err)
	}
	var result map[string]interface{}
	if err := json.Unmarshal(respBytes, &result); err != nil {
		return "", fmt.Errorf("failed to parse create SSL certificate response: %w", err)
	}
	if result["status"] != "success" {
		return "", fmt.Errorf("create SSL certificate failed: %s", result["message"])
	}
	id, _ := result["id"].(string)
	return id, nil
}

func (c *Client) GetSSLCertificate(certID string) (*SSLCertificateInstance, error) {
	endpoint := fmt.Sprintf("/certificates/certificate?id=%s&ssl_view=details&include_sensitive=0", certID)
	respBytes, err := c.Get(endpoint)
	if err != nil {
		return nil, fmt.Errorf("failed to get SSL certificate: %w", err)
	}
	var resp SSLCertificateListResponse
	if err := json.Unmarshal(respBytes, &resp); err != nil {
		return nil, fmt.Errorf("failed to parse get SSL certificate response: %w", err)
	}
	if len(resp.Certificates) == 0 {
		return nil, nil
	}
	return &resp.Certificates[0], nil
}

func (c *Client) ListSSLCertificates() ([]SSLCertificateInstance, error) {
	respBytes, err := c.Get("/certificates")
	if err != nil {
		return nil, fmt.Errorf("failed to list SSL certificates: %w", err)
	}
	var resp SSLCertificateListResponse
	if err := json.Unmarshal(respBytes, &resp); err != nil {
		return nil, fmt.Errorf("failed to parse list SSL certificates response: %w", err)
	}
	return resp.Certificates, nil
}

func (c *Client) DeleteSSLCertificate(certID string) error {
	endpoint := fmt.Sprintf("/certificates/certificate?id=%s", certID)
	respBytes, err := c.Delete(endpoint)
	if err != nil {
		return fmt.Errorf("failed to delete SSL certificate: %w", err)
	}
	var result map[string]string
	if err := json.Unmarshal(respBytes, &result); err != nil {
		return fmt.Errorf("failed to parse delete SSL certificate response: %w", err)
	}
	if result["status"] != "success" {
		return fmt.Errorf("delete SSL certificate failed: %s", result["message"])
	}
	return nil
}
