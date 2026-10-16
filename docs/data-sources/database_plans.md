---
page_title: "Utho: utho_database_plans"
subcategory: "Database (DBaaS)"
description: |-
  List the plans available for Utho managed database clusters.
---

# utho_database_plans

Lists the plans available for Utho managed database clusters. Use the plan `id` as the `size` argument of [`utho_database`](../resources/database) instead of hard-coding plan IDs.

## Example Usage

### List all plans

```hcl
data "utho_database_plans" "all" {}

output "database_plans" {
  value = [
    for p in data.utho_database_plans.all.plans :
    "${p.id}: ${p.name} (${p.cpu} vCPU, ${p.ram} MB RAM, ${p.disk} GB disk)"
  ]
}
```

### Pick a plan by size

```hcl
data "utho_database_plans" "all" {}

locals {
  # First plan with at least 2 vCPU and 4 GB RAM
  db_plan = [
    for p in data.utho_database_plans.all.plans :
    p if tonumber(p.cpu) >= 2 && tonumber(p.ram) >= 4096
  ][0]
}

resource "utho_database" "main" {
  cluster_name  = "production-db"
  dcslug        = "inmumbaizone2"
  engine        = "pg"
  version       = "17"
  size          = local.db_plan.id
  network_type  = "private"
  billing       = "monthly"
  replica_count = "0"
}
```

## Argument Reference

This data source has no arguments.

## Attribute Reference

### Top-level

| Attribute | Type | Description |
|-----------|------|-------------|
| `plans`   | List | List of available database plans. See [plans](#plans). |

### plans

| Attribute | Type   | Description |
|-----------|--------|-------------|
| `id`      | String | Plan ID. Use this as `size` in `utho_database`. |
| `name`    | String | Plan name. |
| `cpu`     | String | Number of vCPUs. |
| `ram`     | String | RAM in MB. |
| `disk`    | String | Disk size in GB. |
| `cost`    | String | Monthly cost. |

## Notes

- All numeric values are returned as strings. Wrap them in `tonumber()` before comparing, as in the example above.
