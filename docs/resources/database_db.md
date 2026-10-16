---
page_title: "Utho: utho_database_db"
subcategory: "Database (DBaaS)"
description: |-
  Create and manage logical databases inside a Utho managed database cluster.
---

# utho_database_db

Creates a logical database inside an existing Utho managed database cluster (`utho_database`). Use this to give each application or service its own database on a shared cluster.

Terraform waits for the cluster to finish provisioning before creating the database, so you can create the cluster and its databases in the same `terraform apply`.

## Example Usage

### Single database

```hcl
resource "utho_database" "main" {
  cluster_name  = "production-db"
  dcslug        = "inmumbaizone2"
  engine        = "pg"
  version       = "17"
  size          = "10151"
  network_type  = "private"
  billing       = "monthly"
  replica_count = "0"
}

resource "utho_database_db" "app" {
  cluster_id = utho_database.main.id
  name       = "app"
}
```

### One database per service

```hcl
locals {
  services = ["orders", "payments", "inventory"]
}

resource "utho_database_db" "service" {
  for_each   = toset(local.services)
  cluster_id = utho_database.main.id
  name       = each.key
}

resource "utho_database_user" "service" {
  for_each   = toset(local.services)
  cluster_id = utho_database.main.id
  name       = "${each.key}_user"
  password   = var.db_passwords[each.key]
}
```

## Argument Reference

| Argument     | Type   | Required | Description |
|--------------|--------|----------|-------------|
| `cluster_id` | String | Yes      | ID of the `utho_database` cluster to create the database in. Changing this forces a new resource. |
| `name`       | String | Yes      | Database name. Changing this forces a new resource. |

## Attribute Reference

| Attribute | Type   | Description |
|-----------|--------|-------------|
| `id`      | String | Identifier in the format `{cluster_id}:{name}`. |

## Notes

- Databases cannot be renamed or moved in place. Changing `name` or `cluster_id` destroys the database and creates a new one, which **deletes its data**. Use `lifecycle { prevent_destroy = true }` on production databases.
- Deleting the parent `utho_database` cluster removes all databases inside it.
- Import is not currently supported for this resource.
