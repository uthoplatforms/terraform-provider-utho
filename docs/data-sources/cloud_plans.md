---
page_title: "Utho: utho_cloud_plans"
subcategory: "Cloud Instances"
description: |-
  List available plans for Utho Cloud instances.
---

# utho_cloud_plans

Fetches available plans (instance sizes) for cloud instances. Plans define the CPU, RAM, disk, and price for a cloud instance. Filter by `dcslug` to see only plans available in a specific data center.

~> **Important:** Always filter by `dcslug` when creating instances — plan availability varies by data center. Only plans with `is_available = "YES"` are returned.

## Example Usage

### List all plans in Mumbai

```hcl
data "utho_cloud_plans" "mumbai" {
  dcslug = "inmumbaizone2"
}

output "plans" {
  value = data.utho_cloud_plans.mumbai.plans
}
```

### Find the cheapest plan with disk included

```hcl
data "utho_cloud_plans" "mumbai" {
  dcslug = "inmumbaizone2"
}

locals {
  # Only plans that include disk (disk != "0")
  plans_with_disk = [
    for p in data.utho_cloud_plans.mumbai.plans :
    p if p.disk != "0"
  ]

  # Sort by price and pick cheapest
  cheapest = local.plans_with_disk[
    index(
      local.plans_with_disk[*].price,
      min(local.plans_with_disk[*].price...)
    )
  ]
}

output "cheapest_plan_id" {
  value = local.cheapest.id
}
```

### Use plan ID in a cloud instance

```hcl
data "utho_cloud_plans" "mumbai" {
  dcslug = "inmumbaizone2"
}

locals {
  # Find the basic 2vCPU/4GB plan with disk
  plan = one([
    for p in data.utho_cloud_plans.mumbai.plans :
    p if p.id == "10308"
  ])
}

resource "utho_cloud" "web" {
  hostname        = "web-01.mhc"
  dcslug          = "inmumbaizone2"
  planid          = local.plan.id
  billingcycle    = "hourly"
  auth            = "option1"
  root_password   = var.root_password
  image           = "ubuntu-22.04-x86_64"
  enable_publicip = "true"
}
```

### Find dedicated CPU plans only

```hcl
data "utho_cloud_plans" "mumbai" {
  dcslug = "inmumbaizone2"
}

output "dedicated_plans" {
  value = [
    for p in data.utho_cloud_plans.mumbai.plans :
    { id = p.id, cpu = p.cpu, ram = p.ram, disk = p.disk, price = p.price }
    if p.slug == "dedicated-cpu" && p.disk != "0"
  ]
}
```

## Argument Reference

| Argument | Type   | Required | Description |
|----------|--------|----------|-------------|
| `dcslug` | String | No       | Filter plans by data center slug. Recommended — plan availability varies by DC. |

## Attribute Reference

### Top-level

| Attribute | Type | Description |
|-----------|------|-------------|
| `plans`   | List | List of available plans. Only includes plans where `is_available = "YES"`. |

### plans

| Attribute        | Type   | Description |
|------------------|--------|-------------|
| `id`             | String | Plan ID. Use this as `planid` in `utho_cloud`. |
| `name`           | String | Plan name. |
| `slug`           | String | Plan type: `basic`, `dedicated-cpu`, `dedicated-memory`, `gpu`. |
| `cpu`            | String | Number of vCPUs. |
| `ram`            | String | RAM in MB (divide by 1024 for GB). |
| `disk`           | String | Included disk in GB. `"0"` means no disk — attach EBS separately. |
| `bandwidth`      | String | Bandwidth in GB. |
| `dedicated_vcore`| String | `"1"` for dedicated vCPU, `"0"` for shared. |
| `price`          | Float  | Monthly price in INR. |
| `is_available`   | String | Always `"YES"` (unavailable plans are filtered out). |

## Plan Types

| Slug | Description |
|------|-------------|
| `basic` | Shared CPU — affordable for dev and low-traffic workloads |
| `dedicated-cpu` | Dedicated vCPU — consistent performance for production |
| `dedicated-memory` | High RAM-to-CPU ratio — databases, in-memory caches |
| `gpu` | GPU-enabled — ML/AI workloads |

## Notes

- Plans with `disk = "0"` require an EBS volume for the OS disk. Use the `ebs` block in `utho_cloud` for these plans.
- Plans with `disk != "0"` include the OS disk in the plan price — no EBS needed.
- Always filter by `dcslug` — not all plans are available in every data center.
