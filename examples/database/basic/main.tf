# ============================================================
# Example: Basic MySQL Database Cluster
# ============================================================
# Creates a managed MySQL cluster. Terraform waits until
# the cluster is Active before completing — all connection
# strings are available in outputs immediately after apply.
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

variable "cluster_name" {
  type    = string
  default = "my-mysql-cluster"
}

resource "utho_database" "main" {
  cluster_name  = var.cluster_name
  dcslug        = "inmumbaizone2"
  engine        = "mysql"
  version       = "8.0"
  size          = "10151"
  network_type  = "public"
  billing       = "hourly"
  replica_count = "0"
}

output "cluster_id"    { value = utho_database.main.id }
output "status"        { value = utho_database.main.status }
output "host"          { value = utho_database.main.host }
output "port"          { value = utho_database.main.port }
output "default_user"  { value = utho_database.main.default_user }
output "connect_cmd"   { value = "mysql -h ${utho_database.main.host} -P ${utho_database.main.port} -u ${utho_database.main.default_user} -p" }
output "password" {
  value     = utho_database.main.default_pass
  sensitive = true
}
output "uri" {
  value     = utho_database.main.uri
  sensitive = true
}
