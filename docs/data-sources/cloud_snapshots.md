---
page_title: "Utho: utho_cloud_snapshots"
subcategory: "Cloud Instances"
description: |-
  List available snapshots in your Utho account.
---

# utho_cloud_snapshots

Fetches all snapshots in your account. Use snapshot IDs to deploy new cloud instances from a snapshot — useful for golden images, disaster recovery, or cloning environments.

## Example Usage

### List all snapshots

```hcl
data "utho_cloud_snapshots" "all" {}

output "snapshots" {
  value = data.utho_cloud_snapshots.all.snapshots
}
```

### Deploy from a specific snapshot

```hcl
data "utho_cloud_snapshots" "all" {}

locals {
  # Find snapshot by name
  golden = one([
    for s in data.utho_cloud_snapshots.all.snapshots :
    s if s.name == "golden-image-v2"
  ])
}

resource "utho_cloud" "restored" {
  hostname        = "restored-01.mhc"
  dcslug          = "inmumbaizone2"
  planid          = "10308"
  billingcycle    = "hourly"
  auth            = "option1"
  root_password   = var.root_password
  snapshotid      = local.golden.id
  enable_publicip = "true"
}
```

### List active snapshots only

```hcl
data "utho_cloud_snapshots" "all" {}

output "active_snapshots" {
  value = [
    for s in data.utho_cloud_snapshots.all.snapshots :
    { id = s.id, name = s.name, size = s.size, created = s.create_date }
    if s.status == "Active"
  ]
}
```

## Attribute Reference

### Top-level

| Attribute   | Type | Description |
|-------------|------|-------------|
| `snapshots` | List | List of all snapshots in your account. |

### snapshots

| Attribute     | Type   | Description |
|---------------|--------|-------------|
| `id`          | String | Snapshot ID. Use this as `snapshotid` in `utho_cloud`. |
| `cloud_id`    | String | ID of the source cloud instance. |
| `name`        | String | Snapshot name. |
| `image`       | String | Base OS image the source instance was running. |
| `size`        | String | Snapshot size in GB. |
| `status`      | String | Snapshot status. Only use `Active` snapshots. |
| `storage`     | String | Storage backend: `general`, `ebs`, or `other`. |
| `create_date` | String | Timestamp when the snapshot was created. |
