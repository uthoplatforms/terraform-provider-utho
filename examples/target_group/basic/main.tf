# ============================================================
# Example: Basic Target Group
# ============================================================

terraform {
  required_providers {
    utho = {
      source  = "dharmajputho/utho-dev"
      version = ">= 0.2.63"
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

# ── Target Group ───────────────────────────────────────────
resource "utho_target_group" "api" {
  name                  = "api-targets"
  protocol              = "HTTP"
  port                  = "80"
  health_check_path     = "/health"
  health_check_protocol = "HTTP"
  health_check_interval = "30"
  health_check_timeout  = "5"
  healthy_threshold     = "3"
  unhealthy_threshold   = "3"
}

# ── Cloud instance target ──────────────────────────────────
resource "utho_target_group_target" "web" {
  target_group_id = utho_target_group.api.id
  type            = "cloud"
  backend_port    = "80"
  cloud_id        = var.cloud_id
  ip              = var.cloud_ip
}

variable "cloud_id"  { type = string }
variable "cloud_ip"  { type = string }

output "target_group_id" { value = utho_target_group.api.id }
output "target_id"       { value = utho_target_group_target.web.id }
