package client

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"
)

// ── Request structs ───────────────────────────────────────────────────────

type AutoScalingPolicy struct {
	Name     string `json:"name"`
	Type     string `json:"type"`
	Compare  string `json:"compare"`
	Value    string `json:"value"`
	Adjust   int    `json:"adjust"`
	Period   string `json:"period"`
	Cooldown string `json:"cooldown"`
}

type AutoScalingSchedule struct {
	Name               string `json:"name"`
	DesiredSize        string `json:"desiredsize"`
	Timezone           string `json:"timezone"`
	Recurrence         string `json:"recurrence"`
	RecurrenceDuration string `json:"recurrence_duration"`
	RecurrenceWeek     string `json:"recurrence_week"`
	SelectedTime       string `json:"selectedTime"`
	SelectedDate       string `json:"selectedDate"`
	StartDate          string `json:"start_date"`
}

type AutoScalingCreateRequest struct {
	Name               string                `json:"name"`
	OSDiskSize         int                   `json:"os_disk_size"`
	DCSlug             string                `json:"dcslug"`
	MinSize            string                `json:"minsize"`
	MaxSize            string                `json:"maxsize"`
	DesiredSize        string                `json:"desiredsize"`
	PlanID             string                `json:"planid"`
	PlanName           string                `json:"planname"`
	InstanceTemplateID string                `json:"instance_templateid"`
	ImageID            string                `json:"image_id"`
	ImageName          string                `json:"image_name"`
	PublicIPEnabled    int                   `json:"public_ip_enabled"`
	VPC                string                `json:"vpc,omitempty"`
	LoadBalancers      string                `json:"load_balancers,omitempty"`
	SecurityGroups     string                `json:"security_groups,omitempty"`
	TargetGroups       string                `json:"target_groups,omitempty"`
	SnapshotID         string                `json:"snapshotid,omitempty"`
	Stack              string                `json:"stack,omitempty"`
	StackID            string                `json:"stackid,omitempty"`
	StackImage         string                `json:"stackimage,omitempty"`
	BackupID           string                `json:"backupid,omitempty"`
	CPUModel           string                `json:"cpumodel,omitempty"`
	Policies           []AutoScalingPolicy   `json:"policies,omitempty"`
	Schedules          []AutoScalingSchedule `json:"schedules,omitempty"`
}

type AutoScalingUpdateRequest struct {
	Name            string `json:"name"`
	MinSize         string `json:"minsize"`
	MaxSize         string `json:"maxsize"`
	DesiredSize     string `json:"desiredsize"`
	PublicIPEnabled string `json:"public_ip_enabled,omitempty"`
	SnapshotID      string `json:"snapshotid,omitempty"`
	Stack           string `json:"stack,omitempty"`
	StackID         string `json:"stackid,omitempty"`
	StackImage      string `json:"stackimage,omitempty"`
}

type ScalingPolicyCreateRequest struct {
	Name        string `json:"name"`
	Type        string `json:"type"`
	Compare     string `json:"compare"`
	Value       string `json:"value"`
	Adjust      string `json:"adjust"`
	Period      string `json:"period"`
	Cooldown    string `json:"cooldown"`
	Product     string `json:"product"`
	ProductID   string `json:"productid"`
	ScalingType string `json:"scaling_type"`
}

type ScalingPolicyUpdateRequest struct {
	Name     string `json:"name"`
	Type     string `json:"type"`
	Compare  string `json:"compare"`
	Value    string `json:"value"`
	Adjust   int    `json:"adjust"`
	Period   string `json:"period"`
	Cooldown string `json:"cooldown"`
}

type ScalingScheduleUpdateRequest struct {
	GroupID     string `json:"groupid"`
	ID          string `json:"id"`
	Name        string `json:"name"`
	DesiredSize string `json:"desiredsize"`
	Timezone    string `json:"timezone"`
	Recurrence  string `json:"recurrence"`
	StartDate   string `json:"start_date"`
	Status      int    `json:"status"`
}

// ── Response structs ──────────────────────────────────────────────────────

