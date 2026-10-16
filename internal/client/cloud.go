package client

import (
	"encoding/json"
	"fmt"
	"strings"
)

// parseAPIResponse parses a standard Utho API response, handling empty bodies gracefully.
func parseAPIResponse(respBytes []byte, result interface{}) error {
	trimmed := strings.TrimSpace(string(respBytes))
	if trimmed == "" || trimmed == "null" {
		return nil
	}
	return json.Unmarshal(respBytes, result)
}

// ── Request structs ───────────────────────────────────────────────────────
type CloudDeployRequest struct {
	DCSlug         string        `json:"dcslug"`
	PlanID         string        `json:"planid"`
	BillingCycle   string        `json:"billingcycle"`
	Auth           string        `json:"auth"`
	EnablePublicIP string        `json:"enable_publicip"`
	SubnetRequired string        `json:"subnetRequired,omitempty"`
	CPUModel       string        `json:"cpumodel,omitempty"`
	EnableBackup   string        `json:"enablebackup,omitempty"`
	Support        string        `json:"support,omitempty"`
	Firewall       string        `json:"firewall,omitempty"`
	RootPassword   string        `json:"root_password,omitempty"`
	SSHKeys        string        `json:"sshkeys,omitempty"`
	Image          string        `json:"image,omitempty"`
	VPC            string        `json:"vpc,omitempty"`
	SnapshotID     string        `json:"snapshotid,omitempty"`
	BackupID       string        `json:"backupid,omitempty"`
	ISO            string        `json:"iso,omitempty"`
	Stack          string        `json:"stack,omitempty"`
	EBS            []EBSVolume   `json:"ebs,omitempty"`
	Cloud          []CloudServer `json:"cloud"`
}

type EBSVolume struct {
	ID   string `json:"id"`
	Disk int    `json:"disk"`
	Type string `json:"type"`
}

type CloudServer struct {
	Hostname string `json:"hostname"`
}

// ── Response structs ──────────────────────────────────────────────────────

type CloudListResponse struct {
	Cloud []CloudInstance `json:"cloud"`
	Meta  CloudMeta       `json:"meta"`
}

type CloudMeta struct {
	Total       int `json:"total"`
	TotalPages  int `json:"totalpages"`
	CurrentPage int `json:"currentpage"`
}

type CloudStorage struct {
	ID   string  `json:"id"`
	Size float64 `json:"size"`
	Type string  `json:"type"`
	Bus  string  `json:"bus"`
}

type CloudInstance struct {
	CloudID      string          `json:"cloudid"`
	Hostname     string          `json:"hostname"`
	Status       string          `json:"status"`
	PowerStatus  string          `json:"powerstatus"`
	IP           string          `json:"ip"`
	CPU          string          `json:"cpu"`
	RAM          string          `json:"ram"`
	DiskSize     int             `json:"disksize"`
	DCSlug       string          `json:"dcslug"`
	BillingCycle string          `json:"billingcycle"`
	PlanSlug     string          `json:"planslug"`
	Image        CloudImage      `json:"image"`
	DCLocation   CloudDCLocation `json:"dclocation"`
	Plan         CloudPlan       `json:"plan"`
	CreatedAt    string          `json:"created_at"`
	Storages     []CloudStorage  `json:"storages"`
}

type CloudImage struct {
	Name         string `json:"name"`
	Distribution string `json:"distribution"`
	Version      string `json:"version"`
	Image        string `json:"image"`
}

type CloudDCLocation struct {
	Location string `json:"location"`
	Country  string `json:"country"`
	DC       string `json:"dc"`
}

type CloudPlan struct {
	ID   string `json:"id"`
	Slug string `json:"slug"`
	Name string `json:"name"`
	CPU  string `json:"cpu"`
	RAM  string `json:"ram"`
	Disk string `json:"disk"`
}

type CloudDeleteResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}

type CloudDeploySuccessResponse struct {
	Status  string `json:"status"`
	CloudID string `json:"cloudid"`
	Message string `json:"message"`
	IPv4    string `json:"ipv4"`
}

// ── API methods ───────────────────────────────────────────────────────────

