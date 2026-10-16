package client

import (
	"encoding/json"
	"fmt"
)

// ── Request structs ───────────────────────────────────────────────────────

type ContainerRegistryCreateRequest struct {
	DCSlug       string `json:"dcslug"`
	PlanID       int    `json:"planid"`
	BillingCycle string `json:"billingcycle"`
	Public       string `json:"public"`
	ProjectName  string `json:"project_name"`
}

type ContainerRegistryUpdateRequest struct {
	Public string `json:"public"`
}

type RegistryRobotAccess struct {
	Resource string `json:"resource"`
	Action   string `json:"action"`
}

type RegistryRobotCreateRequest struct {
	Name        string                `json:"name"`
	Description string                `json:"description"`
	Duration    int                   `json:"duration"`
	Access      []RegistryRobotAccess `json:"access"`
}

type RegistryRobotDeleteRequest struct {
	RobotID int `json:"robot_id"`
}

type RegistryWebhookCreateRequest struct {
	Name           string   `json:"name"`
	Address        string   `json:"address"`
	EventTypes     []string `json:"event_types"`
	AuthHeader     string   `json:"auth_header,omitempty"`
	SkipCertVerify bool     `json:"skip_cert_verify"`
	PayloadFormat  string   `json:"payload_format"`
}

type RegistryWebhookDeleteRequest struct {
	WebhookID int `json:"webhook_id"`
}

type RegistryImmutableCreateRequest struct {
	TagPattern  string `json:"tag_pattern"`
	RepoPattern string `json:"repo_pattern"`
}

type RegistryImmutableDeleteRequest struct {
	RuleID int `json:"rule_id"`
}

type RegistryBuildRequest struct {
	SourceType string `json:"source_type"`
	Dockerfile string `json:"dockerfile"`
	TargetRepo string `json:"target_repo"`
	TargetTag  string `json:"target_tag"`
}

type RegistryTagAddRequest struct {
	Tag string `json:"tag"`
}

type RegistryTagPromoteRequest struct {
	SourceRepo string `json:"source_repo"`
	SourceRef  string `json:"source_ref"`
	TargetRepo string `json:"target_repo"`
	TargetTag  string `json:"target_tag"`
}

// ── Response structs ──────────────────────────────────────────────────────

type ContainerRegistry struct {
	ID           int    `json:"id"`
	HarborPID    int    `json:"harbor_pid"`
	PlanID       int    `json:"plan_id"`
	BillingCycle string `json:"billingcycle"`
	CreatedAt    string `json:"created_at"`
}

type ContainerRegistryListItem struct {
	Name         string `json:"name"`
	ProjectID    int    `json:"project_id"`
	RepoCount    int    `json:"repo_count"`
	CreationTime string `json:"creation_time"`
}

type ContainerRegistryListResponse struct {
	Status string                      `json:"status"`
	Data   []ContainerRegistryListItem `json:"data"`
}

type ContainerRegistryDetailResponse struct {
	Status string `json:"status"`
	Data   struct {
		Project ContainerRegistry `json:"project"`
	} `json:"data"`
}

type RegistryRobot struct {
	ID           int    `json:"id"`
	Name         string `json:"name"`
	Description  string `json:"description"`
	Duration     int    `json:"duration"`
	ExpiresAt    int64  `json:"expires_at"`
	Disable      bool   `json:"disable"`
	CreationTime string `json:"creation_time"`
	Secret       string `json:"secret"`
}

type RegistryRobotCreateResponse struct {
	Status string        `json:"status"`
	Data   RegistryRobot `json:"data"`
}

type RegistryRobotsListResponse struct {
	Status string          `json:"status"`
	Data   []RegistryRobot `json:"data"`
}

// ── Registry methods ──────────────────────────────────────────────────────

func (c *Client) CreateContainerRegistry(req *ContainerRegistryCreateRequest) (string, error) {
	respBytes, err := c.Post("/registry", req)
	if err != nil {
		return "", fmt.Errorf("failed to create container registry: %w", err)
	}
	var result map[string]interface{}
	if err := json.Unmarshal(respBytes, &result); err != nil {
		return "", fmt.Errorf("failed to parse create registry response: %w", err)
	}
	if result["status"] != "success" {
		return "", fmt.Errorf("create registry failed: %s", result["message"])
	}
	return req.ProjectName, nil
}

func (c *Client) GetContainerRegistry(projectName string) (*ContainerRegistry, error) {
	respBytes, err := c.Get(fmt.Sprintf("/registry/cr2/summary?project_name=%s", projectName))
	if err != nil {
		return nil, fmt.Errorf("failed to get container registry: %w", err)
	}
	var resp ContainerRegistryDetailResponse
	if err := json.Unmarshal(respBytes, &resp); err != nil {
		return nil, fmt.Errorf("failed to parse get registry response: %w", err)
	}
	if resp.Data.Project.ID == 0 {
		return nil, nil
	}
	return &resp.Data.Project, nil
}

