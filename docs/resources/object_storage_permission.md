---
page_title: "Utho: utho_object_storage_permission"
subcategory: "Object Storage"
description: |-
  Grant permissions on a Utho Object Storage bucket to an access key.
---

# utho_object_storage_permission

Grants a specific permission level on an Object Storage bucket to an access key. Changing the `permission` value updates the access key's rights in place — no destroy/recreate needed. Destroying this resource revokes the permission.

~> **Note:** The `read` permission is not currently supported by the Utho platform. Use `write` or `full` instead.

## Example Usage

### Grant full access to app service

```hcl
resource "utho_object_storage" "data" {
  name         = "myapp-data"
  dcslug       = "innoida"
  access       = "private"
  billingcycle = "monthly"
}

resource "utho_object_storage_key" "app" {
  dcslug = "innoida"
  name   = "app-key"
}

resource "utho_object_storage_permission" "app" {
  dcslug      = "innoida"
  bucket_name = utho_object_storage.data.name
  access_key  = utho_object_storage_key.app.id
  permission  = "full"
}
```

### Grant write-only access to backup service

```hcl
resource "utho_object_storage_key" "backup" {
  dcslug = "innoida"
  name   = "backup-key"
}

resource "utho_object_storage_permission" "backup" {
  dcslug      = "innoida"
  bucket_name = utho_object_storage.data.name
  access_key  = utho_object_storage_key.backup.id
  permission  = "write"
}
```

### Multi-service bucket access

```hcl
resource "utho_object_storage" "shared" {
  name         = "shared-bucket"
  dcslug       = "innoida"
  access       = "private"
  billingcycle = "monthly"
}

resource "utho_object_storage_key" "app" {
  dcslug = "innoida"
  name   = "app-key"
}

resource "utho_object_storage_key" "backup" {
  dcslug = "innoida"
  name   = "backup-key"
}

# App gets full access
resource "utho_object_storage_permission" "app_full" {
  dcslug      = "innoida"
  bucket_name = utho_object_storage.shared.name
  access_key  = utho_object_storage_key.app.id
  permission  = "full"
}

# Backup service gets write-only
resource "utho_object_storage_permission" "backup_write" {
  dcslug      = "innoida"
  bucket_name = utho_object_storage.shared.name
  access_key  = utho_object_storage_key.backup.id
  permission  = "write"
}
```

### Update permission in place

```hcl
resource "utho_object_storage_permission" "app" {
  dcslug      = "innoida"
  bucket_name = utho_object_storage.data.name
  access_key  = utho_object_storage_key.app.id
  permission  = "write"  # change to "full" and apply — updates in place
}
```

## Argument Reference

| Argument      | Type   | Required | Description |
|---------------|--------|----------|-------------|
| `dcslug`      | String | Yes      | Data center slug. Changing this forces a new resource. |
| `bucket_name` | String | Yes      | Bucket name to grant permission on. Changing this forces a new resource. |
| `access_key`  | String | Yes      | Access key ID (`utho_object_storage_key.name.id`). Changing this forces a new resource. |
| `permission`  | String | Yes      | Permission level: `write`, `full`. Can be updated in place. |

## Attribute Reference

| Attribute | Type   | Description |
|-----------|--------|-------------|
| `id`      | String | Identifier in the format `{bucket_name}:{access_key}:{permission}`. |

## Permission Levels

| Value   | Description |
|---------|-------------|
| `write` | Write objects to the bucket. |
| `full`  | Full read and write access. |

## Related Resources

| Resource | Purpose |
|----------|---------|
| [utho_object_storage](object_storage) | The bucket to grant access to |
| [utho_object_storage_key](object_storage_key) | The access key to grant permissions to |

## Notes

- Remove this resource from config to revoke the permission — no need to update `permission` to `none`
- One key can have different permissions on different buckets — create one permission resource per bucket
- `read` permission is not currently supported by the Utho platform
