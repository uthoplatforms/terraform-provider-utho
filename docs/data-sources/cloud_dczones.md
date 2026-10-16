---
page_title: "Utho: utho_cloud_dczones"
subcategory: "Cloud Instances"
description: |-
  List all available data center zones for Utho Cloud instances.
---

# utho_cloud_dczones

Fetches all available data center zones for cloud instances. Use this to discover valid `dcslug` values before creating resources, and to check which zones support EBS volumes or VPCs.

## Example Usage

### List all active zones

```hcl
data "utho_cloud_dczones" "all" {}

output "active_zones" {
  value = [
    for z in data.utho_cloud_dczones.all.zones :
    z.slug if z.status == "active"
  ]
}
```

### Find zones with EBS support

```hcl
data "utho_cloud_dczones" "all" {}

output "ebs_zones" {
  value = [
    for z in data.utho_cloud_dczones.all.zones :
    { slug = z.slug, city = z.city }
    if z.ebs_available && z.status == "active"
  ]
}
```

### Use in a cloud instance

```hcl
data "utho_cloud_dczones" "all" {}

locals {
  # Pick Mumbai zone
  mumbai = one([
    for z in data.utho_cloud_dczones.all.zones :
    z if z.slug == "inmumbaizone2"
  ])
}

resource "utho_cloud" "web" {
  hostname        = "web-01.mhc"
  dcslug          = local.mumbai.slug
  planid          = "10308"
  billingcycle    = "hourly"
  auth            = "option1"
  root_password   = var.root_password
  image           = "ubuntu-22.04-x86_64"
  enable_publicip = "true"
}
```

## Attribute Reference

### Top-level

| Attribute | Type | Description |
|-----------|------|-------------|
| `zones`   | List | List of all data center zones. |

### zones

| Attribute            | Type   | Description |
|----------------------|--------|-------------|
| `id`                 | String | DC zone ID. |
| `slug`               | String | DC slug — use this as `dcslug` in all resources. |
| `city`               | String | City name (e.g. `Mumbai`). |
| `country`            | String | Country name (e.g. `India`). |
| `cc`                 | String | Country code (e.g. `in`, `de`, `us`). |
| `status`             | String | `active` or `inactive`. Only use `active` zones. |
| `default_cpu`        | String | Default CPU model for this zone: `amd` or `intel`. |
| `ebs_available`      | Bool   | `true` if EBS block volumes are supported in this zone. |
| `public_ip_available`| Bool   | `true` if public IPs can be assigned. |

## Available Zones

| Slug | City | Country | EBS | Status |
|------|------|---------|-----|--------|
| `innoida` | Delhi (Noida) | India | ✓ | active |
| `inmumbaizone2` | Mumbai | India | ✓ | active |
| `inbangalore` | Bangalore | India | ✗ | active |
| `defra1` | Frankfurt | Germany | ✗ | active |
| `uslosangeles` | Los Angeles | United States | ✗ | active |
