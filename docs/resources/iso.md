---
page_title: "Utho: utho_iso"
subcategory: "Compute"
description: |-
  Upload and manage ISO images on Utho by providing a public download URL.
---

# utho_iso

Uploads and manages an ISO image on Utho by providing a public download URL. The ISO is downloaded asynchronously in the background. Use the computed `file` attribute as input to `utho_cloud_iso` to mount the ISO on a cloud instance.

~> **Note:** ISO downloads are asynchronous. The `file` attribute is populated immediately after the download request is accepted, but the ISO may still be downloading. Check download progress in the Utho console before mounting.

## Example Usage

```hcl
# Upload an ISO
resource "utho_iso" "alpine" {
  name   = "alpine-linux"
  dcslug = "innoida"
  url    = "https://dl-cdn.alpinelinux.org/alpine/latest-stable/releases/x86_64/alpine-virt-3.24.3-x86_64.iso"
}

# Mount it on a cloud instance
resource "utho_cloud_iso" "mount" {
  cloud_id = "990001160"
  iso      = utho_iso.alpine.file
}

output "iso_file" { value = utho_iso.alpine.file }
```

## Argument Reference

| Argument | Type   | Required | Description |
|----------|--------|----------|-------------|
| `name`   | String | Yes      | ISO name. Changing forces new resource. |
| `dcslug` | String | Yes      | Data center slug where ISO will be stored (e.g. `innoida`, `inmumbaizone2`). Changing forces new resource. |
| `url`    | String | Yes      | Public download URL of the ISO file. Changing forces new resource. |

## Attribute Reference

| Attribute  | Type   | Description |
|------------|--------|-------------|
| `id`       | String | ISO upload ID returned by Utho. |
| `file`     | String | ISO filename assigned by Utho (e.g. `PzFZp-194454.iso`). Use this as input to `utho_cloud_iso`. |
| `added_at` | String | Timestamp when the ISO was added. |
| `size`     | Number | ISO size in GB. |
| `location` | String | Data center location name. |

## Related Resources

- [utho_cloud_iso](cloud_iso) — Mount this ISO on a cloud instance