// CreateCloud deploys a new cloud instance
func (c *Client) CreateCloud(req *CloudDeployRequest) (*CloudInstance, error) {
	respBytes, err := c.Post("/cloud/deploy", req)
	if err != nil {
		return nil, fmt.Errorf("failed to create cloud instance: %w", err)
	}

	// Parse the deploy response
	var deployResp CloudDeploySuccessResponse
	if err := parseAPIResponse(respBytes, &deployResp); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	if deployResp.Status != "success" {
		return nil, fmt.Errorf("deploy failed: %s", deployResp.Message)
	}

	if deployResp.CloudID == "" {
		return nil, fmt.Errorf("no cloudid returned in response")
	}

	// Return a CloudInstance with what we have from deploy response
	// IP and full details will be fetched by Read() after creation
	return &CloudInstance{
		CloudID: deployResp.CloudID,
		IP:      deployResp.IPv4,
		Status:  "Active",
	}, nil
}

// GetCloud fetches a specific cloud instance by ID
func (c *Client) GetCloud(cloudID string) (*CloudInstance, error) {
	endpoint := fmt.Sprintf("/cloud/%s", cloudID)
	respBytes, err := c.Get(endpoint)
	if err != nil {
		return nil, fmt.Errorf("failed to get cloud instance: %w", err)
	}

	var listResp CloudListResponse
	if err := parseAPIResponse(respBytes, &listResp); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	if len(listResp.Cloud) == 0 {
		return nil, nil // not found — deleted outside terraform
	}

	return &listResp.Cloud[0], nil
}

// DeleteCloud destroys a cloud instance
func (c *Client) DeleteCloud(cloudID string, hostname string, deleteEBS bool) error {
	ebsParam := "no"
	if deleteEBS {
		ebsParam = "yes"
	}

	endpoint := fmt.Sprintf("/cloud/%s/destroy?ebs_delete=%s", cloudID, ebsParam)
	payload := map[string]string{"confirm": hostname}
	respBytes, err := c.doRequest("DELETE", endpoint, payload)
	if err != nil {
		return fmt.Errorf("failed to delete cloud instance: %w", err)
	}

	var deleteResp CloudDeleteResponse
	if err := parseAPIResponse(respBytes, &deleteResp); err != nil {
		return nil
	}

	if deleteResp.Status != "success" {
		return fmt.Errorf("delete failed: %s", deleteResp.Message)
	}

	return nil
}

// PowerAction performs a power action on a cloud instance
func (c *Client) PowerAction(cloudID string, action string) error {
	endpoint := fmt.Sprintf("/cloud/%s/%s", cloudID, action)
	respBytes, err := c.Post(endpoint, nil)
	if err != nil {
		return fmt.Errorf("power action '%s' failed: %w", action, err)
	}

	// API sometimes returns empty body on success
	trimmed := strings.TrimSpace(string(respBytes))
	if trimmed == "" || trimmed == "null" {
		return nil
	}

	var result map[string]string
	if err := json.Unmarshal(respBytes, &result); err != nil {
		return nil
	}

	if result["status"] != "" && result["status"] != "success" {
		return fmt.Errorf("power action failed: %s", result["message"])
	}

	return nil
}

// ListClouds returns all cloud instances in the account
func (c *Client) ListClouds() ([]CloudInstance, error) {
	var allInstances []CloudInstance
	page := 1

	for {
		endpoint := fmt.Sprintf("/cloud?perpage=100&page=%d", page)
		respBytes, err := c.Get(endpoint)
		if err != nil {
			return nil, fmt.Errorf("failed to list cloud instances: %w", err)
		}

		var listResp CloudListResponse
		if err := parseAPIResponse(respBytes, &listResp); err != nil {
			return nil, fmt.Errorf("failed to parse response: %w", err)
		}

		allInstances = append(allInstances, listResp.Cloud...)

		if page >= listResp.Meta.TotalPages || listResp.Meta.TotalPages == 0 {
			break
		}
		page++
	}

	return allInstances, nil
}