type AutoScalingInstance struct {
	ID              string      `json:"id"`
	Name            string      `json:"name"`
	DCSlug          string      `json:"dcslug"`
	MinSize         string      `json:"minsize"`
	MaxSize         string      `json:"maxsize"`
	DesiredSize     string      `json:"desiredsize"`
	PlanID          string      `json:"planid"`
	PlanName        string      `json:"planname"`
	Status          string      `json:"status"`
	CreatedAt       string      `json:"created_at"`
	OSDiskSize      interface{} `json:"os_disk_size"`
	PublicIPEnabled interface{} `json:"public_ip_enabled"`
	ImageName       string      `json:"image_name"`
	SnapshotID      string      `json:"snapshotid"`
	Stack           string      `json:"stack"`
	StackID         string      `json:"stackid"`
	StackImage      string      `json:"stackimage"`
	VPC             interface{} `json:"vpc"`
	LoadBalancers   interface{} `json:"load_balancers"`
	SecurityGroups  interface{} `json:"security_groups"`
	TargetGroups    interface{} `json:"target_groups"`
	BackupID        string      `json:"backupid"`
	CPUModel        string      `json:"cpumodel"`
}

type AutoScalingListResponse struct {
	Status      string                `json:"status"`
	Autoscaling []AutoScalingInstance `json:"autoscaling"`
}

type AutoScalingDetailResponse struct {
	Status      string                `json:"status"`
	Autoscaling []AutoScalingInstance `json:"autoscaling"`
}

// ── Autoscaling Group methods ─────────────────────────────────────────────

func (c *Client) CreateAutoScaling(req *AutoScalingCreateRequest) (string, error) {
	respBytes, err := c.Post("/autoscaling", req)
	if err != nil {
		return "", fmt.Errorf("failed to create autoscaling group: %w", err)
	}
	var result map[string]interface{}
	if err := json.Unmarshal(respBytes, &result); err != nil {
		return "", fmt.Errorf("failed to parse create autoscaling response: %w", err)
	}
	if result["status"] != "success" {
		return "", fmt.Errorf("create autoscaling failed: %s", result["message"])
	}
	// ID comes back as float64 from JSON — format as integer
	switch v := result["id"].(type) {
	case float64:
		return fmt.Sprintf("%.0f", v), nil
	default:
		return fmt.Sprintf("%v", v), nil
	}
}

func (c *Client) WaitForASGReady(id string) error {
	for attempt := 0; attempt < 36; attempt++ { // 6 minutes max
		if attempt > 0 {
			time.Sleep(10 * time.Second)
		}
		asg, err := c.GetAutoScaling(id)
		if err != nil || asg == nil {
			continue
		}
		status := strings.ToLower(asg.Status)
		if status == "active" || status == "running" {
			return nil
		}
		if status == "failed" || status == "error" {
			return fmt.Errorf("auto scaling group entered failed state")
		}
	}
	return fmt.Errorf("STILL_PROVISIONING")
}

func (c *Client) UpdateAutoScaling(id string, req *AutoScalingUpdateRequest) error {
	respBytes, err := c.Put(fmt.Sprintf("/autoscaling/%s", id), req)
	if err != nil {
		return fmt.Errorf("failed to update auto scaling group: %w", err)
	}
	var result map[string]string
	if err := json.Unmarshal(respBytes, &result); err != nil {
		return fmt.Errorf("failed to parse update auto scaling response: %w", err)
	}
	if result["status"] != "success" {
		return fmt.Errorf("update auto scaling failed: %s", result["message"])
	}
	return nil
}

func (c *Client) GetAutoScaling(id string) (*AutoScalingInstance, error) {
	respBytes, err := c.Get(fmt.Sprintf("/autoscaling/%s", id))
	if err != nil {
		return nil, fmt.Errorf("failed to get autoscaling group: %w", err)
	}
	// Use map to avoid unmarshal failures from nested objects
	var raw map[string]interface{}
	if err := json.Unmarshal(respBytes, &raw); err != nil {
		return nil, fmt.Errorf("failed to parse get autoscaling response: %w", err)
	}
	// Try groups key first, then autoscaling
	var groups []interface{}
	if g, ok := raw["groups"].([]interface{}); ok && len(g) > 0 {
		groups = g
	} else if a, ok := raw["autoscaling"].([]interface{}); ok && len(a) > 0 {
		groups = a
	}
	if len(groups) == 0 {
		return nil, nil
	}
	// Marshal first group back to JSON then unmarshal into struct
	groupBytes, err := json.Marshal(groups[0])
	if err != nil {
		return nil, fmt.Errorf("failed to re-marshal autoscaling group: %w", err)
	}
	var asg AutoScalingInstance
	if err := json.Unmarshal(groupBytes, &asg); err != nil {
		return nil, fmt.Errorf("failed to parse autoscaling group: %w", err)
	}
	return &asg, nil
}

