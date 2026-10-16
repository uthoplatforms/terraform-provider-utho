---
page_title: "Utho: utho_cloud_power"
subcategory: "Cloud Instances"
description: |-
  Control the power state of a Utho Cloud instance.
---

# utho_cloud_power

Controls the power state of an existing Utho Cloud instance. Use this resource to start, stop, or reboot an instance without destroying and recreating it.

~> **Note:** This resource manages the power state only. The instance must already exist — either created via `utho_cloud` or imported.

## Example Usage

### Stop an instance

```hcl
resource "utho_cloud_power" "web" {
  cloud_id = utho_cloud.web.id
  action   = "poweroff"
}
```

### Start an instance

```hcl
resource "utho_cloud_power" "web" {
  cloud_id = utho_cloud.web.id
  action   = "poweron"
}
```

### Reboot an instance

```hcl
resource "utho_cloud_power" "web" {
  cloud_id = utho_cloud.web.id
  action   = "reboot"
}
```

### Stop before maintenance, start after

```hcl
# Stop
resource "utho_cloud_power" "maintenance" {
  cloud_id = utho_cloud.web.id
  action   = "poweroff"
}

# After maintenance — change action to poweron and apply
# resource "utho_cloud_power" "maintenance" {
#   cloud_id = utho_cloud.web.id
#   action   = "poweron"
# }
```

## Argument Reference

| Argument   | Type   | Required | Description |
|------------|--------|----------|-------------|
| `cloud_id` | String | Yes      | Cloud instance ID. Get from `utho_cloud.name.id`. Changing this forces a new resource. |
| `action`   | String | Yes      | Power action: `poweron`, `poweroff`, or `reboot`. |

## Attribute Reference

| Attribute | Type   | Description |
|-----------|--------|-------------|
| `id`      | String | Same as `cloud_id`. |

## Power Actions

| Action     | Description |
|------------|-------------|
| `poweron`  | Start a stopped instance. |
| `poweroff` | Gracefully stop a running instance. |
| `reboot`   | Restart a running instance. |

## Notes

- Changing `action` updates the power state in place — no destroy/recreate.
- After a `reboot` or `poweron`, the instance takes 30–60 seconds to be fully ready.
- Use `poweroff` before resizing with `utho_cloud_resize` to avoid data corruption.
- Destroying this resource does NOT destroy the cloud instance — it only removes the power state management from Terraform state.
