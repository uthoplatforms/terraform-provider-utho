package client

import (
	"encoding/json"
	"fmt"
)

// ── Response structs ──────────────────────────────────────────────────────

type CloudDCZone struct {
	ID                 string   `json:"id"`
	Slug               string   `json:"slug"`
	CC                 string   `json:"cc"`
	Country            string   `json:"country"`
	City               string   `json:"city"`
	Status             string   `json:"status"`
	DefaultCPU         string   `json:"default_cpu"`
	AvailableCPUModels []string `json:"available_cpumodels"`
	EBSAvailable       bool     `json:"ebs_available"`
	PublicIPAvailable  bool     `json:"public_ip_available"`
}

type DeployPlan struct {
	ID             string  `json:"id"`
	Type           string  `json:"type"`
	Slug           string  `json:"slug"`
	Name           string  `json:"name"`
	Disk           string  `json:"disk"`
	RAM            string  `json:"ram"`
	CPU            string  `json:"cpu"`
	Bandwidth      string  `json:"bandwidth"`
	DedicatedVCore string  `json:"dedicated_vcore"`
	Price          float64 `json:"price"`
	IsAvailable    string  `json:"is_available"`
}

type DeployImage struct {
	ID           string `json:"id"`
	Distro       string `json:"distro"`
	Image        string `json:"image"`
	Version      string `json:"version"`
	Distribution string `json:"distribution"`
	Category     string `json:"category"`
	Type         string `json:"type"`
}

type DeployDistro struct {
	Distro       string        `json:"distro"`
	Distribution string        `json:"distribution"`
	Images       []DeployImage `json:"images"`
}

type DeploySnapshot struct {
	ID         string `json:"id"`
	CloudID    string `json:"cloudid"`
	Image      string `json:"image"`
	Size       string `json:"size"`
	Name       string `json:"name"`
	Status     string `json:"status"`
	Storage    string `json:"storage"`
	CreateDate string `json:"createdate"`
}

type DeployISO struct {
	Name    string `json:"name"`
	File    string `json:"file"`
	DC      string `json:"dc"`
	AddedAt string `json:"added_at"`
}

type DeploySSHKey struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	Value     string `json:"value"`
	CreatedAt string `json:"created_at"`
}

type DeployFirewall struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	CreatedAt string `json:"created_at"`
}

type CloudDeployData struct {
	DCZones   []CloudDCZone    `json:"dczones"`
	Plans     []DeployPlan     `json:"plans"`
	Distros   []DeployDistro   `json:"distro"`
	Snapshots []DeploySnapshot `json:"snapshots"`
	ISOs      []DeployISO      `json:"isos"`
	Keys      []DeploySSHKey   `json:"keys"`
	Firewalls []DeployFirewall `json:"firewalls"`
}

type VPCSubnetItem struct {
	ID             string `json:"id"`
	UUID           string `json:"uuid"`
	DCSlug         string `json:"dcslug"`
	Name           string `json:"name"`
	Network        string `json:"network"`
	Size           string `json:"size"`
	SubnetType     string `json:"subnet_type"`
	AssignPublicIP string `json:"assign_publicip"`
	Status         string `json:"status"`
	Gateway        string `json:"gateway"`
}

type VPCItemFull struct {
	ID        string          `json:"id"`
	Name      string          `json:"name"`
	Network   string          `json:"network"`
	Size      string          `json:"size"`
	DCSlug    string          `json:"dcslug"`
	Total     int             `json:"total"`
	Available int             `json:"available"`
	Subnets   []VPCSubnetItem `json:"subnets"`
}

type VPCListFullResponse struct {
	Status string        `json:"status"`
	VPC    []VPCItemFull `json:"vpc"`
}

// ── API methods ───────────────────────────────────────────────────────────

func (c *Client) GetCloudDeployData(dcslug string) (*CloudDeployData, error) {
	endpoint := "/cloud/getdeploy?product=cloud"
	if dcslug != "" {
		endpoint = fmt.Sprintf("/cloud/getdeploy?product=cloud&dcslug=%s", dcslug)
	}

	respBytes, err := c.Get(endpoint)
	if err != nil {
		return nil, fmt.Errorf("failed to get cloud deploy data: %w", err)
	}

	var data CloudDeployData
	if err := json.Unmarshal(respBytes, &data); err != nil {
		return nil, fmt.Errorf("failed to parse cloud deploy data: %w", err)
	}

	return &data, nil
}

func (c *Client) ListVPCsFull() ([]VPCItemFull, error) {
	respBytes, err := c.Get("/vpc")
	if err != nil {
		return nil, fmt.Errorf("failed to list VPCs: %w", err)
	}

	var resp VPCListFullResponse
	if err := json.Unmarshal(respBytes, &resp); err != nil {
		return nil, fmt.Errorf("failed to parse VPC list: %w", err)
	}

	return resp.VPC, nil
}
