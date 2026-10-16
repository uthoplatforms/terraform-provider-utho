---
page_title: "Utho: utho_target_groups"
subcategory: "Auto Scaling"
description: |-
  List all target groups in your Utho account.
---

# utho_target_groups

Lists all target groups in your Utho account. Target groups are attached to auto scaling groups through the `target_groups` argument of [`utho_autoscaling`](../resources/autoscaling).

## Example Usage

### List all target groups

```hcl
data "utho_target_groups" "all" {}

output "target_groups" {
  value = [
    for tg in data.utho_target_groups.all.target_groups :
    "${tg.name} (${tg.protocol}:${tg.port})"
  ]
}
```

### Attach an existing target group to an auto scaling group

```hcl
data "utho_target_groups" "all" {}

locals {
  web_tg = one([
    for tg in data.utho_target_groups.all.target_groups :
    tg if tg.name == "web-tg"
  ])
}

resource "utho_autoscaling" "web" {
  # ... other configuration ...
  target_groups = local.web_tg.id
}
```

## Argument Reference

This data source has no arguments.

## Attribute Reference

### Top-level

| Attribute       | Type | Description |
|-----------------|------|-------------|
| `target_groups` | List | List of target groups. See [target_groups](#target_groups). |

### target_groups

| Attribute  | Type   | Description |
|------------|--------|-------------|
| `id`       | String | Target group ID. |
| `name`     | String | Target group name. |
| `protocol` | String | Protocol. |
| `port`     | String | Port. |
