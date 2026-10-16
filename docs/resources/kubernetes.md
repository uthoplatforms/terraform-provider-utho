---
page_title: "Utho: utho_kubernetes"
subcategory: "Kubernetes"
description: |-
  Create and manage Utho Managed Kubernetes clusters.
---

# utho_kubernetes

Creates and manages a Utho Managed Kubernetes cluster. Utho provisions and manages the control plane — you just define your node pools and the cluster version. Worker nodes run in your account and are visible as cloud instances.

After creation, fetch the kubeconfig with the `utho_kubernetes_config` data source to authenticate `kubectl` or configure Kubernetes providers.

## Example Usage

### Single-zone cluster with one node pool

```hcl
resource "utho_kubernetes" "main" {
  dcslug          = "inmumbaizone2"
  cluster_label   = "production"
  cluster_version = "1.30.0-utho"
  network_type    = "public"

  nodepools = [
    {
      label     = "workers"
      size      = "10355"   # 4 vCPU / 8 GB RAM
      count     = "3"
      min_nodes = "2"
      max_nodes = "10"
      disk_size = "30"
      disk_type = "nvme"
    }
  ]
}

output "cluster_id" { value = utho_kubernetes.main.id }
output "cluster_dns" { value = utho_kubernetes.main.dns }
```

### Private cluster inside a VPC

Suitable for production workloads where the control plane endpoint should not be publicly accessible.

```hcl
resource "utho_kubernetes" "private" {
  dcslug          = "inmumbaizone2"
  cluster_label   = "private-cluster"
  cluster_version = "1.30.0-utho"
  network_type    = "publicprivate"
  vpc             = utho_subnet.private.id
  cpumodel        = "intel"

  nodepools = [
    {
      label     = "app-pool"
      size      = "10355"
      count     = "3"
      min_nodes = "2"
      max_nodes = "8"
    }
  ]
}
```

### Multi-pool cluster for workload separation

Use separate node pools for different workload types — general app servers, memory-intensive jobs, etc.

```hcl
resource "utho_kubernetes" "main" {
  dcslug          = "inmumbaizone2"
  cluster_label   = "multi-pool"
  cluster_version = "1.30.0-utho"
  network_type    = "public"

  nodepools = [
    {
      label     = "general"
      size      = "10355"
      count     = "3"
      min_nodes = "2"
      max_nodes = "10"
    }
  ]
}

# Add a high-memory pool for batch workloads
resource "utho_kubernetes_node_pool" "batch" {
  cluster_id = utho_kubernetes.main.id
  label      = "batch"
  size       = "10400"   # higher memory plan
  count      = "2"
  min_nodes  = "1"
  max_nodes  = "5"
}
```

### Get kubeconfig and use with kubectl

```hcl
data "utho_kubernetes_config" "main" {
  cluster_id = utho_kubernetes.main.id
}

# Save to file for kubectl
resource "local_file" "kubeconfig" {
  content         = data.utho_kubernetes_config.main.raw_config
  filename        = "${path.module}/kubeconfig.yaml"
  file_permission = "0600"
}

output "kubeconfig_command" {
  value = "export KUBECONFIG=${path.module}/kubeconfig.yaml"
}
```

## Argument Reference

### Required

| Argument           | Type   | Description |
|--------------------|--------|-------------|
| `dcslug`           | String | Data center for the cluster. Changing this forces a new resource. |
| `cluster_label`    | String | Cluster name. Changing this forces a new resource. |
| `cluster_version`  | String | Kubernetes version (e.g. `1.30.0-utho`). Changing this forces a new resource. |
| `network_type`     | String | `public`, `private`, or `publicprivate`. Changing this forces a new resource. |
| `nodepools`        | List   | One or more node pool definitions. See [Node Pool Block](#node-pool-block). |

### Optional

| Argument    | Type   | Description |
|-------------|--------|-------------|
| `vpc`       | String | VPC subnet ID. Required when `network_type` is `private` or `publicprivate`. Use `utho_subnet.name.id`. Changing this forces a new resource. |
| `cpumodel`  | String | CPU preference: `amd` or `intel`. Changing this forces a new resource. |

## Available Kubernetes Versions

Use one of these values for `cluster_version`:

| Version | Status |
|---------|--------|
| `1.28.9-utho` | Older |
| `1.29.4-utho` | Older |
| `1.30.0-utho` | Stable |
| `1.31.14-utho` | Stable |
| `1.32.10-utho` | Stable |
| `1.34.2-utho` | Latest |
| `1.35.6-utho` | Latest |
| `1.36.2-utho` | Latest |

~> **Tip:** Use the latest stable version (`1.32.10-utho`) for new clusters unless you have a specific version requirement.

### Node Pool Block

```hcl
nodepools = [
  {
    label     = "workers"
    size      = "10355"
    count     = "3"
    min_nodes = "2"
    max_nodes = "10"
    disk_size = "30"    # optional
    disk_type = "nvme"  # optional
  }
]
```

| Argument    | Type   | Required | Description |
|-------------|--------|----------|-------------|
| `label`     | String | Yes      | Node pool name. |
| `size`      | String | Yes      | Plan ID for the worker node size. |
| `count`     | String | Yes      | Initial number of worker nodes. |
| `min_nodes` | String | Yes      | Minimum nodes (for autoscaling). |
| `max_nodes` | String | Yes      | Maximum nodes (for autoscaling). |
| `disk_size` | String | No       | Additional EBS disk size in GB. |
| `disk_type` | String | No       | EBS disk type: `nvme` or `ssd`. |

## Attribute Reference

| Attribute    | Type   | Description |
|--------------|--------|-------------|
| `id`         | String | Unique cluster ID. Use this in `utho_kubernetes_node_pool` and `utho_kubernetes_config`. |
| `status`     | String | Cluster status (`Active`, `Pending`, etc.). |
| `ip`         | String | Control plane IP address. |
| `dns`        | String | Control plane DNS endpoint (used in kubeconfig). |
| `created_at` | String | Creation timestamp. |

## Notes

- Cluster creation takes 5–15 minutes. The cluster status will be `Pending` during this time.
- Node pools defined in the `nodepools` block at creation are the initial pools. Use `utho_kubernetes_node_pool` to add more later.
- Updating the cluster (changing version, network type, etc.) requires destroy and recreate — plan changes carefully for production clusters.
- Always save the kubeconfig output securely — treat it like a root password.