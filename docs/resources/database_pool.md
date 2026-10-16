---
page_title: "Utho: utho_database_pool"
subcategory: "Database (DBaaS)"
description: |-
  Create and manage connection pools for a Utho database cluster.
---

# utho_database_pool

Creates a connection pool for a Utho database cluster. Connection pooling reduces the overhead of opening and closing database connections, improving performance for high-traffic applications.

## Example Usage

```hcl
resource "utho_database_pool" "app" {
  cluster_id = utho_database.main.id
  cloud_id   = utho_database.main.cloud_id
  name       = "app-pool"
  db         = utho_database_db.app.name
  user       = utho_database_user.app.name
  mode       = "transaction"
  size       = 10
}
```

## Argument Reference

| Argument     | Type   | Required | Description |
|--------------|--------|----------|-------------|
| `cluster_id` | String | Yes      | Database cluster ID. Changing this forces a new resource. |
| `cloud_id`   | String | Yes      | Primary node cloud ID (`utho_database.name.cloud_id`). Changing this forces a new resource. |
| `name`       | String | Yes      | Pool name. Changing this forces a new resource. |
| `db`         | String | Yes      | Database name to pool connections for. |
| `user`       | String | Yes      | Database user for the pool. |
| `mode`       | String | Yes      | Pooling mode: `transaction`, `session`, or `statement`. |
| `size`       | Number | Yes      | Maximum number of connections in the pool. |

## Attribute Reference

| Attribute | Type   | Description |
|-----------|--------|-------------|
| `id`      | String | Connection pool ID. |