func (c *Client) DeleteAutoScaling(id string, name string) error {
	respBytes, err := c.Delete(fmt.Sprintf("/autoscaling/%s/?name=%s", id, name))
	if err != nil {
		return fmt.Errorf("failed to delete autoscaling group: %w", err)
	}
	var result map[string]string
	if err := json.Unmarshal(respBytes, &result); err != nil {
		return fmt.Errorf("failed to parse delete autoscaling response: %w", err)
	}
	if result["status"] != "success" {
		return fmt.Errorf("delete autoscaling failed: %s", result["message"])
	}
	return nil
}

// ── Security Group methods ────────────────────────────────────────────────

func (c *Client) AttachASGSecurityGroup(asgID string, sgID string) error {
	respBytes, err := c.Post(fmt.Sprintf("/autoscaling/%s/securitygroup/%s", asgID, sgID), nil)
	if err != nil {
		return fmt.Errorf("failed to attach security group: %w", err)
	}
	var result map[string]string
	if err := json.Unmarshal(respBytes, &result); err != nil {
		return fmt.Errorf("failed to parse attach security group response: %w", err)
	}
	if result["status"] != "success" {
		return fmt.Errorf("attach security group failed: %s", result["message"])
	}
	return nil
}

func (c *Client) DetachASGSecurityGroup(asgID string, sgID string) error {
	respBytes, err := c.Delete(fmt.Sprintf("/autoscaling/%s/securitygroup/%s", asgID, sgID))
	if err != nil {
		return fmt.Errorf("failed to detach security group: %w", err)
	}
	var result map[string]string
	if err := json.Unmarshal(respBytes, &result); err != nil {
		return fmt.Errorf("failed to parse detach security group response: %w", err)
	}
	if result["status"] != "success" {
		return fmt.Errorf("detach security group failed: %s", result["message"])
	}
	return nil
}

// ── Load Balancer methods ─────────────────────────────────────────────────

func (c *Client) AttachASGLoadBalancer(asgID string, lbID string) error {
	respBytes, err := c.Post(fmt.Sprintf("/autoscaling/%s/loadbalancer/%s", asgID, lbID), nil)
	if err != nil {
		return fmt.Errorf("failed to attach load balancer: %w", err)
	}
	var result map[string]string
	if err := json.Unmarshal(respBytes, &result); err != nil {
		return fmt.Errorf("failed to parse attach load balancer response: %w", err)
	}
	if result["status"] != "success" {
		return fmt.Errorf("attach load balancer failed: %s", result["message"])
	}
	return nil
}

func (c *Client) DetachASGLoadBalancer(asgID string, lbID string) error {
	respBytes, err := c.Delete(fmt.Sprintf("/autoscaling/%s/loadbalancer/%s", asgID, lbID))
	if err != nil {
		return fmt.Errorf("failed to detach load balancer: %w", err)
	}
	var result map[string]string
	if err := json.Unmarshal(respBytes, &result); err != nil {
		return fmt.Errorf("failed to parse detach load balancer response: %w", err)
	}
	if result["status"] != "success" {
		return fmt.Errorf("detach load balancer failed: %s", result["message"])
	}
	return nil
}

// ── Target Group methods ──────────────────────────────────────────────────

func (c *Client) AttachASGTargetGroup(asgID string, tgID string) error {
	respBytes, err := c.Post(fmt.Sprintf("/autoscaling/%s/targetgroup/%s", asgID, tgID), nil)
	if err != nil {
		return fmt.Errorf("failed to attach target group: %w", err)
	}
	var result map[string]string
	if err := json.Unmarshal(respBytes, &result); err != nil {
		return fmt.Errorf("failed to parse attach target group response: %w", err)
	}
	if result["status"] != "success" {
		return fmt.Errorf("attach target group failed: %s", result["message"])
	}
	return nil
}