// AttachFirewall attaches a security group to a cloud instance
func (c *Client) AttachFirewall(firewallID string, cloudID string) error {
	endpoint := fmt.Sprintf("/firewall/%s/server/add", firewallID)
	payload := map[string]string{
		"cloudid": cloudID,
	}

	respBytes, err := c.Post(endpoint, payload)
	if err != nil {
		return fmt.Errorf("failed to attach firewall: %w", err)
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
		return fmt.Errorf("attach firewall failed: %s", result["message"])
	}
	return nil
}

// DetachFirewall detaches a security group from a cloud instance
func (c *Client) DetachFirewall(firewallID string, cloudID string) error {
	endpoint := fmt.Sprintf("/firewall/%s/server/%s/delete", firewallID, cloudID)
	respBytes, err := c.Delete(endpoint)
	if err != nil {
		return fmt.Errorf("failed to detach firewall: %w", err)
	}

	trimmed2 := strings.TrimSpace(string(respBytes))
	if trimmed2 == "" || trimmed2 == "null" {
		return nil
	}
	var result2 map[string]string
	if err := json.Unmarshal(respBytes, &result2); err != nil {
		return nil
	}
	if result2["status"] != "" && result2["status"] != "success" {
		return fmt.Errorf("detach firewall failed: %s", result2["message"])
	}
	return nil
}

// ── General Storage ───────────────────────────────────────────────────────

// AddGeneralStorage adds a general disk to a cloud instance
func (c *Client) AddGeneralStorage(cloudID string, sizeGB int) error {
	endpoint := fmt.Sprintf("/cloud/%s/storage/add", cloudID)
	payload := map[string]int{
		"size": sizeGB,
	}

	respBytes, err := c.Post(endpoint, payload)
	if err != nil {
		return fmt.Errorf("failed to add storage: %w", err)
	}

	var result map[string]string
	if err := parseAPIResponse(respBytes, &result); err != nil {
		return nil
	}
	if result["status"] != "" && result["status"] != "success" {
		return fmt.Errorf("add storage failed: %s", result["message"])
	}

	return nil
}

// DeleteGeneralStorage deletes a general disk from a cloud instance
func (c *Client) DeleteGeneralStorage(cloudID string, diskID string) error {
	endpoint := fmt.Sprintf("/cloud/%s/storage/%s/delete", cloudID, diskID)
	respBytes, err := c.Delete(endpoint)
	if err != nil {
		return fmt.Errorf("failed to delete storage: %w", err)
	}

	var result map[string]string
	if err := parseAPIResponse(respBytes, &result); err != nil {
		return nil
	}
	if result["status"] != "" && result["status"] != "success" {
		return fmt.Errorf("delete storage failed: %s", result["message"])
	}

	return nil
}

// UpdateGeneralStorage updates a general disk on a cloud instance
func (c *Client) UpdateGeneralStorage(cloudID string, diskID string, bus string, diskType string) error {
	endpoint := fmt.Sprintf("/cloud/%s/storage/%s/update", cloudID, diskID)
	payload := map[string]string{
		"bus":  bus,
		"type": diskType,
	}

	respBytes, err := c.Post(endpoint, payload)
	if err != nil {
		return fmt.Errorf("failed to update storage: %w", err)
	}

	var result map[string]string
	if err := parseAPIResponse(respBytes, &result); err != nil {
		return nil
	}
	if result["status"] != "" && result["status"] != "success" {
		return fmt.Errorf("update storage failed: %s", result["message"])
	}

	return nil
}

// ── EBS Storage ───────────────────────────────────────────────────────────

// AttachEBS attaches an EBS volume to a cloud instance
func (c *Client) AttachEBS(ebsID string, cloudID string) error {
	endpoint := fmt.Sprintf("/ebs/%s/attach", ebsID)
	payload := map[string]string{
		"type":       "cloud",
		"resourceid": cloudID,
	}

	respBytes, err := c.Put(endpoint, payload)
	if err != nil {
		return fmt.Errorf("failed to attach EBS: %w", err)
	}

	var result map[string]string
	if err := parseAPIResponse(respBytes, &result); err != nil {
		return nil
	}
	if result["status"] != "" && result["status"] != "success" {
		return fmt.Errorf("attach EBS failed: %s", result["message"])
	}

	return nil
}

