---
page_title: "EBS DC Zones - Utho"
subcategory: "Elastic Block Storage"
description: |-
  Retrieve available data center zones for Utho EBS deployment.
---

# data.utho_ebs_dczones

Retrieves the list of data center zones where Utho Elastic Block Storage (EBS) is available. Use the `slug` value as the `dcslug` argument when creating an `utho_ebs` resource.

~> **Note:** EBS is currently available in **Delhi (Noida)**, **Mumbai**, and **Bangalore**.

## Example Usage

```hcl
data "utho_ebs_dczones" "available" {}

output "ebs_zones" { value = data.utho_ebs_dczones.available.dczones }

# Use the first available zone for an EBS volume
resource "utho_ebs" "data" {
  name       = "app-data"
  dcslug     = data.utho_ebs_dczones.available.dczones[0].slug
  disk_type  = "SSD"
  disk       = "50"
  iops       = "250"
  throughput = "125"
}
```

## Attribute Reference

| Attribute  | Type | Description |
|------------|------|-------------|
| `id`       | String | Static identifier `ebs_dczones`. |
| `dczones`  | List   | List of DC zones where EBS is available. |

### dczones

| Attribute | Type   | Description |
|-----------|--------|-------------|
| `id`      | String | Data center ID. |
| `slug`    | String | Data center slug (use this as `dcslug` in `utho_ebs`). |
| `city`    | String | City name (e.g. `Delhi (Noida)`, `Mumbai`, `Bangalore`). |
| `country` | String | Country name. |
| `cc`      | String | Country code (e.g. `in`). |
| `status`  | String | Data center status (`active`). |

## Available Zones

| Slug           | City          |
|----------------|---------------|
| `innoida`      | Delhi (Noida) |
| `inmumbaizone2`| Mumbai        |
| `inbangalore`  | Bangalore     |