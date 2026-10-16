package client

import (
	"encoding/json"
	"fmt"
	"time"
)

// ── Request structs ───────────────────────────────────────────────────────

type KubernetesNodePool struct {
	Label    string       `json:"label"`
	Size     string       `json:"size"`
	Count    string       `json:"count"`
	MinNodes string       `json:"min_nodes"`
	MaxNodes string       `json:"max_nodes"`
	EBS      []K8sEBSDisk `json:"ebs,omitempty"`
}

type K8sEBSDisk struct {
	Disk string `json:"disk"`
	Type string `json:"type"`
}

type KubernetesDeployRequest struct {
	DCSlug         string               `json:"dcslug"`
	ClusterLabel   string               `json:"cluster_label"`
	ClusterVersion string               `json:"cluster_version"`
	NodePools      []KubernetesNodePool `json:"nodepools"`
	VPC            string               `json:"vpc,omitempty"`
	NetworkType    string               `json:"network_type"`
	CPUModel       string               `json:"cpumodel,omitempty"`
}

type KubernetesDestroyRequest struct {
	Confirm string `json:"confirm"`
}

type NodePoolUpdateRequest struct {
	Count    string `json:"count"`
	MinNodes string `json:"min_nodes,omitempty"`
	MaxNodes string `json:"max_nodes,omitempty"`
}

type AddNodePoolRequest struct {
	NodePools []KubernetesNodePool `json:"nodepools"`
}

// ── Response structs ──────────────────────────────────────────────────────

type KubernetesCluster struct {
	ID          string  `json:"id"`
	CloudID     string  `json:"cloudid"`
	Hostname    string  `json:"hostname"`
	DCSlug      string  `json:"dcslug"`
	Status      string  `json:"status"`
	AppStatus   string  `json:"app_status"`
	IP          string  `json:"ip"`
	WorkerCount string  `json:"worker_count"`
	CreatedAt   string  `json:"created_at"`
	PowerStatus string  `json:"powerstatus"`
	Cost        float64 `json:"cost"`
}

type KubernetesListResponse struct {
	K8s   []KubernetesCluster `json:"k8s"`
	RCode string              `json:"rcode"`
}

type KubernetesClusterDetail struct {
	ID          string `json:"id"`
	Version     string `json:"version"`
	Label       string `json:"label"`
	DCSlug      string `json:"dcslug"`
	NetworkType string `json:"network_type"`
	Status      string `json:"status"`
	AppStatus   string `json:"app_status"`
	VPC         string `json:"vpc"`
	DNS         string `json:"dns"`
	IP          string `json:"ip"`
	CreatedAt   string `json:"created_at"`
}

type KubernetesDetailResponse struct {
	Info struct {
		Cluster KubernetesClusterDetail `json:"cluster"`
	} `json:"info"`
	K8s   []KubernetesCluster `json:"k8s"`
	RCode string              `json:"rcode"`
}

// ── API methods ───────────────────────────────────────────────────────────

// waitForK8sReady polls until cluster is ready (status=Active, app_status=Installed)
// Uses the list API since the detail API doesn't always return status fields
// K8s clusters take 5-15 minutes to provision
func (c *Client) waitForK8sReady(clusterID string) error {
	for attempt := 0; attempt < 60; attempt++ { // max 10 minutes
		if attempt > 0 {
			time.Sleep(10 * time.Second)
		}
		respBytes, err := c.Get("/kubernetes")
		if err != nil {
			continue
		}
		var resp KubernetesListResponse
		if json.Unmarshal(respBytes, &resp) != nil {
			continue
		}
		for _, k := range resp.K8s {
			if k.ID == clusterID {
				if k.Status == "Active" && k.AppStatus == "Installed" {
					return nil
				}
				if k.AppStatus == "Failed" {
					return fmt.Errorf("kubernetes cluster entered a Failed state — please check the Utho Console")
				}
				break
			}
		}
	}
	return nil // proceed after timeout
}

// getK8sIP fetches the cluster IP from the list API
func (c *Client) GetK8sIP(clusterID string) string {
	respBytes, err := c.Get("/kubernetes")
	if err != nil {
		return ""
	}
	var resp KubernetesListResponse
	if json.Unmarshal(respBytes, &resp) != nil {
		return ""
	}
	for _, k := range resp.K8s {
		if k.ID == clusterID {
			return k.IP
		}
	}
	return ""
}

