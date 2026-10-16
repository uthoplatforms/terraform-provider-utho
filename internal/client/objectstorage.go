package client

import (
	"encoding/json"
	"fmt"
)

// ── Request structs ───────────────────────────────────────────────────────

type BucketCreateRequest struct {
	DCSlug       string `json:"dcslug"`
	Name         string `json:"name"`
	Size         string `json:"size"`
	PlanID       string `json:"planid"`
	Cycle        string `json:"cycle"`
	Billing      string `json:"billing"`
	BillingCycle string `json:"billingcycle"`
	Price        string `json:"price"`
}

type AccessKeyCreateRequest struct {
	AccessKey string `json:"accesskey"`
}

type AccessKeyStatusRequest struct {
	AccessKey string `json:"accesskey"`
	Status    string `json:"status"`
}

// ── Response structs ──────────────────────────────────────────────────────

type BucketInstance struct {
	Name           string  `json:"name"`
	Access         string  `json:"access"`
	DCSlug         string  `json:"dcslug"`
	Size           string  `json:"size"`
	Status         string  `json:"status"`
	AccessKey      string  `json:"access_key"`
	SecretKey      string  `json:"secret_key"`
	CreatedAt      string  `json:"created_at"`
	ObjectCount    int     `json:"object_count"`
	VersionEnabled bool    `json:"version_enabled"`
	BillingCycle   string  `json:"billingcycle"`
	Cost           float64 `json:"cost"`
	PlanGB         int     `json:"plan_gb"`
	UsedGB         float64 `json:"used_gb"`
}

type AccessKeyInstance struct {
	AccessKey string `json:"accesskey"`
	SecretKey string `json:"secretkey"`
}

// ── Bucket methods ────────────────────────────────────────────────────────

func (c *Client) CreateBucket(req *BucketCreateRequest) (string, error) {
	respBytes, err := c.Post("/objectstorage/bucket/create", req)
	if err != nil {
		return "", fmt.Errorf("failed to create bucket: %w", err)
	}
	var result map[string]interface{}
	if err := json.Unmarshal(respBytes, &result); err != nil {
		return "", fmt.Errorf("failed to parse create bucket response: %w", err)
	}
	if result["status"] != "success" {
		return "", fmt.Errorf("create bucket failed: %s", result["message"])
	}
	return req.Name, nil
}

func (c *Client) GetBucket(dcslug string, name string) (*BucketInstance, error) {
	respBytes, err := c.Get(fmt.Sprintf("/objectstorage/%s/bucket", dcslug))
	if err != nil {
		return nil, fmt.Errorf("failed to get bucket: %w", err)
	}
	var result struct {
		Buckets []BucketInstance `json:"buckets"`
	}
	if err := json.Unmarshal(respBytes, &result); err != nil {
		return nil, fmt.Errorf("failed to parse get bucket response: %w", err)
	}
	for _, b := range result.Buckets {
		if b.Name == name {
			return &b, nil
		}
	}
	return nil, nil
}

func (c *Client) DeleteBucket(dcslug string, name string) error {
	respBytes, err := c.Delete(fmt.Sprintf("/objectstorage/%s/bucket/%s/delete", dcslug, name))
	if err != nil {
		return fmt.Errorf("failed to delete bucket: %w", err)
	}
	var result map[string]string
	if err := json.Unmarshal(respBytes, &result); err != nil {
		return fmt.Errorf("failed to parse delete bucket response: %w", err)
	}
	if result["status"] != "success" {
		return fmt.Errorf("delete bucket failed: %s", result["message"])
	}
	return nil
}

func (c *Client) EnableBucketVersioning(dcslug string, name string) error {
	respBytes, err := c.Post(fmt.Sprintf("/objectstorage/%s/bucket/%s/version", dcslug, name), nil)
	if err != nil {
		return fmt.Errorf("failed to enable versioning: %w", err)
	}
	var result map[string]string
	if err := json.Unmarshal(respBytes, &result); err != nil {
		return fmt.Errorf("failed to parse enable versioning response: %w", err)
	}
	if result["status"] != "success" {
		return fmt.Errorf("enable versioning failed: %s", result["message"])
	}
	return nil
}

