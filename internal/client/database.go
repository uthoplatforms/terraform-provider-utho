package client

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"
)

// ── Request structs ───────────────────────────────────────────────────────

type DatabaseCreateRequest struct {
	DCSlug         string `json:"dcslug"`
	NetworkType    string `json:"network_type"`
	VPC            string `json:"vpc,omitempty"`
	Firewall       string `json:"firewall,omitempty"`
	ClusterLabel   string `json:"cluster_label"`
	ClusterEngine  string `json:"cluster_engine"`
	ClusterVersion string `json:"cluster_version"`
	Size           string `json:"size"`
	Billing        string `json:"billing"`
	PITREnabled    string `json:"pitr_enabled"`
	ReplicaCount   string `json:"replica_count"`
}

type DatabaseReadonlyRequest struct {
	DCSlug       string `json:"dcslug"`
	SizeValue    string `json:"size_value"`
	Size         string `json:"size"`
	PlanCost     string `json:"plan_cost"`
	VPC          string `json:"vpc,omitempty"`
	ClusterLabel string `json:"cluster_label"`
	Subnet       string `json:"subnet,omitempty"`
}

type DatabaseDBRequest struct {
	Name string `json:"name"`
}

type DatabaseUserRequest struct {
	Name     string `json:"name"`
	Password string `json:"password"`
}

type DatabasePermissionRequest struct {
	DBUser     string `json:"db_user"`
	Permission string `json:"permission"`
}

type DatabasePoolRequest struct {
	Name string `json:"name"`
	DB   string `json:"db"`
	User string `json:"user"`
	Mode string `json:"mode"`
	Size int    `json:"size"`
}

type DatabaseBackupRequest struct {
	ClusterName string `json:"clustername"`
}

type DatabaseBackupRestoreRequest struct {
	DBName string `json:"dbname"`
}

type DatabaseSnapshotRestoreRequest struct {
	SnapshotID string `json:"snapshotid"`
	CloudID    string `json:"cloudid"`
}

type DatabaseResizeRequest struct {
	CloudID string `json:"cloudid"`
	Plan    string `json:"plan"`
}

type DatabaseTrustedHostRequest struct {
	TrustedHost string `json:"trusted_host"`
	Type        string `json:"type"`
}

// ── Response structs ──────────────────────────────────────────────────────

type DatabaseConnectionStrings struct {
	URI         string `json:"uri"`
	URIPrivate  string `json:"uri_private"`
	Host        string `json:"host"`
	HostPrivate string `json:"host_private"`
	Port        string `json:"port"`
	User        string `json:"user"`
	Database    string `json:"database"`
}

type DatabaseNode struct {
	CloudID           string                    `json:"cloudid"`
	ConnectionStrings DatabaseConnectionStrings `json:"connection_strings"`
}

type DatabaseNodes struct {
	Primary []DatabaseNode `json:"primary"`
}

type DatabaseInstance struct {
	ID              string        `json:"id"`
	ClusterName     string        `json:"cluster_name"`
	Engine          string        `json:"engine"`
	Version         string        `json:"version"`
	Port            string        `json:"port"`
	DCSlug          string        `json:"dcslug"`
	Status          string        `json:"status"`
	CreatedAt       string        `json:"created_at"`
	DefaultUser     string        `json:"dbdefault_user"`
	DefaultPass     string        `json:"dbdefault_pass"`
	DefaultDBName   string        `json:"dbdefault_dbname"`
	IsSSL           string        `json:"is_ssl"`
	PITREnabled     string        `json:"pitr_enabled"`
	AutomatedBackup string        `json:"automated_backup"`
	Nodes           DatabaseNodes `json:"nodes"`
}

// GetCloudID returns the primary node cloud ID
func (d *DatabaseInstance) GetCloudID() string {
	if len(d.Nodes.Primary) > 0 {
		return d.Nodes.Primary[0].CloudID
	}
	return ""
}

// GetHost returns the public connection host
func (d *DatabaseInstance) GetHost() string {
	if len(d.Nodes.Primary) > 0 {
		return d.Nodes.Primary[0].ConnectionStrings.Host
	}
	return ""
}

