# ============================================================
# Example: Multiple IAM Users with Different Roles
# ============================================================
# Creates three IAM users:
# - admin: full access to all services
# - developer: compute, database, kubernetes access
# - readonly: read-only access to all services
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

# ── Admin User ─────────────────────────────────────────────
resource "utho_iam_user" "admin" {
  fullname    = "Admin User"
  email       = "admin@mycompany.com"
  mobilecc    = "91"
  mobile      = "9999999999"
  permissions = "compute_read,compute_write,compute_delete,database_read,database_write,database_delete,objectstorage_read,objectstorage_write,objectstorage_delete,kubernetes_read,kubernetes_write,kubernetes_delete,monitoring_read,monitoring_write,monitoring_delete,loadbalancer_read,loadbalancer_write,loadbalancer_delete,dns_read,dns_write,dns_delete,vpc_read,vpc_write,vpc_delete,firewall_read,firewall_write,firewall_delete,autoscaling_read,autoscaling_write,autoscaling_delete,snapshot_read,snapshot_write,snapshot_delete,sshkey_read,sshkey_write,sshkey_delete,billing_read,api_read,api_write,api_delete"
}

# ── Developer User ─────────────────────────────────────────
resource "utho_iam_user" "developer" {
  fullname    = "Developer User"
  email       = "dev@mycompany.com"
  mobilecc    = "91"
  mobile      = "8888888888"
  permissions = "compute_read,compute_write,compute_delete,database_read,database_write,kubernetes_read,kubernetes_write,objectstorage_read,objectstorage_write,monitoring_read,snapshot_read,snapshot_write,sshkey_read,sshkey_write"
}

# ── Read-Only User ─────────────────────────────────────────
resource "utho_iam_user" "readonly" {
  fullname    = "Read Only User"
  email       = "readonly@mycompany.com"
  mobilecc    = "91"
  mobile      = "7777777777"
  permissions = "compute_read,database_read,objectstorage_read,kubernetes_read,monitoring_read,loadbalancer_read,dns_read,vpc_read,firewall_read,autoscaling_read,snapshot_read,billing_read"
}

output "admin_id"     { value = utho_iam_user.admin.id }
output "developer_id" { value = utho_iam_user.developer.id }
output "readonly_id"  { value = utho_iam_user.readonly.id }
