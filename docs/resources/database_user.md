---
page_title: "Utho: utho_database_user"
subcategory: "Database (DBaaS)"
description: |-
  Create and manage users inside a Utho database cluster.
---

# utho_database_user

Creates a database user inside an existing Utho database cluster.

## Example Usage

```hcl
resource "utho_database_user" "app" {
  cluster_id = utho_database.main.id
  name       = "appuser"
  password   = var.db_password
}

output "db_username" { value = utho_database_user.app.name }
```

## Argument Reference

| Argument     | Type   | Required | Description |
|--------------|--------|----------|-------------|
| `cluster_id` | String | Yes      | Database cluster ID. Changing this forces a new resource. |
| `name`       | String | Yes      | Username. Changing this forces a new resource. |
| `password`   | String | Yes      | User password. **Sensitive.** |

## Attribute Reference

| Attribute            | Type   | Description |
|----------------------|--------|-------------|
| `id`                 | String | Identifier in the format `{cluster_id}:{name}`. |
| `generated_password` | String | Generated password if none provided. **Sensitive.** |
