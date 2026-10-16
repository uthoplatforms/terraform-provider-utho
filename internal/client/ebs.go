package client

import (
	"encoding/json"
	"fmt"
)

// ── Request structs ───────────────────────────────────────────────────────

type EBSCreateRequest struct {
	DCSlug     string `json:"dcslug"`
	DiskType   string `json:"disk_type"`
	Name       string `json:"name"`
	Disk       string `json:"disk"`
	IOPS       string `json:"iops"`
	Throughput string `json:"throughput"`
}

type EBSUpdateNameRequest struct {
	Name string `json:"name"`
}

type EBSResizeRequest struct {
	Disk       string `json:"disk"`
	IOPS       string `json:"iops"`
	Throughput string `json:"throughput"`
}

type EBSAttachRequest struct {
	ResourceID string `json:"resourceid"`
	Type       string `json:"type"`
}

// ── Response structs ──────────────────────────────────────────────────────

type EBSInstance struct {
	DID        string `json:"did"`
	CloudID    string `json:"cloudid"`
	Name       string `json:"name"`
	Size       string `json:"size"`
	Status     string `json:"status"`
	IOPS       string `json:"iops"`
	Throughput string `json:"throughput"`
	CreatedAt  string `json:"created_at"`
	DCSlug     string `json:"dc"`
	Location   struct {
		City    string `json:"city"`
		Country string `json:"country"`
		DC      string `json:"dc"`
		DCCC    string `json:"dccc"`
	} `json:"location"`
}

// ── Client methods ────────────────────────────────────────────────────────

func (c *Client) CreateEBS(req *EBSCreateRequest) (string, error) {
	respBytes, err := c.Post("/ebs", req)
	if err != nil {
		return "", fmt.Errorf("failed to create EBS volume: %w", err)
	}
	var result map[string]interface{}
	if err := json.Unmarshal(respBytes, &result); err != nil {
		return "", fmt.Errorf("failed to parse create EBS response: %w", err)
	}
	if result["status"] != "success" {
		return "", fmt.Errorf("create EBS failed: %s", result["message"])
	}
	switch v := result["id"].(type) {
	case float64:
		return fmt.Sprintf("%.0f", v), nil
	default:
		return fmt.Sprintf("%v", v), nil
	}
}

func (c *Client) GetEBS(ebsID string) (*EBSInstance, error) {
	respBytes, err := c.Get(fmt.Sprintf("/ebs/%s", ebsID))
	if err != nil {
		return nil, fmt.Errorf("failed to get EBS volume: %w", err)
	}
	var result struct {
		Status string        `json:"status"`
		EBS    []EBSInstance `json:"ebs"`
	}
	if err := json.Unmarshal(respBytes, &result); err != nil {
		return nil, fmt.Errorf("failed to parse get EBS response: %w", err)
	}
	if len(result.EBS) == 0 {
		return nil, nil
	}
	return &result.EBS[0], nil
}

func (c *Client) UpdateEBSName(ebsID string, req *EBSUpdateNameRequest) error {
	respBytes, err := c.Put(fmt.Sprintf("/ebs/%s/updatename", ebsID), req)
	if err != nil {
		return fmt.Errorf("failed to update EBS name: %w", err)
	}
	// API has no response body — ignore empty response
	if len(respBytes) == 0 {
		return nil
	}
	var result map[string]interface{}
	if err := json.Unmarshal(respBytes, &result); err != nil {
		return nil
	}
	if status, ok := result["status"].(string); ok && status != "success" {
		return fmt.Errorf("update EBS name failed: %s", result["message"])
	}
	return nil
}

func (c *Client) ResizeEBS(ebsID string, req *EBSResizeRequest) error {
	respBytes, err := c.Put(fmt.Sprintf("/ebs/%s/resize", ebsID), req)
	if err != nil {
		return fmt.Errorf("failed to resize EBS volume: %w", err)
	}
	if len(respBytes) == 0 {
		return nil
	}
	var result map[string]interface{}
	if err := json.Unmarshal(respBytes, &result); err != nil {
		return nil
	}
	if status, ok := result["status"].(string); ok && status != "success" {
		return fmt.Errorf("resize EBS failed: %s", result["message"])
	}
	return nil
}

func (c *Client) AttachEBSVolume(ebsID string, req *EBSAttachRequest) error {
	respBytes, err := c.Put(fmt.Sprintf("/ebs/%s/attach", ebsID), req)
	if err != nil {
		return fmt.Errorf("failed to attach EBS volume: %w", err)
	}
	var result map[string]interface{}
	if err := json.Unmarshal(respBytes, &result); err != nil {
		return fmt.Errorf("failed to parse attach EBS response: %w", err)
	}
	if result["status"] != "success" {
		return fmt.Errorf("attach EBS failed: %s", result["message"])
	}
	return nil
}

func (c *Client) DetachEBSVolume(ebsID string, req *EBSAttachRequest) error {
	respBytes, err := c.Put(fmt.Sprintf("/ebs/%s/dettach", ebsID), req)
	if err != nil {
		return fmt.Errorf("failed to detach EBS volume: %w", err)
	}
	var result map[string]interface{}
	if err := json.Unmarshal(respBytes, &result); err != nil {
		return fmt.Errorf("failed to parse detach EBS response: %w", err)
	}
	if result["status"] != "success" {
		return fmt.Errorf("detach EBS failed: %s", result["message"])
	}
	return nil
}

func (c *Client) DeleteEBS(ebsID string) error {
	respBytes, err := c.Delete(fmt.Sprintf("/ebs/%s/destroy", ebsID))
	if err != nil {
		return fmt.Errorf("failed to delete EBS volume: %w", err)
	}
	var result map[string]interface{}
	if err := json.Unmarshal(respBytes, &result); err != nil {
		return fmt.Errorf("failed to parse delete EBS response: %w", err)
	}
	if result["status"] != "success" {
		return fmt.Errorf("delete EBS failed: %s", result["message"])
	}
	return nil
}
