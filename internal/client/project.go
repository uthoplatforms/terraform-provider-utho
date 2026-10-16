package client

import (
	"encoding/json"
	"fmt"
)

// ── Request structs ───────────────────────────────────────────────────────

type ProjectCreateRequest struct {
	Name        string `json:"name"`
	Description string `json:"description,omitempty"`
	Environment string `json:"environment"`
}

type ProjectMemberRequest struct {
	UserID int `json:"userid"`
}

// ── Response structs ──────────────────────────────────────────────────────

type ProjectInstance struct {
	ID            int    `json:"id"`
	UUID          string `json:"uuid"`
	Name          string `json:"name"`
	Description   string `json:"description"`
	Environment   string `json:"environment"`
	Status        string `json:"status"`
	IsDefault     bool   `json:"is_default"`
	MemberCount   int    `json:"member_count"`
	ResourceCount int    `json:"resource_count"`
	CreatedAt     string `json:"created_at"`
	UpdatedAt     string `json:"updated_at"`
}

type ProjectDetailResponse struct {
	Status  string          `json:"status"`
	Project ProjectInstance `json:"project"`
}

type ProjectListResponse struct {
	Status   string            `json:"status"`
	Projects []ProjectInstance `json:"projects"`
}

// ── API methods ───────────────────────────────────────────────────────────

func (c *Client) CreateProject(req *ProjectCreateRequest) (*ProjectInstance, error) {
	respBytes, err := c.Post("/projects/create", req)
	if err != nil {
		return nil, fmt.Errorf("failed to create project: %w", err)
	}
	var resp ProjectDetailResponse
	if err := json.Unmarshal(respBytes, &resp); err != nil {
		return nil, fmt.Errorf("failed to parse create project response: %w", err)
	}
	if resp.Status != "success" {
		return nil, fmt.Errorf("create project failed")
	}
	return &resp.Project, nil
}

func (c *Client) GetProject(projectID string) (*ProjectInstance, error) {
	respBytes, err := c.Get(fmt.Sprintf("/projects/%s", projectID))
	if err != nil {
		return nil, fmt.Errorf("failed to get project: %w", err)
	}
	var resp ProjectDetailResponse
	if err := json.Unmarshal(respBytes, &resp); err != nil {
		return nil, fmt.Errorf("failed to parse get project response: %w", err)
	}
	if resp.Project.ID == 0 {
		return nil, nil
	}
	return &resp.Project, nil
}

func (c *Client) UpdateProject(projectID string, req *ProjectCreateRequest) (*ProjectInstance, error) {
	respBytes, err := c.Post(fmt.Sprintf("/projects/%s/update", projectID), req)
	if err != nil {
		return nil, fmt.Errorf("failed to update project: %w", err)
	}
	var resp ProjectDetailResponse
	if err := json.Unmarshal(respBytes, &resp); err != nil {
		return nil, fmt.Errorf("failed to parse update project response: %w", err)
	}
	if resp.Status != "success" {
		return nil, fmt.Errorf("update project failed")
	}
	return &resp.Project, nil
}

func (c *Client) DeleteProject(projectID string) error {
	respBytes, err := c.Delete(fmt.Sprintf("/projects/%s/delete", projectID))
	if err != nil {
		return fmt.Errorf("failed to delete project: %w", err)
	}
	var result map[string]string
	if err := json.Unmarshal(respBytes, &result); err != nil {
		return fmt.Errorf("failed to parse delete project response: %w", err)
	}
	if result["status"] != "success" {
		return fmt.Errorf("delete project failed: %s", result["message"])
	}
	return nil
}

func (c *Client) AddProjectMember(projectID string, userID int) error {
	respBytes, err := c.Post(fmt.Sprintf("/projects/%s/members/add", projectID),
		&ProjectMemberRequest{UserID: userID})
	if err != nil {
		return fmt.Errorf("failed to add project member: %w", err)
	}
	var result map[string]string
	if err := json.Unmarshal(respBytes, &result); err != nil {
		return fmt.Errorf("failed to parse add project member response: %w", err)
	}
	if result["status"] != "success" {
		return fmt.Errorf("add project member failed: %s", result["message"])
	}
	return nil
}

func (c *Client) RemoveProjectMember(projectID string, userID string) error {
	respBytes, err := c.Delete(fmt.Sprintf("/projects/%s/members/%s/remove", projectID, userID))
	if err != nil {
		return fmt.Errorf("failed to remove project member: %w", err)
	}
	var result map[string]string
	if err := json.Unmarshal(respBytes, &result); err != nil {
		return fmt.Errorf("failed to parse remove project member response: %w", err)
	}
	if result["status"] != "success" {
		return fmt.Errorf("remove project member failed: %s", result["message"])
	}
	return nil
}