// GetHostPrivate returns the private connection host
func (d *DatabaseInstance) GetHostPrivate() string {
	if len(d.Nodes.Primary) > 0 {
		return d.Nodes.Primary[0].ConnectionStrings.HostPrivate
	}
	return ""
}

// GetURI returns the public connection URI
func (d *DatabaseInstance) GetURI() string {
	if len(d.Nodes.Primary) > 0 {
		return d.Nodes.Primary[0].ConnectionStrings.URI
	}
	return ""
}

// GetURIPrivate returns the private connection URI
func (d *DatabaseInstance) GetURIPrivate() string {
	if len(d.Nodes.Primary) > 0 {
		return d.Nodes.Primary[0].ConnectionStrings.URIPrivate
	}
	return ""
}

type DatabaseListResponse struct {
	Status    string             `json:"status"`
	Databases []DatabaseInstance `json:"databases"`
}

type DatabaseDetailResponse struct {
	Status    string             `json:"status"`
	Databases []DatabaseInstance `json:"databases"`
}

type DatabaseUserResponse struct {
	Status            string `json:"status"`
	Message           string `json:"message"`
	Username          string `json:"username"`
	GeneratedPassword string `json:"generated_password"`
}

// ── Cluster methods ───────────────────────────────────────────────────────

var ErrDatabaseStillProvisioning = fmt.Errorf("STILL_PROVISIONING")

func (c *Client) WaitForDatabaseReady(clusterID string) error {
	for attempt := 0; attempt < 60; attempt++ {
		if attempt > 0 {
			time.Sleep(10 * time.Second)
		}
		db, err := c.GetDatabase(clusterID, "")
		if err != nil {
			continue
		}
		if db == nil {
			continue
		}
		status := strings.ToLower(db.Status)
		if status == "active" || status == "running" {
			return nil
		}
		if status == "failed" || status == "error" {
			return fmt.Errorf("database cluster entered a failed state — check the Utho Console")
		}
	}
	return ErrDatabaseStillProvisioning
}

func (c *Client) CreateDatabase(req *DatabaseCreateRequest) (string, error) {
	respBytes, err := c.Post("/databases", req)
	if err != nil {
		return "", fmt.Errorf("failed to create database cluster: %w", err)
	}
	var result map[string]interface{}
	if err := json.Unmarshal(respBytes, &result); err != nil {
		return "", fmt.Errorf("failed to parse create database response: %w", err)
	}
	if result["status"] != "success" {
		return "", fmt.Errorf("create database failed: %s", result["message"])
	}
	id := fmt.Sprintf("%v", result["id"])
	return id, nil
}

func (c *Client) GetDatabase(clusterID string, cloudID string) (*DatabaseInstance, error) {
	// Use list API to find cluster by ID
	respBytes, err := c.Get("/databases")
	if err != nil {
		return nil, fmt.Errorf("failed to list databases: %w", err)
	}
	var resp DatabaseListResponse
	if err := json.Unmarshal(respBytes, &resp); err != nil {
		return nil, fmt.Errorf("failed to parse database list response: %w", err)
	}
	for _, db := range resp.Databases {
		if db.ID == clusterID {
			return &db, nil
		}
	}
	return nil, nil
}

func (c *Client) DeleteDatabase(clusterID string, cloudID string, clusterName string) error {
	endpoint := fmt.Sprintf("/databases/%s/%s?confirm=%s", clusterID, cloudID, clusterName)
	respBytes, err := c.Delete(endpoint)
	if err != nil {
		return fmt.Errorf("failed to delete database cluster: %w", err)
	}
	var result map[string]string
	if err := json.Unmarshal(respBytes, &result); err != nil {
		return fmt.Errorf("failed to parse delete database response: %w", err)
	}
	if result["status"] != "success" {
		return fmt.Errorf("delete database failed: %s", result["message"])
	}
	return nil
}

