# ============================================================
# Example: Multi-Environment Projects
# ============================================================
# Creates dev, staging and production projects with
# shared team members.
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

variable "dev_user_id" {
  description = "Developer IAM user ID."
  type        = number
}

variable "devops_user_id" {
  description = "DevOps IAM user ID."
  type        = number
}

# ── Projects ───────────────────────────────────────────────
resource "utho_project" "dev" {
  name        = "my-app-dev"
  description = "Development environment"
  environment = "development"
}

resource "utho_project" "staging" {
  name        = "my-app-staging"
  description = "Staging environment"
  environment = "staging"
}

resource "utho_project" "prod" {
  name        = "my-app-prod"
  description = "Production environment"
  environment = "production"
}

# ── Members — dev has access to all, devops to staging+prod ─
resource "utho_project_member" "dev_in_dev" {
  project_id = utho_project.dev.id
  user_id    = var.dev_user_id
}

resource "utho_project_member" "dev_in_staging" {
  project_id = utho_project.staging.id
  user_id    = var.dev_user_id
}

resource "utho_project_member" "devops_in_staging" {
  project_id = utho_project.staging.id
  user_id    = var.devops_user_id
}

resource "utho_project_member" "devops_in_prod" {
  project_id = utho_project.prod.id
  user_id    = var.devops_user_id
}

# ── Outputs ────────────────────────────────────────────────
output "dev_project_id"     { value = utho_project.dev.id }
output "staging_project_id" { value = utho_project.staging.id }
output "prod_project_id"    { value = utho_project.prod.id }
