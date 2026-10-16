---
page_title: "Utho: utho_kubernetes_node_pool"
subcategory: "Kubernetes"
description: |-
  Add and manage node pools in a Utho Kubernetes cluster.
---

# utho_kubernetes_node_pool

Adds a new node pool to an existing Utho Kubernetes cluster and manages its worker count and autoscaling settings. Use this resource to add pools after cluster creation.

Destroying this resource scales the node pool to 0 workers (no delete API exists).

## Example Usage

### Add a node pool to an existing cluster

```hcl
resource "utho_kubernetes_node_pool" "gpu_pool" {
  cluster_id = utho_kubernetes.main.id
  label      = "gpu-pool"
  size       = "10400"
  count      = "2"
  min_nodes  = "1"
  max_nodes  = "5"
}
```

### Node pool with additional EBS storage

```hcl
resource "utho_kubernetes_node_pool" "storage_pool" {
  cluster_id = utho_kubernetes.main.id
  label      = "storage-pool"
  size       = "10355"
  count      = "3"
  min_nodes  = "2"
  max_nodes  = "10"
  disk_size  = "50"
  disk_type  = "nvme"
}
```

### Update node count

```hcl
resource "utho_kubernetes_node_pool" "workers" {
  cluster_id = utho_kubernetes.main.id
  label      = "worker-pool"
  size       = "10355"
  count      = "5"      # change this and run terraform apply
  min_nodes  = "2"
  max_nodes  = "10"
}
```

## Argument Reference

| Argument     | Type   | Required | Description |
|--------------|--------|----------|-------------|
| `cluster_id` | String | Yes      | Kubernetes cluster ID. Changing this forces a new resource. |
| `label`      | String | Yes      | Node pool label. Changing this forces a new resource. |
| `size`       | String | Yes      | Plan ID for worker node size. Changing this forces a new resource. |
| `count`      | String | Yes      | Desired number of worker nodes. Updatable in place. |
| `min_nodes`  | String | Yes      | Minimum nodes for autoscaling. Updatable in place. |
| `max_nodes`  | String | Yes      | Maximum nodes for autoscaling. Updatable in place. |
| `disk_size`  | String | No       | Additional EBS disk size in GB. Changing this forces a new resource. |
| `disk_type`  | String | No       | EBS disk type: `nvme` or `ssd`. Changing this forces a new resource. |

## Attribute Reference

| Attribute | Type   | Description |
|-----------|--------|-------------|
| `id`      | String | Identifier in the format `{cluster_id}:{label}`. |
| `pool_id` | String | Node pool ID assigned by Utho. |

## Notes

- Changing `count`, `min_nodes`, or `max_nodes` updates the pool in place.
- Destroying this resource scales the pool to 0 nodes — it is not physically removed until the cluster is deleted.
