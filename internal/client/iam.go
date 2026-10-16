package client

import (
	"encoding/json"
	"fmt"
)

type IAMUserCreateRequest struct {
	FullName    string `json:"fullname"`
	Email       string `json:"email"`
	MobileCC    string `json:"mobilecc"`
	Mobile      string `json:"mobile"`
	Permissions string `json:"permissions"`
}

type IAMUserUpdateRequest struct {
	Permissions string `json:"permissions"`
}

type IAMUserInstance struct {
	ID          string `json:"id"`
	FullName    string `json:"fullname"`
	Email       string `json:"email"`
	Permissions string `json:"permissions"`
	Status      string `json:"status"`
	Resources   string `json:"resources"`
	DateAdded   string `json:"dateadded"`
}

func (c *Client) CreateIAMUser(req *IAMUserCreateRequest) (string, error) {
	respBytes, err := c.Post("/user", req)
	if err != nil {
		return "", fmt.Errorf("failed to create IAM user: %w", err)
	}
	var result map[string]interface{}
	if err := json.Unmarshal(respBytes, &result); err != nil {
		return "", fmt.Errorf("failed to parse create IAM user response: %w", err)
	}
	if result["status"] != "success" {
		return "", fmt.Errorf("create IAM user failed: %s", result["message"])
	}
	// API returns subuser ID — find account access ID from list
	var subuserID string
	switch v := result["id"].(type) {
	case float64:
		subuserID = fmt.Sprintf("%.0f", v)
	default:
		subuserID = fmt.Sprintf("%v", v)
	}
	// Find account access ID by matching subuser ID
	users, err := c.ListIAMUsers()
	if err == nil {
		for _, u := range users {
			if u.Email == req.Email {
				return u.ID, nil
			}
		}
	}
	return subuserID, nil
}

func (c *Client) GetIAMUser(userID string) (*IAMUserInstance, error) {
	respBytes, err := c.Get(fmt.Sprintf("/user/%s", userID))
	if err != nil {
		return nil, fmt.Errorf("failed to get IAM user: %w", err)
	}
	var result struct {
		Status string          `json:"status"`
		Info   IAMUserInstance `json:"info"`
	}
	if err := json.Unmarshal(respBytes, &result); err != nil {
		return nil, fmt.Errorf("failed to parse get IAM user response: %w", err)
	}
	if result.Info.ID == "" {
		return nil, nil
	}
	return &result.Info, nil
}

func (c *Client) ListIAMUsers() ([]IAMUserInstance, error) {
	respBytes, err := c.Get("/user")
	if err != nil {
		return nil, fmt.Errorf("failed to list IAM users: %w", err)
	}
	var result struct {
		Status        string `json:"status"`
		AccountAccess []struct {
			ID          string `json:"id"`
			FullName    string `json:"fullname"`
			Email       string `json:"email"`
			Permissions string `json:"permissions"`
			Status      string `json:"status"`
			Resources   string `json:"resources"`
			DateAdded   string `json:"dateadded"`
		} `json:"accountaccess"`
	}
	if err := json.Unmarshal(respBytes, &result); err != nil {
		return nil, fmt.Errorf("failed to parse list IAM users response: %w", err)
	}
	var users []IAMUserInstance
	for _, u := range result.AccountAccess {
		users = append(users, IAMUserInstance{
			ID:          u.ID,
			FullName:    u.FullName,
			Email:       u.Email,
			Permissions: u.Permissions,
			Status:      u.Status,
			Resources:   u.Resources,
			DateAdded:   u.DateAdded,
		})
	}
	return users, nil
}

func (c *Client) UpdateIAMUser(userID string, req *IAMUserUpdateRequest) error {
	respBytes, err := c.Put(fmt.Sprintf("/user/%s", userID), req)
	if err != nil {
		return fmt.Errorf("failed to update IAM user: %w", err)
	}
	var result map[string]string
	if err := json.Unmarshal(respBytes, &result); err != nil {
		return fmt.Errorf("failed to parse update IAM user response: %w", err)
	}
	if result["status"] != "success" {
		return fmt.Errorf("update IAM user failed: %s", result["message"])
	}
	return nil
}

func (c *Client) DeleteIAMUser(userID string) error {
	respBytes, err := c.Delete(fmt.Sprintf("/user/%s", userID))
	if err != nil {
		return fmt.Errorf("failed to delete IAM user: %w", err)
	}
	var result map[string]string
	if err := json.Unmarshal(respBytes, &result); err != nil {
		return fmt.Errorf("failed to parse delete IAM user response: %w", err)
	}
	if result["status"] != "success" {
		return fmt.Errorf("delete IAM user failed: %s", result["message"])
	}
	return nil
}
