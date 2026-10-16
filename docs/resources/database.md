---
page_title: "Utho: utho_database"
subcategory: "Database (DBaaS)"
description: |-
  Create and manage Utho managed database clusters.
---

# utho_database

Creates and manages a fully managed database cluster on Utho. Supports MySQL, PostgreSQL, and MariaDB. Utho handles provisioning, backups, failover, and maintenance. You connect using standard database clients with the provided connection strings.

~> **Note:** Database clusters take 4-10 minutes to provision. Terraform polls until the cluster is Active and all connection strings are available.

-> **Required:** The `billing` argument is required and must be set to `hourly` or `monthly`. Omitting it will cause the plan to fail immediately.

## Example Usage

### MySQL cluster

```hcl
resource "utho_database" "mysql" {
  cluster_name  = "production-mysql"
  dcslug        = "inmumbaizone2"
  engine        = "mysql"
  version       = "8.0"
  size          = "10151"
  network_type  = "public"
  billing       = "hourly"
  replica_count = "0"
}

output "mysql_host" { value = utho_database.mysql.host }
output "mysql_port" { value = utho_database.mysql.port }
output "mysql_user" { value = utho_database.mysql.default_user }
output "mysql_uri" {
  value     = utho_database.mysql.uri
  sensitive = true
}
```

### PostgreSQL cluster

```hcl
resource "utho_database" "postgres" {
  cluster_name  = "production-postgres"
  dcslug        = "inmumbaizone2"
  engine        = "pg"
  version       = "17"
  size          = "10151"
  network_type  = "public"
  billing       = "hourly"
  replica_count = "0"
  pitr_enabled  = "1"
}
```

### MariaDB cluster

```hcl
resource "utho_database" "mariadb" {
  cluster_name  = "production-mariadb"
  dcslug        = "inmumbaizone2"
  engine        = "mariadb"
  version       = "10.11"
  size          = "10151"
  network_type  = "public"
  billing       = "hourly"
  replica_count = "0"
}
```

### Cluster with security group

```hcl
resource "utho_database" "secured" {
  cluster_name  = "production-db"
  dcslug        = "inmumbaizone2"
  engine        = "mysql"
  version       = "8.0"
  size          = "10151"
  network_type  = "public"
  firewall      = utho_firewall.db.id
  billing       = "hourly"
  replica_count = "0"
}
```

### Connection and TLS guidance

~> **Note:** Utho managed databases do not support SSL/TLS connections. Use `sslmode=disable` (PostgreSQL) or omit SSL flags (MySQL/MariaDB) when connecting. Using `sslmode=require` will result in `server does not support SSL` errors.

Connect to PostgreSQL without SSL:

```bash
psql "host=<host> port=5432 dbname=defaultdb user=dbadmin sslmode=disable"
```

Connect to MySQL without SSL:

```bash
mysql -h <host> -P 3306 -u dbadmin -p defaultdb
```

### Use connection string in app server

```hcl
resource "utho_database" "main" {
  cluster_name  = "app-db"
  dcslug        = "inmumbaizone2"
  engine        = "mysql"
  version       = "8.0"
  size          = "10151"
  network_type  = "public"
  billing       = "hourly"
  replica_count = "0"
}

resource "utho_cloud" "app" {
  hostname = "app-server.mhc"
  dcslug   = "inmumbaizone2"
  # ...
}

output "connect_cmd" {
  value = "mysql -h ${utho_database.main.host} -P ${utho_database.main.port} -u ${utho_database.main.default_user} -p"
}
```

## Argument Reference

| Argument        | Type   | Required | Description |
|-----------------|--------|----------|-------------|
| `cluster_name`  | String | Yes      | Unique cluster name. Changing this forces a new resource. |
| `dcslug`        | String | Yes      | Data center slug. Changing this forces a new resource. |
| `engine`        | String | Yes      | Database engine: `mysql`, `pg`, `mariadb`. Changing this forces a new resource. |
| `version`       | String | Yes      | Engine version (e.g. `8.0`, `17`, `10.11`). Changing this forces a new resource. |
| `size`          | String | Yes      | Plan ID. Use `data.utho_database_plans` to list available plans. |
| `network_type`  | String | Yes      | `public` or `private`. Changing this forces a new resource. |
| `billing`       | String | Yes      | Billing cycle: `hourly` or `monthly`. |
| `replica_count` | String | Yes      | Number of replica nodes. Use `"0"` for standalone. |
| `firewall`      | String | No       | Security group ID to restrict access. |
| `vpc`           | String | No       | VPC subnet ID for private network clusters. |
| `pitr_enabled`  | String | No       | Enable point-in-time recovery: `"1"` or `"0"`. |

## Attribute Reference

| Attribute       | Type   | Description |
|-----------------|--------|-------------|
| `id`            | String | Cluster ID. |
| `cloud_id`      | String | Primary node cloud ID. |
| `status`        | String | Cluster status (`Active`, `Pending`). |
| `default_user`  | String | Default admin username. |
| `default_pass`  | String | Default admin password. **Sensitive.** |
| `default_dbname`| String | Default database name. |
| `port`          | String | Database port. |
| `host`          | String | Public connection hostname. |
| `host_private`  | String | Private connection hostname. |
| `uri`           | String | Public connection URI. **Sensitive.** |
| `uri_private`   | String | Private connection URI. **Sensitive.** |
| `created_at`    | String | Creation timestamp. |

## Supported Engines and Versions

| Engine    | `engine` value | Versions |
|-----------|---------------|---------|
| MySQL     | `mysql`       | `8.0` |
| PostgreSQL| `pg`          | `16`, `17` |
| MariaDB   | `mariadb`     | `10.11` |

## Related Resources

| Resource | Purpose |
|----------|---------|
| [utho_database_db](database_db) | Create databases inside the cluster |
| [utho_database_user](database_user) | Create database users |
| [utho_database_pool](database_pool) | Create connection pools |
| [data.utho_database_plans](../data-sources/database_plans) | List available plans |