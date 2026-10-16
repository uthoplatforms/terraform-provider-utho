# ============================================================
# Example: Basic Project with Member
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

variable "member_user_id" {
  description = "IAM sub-user ID to add as project member."
  type        = number
}

resource "utho_project" "main" {
  name        = "my-app"
  description = "Main application project"
  environment = "production"
}

resource "utho_project_member" "dev" {
  project_id = utho_project.main.id
  user_id    = var.member_user_id
}

output "project_id"     { value = utho_project.main.id }
output "project_status" { value = utho_project.main.status }
