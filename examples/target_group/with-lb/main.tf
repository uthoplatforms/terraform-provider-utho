# ============================================================
# Example: Target Group with Load Balancer
# ============================================================
# Creates a target group with cloud instance targets and
# attaches it to a load balancer frontend.
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

variable "cloud_id"  { description = "Cloud instance ID to add as target." type = string }
variable "cloud_ip"  { description = "Cloud instance public IP." type = string }
variable "subnet_id" { description = "Subnet UUID for the load balancer." type = string }

# ── Target Group ───────────────────────────────────────────
resource "utho_target_group" "web" {
  name                  = "web-targets"
  protocol              = "HTTP"
  port                  = "80"
  health_check_path     = "/health"
  health_check_protocol = "HTTP"
  health_check_interval = "30"
  health_check_timeout  = "5"
  healthy_threshold     = "3"
  unhealthy_threshold   = "3"
}

# ── Target ─────────────────────────────────────────────────
resource "utho_target_group_target" "web" {
  target_group_id = utho_target_group.web.id
  type            = "cloud"
  backend_port    = "80"
  cloud_id        = var.cloud_id
  ip              = var.cloud_ip
}

# ── Load Balancer ──────────────────────────────────────────
resource "utho_loadbalancer" "main" {
  name            = "prod-lb"
  dcslug          = "innoida"
  type            = "application"
  enable_publicip = "true"
  vpc             = var.subnet_id
}

# ── Frontend ───────────────────────────────────────────────
resource "utho_loadbalancer_frontend" "http" {
  loadbalancer_id  = utho_loadbalancer.main.id
  name             = "http-frontend"
  proto            = "http"
  port             = "80"
  algorithm        = "roundrobin"
  target_group_id  = utho_target_group.web.id
  redirect_https   = "0"
  cookie           = "0"
}

# ── Outputs ────────────────────────────────────────────────
output "lb_ip"           { value = utho_loadbalancer.main.ip }
output "target_group_id" { value = utho_target_group.web.id }
output "target_id"       { value = utho_target_group_target.web.id }