// DetachEBS detaches an EBS volume from a cloud instance
func (c *Client) DetachEBS(ebsID string, cloudID string) error {
	endpoint := fmt.Sprintf("/ebs/%s/dettach", ebsID)
	payload := map[string]string{
		"type":       "cloud",
		"resourceid": cloudID,
	}

	respBytes, err := c.Put(endpoint, payload)
	if err != nil {
		return fmt.Errorf("failed to detach EBS: %w", err)
	}

	var result map[string]string
	if err := parseAPIResponse(respBytes, &result); err != nil {
		return nil
	}
	if result["status"] != "" && result["status"] != "success" {
		return fmt.Errorf("detach EBS failed: %s", result["message"])
	}

	return nil
}

// UpdateEBS updates an attached EBS volume on a cloud instance
func (c *Client) UpdateEBS(cloudID string, ebsID string, bus string, diskType string) error {
	endpoint := fmt.Sprintf("/cloud/%s/storage/%s/update", cloudID, ebsID)
	payload := map[string]string{
		"bus":  bus,
		"type": diskType,
	}

	respBytes, err := c.Post(endpoint, payload)
	if err != nil {
		return fmt.Errorf("failed to update EBS: %w", err)
	}

	var result map[string]string
	if err := parseAPIResponse(respBytes, &result); err != nil {
		return nil
	}
	if result["status"] != "" && result["status"] != "success" {
		return fmt.Errorf("update EBS failed: %s", result["message"])
	}

	return nil
}

// ── Snapshot ──────────────────────────────────────────────────────────────

// CreateSnapshot creates a snapshot of a cloud instance
func (c *Client) CreateSnapshot(cloudID string, name string) error {
	endpoint := fmt.Sprintf("/cloud/%s/snapshot/create", cloudID)
	payload := map[string]string{
		"name": name,
	}

	respBytes, err := c.Post(endpoint, payload)
	if err != nil {
		return fmt.Errorf("failed to create snapshot: %w", err)
	}

	var result map[string]string
	if err := parseAPIResponse(respBytes, &result); err != nil {
		return nil
	}
	if result["status"] != "" && result["status"] != "success" {
		return fmt.Errorf("create snapshot failed: %s", result["message"])
	}

	return nil
}

// DeleteSnapshot deletes a snapshot of a cloud instance
func (c *Client) DeleteSnapshot(cloudID string, snapshotID string) error {
	endpoint := fmt.Sprintf("/cloud/%s/snapshot/%s/delete", cloudID, snapshotID)
	respBytes, err := c.Delete(endpoint)
	if err != nil {
		return fmt.Errorf("failed to delete snapshot: %w", err)
	}

	var result map[string]string
	if err := parseAPIResponse(respBytes, &result); err != nil {
		return nil
	}
	if result["status"] != "" && result["status"] != "success" {
		return fmt.Errorf("delete snapshot failed: %s", result["message"])
	}

	return nil
}

// RestoreSnapshot restores a cloud instance from a snapshot
func (c *Client) RestoreSnapshot(cloudID string, snapshotID string) error {
	endpoint := fmt.Sprintf("/cloud/%s/snapshot/%s/restore", cloudID, snapshotID)
	respBytes, err := c.Post(endpoint, nil)
	if err != nil {
		return fmt.Errorf("failed to restore snapshot: %w", err)
	}

	var result map[string]string
	if err := parseAPIResponse(respBytes, &result); err != nil {
		return nil
	}
	if result["status"] != "" && result["status"] != "success" {
		return fmt.Errorf("restore snapshot failed: %s", result["message"])
	}

	return nil
}

// ── ISO ───────────────────────────────────────────────────────────────────

