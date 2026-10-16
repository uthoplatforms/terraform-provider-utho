# ============================================================
# Example: Full Production Database Stack
# ============================================================
# Creates a complete production database setup:
#   - MySQL cluster
#   - Application database
#   - Application user
#   - Connection pool
# ============================================================

terraform {
  required_providers {
    utho = {
      source  = "nitinuthocloud/utho"
      version = "~> 0.7"
    }
  }
  required_version = ">= 1.0"
}

provider "utho" {
  api_key = var.utho_api_key
}

variable "utho_api_key" {
  type      = string
  sensitive = true
}

variable "project" {
  type    = string
  default = "myapp"
}

variable "db_password" {
  type      = string
  sensitive = true
}

# ── Database Cluster ───────────────────────────────────────

resource "utho_database" "main" {
  cluster_name  = "${var.project}-db"
  dcslug        = "inmumbaizone2"
  engine        = "mysql"
  version       = "8.0"
  size          = "10151"
  network_type  = "public"
  billing       = "hourly"
  replica_count = "0"
}

# ── Application Database ───────────────────────────────────

resource "utho_database_db" "app" {
  cluster_id = utho_database.main.id
  name       = "${var.project}_db"
}

# ── Application User ───────────────────────────────────────

resource "utho_database_user" "app" {
  cluster_id = utho_database.main.id
  name       = "${var.project}user"
  password   = var.db_password
}

# ── Connection Pool ────────────────────────────────────────

resource "utho_database_pool" "app" {
  cluster_id = utho_database.main.id
  cloud_id   = utho_database.main.cloud_id
  name       = "${var.project}-pool"
  db         = utho_database_db.app.name
  user       = utho_database_user.app.name
  mode       = "transaction"
  size       = 20
}

# ── Outputs ────────────────────────────────────────────────

output "cluster_id"   { value = utho_database.main.id }
output "host"         { value = utho_database.main.host }
output "port"         { value = utho_database.main.port }
output "database"     { value = utho_database_db.app.name }
output "username"     { value = utho_database_user.app.name }
output "pool_id"      { value = utho_database_pool.app.id }
output "connect_cmd"  { value = "mysql -h ${utho_database.main.host} -P ${utho_database.main.port} -u ${utho_database_user.app.name} -p ${utho_database_db.app.name}" }
output "uri" {
  value     = utho_database.main.uri
  sensitive = true
}
