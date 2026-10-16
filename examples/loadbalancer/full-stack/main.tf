# ============================================================
# Example: Full Production Load Balancer Stack
# ============================================================
# Complete production-ready setup:
#   - VPC + subnet
#   - Security group with SSH + HTTP rules
#   - 3 cloud instances as backends
#   - Application load balancer
#   - HTTP frontend with roundrobin
#   - 3 backends
#   - Optimized LB settings (http2, compression, timeouts)
#
# Everything deployed in a single terraform apply.
#
# Usage:
#   terraform init
#   terraform apply
#   curl http://<lb_ip>
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

# ── Variables ─────────────────────────────────────────────

variable "utho_api_key" {
  description = "Your Utho API key."
  type        = string
  sensitive   = true
}

variable "project_name" {
  description = "Project name used as prefix for all resources."
  type        = string
  default     = "myapp"
}

variable "dcslug" {
  description = "Data center slug."
  type        = string
  default     = "inmumbaizone2"
}

variable "backend_count" {
  description = "Number of backend instances. Change to scale."
  type        = number
  default     = 3
}

# ── Network ────────────────────────────────────────────────

resource "utho_vpc" "main" {
  name    = "${var.project_name}-vpc"
  network = "10.0.0.0"
  size    = "24"
  dcslug  = var.dcslug
  planid  = "1008"
}

resource "utho_subnet" "public" {
  name            = "${var.project_name}-subnet"
  vpc_id          = utho_vpc.main.id
  network         = "10.0.0.0"
  size            = 24
  type            = "public"
  assign_publicip = 1
}

# ── Security Group ─────────────────────────────────────────

resource "utho_firewall" "backend" {
  name = "${var.project_name}-backend-sg"
}

resource "utho_firewall_rule" "ssh" {
  firewall_id  = utho_firewall.backend.id
  type         = "incoming"
  service      = "SSH"
  protocol     = "tcp"
  port         = "22"
  port_range   = "22"
  addresses    = "0.0.0.0/0"
  source_range = "0.0.0.0/0"
}

resource "utho_firewall_rule" "http" {
  firewall_id  = utho_firewall.backend.id
  type         = "incoming"
  service      = "HTTP"
  protocol     = "tcp"
  port         = "80"
  port_range   = "80"
  addresses    = "0.0.0.0/0"
  source_range = "0.0.0.0/0"
}

# ── SSH Key ────────────────────────────────────────────────

resource "utho_ssh_key" "deploy" {
  name   = "${var.project_name}-deploy-key"
  sshkey = file("~/.ssh/id_ed25519.pub")
}

# ── Backend Instances ──────────────────────────────────────

resource "utho_cloud" "backend" {
  count = var.backend_count

  hostname        = "${var.project_name}-backend-${count.index + 1}.mhc"
  dcslug          = var.dcslug
  planid          = "10308"
  billingcycle    = "hourly"
  auth            = "option2"
  sshkeys         = utho_ssh_key.deploy.id
  image           = "ubuntu-22.04-x86_64"
  enable_publicip = "true"
  cpumodel        = "amd"
  firewall        = utho_firewall.backend.id
}

# ── Load Balancer ──────────────────────────────────────────

resource "utho_loadbalancer" "main" {
  name            = "${var.project_name}-lb"
  type            = "application"
  dcslug          = var.dcslug
  vpc             = utho_subnet.public.id
  enable_publicip = "true"
}

# ── Frontend ───────────────────────────────────────────────

resource "utho_loadbalancer_frontend" "http" {
  loadbalancer_id = utho_loadbalancer.main.id
  name            = "http"
  algorithm       = "roundrobin"
  proto           = "http"
  port            = "80"
  cookie          = "0"
  redirecthttps   = "0"
  certificate_id  = "0"
}

# ── Backends ───────────────────────────────────────────────

resource "utho_loadbalancer_backend" "app" {
  count = var.backend_count

  loadbalancer_id = utho_loadbalancer.main.id
  frontend_id     = utho_loadbalancer_frontend.http.id
  backend_port    = "80"
  weight          = "1"
  type            = "cloud"
  cloudid         = utho_cloud.backend[count.index].id
}

# ── LB Settings ────────────────────────────────────────────
# All timeouts in milliseconds

resource "utho_loadbalancer_settings" "main" {
  loadbalancer_id        = utho_loadbalancer.main.id
  timeout_connect        = "5000"     # 5 seconds
  timeout_client         = "50000"    # 50 seconds
  timeout_server         = "50000"    # 50 seconds
  timeout_http_request   = "5000"     # 5 seconds
  timeout_http_keepalive = "60000"    # 60 seconds
  timeout_tunnel         = "3600000"  # 1 hour
  max_connections        = "2000"
  http2                  = "1"
  compression            = "1"
  hsts                   = "0"
}

# ── Outputs ────────────────────────────────────────────────

output "lb_ip" {
  description = "Load balancer public IP."
  value       = utho_loadbalancer.main.ip
}

output "lb_dns" {
  description = "Load balancer DNS hostname."
  value       = utho_loadbalancer.main.dns
}

output "backend_ips" {
  description = "Backend server IPs."
  value       = utho_cloud.backend[*].ip
}

output "ssh_commands" {
  description = "SSH into each backend."
  value = [
    for instance in utho_cloud.backend :
    "ssh root@${instance.ip}"
  ]
}

output "curl_test" {
  description = "Test the load balancer."
  value       = "curl http://${utho_loadbalancer.main.ip}"
}
