---
page_title: "Utho: utho_autoscalings"
subcategory: "Auto Scaling"
description: |-
  List all auto scaling groups in your Utho account.
---

# utho_autoscalings

Lists all auto scaling groups in your Utho account. Use this to look up groups created outside of the current configuration, or to report on capacity across environments.

## Example Usage

### List all auto scaling groups

```hcl
data "utho_autoscalings" "all" {}

output "autoscaling_groups" {
  value = data.utho_autoscalings.all.autoscalings[*].name
}
```

### Look up a group by name

```hcl
data "utho_autoscalings" "all" {}

locals {
  web_asg = one([
    for g in data.utho_autoscalings.all.autoscalings :
    g if g.name == "web-asg"
  ])
}

output "web_asg_capacity" {
  value = {
    min     = local.web_asg.minsize
    desired = local.web_asg.desiredsize
    max     = local.web_asg.maxsize
  }
}
```

## Argument Reference

This data source has no arguments.

## Attribute Reference

### Top-level

| Attribute      | Type | Description |
|----------------|------|-------------|
| `autoscalings` | List | List of auto scaling groups. See [autoscalings](#autoscalings). |

### autoscalings

| Attribute     | Type   | Description |
|---------------|--------|-------------|
| `id`          | String | Auto scaling group ID. |
| `name`        | String | Auto scaling group name. |
| `dcslug`      | String | Data center slug. |
| `minsize`     | String | Minimum number of instances. |
| `maxsize`     | String | Maximum number of instances. |
| `desiredsize` | String | Desired number of instances. |
| `status`      | String | Group status. |
