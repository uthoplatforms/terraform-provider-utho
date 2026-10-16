---
page_title: "Utho: utho_cloud_ebs"
subcategory: "Cloud Instances"
description: |-
  Attach an existing EBS block volume to a Utho Cloud instance.
---

# utho_cloud_ebs

Attaches an existing EBS (Elastic Block Storage) volume to a cloud instance. Use this to connect pre-created EBS volumes to instances for additional persistent storage.

~> **Note:** This resource attaches an **existing** EBS volume. To create EBS volumes at instance creation time, use the `ebs` block inside `utho_cloud` instead. To create standalone EBS volumes, use the Utho Dashboard and note the volume ID.

## Example Usage

### Attach an existing EBS volume

```hcl
resource "utho_cloud_ebs" "attach" {
  cloud_id = utho_cloud.db.id
  ebs_id   = "your-ebs-volume-id"   # Get from Utho Dashboard
}
```

### Full setup — instance with EBS at creation

For attaching EBS at creation time, use the `ebs` block inside `utho_cloud` instead:

```hcl
resource "utho_cloud" "db" {
  hostname        = "db-01.mhc"
  dcslug          = "inmumbaizone2"
  planid          = "10355"   # EBS plan — no included disk
  billingcycle    = "monthly"
  auth            = "option2"
  sshkeys         = utho_ssh_key.deploy.id
  image           = "ubuntu-22.04-x86_64"
  enable_publicip = "false"
  cpumodel        = "amd"
  delete_ebs      = true

  ebs = [
    { id = "1", disk = 80, type = "nvme" },
    { id = "2", disk = 500, type = "nvme" },
  ]
}
```

## Argument Reference

| Argument   | Type   | Required | Description |
|------------|--------|----------|-------------|
| `cloud_id` | String | Yes      | Cloud instance ID. Changing this forces a new resource. |
| `ebs_id`   | String | Yes      | EBS volume ID to attach. Get this from the Utho Dashboard. |

## Attribute Reference

| Attribute | Type   | Description |
|-----------|--------|-------------|
| `id`      | String | Unique attachment ID. |

## Notes

- This resource attaches an existing EBS volume — it does not create one.
- To attach EBS volumes at instance creation, use the `ebs` block in `utho_cloud`.
- Destroying this resource detaches the volume from the instance but does NOT delete the volume itself.
- EBS volumes are DC-specific — the volume and instance must be in the same data center.