func (c *Client) CreateKubernetesCluster(req *KubernetesDeployRequest) (string, error) {
	respBytes, err := c.Post("/kubernetes/deploy", req)
	if err != nil {
		return "", fmt.Errorf("failed to create Kubernetes cluster: %w", err)
	}
	var result map[string]interface{}
	if err := json.Unmarshal(respBytes, &result); err != nil {
		return "", fmt.Errorf("failed to parse create Kubernetes cluster response: %w", err)
	}
	if result["status"] != "success" {
		return "", fmt.Errorf("create Kubernetes cluster failed: %s", result["message"])
	}
	id := fmt.Sprintf("%v", result["id"])
	return id, nil
}

func (c *Client) GetKubernetesCluster(clusterID string) (*KubernetesClusterDetail, error) {
	respBytes, err := c.Get(fmt.Sprintf("/kubernetes/%s", clusterID))
	if err != nil {
		return nil, fmt.Errorf("failed to get Kubernetes cluster: %w", err)
	}
	var resp KubernetesDetailResponse
	if err := json.Unmarshal(respBytes, &resp); err != nil {
		return nil, fmt.Errorf("failed to parse get Kubernetes cluster response: %w", err)
	}
	if resp.Info.Cluster.ID == "" {
		return nil, nil
	}
	return &resp.Info.Cluster, nil
}

func (c *Client) DeleteKubernetesCluster(clusterID string) error {
	respBytes, err := c.doRequest("DELETE", fmt.Sprintf("/kubernetes/%s/destroy", clusterID),
		&KubernetesDestroyRequest{Confirm: "I am aware this action will delete data and cluster permanently"})
	if err != nil {
		return fmt.Errorf("failed to delete Kubernetes cluster: %w", err)
	}
	var result map[string]string
	if err := json.Unmarshal(respBytes, &result); err != nil {
		return fmt.Errorf("failed to parse delete Kubernetes cluster response: %w", err)
	}
	if result["status"] != "success" {
		return fmt.Errorf("delete Kubernetes cluster failed: %s", result["message"])
	}
	return nil
}

func (c *Client) AddNodePool(clusterID string, nodePools []KubernetesNodePool) error {
	respBytes, err := c.Post(fmt.Sprintf("/kubernetes/%s/nodepool/add", clusterID), &AddNodePoolRequest{
		NodePools: nodePools,
	})
	if err != nil {
		return fmt.Errorf("failed to add node pool: %w", err)
	}
	var result map[string]string
	if err := json.Unmarshal(respBytes, &result); err != nil {
		return fmt.Errorf("failed to parse add node pool response: %w", err)
	}
	if result["status"] != "success" {
		return fmt.Errorf("add node pool failed: %s", result["message"])
	}
	return nil
}

func (c *Client) UpdateNodePool(clusterID string, poolID string, req *NodePoolUpdateRequest) error {
	respBytes, err := c.Post(fmt.Sprintf("/kubernetes/%s/nodepool/%s/update", clusterID, poolID), req)
	if err != nil {
		return fmt.Errorf("failed to update node pool: %w", err)
	}
	var result map[string]string
	if err := json.Unmarshal(respBytes, &result); err != nil {
		return fmt.Errorf("failed to parse update node pool response: %w", err)
	}
	if result["status"] != "success" {
		return fmt.Errorf("update node pool failed: %s", result["message"])
	}
	return nil
}

func (c *Client) DeleteNodePool(clusterID string, poolID string) error {
	// Delete nodepool = scale count to 0
	return c.UpdateNodePool(clusterID, poolID, &NodePoolUpdateRequest{Count: "0"})
}

func (c *Client) DeleteWorkerNode(clusterID string, poolID string, workerID string) error {
	endpoint := fmt.Sprintf("/kubernetes/%s/nodepool/%s/%s/delete", clusterID, poolID, workerID)
	respBytes, err := c.Delete(endpoint)
	if err != nil {
		return fmt.Errorf("failed to delete worker node: %w", err)
	}
	var result map[string]string
	if err := json.Unmarshal(respBytes, &result); err != nil {
		return fmt.Errorf("failed to parse delete worker node response: %w", err)
	}
	if result["status"] != "success" {
		return fmt.Errorf("delete worker node failed: %s", result["message"])
	}
	return nil
}

func (c *Client) GetKubeconfig(clusterID string) (string, error) {
	respBytes, err := c.Get(fmt.Sprintf("/kubernetes/%s/download", clusterID))
	if err != nil {
		return "", fmt.Errorf("failed to get kubeconfig: %w", err)
	}
	return string(respBytes), nil
}
