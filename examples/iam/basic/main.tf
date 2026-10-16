# ============================================================
# Example: Basic IAM User
# ============================================================
# Creates an IAM sub-user with read/write access to
# compute and database services.
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

resource "utho_iam_user" "dev" {
  fullname    = "Dev User"
  email       = "dev@mycompany.com"
  mobilecc    = "91"
  mobile      = "9999999999"
  permissions = "compute_read,compute_write,database_read,database_write,monitoring_read"
}

output "user_id"    { value = utho_iam_user.dev.id }
output "user_status" { value = utho_iam_user.dev.status }
