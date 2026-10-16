---
page_title: "Utho: utho_cloud_storage"
subcategory: "Cloud Instances"
description: |-
  Add general-purpose storage disks to a Utho Cloud instance.
---

# utho_cloud_storage

Adds a general-purpose storage disk to an existing cloud instance. Unlike EBS volumes (`utho_cloud_ebs`), general storage is directly attached to the instance's hardware and offers consistent local performance.

## Example Usage

### Add a storage disk

```hcl
resource "utho_cloud_storage" "extra" {
  cloud_id = utho_cloud.web.id
  size_gb  = 100
}
```

### Add large storage for media files

```hcl
resource "utho_cloud_storage" "media" {
  cloud_id = utho_cloud.app.id
  size_gb  = 500
}

output "storage_id" {
  value = utho_cloud_storage.media.id
}
```

## Argument Reference

| Argument   | Type   | Required | Description |
|------------|--------|----------|-------------|
| `cloud_id` | String | Yes      | Cloud instance ID. Changing this forces a new resource. |
| `size_gb`  | Number | Yes      | Disk size in GB. |

## Attribute Reference

| Attribute | Type   | Description |
|-----------|--------|-------------|
| `id`      | String | Unique storage disk ID. |
| `disk_id` | String | Internal disk ID used for management operations. |

## EBS vs General Storage

| | General Storage (`utho_cloud_storage`) | EBS (`utho_cloud_ebs`) |
|---|---|---|
| Type | Local disk | Network-attached block volume |
| Performance | Consistent local IOPS | High performance NVMe or SSD |
| Persistence | Tied to instance | Independent of instance |
| Use case | Extra disk space, media files | Databases, critical data |

## Notes

- General storage disks are mounted automatically. Use standard Linux tools to format and mount them.
- Destroying this resource removes the disk and all data on it permanently.
- For databases or critical data, use `utho_cloud_ebs` instead — EBS volumes can be detached and reattached independently of the instance.
