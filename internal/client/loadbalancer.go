package client

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"
)

// ── Request structs ───────────────────────────────────────────────────────

type LoadBalancerCreateRequest struct {
	Type           string `json:"type"`
	Name           string `json:"name"`
	DCSlug         string `json:"dcslug"`
	VPC            string `json:"vpc,omitempty"`
	EnablePublicIP string `json:"enable_publicip"`
	Firewall       string `json:"firewall,omitempty"`
}

type LoadBalancerFrontendRequest struct {
	Name          string `json:"name"`
	Algorithm     string `json:"algorithm"`
	Proto         string `json:"proto"`
	Port          string `json:"port"`
	Cookie        string `json:"cookie"`
	RedirectHTTPS string `json:"redirecthttps"`
	CertificateID string `json:"certificate_id"`
	CookieName    string `json:"cookiename,omitempty"`
}

type LoadBalancerBackendRequest struct {
	FrontendID  string `json:"frontend_id"`
	BackendPort string `json:"backend_port"`
	Weight      string `json:"weight"`
	Type        string `json:"type"`
	CloudID     string `json:"cloudid,omitempty"`
	IP          string `json:"ip,omitempty"`
}

type LoadBalancerACLRequest struct {
	FrontendID    string `json:"frontend_id"`
	Name          string `json:"name"`
	ConditionType string `json:"conditionType"`
	Value         string `json:"value"`
}

type LoadBalancerSettingsRequest struct {
	TimeoutConnect       string `json:"timeout_connect"`
	TimeoutClient        string `json:"timeout_client"`
	TimeoutServer        string `json:"timeout_server"`
	TimeoutHTTPRequest   string `json:"timeout_http_request"`
	TimeoutHTTPKeepalive string `json:"timeout_http_keepalive"`
	TimeoutTunnel        string `json:"timeout_tunnel"`
	MaxConnections       string `json:"max_connections"`
	HTTP2                string `json:"http2"`
	Compression          string `json:"compression"`
	HSTS                 string `json:"hsts"`
}

// ── Response structs ──────────────────────────────────────────────────────

type LoadBalancerInstance struct {
	ID           string `json:"id"`
	Name         string `json:"name"`
	Type         string `json:"type"`
	IP           string `json:"ip"`
	DNS          string `json:"dns"`
	DCSlug       string `json:"dcslug"`
	Status       string `json:"status"`
	BackendCount string `json:"backendcount"`
	CreatedAt    string `json:"created_at"`
}

type LoadBalancerListResponse struct {
	LoadBalancers []LoadBalancerInstance `json:"loadbalancers"`
}

// ── API methods ───────────────────────────────────────────────────────────

// checkLBReady polls until the LB is ready (status=Active, app_status=Installed).
// Polls every 10 seconds for up to 5 minutes before giving up.
func (c *Client) checkLBReady(lbID string) error {
	for attempt := 0; attempt < 30; attempt++ {
		if attempt > 0 {
			time.Sleep(10 * time.Second)
		}

		respBytes, err := c.Get(fmt.Sprintf("/loadbalancer/%s", lbID))
		if err != nil {
			continue
		}
		var resp map[string]interface{}
		if json.Unmarshal(respBytes, &resp) != nil {
			continue
		}
		lbs, _ := resp["loadbalancers"].([]interface{})
		if len(lbs) == 0 {
			continue
		}
		lb, _ := lbs[0].(map[string]interface{})
		status := fmt.Sprintf("%v", lb["status"])
		appStatus := fmt.Sprintf("%v", lb["app_status"])

		// Ready
		if status == "Active" && appStatus == "Installed" {
			return nil
		}

		// Failed — no point waiting
		if appStatus == "Failed" {
			return fmt.Errorf(
				"load balancer entered a Failed state — please check the Utho Console and resolve the issue, then run terraform apply again",
			)
		}

		// Still pending — continue polling
	}

	return fmt.Errorf(
		"load balancer did not become ready after 5 minutes — please check the Utho Console and try again",
	)
}