func (c *Client) DisableBucketVersioning(dcslug string, name string) error {
	respBytes, err := c.Delete(fmt.Sprintf("/objectstorage/%s/bucket/%s/version", dcslug, name))
	if err != nil {
		return fmt.Errorf("failed to disable versioning: %w", err)
	}
	var result map[string]string
	if err := json.Unmarshal(respBytes, &result); err != nil {
		return fmt.Errorf("failed to parse disable versioning response: %w", err)
	}
	if result["status"] != "success" {
		return fmt.Errorf("disable versioning failed: %s", result["message"])
	}
	return nil
}

func (c *Client) UpdateBucketPolicy(dcslug string, name string, policy string) error {
	endpoint := fmt.Sprintf("/objectstorage/%s/bucket/%s/policy/%s", dcslug, name, policy)
	respBytes, err := c.Post(endpoint, nil)
	if err != nil {
		return fmt.Errorf("failed to update bucket policy: %w", err)
	}
	var result map[string]string
	if err := json.Unmarshal(respBytes, &result); err != nil {
		return fmt.Errorf("failed to parse update bucket policy response: %w", err)
	}
	if result["status"] != "success" {
		return fmt.Errorf("update bucket policy failed: %s", result["message"])
	}
	return nil
}

func (c *Client) GrantBucketPermission(dcslug string, name string, permission string, accessKey string) error {
	endpoint := fmt.Sprintf("/objectstorage/%s/bucket/%s/permission/%s/accesskey/%s", dcslug, name, permission, accessKey)
	respBytes, err := c.Post(endpoint, nil)
	if err != nil {
		return fmt.Errorf("failed to grant bucket permission: %w", err)
	}
	var result map[string]string
	if err := json.Unmarshal(respBytes, &result); err != nil {
		return fmt.Errorf("failed to parse grant permission response: %w", err)
	}
	if result["status"] != "success" {
		return fmt.Errorf("grant bucket permission failed: %s", result["message"])
	}
	return nil
}

// ── Access Key methods ────────────────────────────────────────────────────

func (c *Client) CreateAccessKey(dcslug string, name string) (*AccessKeyInstance, error) {
	endpoint := fmt.Sprintf("/objectstorage/%s/accesskey/create", dcslug)
	respBytes, err := c.Post(endpoint, &AccessKeyCreateRequest{AccessKey: name})
	if err != nil {
		return nil, fmt.Errorf("failed to create access key: %w", err)
	}
	var result map[string]interface{}
	if err := json.Unmarshal(respBytes, &result); err != nil {
		return nil, fmt.Errorf("failed to parse create access key response: %w", err)
	}
	if result["status"] != "success" {
		return nil, fmt.Errorf("create access key failed: %s", result["message"])
	}
	return &AccessKeyInstance{
		AccessKey: result["accesskey"].(string),
		SecretKey: result["secretkey"].(string),
	}, nil
}

func (c *Client) UpdateAccessKeyStatus(dcslug string, accessKey string, status string) error {
	endpoint := fmt.Sprintf("/objectstorage/%s/accesskey/%s/status", dcslug, accessKey)
	respBytes, err := c.Post(endpoint, &AccessKeyStatusRequest{
		AccessKey: accessKey,
		Status:    status,
	})
	if err != nil {
		return fmt.Errorf("failed to update access key status: %w", err)
	}
	var result map[string]string
	if err := json.Unmarshal(respBytes, &result); err != nil {
		return fmt.Errorf("failed to parse update access key status response: %w", err)
	}
	if result["status"] != "success" {
		return fmt.Errorf("update access key status failed: %s", result["message"])
	}
	return nil
}

func (c *Client) DeleteAccessKey(dcslug string, accessKey string) error {
	return c.UpdateAccessKeyStatus(dcslug, accessKey, "remove")
}
