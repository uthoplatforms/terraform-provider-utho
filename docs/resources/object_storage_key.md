---
page_title: "Utho: utho_object_storage_key"
subcategory: "Object Storage"
description: |-
  Create and manage access keys for Utho Object Storage.
---

# utho_object_storage_key

Creates and manages an S3-compatible access key for Utho Object Storage. Use access keys to give services scoped access to buckets — instead of using the bucket's root credentials.

~> **Important:** The `secret_key` is only available immediately after creation. Save it to a secrets manager right away — it cannot be retrieved again.

## Example Usage

### Create an access key

```hcl
resource "utho_object_storage_key" "app" {
  dcslug = "innoida"
  name   = "app-service-key"
}

output "access_key" {
  value     = utho_object_storage_key.app.access_key
  sensitive = true
}

output "secret_key" {
  value     = utho_object_storage_key.app.secret_key
  sensitive = true
}
```

### Create multiple keys for different services

```hcl
resource "utho_object_storage_key" "app" {
  dcslug = "innoida"
  name   = "app-key"
}

resource "utho_object_storage_key" "backup" {
  dcslug = "innoida"
  name   = "backup-key"
}

resource "utho_object_storage_permission" "app_full" {
  dcslug      = "innoida"
  bucket_name = utho_object_storage.main.name
  access_key  = utho_object_storage_key.app.id
  permission  = "full"
}

resource "utho_object_storage_permission" "backup_write" {
  dcslug      = "innoida"
  bucket_name = utho_object_storage.main.name
  access_key  = utho_object_storage_key.backup.id
  permission  = "write"
}
```

### Disable an access key

```hcl
resource "utho_object_storage_key" "app" {
  dcslug = "innoida"
  name   = "app-key"
  status = "disable"
}
```

Change `status` from `disable` to `enable` and apply to re-enable — updates in place without recreating.

## Argument Reference

| Argument | Type   | Required | Description |
|----------|--------|----------|-------------|
| `dcslug` | String | Yes      | Data center slug. Changing this forces a new resource. |
| `name`   | String | Yes      | Access key name for identification. Changing this forces a new resource. |
| `status` | String | No       | Key status: `enable` (default) or `disable`. Can be updated in place. |

## Attribute Reference

| Attribute    | Type   | Description |
|--------------|--------|-------------|
| `id`         | String | The access key string — use this as `access_key` in `utho_object_storage_permission`. |
| `access_key` | String | Generated S3 access key. **Sensitive.** |
| `secret_key` | String | Generated S3 secret key. **Sensitive — only available at creation time.** |
| `status`     | String | Current key status: `enable` or `disable`. |

## Related Resources

| Resource | Purpose |
|----------|---------|
| [utho_object_storage](object_storage) | The bucket to grant access to |
| [utho_object_storage_permission](object_storage_permission) | Grant this key access to a bucket |

## Notes

- One key can be granted different permissions on different buckets via `utho_object_storage_permission`
- Disabling a key immediately revokes all S3 API access without deleting the key
- Delete the resource to permanently remove the key