func (c *Client) CreateLoadBalancer(req *LoadBalancerCreateRequest) (string, error) {
	respBytes, err := c.Post("/loadbalancer", req)
	if err != nil {
		// On 504/timeout, LB may have been created — look it up by name
		if strings.Contains(err.Error(), "504") || strings.Contains(err.Error(), "timeout") || strings.Contains(err.Error(), "connection reset") {
			time.Sleep(5 * time.Second)
			id, lookupErr := c.findLBByName(req.Name)
			if lookupErr == nil && id != "" {
				return id, nil
			}
		}
		return "", fmt.Errorf("failed to create load balancer: %w", err)
	}
	var result map[string]interface{}
	if err := json.Unmarshal(respBytes, &result); err != nil {
		return "", fmt.Errorf("failed to parse create load balancer response: %w", err)
	}
	if result["status"] != "success" {
		return "", fmt.Errorf("create load balancer failed: %s", result["message"])
	}
	id := fmt.Sprintf("%v", result["loadbalancerid"])
	return id, nil
}

// findLBByName looks up a load balancer ID by name — used to recover from 504 timeouts
func (c *Client) findLBByName(name string) (string, error) {
	respBytes, err := c.Get("/loadbalancer")
	if err != nil {
		return "", err
	}
	var resp LoadBalancerListResponse
	if err := json.Unmarshal(respBytes, &resp); err != nil {
		return "", err
	}
	for _, lb := range resp.LoadBalancers {
		if lb.Name == name {
			return lb.ID, nil
		}
	}
	return "", fmt.Errorf("LB not found")
}

func (c *Client) GetLoadBalancer(lbID string) (*LoadBalancerInstance, error) {
	respBytes, err := c.Get(fmt.Sprintf("/loadbalancer/%s", lbID))
	if err != nil {
		return nil, fmt.Errorf("failed to get load balancer: %w", err)
	}
	var resp LoadBalancerListResponse
	if err := json.Unmarshal(respBytes, &resp); err != nil {
		return nil, fmt.Errorf("failed to parse get load balancer response: %w", err)
	}
	if len(resp.LoadBalancers) == 0 {
		return nil, nil
	}
	return &resp.LoadBalancers[0], nil
}

func (c *Client) DeleteLoadBalancer(lbID string) error {
	respBytes, err := c.Delete(fmt.Sprintf("/loadbalancer/%s", lbID))
	if err != nil {
		if strings.Contains(err.Error(), "404") || strings.Contains(err.Error(), "not found") {
			return nil
		}
		return fmt.Errorf("failed to delete load balancer: %w", err)
	}
	trimmed := strings.TrimSpace(string(respBytes))
	if trimmed == "" || trimmed == "null" {
		return nil
	}
	var result map[string]string
	if err := json.Unmarshal(respBytes, &result); err != nil {
		return nil
	}
	if result["status"] != "success" {
		return fmt.Errorf("delete load balancer failed: %s", result["message"])
	}
	return nil
}

func (c *Client) AddLBFrontend(lbID string, req *LoadBalancerFrontendRequest) (string, error) {
	// Check LB is ready before adding frontend
	if err := c.checkLBReady(lbID); err != nil {
		return "", err
	}
	endpoint := fmt.Sprintf("/loadbalancer/%s/frontend", lbID)
	for attempt := 0; attempt < 6; attempt++ {
		respBytes, err := c.Post(endpoint, req)
		if err != nil {
			return "", fmt.Errorf("failed to add frontend: %w", err)
		}
		var result map[string]interface{}
		if err := json.Unmarshal(respBytes, &result); err != nil {
			return "", fmt.Errorf("failed to parse add frontend response: %w", err)
		}
		if result["status"] == "success" {
			return fmt.Sprintf("%v", result["id"]), nil
		}
		msg := fmt.Sprintf("%v", result["message"])
		if strings.Contains(msg, "pending action") || strings.Contains(msg, "in process") {
			time.Sleep(10 * time.Second)
			continue
		}
		return "", fmt.Errorf("add frontend failed: %s", msg)
	}
	return "", fmt.Errorf("add frontend failed: LB still processing after retries")
}

