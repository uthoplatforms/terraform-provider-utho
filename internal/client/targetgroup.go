package client

import (
	"encoding/json"
	"fmt"
)

// ── Request structs ───────────────────────────────────────────────────────

type TargetGroupCreateRequest struct {
	Name                string `json:"name"`
	Protocol            string `json:"protocol"`
	Port                string `json:"port"`
	HealthCheckPath     string `json:"health_check_path"`
	HealthCheckProtocol string `json:"health_check_protocol"`
	HealthCheckInterval string `json:"health_check_interval"`
	HealthCheckTimeout  string `json:"health_check_timeout"`
	HealthyThreshold    string `json:"healthy_threshold"`
	UnhealthyThreshold  string `json:"unhealthy_threshold"`
}

type TargetGroupUpdateRequest struct {
	Name                string `json:"name"`
	Protocol            string `json:"protocol"`
	Port                string `json:"port"`
	HealthCheckPath     string `json:"health_check_path"`
	HealthCheckProtocol string `json:"health_check_protocol"`
	HealthCheckInterval string `json:"health_check_interval"`
	HealthCheckTimeout  string `json:"health_check_timeout"`
	HealthyThreshold    string `json:"healthy_threshold"`
	UnhealthyThreshold  string `json:"unhealthy_threshold"`
}

type TargetGroupAddTargetRequest struct {
	BackendPort         string `json:"backend_port"`
	Type                string `json:"type"`
	CloudID             string `json:"cloudid,omitempty"`
	IP                  string `json:"ip,omitempty"`
	KubernetesClusterID string `json:"kubernetes_clusterid,omitempty"`
	PoolName            string `json:"pool_name,omitempty"`
}

// ── Response structs ──────────────────────────────────────────────────────

type TargetGroupTarget struct {
	ID                  string `json:"id"`
	IP                  string `json:"ip"`
	CloudID             string `json:"cloudid"`
	KubernetesClusterID string `json:"kubernetes_clusterid"`
	BackendPort         string `json:"backend_port"`
	BackendProtocol     string `json:"backend_protocol"`
	Status              string `json:"status"`
	TargetGroupID       string `json:"targetgroup_id"`
}

type TargetGroup struct {
	ID                  string              `json:"id"`
	Name                string              `json:"name"`
	Port                string              `json:"port"`
	Protocol            string              `json:"protocol"`
	HealthCheckPath     string              `json:"health_check_path"`
	HealthCheckInterval string              `json:"health_check_interval"`
	HealthCheckProtocol string              `json:"health_check_protocol"`
	HealthCheckTimeout  string              `json:"health_check_timeout"`
	HealthyThreshold    string              `json:"healthy_threshold"`
	UnhealthyThreshold  string              `json:"unhealthy_threshold"`
	CreatedAt           string              `json:"created_at"`
	Targets             []TargetGroupTarget `json:"targets"`
}

// ── Client methods ────────────────────────────────────────────────────────

func (c *Client) CreateTargetGroup(req *TargetGroupCreateRequest) (string, error) {
	respBytes, err := c.Post("/targetgroup", req)
	if err != nil {
		return "", fmt.Errorf("failed to create target group: %w", err)
	}
	var result map[string]interface{}
	if err := json.Unmarshal(respBytes, &result); err != nil {
		return "", fmt.Errorf("failed to parse create target group response: %w", err)
	}
	if result["status"] != "success" {
		return "", fmt.Errorf("create target group failed: %s", result["message"])
	}
	switch v := result["id"].(type) {
	case float64:
		return fmt.Sprintf("%.0f", v), nil
	default:
		return fmt.Sprintf("%v", v), nil
	}
}

func (c *Client) GetTargetGroup(id string) (*TargetGroup, error) {
	respBytes, err := c.Get(fmt.Sprintf("/targetgroup/%s", id))
	if err != nil {
		return nil, fmt.Errorf("failed to get target group: %w", err)
	}
	var result struct {
		Status      string      `json:"status"`
		TargetGroup TargetGroup `json:"targetgroup"`
	}
	if err := json.Unmarshal(respBytes, &result); err != nil {
		return nil, fmt.Errorf("failed to parse get target group response: %w", err)
	}
	if result.TargetGroup.ID == "" {
		return nil, nil
	}
	return &result.TargetGroup, nil
}

func (c *Client) UpdateTargetGroup(id string, req *TargetGroupUpdateRequest) error {
	respBytes, err := c.Put(fmt.Sprintf("/targetgroup/%s", id), req)
	if err != nil {
		return fmt.Errorf("failed to update target group: %w", err)
	}
	var result map[string]interface{}
	if err := json.Unmarshal(respBytes, &result); err != nil {
		return fmt.Errorf("failed to parse update target group response: %w", err)
	}
	if result["status"] != "success" {
		return fmt.Errorf("update target group failed: %s", result["message"])
	}
	return nil
}

func (c *Client) DeleteTargetGroup(id string, name string) error {
	respBytes, err := c.Delete(fmt.Sprintf("/targetgroup/%s?name=%s", id, name))
	if err != nil {
		return fmt.Errorf("failed to delete target group: %w", err)
	}
	var result map[string]interface{}
	if err := json.Unmarshal(respBytes, &result); err != nil {
		return fmt.Errorf("failed to parse delete target group response: %w", err)
	}
	if result["status"] != "success" {
		return fmt.Errorf("delete target group failed: %s", result["message"])
	}
	return nil
}

func (c *Client) AddTargetGroupTarget(tgID string, req *TargetGroupAddTargetRequest) error {
	respBytes, err := c.Post(fmt.Sprintf("/targetgroup/%s/target", tgID), req)
	if err != nil {
		return fmt.Errorf("failed to add target: %w", err)
	}
	var result map[string]interface{}
	if err := json.Unmarshal(respBytes, &result); err != nil {
		return fmt.Errorf("failed to parse add target response: %w", err)
	}
	if result["status"] != "success" {
		return fmt.Errorf("add target failed: %s", result["message"])
	}
	return nil
}

func (c *Client) DeleteTargetGroupTarget(tgID string, targetID string) error {
	respBytes, err := c.Delete(fmt.Sprintf("/targetgroup/%s/target/%s", tgID, targetID))
	if err != nil {
		return fmt.Errorf("failed to delete target: %w", err)
	}
	var result map[string]interface{}
	if err := json.Unmarshal(respBytes, &result); err != nil {
		return fmt.Errorf("failed to parse delete target response: %w", err)
	}
	if result["status"] != "success" {
		return fmt.Errorf("delete target failed: %s", result["message"])
	}
	return nil
}