// MountISO mounts an ISO on a cloud instance and boots from it
func (c *Client) MountISO(cloudID string, isoName string) error {
	endpoint := fmt.Sprintf("/cloud/%s/mountiso", cloudID)
	payload := map[string]string{
		"iso": isoName,
	}

	respBytes, err := c.Post(endpoint, payload)
	if err != nil {
		return fmt.Errorf("failed to mount ISO: %w", err)
	}

	var result map[string]string
	if err := parseAPIResponse(respBytes, &result); err != nil {
		return nil
	}
	if result["status"] != "" && result["status"] != "success" {
		return fmt.Errorf("mount ISO failed: %s", result["message"])
	}

	return nil
}

// ── Resize ────────────────────────────────────────────────────────────────

// ResizeCloud resizes a cloud instance
// resizeType: "ramcpu" (CPU/RAM only) or "full" (CPU/RAM + disk)
func (c *Client) ResizeCloud(cloudID string, planID string, resizeType string) error {
	endpoint := fmt.Sprintf("/cloud/%s/resize", cloudID)
	payload := map[string]string{
		"plan": planID,
		"type": resizeType,
	}

	respBytes, err := c.Post(endpoint, payload)
	if err != nil {
		return fmt.Errorf("failed to resize cloud instance: %w", err)
	}

	var result map[string]string
	if err := parseAPIResponse(respBytes, &result); err != nil {
		return nil
	}
	if result["status"] != "" && result["status"] != "success" {
		return fmt.Errorf("resize failed: %s", result["message"])
	}

	return nil
}

// ── Network ───────────────────────────────────────────────────────────────

// AssignPublicIP assigns an additional public IP to a cloud instance
func (c *Client) AssignPublicIP(cloudID string) (string, error) {
	endpoint := fmt.Sprintf("/cloud/%s/assignpublicip", cloudID)
	respBytes, err := c.Post(endpoint, nil)
	if err != nil {
		return "", fmt.Errorf("failed to assign public IP: %w", err)
	}

	var result map[string]interface{}
	if err := parseAPIResponse(respBytes, &result); err != nil {
		return "", fmt.Errorf("failed to parse response: %w", err)
	}

	if result["status"] != "success" {
		return "", fmt.Errorf("assign public IP failed: %s", result["message"])
	}

	ip, _ := result["ip"].(string)
	return ip, nil
}

// DeletePublicIP releases a public IP from a cloud instance
func (c *Client) DeletePublicIP(cloudID string, ip string) error {
	endpoint := fmt.Sprintf("/cloud/%s/ip/%s/delete", cloudID, ip)
	respBytes, err := c.Delete(endpoint)
	if err != nil {
		return fmt.Errorf("failed to delete public IP: %w", err)
	}

	var result map[string]string
	if err := parseAPIResponse(respBytes, &result); err != nil {
		return nil
	}
	if result["status"] != "" && result["status"] != "success" {
		return fmt.Errorf("delete public IP failed: %s", result["message"])
	}

	return nil
}

// AttachVPC attaches a VPC subnet to a cloud instance
func (c *Client) AttachVPC(cloudID string, subnetID string) (string, error) {
	endpoint := fmt.Sprintf("/cloud/%s/vpc/%s/attach", cloudID, subnetID)
	respBytes, err := c.Post(endpoint, nil)
	if err != nil {
		return "", fmt.Errorf("failed to attach VPC: %w", err)
	}

	var result map[string]interface{}
	if err := parseAPIResponse(respBytes, &result); err != nil {
		return "", fmt.Errorf("failed to parse response: %w", err)
	}

	if result["status"] != "success" {
		return "", fmt.Errorf("attach VPC failed: %s", result["message"])
	}

	ip, _ := result["ip"].(string)
	return ip, nil
}

// DetachVPC detaches a VPC subnet from a cloud instance
func (c *Client) DetachVPC(cloudID string, subnetID string) error {
	endpoint := fmt.Sprintf("/cloud/%s/vpc/%s/detach", cloudID, subnetID)
	respBytes, err := c.Post(endpoint, nil)
	if err != nil {
		return fmt.Errorf("failed to detach VPC: %w", err)
	}

	var result map[string]string
	if err := parseAPIResponse(respBytes, &result); err != nil {
		return nil
	}
	if result["status"] != "" && result["status"] != "success" {
		return fmt.Errorf("detach VPC failed: %s", result["message"])
	}

	return nil
}