func (c *Client) UpdateContainerRegistry(projectName string, public string) error {
	respBytes, err := c.Put(fmt.Sprintf("/registry/project/%s", projectName),
		&ContainerRegistryUpdateRequest{Public: public})
	if err != nil {
		return fmt.Errorf("failed to update container registry: %w", err)
	}
	var result map[string]string
	if err := json.Unmarshal(respBytes, &result); err != nil {
		return fmt.Errorf("failed to parse update registry response: %w", err)
	}
	if result["status"] != "success" {
		return fmt.Errorf("update registry failed: %s", result["message"])
	}
	return nil
}

func (c *Client) DeleteContainerRegistry(projectName string) error {
	respBytes, err := c.Delete(fmt.Sprintf("/registry/project/%s", projectName))
	if err != nil {
		return fmt.Errorf("failed to delete container registry: %w", err)
	}
	var result map[string]string
	if err := json.Unmarshal(respBytes, &result); err != nil {
		return fmt.Errorf("failed to parse delete registry response: %w", err)
	}
	if result["status"] != "success" {
		return fmt.Errorf("delete registry failed: %s", result["message"])
	}
	return nil
}

// ── Robot methods ─────────────────────────────────────────────────────────

func (c *Client) CreateRegistryRobot(projectName string, req *RegistryRobotCreateRequest) (*RegistryRobot, error) {
	respBytes, err := c.Post(fmt.Sprintf("/registry/cr2/robot_create?project_name=%s", projectName), req)
	if err != nil {
		return nil, fmt.Errorf("failed to create registry robot: %w", err)
	}
	var resp RegistryRobotCreateResponse
	if err := json.Unmarshal(respBytes, &resp); err != nil {
		return nil, fmt.Errorf("failed to parse create robot response: %w", err)
	}
	if resp.Status != "success" {
		return nil, fmt.Errorf("create registry robot failed")
	}
	return &resp.Data, nil
}

func (c *Client) DeleteRegistryRobot(projectName string, robotID int) error {
	respBytes, err := c.Post(fmt.Sprintf("/registry/cr2/robot_delete?project_name=%s", projectName),
		&RegistryRobotDeleteRequest{RobotID: robotID})
	if err != nil {
		return fmt.Errorf("failed to delete registry robot: %w", err)
	}
	var result map[string]string
	if err := json.Unmarshal(respBytes, &result); err != nil {
		return fmt.Errorf("failed to parse delete robot response: %w", err)
	}
	if result["status"] != "success" {
		return fmt.Errorf("delete registry robot failed")
	}
	return nil
}

// ── Webhook methods ───────────────────────────────────────────────────────

func (c *Client) CreateRegistryWebhook(projectName string, req *RegistryWebhookCreateRequest) (string, error) {
	respBytes, err := c.Post(fmt.Sprintf("/registry/cr2/webhook_create?project_name=%s", projectName), req)
	if err != nil {
		return "", fmt.Errorf("failed to create registry webhook: %w", err)
	}
	var result map[string]interface{}
	if err := json.Unmarshal(respBytes, &result); err != nil {
		return "", fmt.Errorf("failed to parse create webhook response: %w", err)
	}
	if result["status"] != "success" {
		return "", fmt.Errorf("create registry webhook failed")
	}
	return fmt.Sprintf("%v", result["id"]), nil
}

func (c *Client) DeleteRegistryWebhook(projectName string, webhookID int) error {
	respBytes, err := c.Post(fmt.Sprintf("/registry/cr2/webhook_delete?project_name=%s", projectName),
		&RegistryWebhookDeleteRequest{WebhookID: webhookID})
	if err != nil {
		return fmt.Errorf("failed to delete registry webhook: %w", err)
	}
	var result map[string]string
	if err := json.Unmarshal(respBytes, &result); err != nil {
		return fmt.Errorf("failed to parse delete webhook response: %w", err)
	}
	if result["status"] != "success" {
		return fmt.Errorf("delete registry webhook failed")
	}
	return nil
}

// ── Immutable tag methods ─────────────────────────────────────────────────

func (c *Client) CreateImmutableRule(projectName string, req *RegistryImmutableCreateRequest) error {
	respBytes, err := c.Post(fmt.Sprintf("/registry/cr2/immutable_create?project_name=%s", projectName), req)
	if err != nil {
		return fmt.Errorf("failed to create immutable rule: %w", err)
	}
	var result map[string]string
	if err := json.Unmarshal(respBytes, &result); err != nil {
		return fmt.Errorf("failed to parse create immutable rule response: %w", err)
	}
	if result["status"] != "success" {
		return fmt.Errorf("create immutable rule failed")
	}
	return nil
}

func (c *Client) DeleteImmutableRule(projectName string, ruleID int) error {
	respBytes, err := c.Post(fmt.Sprintf("/registry/cr2/immutable_delete?project_name=%s", projectName),
		&RegistryImmutableDeleteRequest{RuleID: ruleID})
	if err != nil {
		return fmt.Errorf("failed to delete immutable rule: %w", err)
	}
	var result map[string]string
	if err := json.Unmarshal(respBytes, &result); err != nil {
		return fmt.Errorf("failed to parse delete immutable rule response: %w", err)
	}
	if result["status"] != "success" {
		return fmt.Errorf("delete immutable rule failed")
	}
	return nil
}