func (c *Client) DeleteLBFrontend(lbID string, frontendID string) error {
	endpoint := fmt.Sprintf("/loadbalancer/%s/frontend/%s", lbID, frontendID)

	// Retry up to 5 times to handle "pending action" errors
	for attempt := 0; attempt < 5; attempt++ {
		respBytes, err := c.Delete(endpoint)
		if err != nil {
			// 404 = already deleted
			if strings.Contains(err.Error(), "404") || strings.Contains(err.Error(), "not found") {
				return nil
			}
			return fmt.Errorf("failed to delete frontend: %w", err)
		}

		trimmed := strings.TrimSpace(string(respBytes))
		if trimmed == "" || trimmed == "null" {
			return nil
		}

		var result map[string]string
		if err := json.Unmarshal(respBytes, &result); err != nil {
			return nil
		}

		if result["status"] == "success" {
			return nil
		}

		msg := result["message"]
		// Not found = already deleted
		if strings.Contains(msg, "not found") || strings.Contains(msg, "Not Found") {
			return nil
		}
		// Pending action = retry after delay
		if strings.Contains(msg, "pending action") || strings.Contains(msg, "in process") {
			time.Sleep(10 * time.Second)
			continue
		}

		return fmt.Errorf("delete frontend failed: %s", msg)
	}

	return fmt.Errorf("delete frontend failed: LB still processing after retries")
}

// fetchBackendID looks up a backend ID from the LB by matching cloudID or IP
func (c *Client) fetchBackendID(lbID string, frontendID string, cloudID string, ip string) (string, error) {
	lb, err := c.Get(fmt.Sprintf("/loadbalancer/%s", lbID))
	if err != nil {
		return "unknown", nil
	}
	var lbResp map[string]interface{}
	if json.Unmarshal(lb, &lbResp) != nil {
		return "unknown", nil
	}
	lbs, _ := lbResp["loadbalancers"].([]interface{})
	if len(lbs) == 0 {
		return "unknown", nil
	}
	lbData, _ := lbs[0].(map[string]interface{})
	frontends, _ := lbData["frontends"].([]interface{})
	for _, fe := range frontends {
		feMap, _ := fe.(map[string]interface{})
		if fmt.Sprintf("%v", feMap["id"]) == frontendID {
			backends, _ := feMap["backends"].([]interface{})
			for _, be := range backends {
				beMap, _ := be.(map[string]interface{})
				if (cloudID != "" && fmt.Sprintf("%v", beMap["cloudid"]) == cloudID) ||
					(ip != "" && fmt.Sprintf("%v", beMap["ip"]) == ip) {
					return fmt.Sprintf("%v", beMap["id"]), nil
				}
			}
		}
	}
	return "unknown", nil
}

func (c *Client) AddLBBackend(lbID string, req *LoadBalancerBackendRequest) (string, error) {
	// Check LB is ready before adding backend
	if err := c.checkLBReady(lbID); err != nil {
		return "", err
	}
	endpoint := fmt.Sprintf("/loadbalancer/%s/backend", lbID)
	for attempt := 0; attempt < 6; attempt++ {
		respBytes, err := c.Post(endpoint, req)
		if err != nil {
			return "", fmt.Errorf("failed to add backend: %w", err)
		}
		var result map[string]interface{}
		if err := json.Unmarshal(respBytes, &result); err != nil {
			return "", fmt.Errorf("failed to parse add backend response: %w", err)
		}
		if result["status"] == "success" {
			// Try to get ID from response first
			if result["id"] != nil && fmt.Sprintf("%v", result["id"]) != "<nil>" && fmt.Sprintf("%v", result["id"]) != "" {
				return fmt.Sprintf("%v", result["id"]), nil
			}
			// API didn't return ID — fetch it from the LB
			return c.fetchBackendID(lbID, req.FrontendID, req.CloudID, req.IP)
		}
		msg := fmt.Sprintf("%v", result["message"])

		// Backend already exists (duplicate IP) — fetch its ID from the LB
		if duplicates, ok := result["duplicates"]; ok {
			dups, _ := duplicates.([]interface{})
			if len(dups) > 0 {
				return c.fetchBackendID(lbID, req.FrontendID, req.CloudID, req.IP)
			}
		}

		if strings.Contains(msg, "pending action") || strings.Contains(msg, "in process") {
			time.Sleep(10 * time.Second)
			continue
		}
		return "", fmt.Errorf("add backend failed: %s", msg)
	}
	return "", fmt.Errorf("add backend failed: LB still processing after retries")
}