func (c *Client) AddDatabaseReplica(clusterID string) (string, error) {
	respBytes, err := c.Post(fmt.Sprintf("/databases/%s/replica", clusterID), nil)
	if err != nil {
		return "", fmt.Errorf("failed to add replica: %w", err)
	}
	var result map[string]interface{}
	if err := json.Unmarshal(respBytes, &result); err != nil {
		return "", fmt.Errorf("failed to parse add replica response: %w", err)
	}
	if result["status"] != "success" {
		return "", fmt.Errorf("add replica failed: %s", result["message"])
	}
	cloudID := fmt.Sprintf("%v", result["cloudid"])
	return cloudID, nil
}

func (c *Client) DeleteDatabaseReplica(clusterID string, cloudID string, clusterName string) error {
	endpoint := fmt.Sprintf("/databases/%s/replica/%s?confirm=%s", clusterID, cloudID, clusterName)
	respBytes, err := c.Delete(endpoint)
	if err != nil {
		return fmt.Errorf("failed to delete replica: %w", err)
	}
	var result map[string]string
	if err := json.Unmarshal(respBytes, &result); err != nil {
		return fmt.Errorf("failed to parse delete replica response: %w", err)
	}
	if result["status"] != "success" {
		return fmt.Errorf("delete replica failed: %s", result["message"])
	}
	return nil
}

func (c *Client) AddReadonlyReplica(clusterID string, req *DatabaseReadonlyRequest) error {
	respBytes, err := c.Post(fmt.Sprintf("/databases/%s/readonly", clusterID), req)
	if err != nil {
		return fmt.Errorf("failed to add readonly replica: %w", err)
	}
	var result map[string]string
	if err := json.Unmarshal(respBytes, &result); err != nil {
		return fmt.Errorf("failed to parse add readonly replica response: %w", err)
	}
	if result["status"] != "success" {
		return fmt.Errorf("add readonly replica failed: %s", result["message"])
	}
	return nil
}

func (c *Client) ResizeDatabase(clusterID string, req *DatabaseResizeRequest) error {
	respBytes, err := c.Post(fmt.Sprintf("/databases/%s/resize", clusterID), req)
	if err != nil {
		return fmt.Errorf("failed to resize database: %w", err)
	}
	var result map[string]string
	if err := json.Unmarshal(respBytes, &result); err != nil {
		return fmt.Errorf("failed to parse resize database response: %w", err)
	}
	if result["status"] != "success" {
		return fmt.Errorf("resize database failed: %s", result["message"])
	}
	return nil
}

// ── Database methods ──────────────────────────────────────────────────────

func (c *Client) CreateDBDatabase(clusterID string, name string) error {
	respBytes, err := c.Post(fmt.Sprintf("/databases/%s/database", clusterID), &DatabaseDBRequest{Name: name})
	if err != nil {
		return fmt.Errorf("failed to create database: %w", err)
	}
	var result map[string]string
	if err := json.Unmarshal(respBytes, &result); err != nil {
		return fmt.Errorf("failed to parse create database response: %w", err)
	}
	if result["status"] != "success" {
		return fmt.Errorf("create database failed: %s", result["message"])
	}
	return nil
}

func (c *Client) DeleteDBDatabase(clusterID string, dbName string) error {
	endpoint := fmt.Sprintf("/databases/%s/database/%s?confirm=%s", clusterID, dbName, dbName)
	respBytes, err := c.Delete(endpoint)
	if err != nil {
		return fmt.Errorf("failed to delete database: %w", err)
	}
	var result map[string]string
	if err := json.Unmarshal(respBytes, &result); err != nil {
		return fmt.Errorf("failed to parse delete database response: %w", err)
	}
	if result["status"] != "success" {
		return fmt.Errorf("delete database failed: %s", result["message"])
	}
	return nil
}

func (c *Client) AssignDBPermission(clusterID string, dbName string, req *DatabasePermissionRequest) error {
	respBytes, err := c.Post(fmt.Sprintf("/databases/%s/database/%s/user", clusterID, dbName), req)
	if err != nil {
		return fmt.Errorf("failed to assign permission: %w", err)
	}
	var result map[string]string
	if err := json.Unmarshal(respBytes, &result); err != nil {
		return fmt.Errorf("failed to parse assign permission response: %w", err)
	}
	if result["status"] != "success" {
		return fmt.Errorf("assign permission failed: %s", result["message"])
	}
	return nil
}

