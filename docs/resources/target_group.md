---
page_title: "Utho: utho_target_group"
subcategory: "Load Balancing"
description: |-
  Create and manage Utho Target Groups.
---

# utho_target_group

Creates and manages a Utho Target Group. A target group defines a pool of backend servers along with health check configuration. Attach it to a Load Balancer frontend or Auto Scaling group to route traffic.

## Example Usage

```hcl
resource "utho_target_group" "api" {
  name                 = "api-targets"
  protocol             = "HTTP"
  port                 = "80"
  health_check_path     = "/health"
  health_check_protocol = "HTTP"
  health_check_interval = "30"
  health_check_timeout  = "5"
  healthy_threshold    = "3"
  unhealthy_threshold  = "3"
}
```

## Argument Reference

| Argument               | Type   | Required | Description |
|------------------------|--------|----------|-------------|
| `name`                 | String | Yes      | Target group name. Must be unique in your account. |
| `protocol`             | String | Yes      | Traffic protocol: `HTTP`, `HTTPS`, `TCP`, `UDP`. |
| `port`                 | String | Yes      | Port on which registered targets receive traffic (e.g. `80`, `443`). |
| `health_check_path`    | String | Yes      | Health check endpoint path (e.g. `/health`, `/status`). |
| `health_check_protocol`| String | Yes      | Health check protocol: `HTTP`, `HTTPS`, `TCP`, `UDP`. |
| `health_check_interval`| String | Yes      | Health check interval in seconds (e.g. `30`). |
| `health_check_timeout` | String | Yes      | Health check timeout in seconds (e.g. `5`). |
| `healthy_threshold`    | String | Yes      | Consecutive successful checks to mark a target healthy (e.g. `3`). |
| `unhealthy_threshold`  | String | Yes      | Consecutive failed checks to mark a target unhealthy (e.g. `3`). |

## Attribute Reference

| Attribute    | Type   | Description |
|--------------|--------|-------------|
| `id`         | String | Target group ID. |
| `created_at` | String | Creation timestamp. |

## Related Resources

- [utho_target_group_target](target_group_target) — Register backend servers in this target group