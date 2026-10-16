---
page_title: "EBS Volume - Utho"
subcategory: "Elastic Block Storage"
description: |-
  Create and manage Utho Elastic Block Storage (EBS) volumes.
---

# utho_ebs

Creates and manages a Utho Elastic Block Storage (EBS) volume. EBS volumes are persistent block storage that can be attached to cloud instances using `utho_ebs_attachment`.

## Example Usage

```hcl
resource "utho_ebs" "data" {
  name       = "app-data-volume"
  dcslug     = "innoida"
  disk_type  = "SSD"
  disk       = "50"
  iops       = "250"
  throughput = "125"
}
```

## Argument Reference

| Argument     | Type   | Required | Description |
|--------------|--------|----------|-------------|
| `name`       | String | Yes      | Volume name. Can be updated in place. |
| `dcslug`     | String | Yes      | Data center slug (e.g. `innoida`). Changing forces new resource. |
| `disk_type`  | String | Yes      | Disk type: `SSD` or `HDD`. Changing forces new resource. |
| `disk`       | String | Yes      | Volume size in GB (10–10240). Can only be increased, never decreased. |
| `iops`       | String | Yes      | IOPS (1 to 5×disk size, e.g. max 50 for 10 GB). Can be updated in place. |
| `throughput` | String | Yes      | Throughput in MB/s (max 500). Can be updated in place. |

## Attribute Reference

| Attribute    | Type   | Description |
|--------------|--------|-------------|
| `id`         | String | EBS volume ID. |
| `status`     | String | Volume status (`Active`, `Inactive`). |
| `cloudid`    | String | Cloud instance ID the volume is currently attached to. |
| `created_at` | String | Creation timestamp. |

## Notes

- Volume size (`disk`) can only be increased — decreasing is not supported.
- IOPS maximum = 5 × disk size in GB.
- Throughput maximum = 500 MB/s.
- Use `utho_ebs_attachment` to attach/detach volumes from cloud instances.

## Related Resources

- [utho_ebs_attachment](ebs_attachment) — Attach this volume to a cloud instance