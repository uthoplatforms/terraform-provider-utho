---
page_title: "Utho: utho_cloud_snapshot"
subcategory: "Cloud Instances"
description: |-
  Create and manage snapshots of Utho Cloud instances.
---

# utho_cloud_snapshot

Creates and manages a snapshot of a Utho Cloud instance. Snapshots capture the full disk state of an instance at a point in time. Use them for backups, golden images, or deploying identical copies of a configured server.

Once created, use `snapshot_id` attribute to deploy new instances from it via `snapshotid` in `utho_cloud`.

## Example Usage

### Create a snapshot

```hcl
resource "utho_cloud_snapshot" "backup" {
  cloud_id = utho_cloud.web.id
  name     = "web-backup-before-upgrade"
}

output "snapshot_id" {
  value = utho_cloud_snapshot.backup.snapshot_id
}
```

### Create golden image snapshot and deploy from it

```hcl
resource "utho_cloud_snapshot" "golden" {
  cloud_id = utho_cloud.base.id
  name     = "golden-image-v1"
}

resource "utho_cloud" "worker" {
  count = 5

  hostname        = "worker-${count.index + 1}.mhc"
  dcslug          = "inmumbaizone2"
  planid          = "10308"
  billingcycle    = "hourly"
  auth            = "option1"
  root_password   = var.root_password
  snapshotid      = utho_cloud_snapshot.golden.snapshot_id
  enable_publicip = "true"
  cpumodel        = "amd"

  depends_on = [utho_cloud_snapshot.golden]
}
```

### Use existing snapshot from data source

```hcl
data "utho_cloud_snapshots" "all" {}

locals {
  latest = one([
    for s in data.utho_cloud_snapshots.all.snapshots :
    s if s.name == "golden-image-v1" && s.status == "Active"
  ])
}

resource "utho_cloud" "restored" {
  hostname        = "restored-01.mhc"
  dcslug          = "inmumbaizone2"
  planid          = "10308"
  billingcycle    = "hourly"
  auth            = "option1"
  root_password   = var.root_password
  snapshotid      = local.latest.id
  enable_publicip = "true"
  cpumodel        = "amd"
}
```

### Automated backup before risky operations

```hcl
resource "utho_cloud_snapshot" "pre_upgrade" {
  cloud_id = utho_cloud.app.id
  name     = "pre-upgrade-${formatdate("YYYY-MM-DD", timestamp())}"
}
```

## Argument Reference

| Argument    | Type   | Required | Description |
|-------------|--------|----------|-------------|
| `cloud_id`  | String | Yes      | Cloud instance ID to snapshot. Changing this forces a new resource. |
| `name`      | String | Yes      | Name for the snapshot. |

## Attribute Reference

| Attribute     | Type   | Description |
|---------------|--------|-------------|
| `id`          | String | Internal resource ID (`cloud_id:name`). |
| `snapshot_id` | String | Snapshot ID. Use this as `snapshotid` in `utho_cloud` to deploy from this snapshot. |

## Notes

- Snapshot creation takes 2–10 minutes depending on disk size.
- Use `data.utho_cloud_snapshots` to list all your snapshots and their IDs.
- Snapshots are stored in the same data center as the source instance.
- Destroying this resource deletes the snapshot permanently.
- Take snapshots before major changes: OS upgrades, configuration changes, or resizing.
