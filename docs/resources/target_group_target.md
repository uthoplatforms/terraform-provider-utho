---
page_title: "Utho: utho_target_group_target"
subcategory: "Load Balancing"
description: |-
  Register a backend server as a target in a Utho Target Group.
---

# utho_target_group_target

Registers a backend server as a target in a Utho Target Group. Supports cloud instances, Kubernetes node pools, and custom IPs. Destroying this resource removes the target from the group.

## Example Usage

### Cloud instance target

```hcl
resource "utho_target_group_target" "web" {
  target_group_id = utho_target_group.api.id
  type            = "cloud"
  backend_port    = "80"
  cloud_id        = "990001160"
  ip              = "103.127.30.105"
}
```

### Custom IP target

```hcl
resource "utho_target_group_target" "custom" {
  target_group_id = utho_target_group.api.id
  type            = "custom"
  backend_port    = "8080"
  ip              = "192.168.1.100"
}
```

### Kubernetes node pool target

```hcl
resource "utho_target_group_target" "k8s" {
  target_group_id      = utho_target_group.api.id
  type                 = "kubernetes"
  backend_port         = "3000"
  kubernetes_cluster_id = "751036"
  pool_name            = "pool-w7vlhlbx"
}
```

## Argument Reference

| Argument               | Type   | Required | Description |
|------------------------|--------|----------|-------------|
| `target_group_id`      | String | Yes      | Target group ID. Changing forces new resource. |
| `type`                 | String | Yes      | Target type: `cloud`, `kubernetes`, `custom`. Changing forces new resource. |
| `backend_port`         | String | Yes      | Port on which the target receives traffic. Changing forces new resource. |
| `cloud_id`             | String | No       | Cloud instance ID. Required when `type` is `cloud`. |
| `ip`                   | String | No       | IP address. Required when `type` is `custom`. For `cloud` targets, provide the instance public IP. |
| `kubernetes_cluster_id`| String | No       | Kubernetes cluster ID. Required when `type` is `kubernetes`. |
| `pool_name`            | String | No       | Kubernetes node pool name. Required when `type` is `kubernetes`. All node IPs in the pool are registered as targets. |

## Attribute Reference

| Attribute | Type   | Description |
|-----------|--------|-------------|
| `id`      | String | Target ID. |

## Notes

- All arguments force a new resource on change — targets cannot be updated in place.
- For Kubernetes targets, all node IPs in the selected pool are registered as individual targets.
- Destroying this resource removes the target from the group but does not affect the underlying instance.