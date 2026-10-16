---
page_title: "Utho: utho_cloud_iso"
subcategory: "Cloud Instances"
description: |-
  Mount or unmount an ISO on a Utho Cloud instance.
---

# utho_cloud_iso

Mounts or unmounts an ISO image on an existing cloud instance. Use this for custom OS installations, rescue environments, or running bootable tools on a live instance.

## Example Usage

### Mount an ISO

```hcl
data "utho_cloud_isos" "mumbai" {
  dcslug = "inmumbaizone2"
}

locals {
  rescue_iso = one([
    for iso in data.utho_cloud_isos.mumbai.isos :
    iso if iso.name == "rescue-environment"
  ])
}

resource "utho_cloud_iso" "rescue" {
  cloud_id = utho_cloud.web.id
  iso      = local.rescue_iso.name
  action   = "mount"
}
```

### Unmount an ISO

```hcl
resource "utho_cloud_iso" "rescue" {
  cloud_id = utho_cloud.web.id
  iso      = "rescue-environment"
  action   = "unmount"
}
```

## Argument Reference

| Argument   | Type   | Required | Description |
|------------|--------|----------|-------------|
| `cloud_id` | String | Yes      | Cloud instance ID. Changing this forces a new resource. |
| `iso`      | String | Yes      | ISO name. Use [utho_cloud_isos](../data-sources/cloud_isos) to list available ISOs. |
| `action`   | String | Yes      | `mount` or `unmount`. |

## Attribute Reference

| Attribute | Type   | Description |
|-----------|--------|-------------|
| `id`      | String | Same as `cloud_id`. |

## Notes

- ISOs must be uploaded to your account before mounting. Upload from the Utho Dashboard.
- After mounting, reboot the instance to boot from the ISO.
- After OS installation or rescue operation, unmount the ISO and reboot again to boot from disk.
- Use `data.utho_cloud_isos` to list available ISOs and their names.
- ISOs are DC-specific — ensure the ISO is available in the same DC as your instance.
