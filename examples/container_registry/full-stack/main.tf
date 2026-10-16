# ============================================================
# Example: Full Container Registry Stack
# ============================================================
# Registry + CI robot account + webhook + immutable tag rule
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

variable "webhook_url" {
  description = "Webhook URL for registry push events."
  type        = string
}

# ── Registry ───────────────────────────────────────────────
resource "utho_container_registry" "main" {
  project_name = "prod-registry"
  dcslug       = "innoida"
  planid       = 10276
  billingcycle = "monthly"
  public       = false
}

# ── Robot account for CI/CD ────────────────────────────────
resource "utho_container_registry_robot" "ci" {
  project_name = utho_container_registry.main.project_name
  name         = "ci-pipeline"
  description  = "CI/CD robot for image push"
  duration     = 90

  access = [
    {
      resource = "repository"
      action   = "push"
    },
    {
      resource = "repository"
      action   = "pull"
    }
  ]
}

# ── Webhook on push ────────────────────────────────────────
resource "utho_container_registry_webhook" "deploy" {
  project_name     = utho_container_registry.main.project_name
  name             = "deploy-trigger"
  address          = var.webhook_url
  event_types      = ["PUSH_ARTIFACT"]
  payload_format   = "Default"
  skip_cert_verify = false
}

# ── Immutable rule for version tags ───────────────────────
resource "utho_container_registry_immutable_rule" "versions" {
  project_name = utho_container_registry.main.project_name
  tag_pattern  = "v*"
  repo_pattern = "**"
}

# ── Outputs ────────────────────────────────────────────────
output "registry_id"   { value = utho_container_registry.main.id }
output "robot_id"      { value = utho_container_registry_robot.ci.id }
output "webhook_id"    { value = utho_container_registry_webhook.deploy.id }
output "robot_secret" {
  value     = utho_container_registry_robot.ci.secret
  sensitive = true
}