func (c *Client) DeleteLBBackend(lbID string, backendID string) error {
	endpoint := fmt.Sprintf("/loadbalancer/%s/backend/%s", lbID, backendID)
	for attempt := 0; attempt < 5; attempt++ {
		respBytes, err := c.Delete(endpoint)
		if err != nil {
			if strings.Contains(err.Error(), "404") || strings.Contains(err.Error(), "not found") {
				return nil
			}
			return fmt.Errorf("failed to delete backend: %w", err)
		}
		trimmed := strings.TrimSpace(string(respBytes))
		if trimmed == "" || trimmed == "null" {
			return nil
		}
		var result map[string]string
		if err := json.Unmarshal(respBytes, &result); err != nil {
			return nil
		}
		if result["status"] == "success" {
			return nil
		}
		msg := result["message"]
		if strings.Contains(msg, "not found") || strings.Contains(msg, "Not Found") {
			return nil
		}
		if strings.Contains(msg, "pending action") || strings.Contains(msg, "in process") {
			time.Sleep(10 * time.Second)
			continue
		}
		return fmt.Errorf("delete backend failed: %s", msg)
	}
	return fmt.Errorf("delete backend failed: LB still processing after retries")
}

func (c *Client) AddLBACL(lbID string, req *LoadBalancerACLRequest) (string, error) {
	respBytes, err := c.Post(fmt.Sprintf("/loadbalancer/%s/acl", lbID), req)
	if err != nil {
		return "", fmt.Errorf("failed to add ACL rule: %w", err)
	}
	var result map[string]interface{}
	if err := json.Unmarshal(respBytes, &result); err != nil {
		return "", fmt.Errorf("failed to parse add ACL response: %w", err)
	}
	if result["status"] != "success" {
		return "", fmt.Errorf("add ACL failed: %s", result["message"])
	}
	id := fmt.Sprintf("%v", result["id"])
	return id, nil
}

func (c *Client) UpdateLBACL(lbID string, aclID string, req *LoadBalancerACLRequest) error {
	respBytes, err := c.Put(fmt.Sprintf("/loadbalancer/%s/acl/%s", lbID, aclID), req)
	if err != nil {
		return fmt.Errorf("failed to update ACL rule: %w", err)
	}
	var result map[string]string
	if err := json.Unmarshal(respBytes, &result); err != nil {
		return fmt.Errorf("failed to parse update ACL response: %w", err)
	}
	if result["status"] != "success" {
		return fmt.Errorf("update ACL failed: %s", result["message"])
	}
	return nil
}

func (c *Client) DeleteLBACL(lbID string, aclID string) error {
	endpoint := fmt.Sprintf("/loadbalancer/%s/acl/%s", lbID, aclID)
	for attempt := 0; attempt < 5; attempt++ {
		respBytes, err := c.Delete(endpoint)
		if err != nil {
			if strings.Contains(err.Error(), "404") || strings.Contains(err.Error(), "not found") {
				return nil
			}
			return fmt.Errorf("failed to delete ACL rule: %w", err)
		}
		trimmed := strings.TrimSpace(string(respBytes))
		if trimmed == "" || trimmed == "null" {
			return nil
		}
		var result map[string]string
		if err := json.Unmarshal(respBytes, &result); err != nil {
			return nil
		}
		if result["status"] == "success" {
			return nil
		}
		msg := result["message"]
		if strings.Contains(msg, "not found") || strings.Contains(msg, "Not Found") {
			return nil
		}
		if strings.Contains(msg, "pending action") || strings.Contains(msg, "in process") {
			time.Sleep(10 * time.Second)
			continue
		}
		return fmt.Errorf("delete ACL failed: %s", msg)
	}
	return fmt.Errorf("delete ACL failed: LB still processing after retries")
}

func (c *Client) UpdateLBSettings(lbID string, req *LoadBalancerSettingsRequest) error {
	respBytes, err := c.Put(fmt.Sprintf("/loadbalancer/%s/settings", lbID), req)
	if err != nil {
		return fmt.Errorf("failed to update load balancer settings: %w", err)
	}
	trimmed := strings.TrimSpace(string(respBytes))
	if trimmed == "" || trimmed == "null" {
		return nil
	}
	var result map[string]string
	if err := json.Unmarshal(respBytes, &result); err != nil {
		return nil
	}
	if result["status"] != "" && result["status"] != "success" {
		return fmt.Errorf("update settings failed: %s", result["message"])
	}
	return nil
}
