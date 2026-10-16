---
page_title: "EBS Attachment - Utho"
subcategory: "Elastic Block Storage"
description: |-
  Attach a Utho EBS volume to a cloud instance.
---

# utho_ebs_attachment

Attaches an EBS volume to a cloud instance. Destroying this resource detaches the volume without deleting it.

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

resource "utho_cloud" "web" {
  # ... cloud instance config
}

resource "utho_ebs_attachment" "data" {
  ebs_id   = utho_ebs.data.id
  cloud_id = utho_cloud.web.id
}
```

## Argument Reference

| Argument   | Type   | Required | Description |
|------------|--------|----------|-------------|
| `ebs_id`   | String | Yes      | EBS volume ID. Changing forces new resource. |
| `cloud_id` | String | Yes      | Cloud instance ID to attach the volume to. Changing forces new resource. |

## Attribute Reference

| Attribute | Type   | Description |
|-----------|--------|-------------|
| `id`      | String | Composite ID: `{ebs_id}:{cloud_id}`. |
| `device`  | String | Device name assigned by the instance (e.g. `vdb`, `vdc`). |

## Notes

- An EBS volume can only be attached to one instance at a time.
- Destroying this resource detaches the volume — the volume itself is preserved.
- Both volume and instance must be in the same data center.