func (c *Client) DetachASGTargetGroup(asgID string, tgID string) error {
	respBytes, err := c.Delete(fmt.Sprintf("/autoscaling/%s/targetgroup/%s", asgID, tgID))
	if err != nil {
		return fmt.Errorf("failed to detach target group: %w", err)
	}
	var result map[string]string
	if err := json.Unmarshal(respBytes, &result); err != nil {
		return fmt.Errorf("failed to parse detach target group response: %w", err)
	}
	if result["status"] != "success" {
		return fmt.Errorf("detach target group failed: %s", result["message"])
	}
	return nil
}

// ── Scaling Policy methods ────────────────────────────────────────────────

func (c *Client) CreateScalingPolicy(req *ScalingPolicyCreateRequest) (string, error) {
	respBytes, err := c.Post("/autoscaling/policy", req)
	if err != nil {
		return "", fmt.Errorf("failed to create scaling policy: %w", err)
	}
	var result map[string]interface{}
	if err := json.Unmarshal(respBytes, &result); err != nil {
		return "", fmt.Errorf("failed to parse create scaling policy response: %w", err)
	}
	if result["status"] != "success" {
		return "", fmt.Errorf("create scaling policy failed: %s", result["message"])
	}
	// ID comes back as float64 from JSON — format as integer
	switch v := result["id"].(type) {
	case float64:
		return fmt.Sprintf("%.0f", v), nil
	default:
		return fmt.Sprintf("%v", v), nil
	}
}

func (c *Client) UpdateScalingPolicy(policyID string, req *ScalingPolicyUpdateRequest) error {
	respBytes, err := c.Put(fmt.Sprintf("/autoscaling/policy/%s", policyID), req)
	if err != nil {
		return fmt.Errorf("failed to update scaling policy: %w", err)
	}
	var result map[string]string
	if err := json.Unmarshal(respBytes, &result); err != nil {
		return fmt.Errorf("failed to parse update scaling policy response: %w", err)
	}
	if result["status"] != "success" {
		return fmt.Errorf("update scaling policy failed: %s", result["message"])
	}
	return nil
}

func (c *Client) DeleteScalingPolicy(policyID string) error {
	respBytes, err := c.Delete(fmt.Sprintf("/autoscaling/policy/%s", policyID))
	if err != nil {
		if strings.Contains(err.Error(), "404") || strings.Contains(err.Error(), "not found") {
			return nil
		}
		return fmt.Errorf("failed to delete scaling policy: %w", err)
	}
	var result map[string]string
	if err := json.Unmarshal(respBytes, &result); err != nil {
		return nil
	}
	if result["status"] != "success" {
		return fmt.Errorf("delete scaling policy failed: %s", result["message"])
	}
	return nil
}

// ── Scaling Schedule methods ──────────────────────────────────────────────

func (c *Client) UpdateScalingSchedule(asgID string, scheduleID string, req *ScalingScheduleUpdateRequest) error {
	respBytes, err := c.Put(fmt.Sprintf("/autoscaling/%s/schedulepolicy/%s", asgID, scheduleID), req)
	if err != nil {
		return fmt.Errorf("failed to update scaling schedule: %w", err)
	}
	var result map[string]string
	if err := json.Unmarshal(respBytes, &result); err != nil {
		return fmt.Errorf("failed to parse update scaling schedule response: %w", err)
	}
	if result["status"] != "success" {
		return fmt.Errorf("update scaling schedule failed: %s", result["message"])
	}
	return nil
}

func (c *Client) DeleteScalingSchedule(asgID string, scheduleID string) error {
	respBytes, err := c.Delete(fmt.Sprintf("/autoscaling/%s/schedulepolicy/%s", asgID, scheduleID))
	if err != nil {
		return fmt.Errorf("failed to delete scaling schedule: %w", err)
	}
	var result map[string]string
	if err := json.Unmarshal(respBytes, &result); err != nil {
		return fmt.Errorf("failed to parse delete scaling schedule response: %w", err)
	}
	if result["status"] != "success" {
		return fmt.Errorf("delete scaling schedule failed: %s", result["message"])
	}
	return nil
}