// ── User methods ──────────────────────────────────────────────────────────

func (c *Client) CreateDBUser(clusterID string, req *DatabaseUserRequest) (*DatabaseUserResponse, error) {
	respBytes, err := c.Post(fmt.Sprintf("/databases/%s/user", clusterID), req)
	if err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}
	var result DatabaseUserResponse
	if err := json.Unmarshal(respBytes, &result); err != nil {
		return nil, fmt.Errorf("failed to parse create user response: %w", err)
	}
	if result.Status != "success" {
		return nil, fmt.Errorf("create user failed: %s", result.Message)
	}
	return &result, nil
}

func (c *Client) DeleteDBUser(clusterID string, username string) error {
	endpoint := fmt.Sprintf("/databases/%s/user/?confirm=%s", clusterID, username)
	respBytes, err := c.Delete(endpoint)
	if err != nil {
		return fmt.Errorf("failed to delete user: %w", err)
	}
	var result map[string]string
	if err := json.Unmarshal(respBytes, &result); err != nil {
		return fmt.Errorf("failed to parse delete user response: %w", err)
	}
	if result["status"] != "success" {
		return fmt.Errorf("delete user failed: %s", result["message"])
	}
	return nil
}

// ── Connection Pool methods ───────────────────────────────────────────────

func (c *Client) CreateConnectionPool(clusterID string, cloudID string, req *DatabasePoolRequest) (string, error) {
	respBytes, err := c.Post(fmt.Sprintf("/databases/%s/%s/pool", clusterID, cloudID), req)
	if err != nil {
		return "", fmt.Errorf("failed to create connection pool: %w", err)
	}
	var result map[string]interface{}
	if err := json.Unmarshal(respBytes, &result); err != nil {
		return "", fmt.Errorf("failed to parse create connection pool response: %w", err)
	}
	if result["status"] != "success" {
		return "", fmt.Errorf("create connection pool failed: %s", result["message"])
	}
	id := fmt.Sprintf("%v", result["id"])
	return id, nil
}

func (c *Client) DeleteConnectionPool(clusterID string, cloudID string, poolID string) error {
	endpoint := fmt.Sprintf("/databases/%s/%s/pool?id=%s", clusterID, cloudID, poolID)
	respBytes, err := c.Delete(endpoint)
	if err != nil {
		return fmt.Errorf("failed to delete connection pool: %w", err)
	}
	var result map[string]string
	if err := json.Unmarshal(respBytes, &result); err != nil {
		return fmt.Errorf("failed to parse delete connection pool response: %w", err)
	}
	if result["status"] != "success" {
		return fmt.Errorf("delete connection pool failed: %s", result["message"])
	}
	return nil
}

// ── Backup methods ────────────────────────────────────────────────────────

func (c *Client) CreateDBBackup(clusterID string, clusterName string) error {
	respBytes, err := c.Post(fmt.Sprintf("/databases/%s/backup", clusterID), &DatabaseBackupRequest{ClusterName: clusterName})
	if err != nil {
		return fmt.Errorf("failed to create backup: %w", err)
	}
	var result map[string]string
	if err := json.Unmarshal(respBytes, &result); err != nil {
		return fmt.Errorf("failed to parse create backup response: %w", err)
	}
	if result["status"] != "success" {
		return fmt.Errorf("create backup failed: %s", result["message"])
	}
	return nil
}

func (c *Client) DeleteDBBackup(clusterID string, backupID string, dbName string) error {
	respBytes, err := c.doRequest("DELETE",
		fmt.Sprintf("/databases/%s/backup/%s/delete", clusterID, backupID),
		map[string]string{"dbname": dbName})
	if err != nil {
		return fmt.Errorf("failed to delete backup: %w", err)
	}
	var result map[string]string
	if err := json.Unmarshal(respBytes, &result); err != nil {
		return fmt.Errorf("failed to parse delete backup response: %w", err)
	}
	if result["status"] != "success" {
		return fmt.Errorf("delete backup failed: %s", result["message"])
	}
	return nil
}

