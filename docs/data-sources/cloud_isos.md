---
page_title: "Utho: utho_cloud_isos"
subcategory: "Cloud Instances"
description: |-
  List available ISOs in your Utho account.
---

# utho_cloud_isos

Fetches all ISOs uploaded to your account. Use ISO names to deploy cloud instances with custom OS images — useful for installing operating systems not available in the standard image library.

## Example Usage

### List all ISOs

```hcl
data "utho_cloud_isos" "all" {}

output "isos" {
  value = data.utho_cloud_isos.all.isos
}
```

### List ISOs available in a specific DC

```hcl
data "utho_cloud_isos" "mumbai" {
  dcslug = "inmumbaizone2"
}

output "mumbai_isos" {
  value = data.utho_cloud_isos.mumbai.isos[*].name
}
```

### Deploy from an ISO

```hcl
data "utho_cloud_isos" "mumbai" {
  dcslug = "inmumbaizone2"
}

locals {
  my_iso = one([
    for iso in data.utho_cloud_isos.mumbai.isos :
    iso if iso.name == "my-custom-os"
  ])
}

resource "utho_cloud" "custom" {
  hostname        = "custom-01.mhc"
  dcslug          = "inmumbaizone2"
  planid          = "10308"
  billingcycle    = "hourly"
  auth            = "option1"
  root_password   = var.root_password
  iso             = local.my_iso.name
  enable_publicip = "true"
}
```

## Argument Reference

| Argument | Type   | Required | Description |
|----------|--------|----------|-------------|
| `dcslug` | String | No       | Filter ISOs by data center slug. ISOs are uploaded to specific DCs. |

## Attribute Reference

### Top-level

| Attribute | Type | Description |
|-----------|------|-------------|
| `isos`    | List | List of available ISOs. |

### isos

| Attribute  | Type   | Description |
|------------|--------|-------------|
| `name`     | String | ISO name — use this as `iso` in `utho_cloud`. |
| `file`     | String | Internal ISO filename. |
| `dc`       | String | Data center where this ISO is available. |
| `added_at` | String | Upload timestamp. |

## Notes

- ISOs are DC-specific — an ISO uploaded to `inmumbaizone2` is not available in `innoida`.
- Always filter by `dcslug` to ensure the ISO is available in your target DC.
- Upload ISOs from the [Utho Dashboard](https://console.utho.com) before referencing them here.
