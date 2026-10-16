---
page_title: "Running a Managed Database"
subcategory: "Tutorials"
description: |-
  Create a Utho managed PostgreSQL, MySQL, or MariaDB cluster with databases, users, connection pooling, and private networking.
---

# Running a Managed Database

Utho managed databases take care of provisioning, patching, backups, and replication. This tutorial creates a private PostgreSQL cluster, adds a database and an application user, sets up connection pooling, and passes connection details to your application servers.

## Prerequisites

* Complete [Getting Started](getting-started).
* A VPC with a private subnet. See [Private Networking with VPC and NAT](private-networking).

## Choose a plan

Look up the available plans instead of guessing an ID:

```hcl
data "utho_database_plans" "all" {}

output "database_plans" {
  value = [
    for p in data.utho_database_plans.all.plans :
    "${p.id}: ${p.cpu} vCPU / ${p.ram} MB RAM / ${p.disk} GB"
  ]
}
```

Use the chosen plan's `id` as the `size` of the cluster.

## Create the cluster

```hcl
variable "dcslug" {
  type    = string
  default = "inmumbaizone2"
}

resource "utho_database" "main" {
  cluster_name  = "production-db"
  dcslug        = var.dcslug
  engine        = "pg"
  version       = "17"
  size          = "10157"
  network_type  = "private"
  vpc           = utho_subnet.private.id
  billing       = "monthly"
  replica_count = "1"
  pitr_enabled  = "1"

  lifecycle {
    prevent_destroy = true
  }
}
```

A few decisions to make here:

* **Engine and version**: `engine` is `pg`, `mysql`, or `mariadb`. Both `engine` and `version` force a new cluster when changed, so choose carefully.
* **Network**: `network_type = "private"` with a `vpc` subnet keeps the database off the public internet. Use `public` only for development.
* **High availability**: `replica_count = "1"` or more adds standby nodes.
* **Recovery**: `pitr_enabled = "1"` enables point-in-time recovery for PostgreSQL.
* **Protection**: `prevent_destroy` makes Terraform refuse any plan that would delete the cluster and its data.

## Add a database, user, and connection pool

Give each application its own database and credentials rather than using the default admin user:

```hcl
variable "app_db_password" {
  type      = string
  sensitive = true
}

resource "utho_database_db" "app" {
  cluster_id = utho_database.main.id
  name       = "app"
}

resource "utho_database_user" "app" {
  cluster_id = utho_database.main.id
  name       = "app_user"
  password   = var.app_db_password
}
```

For PostgreSQL, add a PgBouncer connection pool so many application processes can share a small number of server connections. See [`utho_database_pool`](../resources/database_pool) for the available arguments.

## Restrict access

Attach a security group to the cluster through its `firewall` argument, and allow the database port only from your application subnet. The rule builds the subnet's CIDR block from its `network` and `size` attributes:

```hcl
resource "utho_firewall" "db" {
  name = "db-sg"
}

resource "utho_firewall_rule" "postgres_from_app" {
  firewall_id  = utho_firewall.db.id
  type         = "incoming"
  service      = "CUSTOM"
  protocol     = "tcp"
  port         = "5432"
  port_range   = "5432"
  addresses    = "${utho_subnet.private.network}/${utho_subnet.private.size}"
  source_range = "${utho_subnet.private.network}/${utho_subnet.private.size}"
}
```

## Connect your application

The cluster exports its connection details as attributes. Prefer the private host and URI from inside the VPC:

```hcl
output "db_host" {
  value = utho_database.main.host_private
}

output "db_port" {
  value = utho_database.main.port
}

output "db_uri" {
  value     = utho_database.main.uri_private
  sensitive = true
}
```

To hand the connection string to a Kubernetes workload, store it in a secret:

```hcl
resource "kubernetes_secret" "db" {
  metadata {
    name      = "database"
    namespace = "app"
  }

  data = {
    DATABASE_URL = "postgres://${utho_database_user.app.name}:${var.app_db_password}@${utho_database.main.host_private}:${utho_database.main.port}/${utho_database_db.app.name}"
  }
}
```

!> **Warning:** Database passwords and connection URIs are stored in Terraform state. Use an encrypted remote backend and restrict who can read it.

## Next steps

* Alert on database health with [`utho_alert`](../resources/alert).
* See [`utho_database`](../resources/database) for all supported engines, versions, and attributes.