func (c *Client) RestoreDBBackup(clusterID string, backupID string, dbName string) error {
	respBytes, err := c.Post(fmt.Sprintf("/databases/%s/backup/%s/restore", clusterID, backupID),
		&DatabaseBackupRestoreRequest{DBName: dbName})
	if err != nil {
		return fmt.Errorf("failed to restore backup: %w", err)
	}
	var result map[string]string
	if err := json.Unmarshal(respBytes, &result); err != nil {
		return fmt.Errorf("failed to parse restore backup response: %w", err)
	}
	if result["status"] != "success" {
		return fmt.Errorf("restore backup failed: %s", result["message"])
	}
	return nil
}

// ── PITR methods ──────────────────────────────────────────────────────────

func (c *Client) EnablePITR(clusterID string, cloudID string) error {
	endpoint := fmt.Sprintf("/databases?action=pitr-enabled&clusterid=%s&cloudid=%s", clusterID, cloudID)
	respBytes, err := c.Post(endpoint, nil)
	if err != nil {
		return fmt.Errorf("failed to enable PITR: %w", err)
	}
	var result map[string]string
	if err := json.Unmarshal(respBytes, &result); err != nil {
		return fmt.Errorf("failed to parse enable PITR response: %w", err)
	}
	if result["status"] != "success" {
		return fmt.Errorf("enable PITR failed: %s", result["message"])
	}
	return nil
}

func (c *Client) DisablePITR(clusterID string, cloudID string) error {
	endpoint := fmt.Sprintf("/databases?action=pitr-disabled&clusterid=%s&cloudid=%s", clusterID, cloudID)
	respBytes, err := c.Post(endpoint, nil)
	if err != nil {
		return fmt.Errorf("failed to disable PITR: %w", err)
	}
	var result map[string]string
	if err := json.Unmarshal(respBytes, &result); err != nil {
		return fmt.Errorf("failed to parse disable PITR response: %w", err)
	}
	if result["status"] != "success" {
		return fmt.Errorf("disable PITR failed: %s", result["message"])
	}
	return nil
}

// ── Security Group methods ────────────────────────────────────────────────

func (c *Client) AttachDBSecurityGroup(clusterID string, sgID string) error {
	respBytes, err := c.Post(fmt.Sprintf("/databases/%s/securitygroup/%s", clusterID, sgID), nil)
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

func (c *Client) DetachDBSecurityGroup(clusterID string, sgID string) error {
	respBytes, err := c.Delete(fmt.Sprintf("/databases/%s/securitygroup/%s", clusterID, sgID))
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

// ── Trusted Host methods ──────────────────────────────────────────────────

func (c *Client) AddTrustedHost(clusterID string, host string) error {
	respBytes, err := c.Post(fmt.Sprintf("/databases/%s/trustedhost", clusterID),
		&DatabaseTrustedHostRequest{TrustedHost: host, Type: "add"})
	if err != nil {
		return fmt.Errorf("failed to add trusted host: %w", err)
	}
	var result map[string]string
	if err := json.Unmarshal(respBytes, &result); err != nil {
		return fmt.Errorf("failed to parse add trusted host response: %w", err)
	}
	if result["status"] != "success" {
		return fmt.Errorf("add trusted host failed: %s", result["message"])
	}
	return nil
}

func (c *Client) DeleteTrustedHost(clusterID string, host string) error {
	endpoint := fmt.Sprintf("/databases/%s/trustedhost/%s?confirm=%s", clusterID, host, host)
	respBytes, err := c.Delete(endpoint)
	if err != nil {
		return fmt.Errorf("failed to delete trusted host: %w", err)
	}
	var result map[string]string
	if err := json.Unmarshal(respBytes, &result); err != nil {
		return fmt.Errorf("failed to parse delete trusted host response: %w", err)
	}
	if result["status"] != "success" {
		return fmt.Errorf("delete trusted host failed: %s", result["message"])
	}
	return nil
}
