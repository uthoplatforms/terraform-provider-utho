# ============================================================
# Example: Basic Load Balancer (HTTP)
# ============================================================
# The simplest load balancer setup — HTTP traffic distributed
# across 2 cloud instances using round-robin algorithm.
#
# What this creates:
#   - 1 VPC + 1 subnet (required for LB)
#   - 1 SSH key
#   - 2 cloud instances (backends)
#   - 1 application load balancer
#   - 1 HTTP frontend (port 80)
#   - 2 backends
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
  description = "Your Utho API key. Get it from https://console.utho.com/api"
  type        = string
  sensitive   = true
}

variable "dcslug" {
  description = "Data center slug."
  type        = string
  default     = "inmumbaizone2"
}

variable "backend_count" {
  description = "Number of backend servers."
  type        = number
  default     = 2
}

# ── Network ────────────────────────────────────────────────
# LB requires a VPC subnet

resource "utho_vpc" "main" {
  name    = "lb-vpc"
  network = "10.0.0.0"
  size    = "24"
  dcslug  = var.dcslug
  planid  = "1008"
}

resource "utho_subnet" "public" {
  name            = "lb-subnet"
  vpc_id          = utho_vpc.main.id
  network         = "10.0.0.0"
  size            = 24
  type            = "public"
  assign_publicip = 1
}

# ── SSH Key ────────────────────────────────────────────────

resource "utho_ssh_key" "deploy" {
  name   = "lb-deploy-key"
  sshkey = file("~/.ssh/id_ed25519.pub")
}

# ── Backend Instances ──────────────────────────────────────

resource "utho_cloud" "backend" {
  count = var.backend_count

  hostname        = "web-${count.index + 1}.mhc"
  dcslug          = var.dcslug
  planid          = "10308"   # 2 vCPU / 4 GB
  billingcycle    = "hourly"
  auth            = "option2"
  sshkeys         = utho_ssh_key.deploy.id
  image           = "ubuntu-22.04-x86_64"
  enable_publicip = "true"
  cpumodel        = "amd"
}

# ── Load Balancer ──────────────────────────────────────────

resource "utho_loadbalancer" "main" {
  name            = "web-lb"
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

resource "utho_loadbalancer_backend" "web" {
  count = var.backend_count

  loadbalancer_id = utho_loadbalancer.main.id
  frontend_id     = utho_loadbalancer_frontend.http.id
  backend_port    = "80"
  weight          = "1"
  type            = "cloud"
  cloudid         = utho_cloud.backend[count.index].id
}

# ── Outputs ────────────────────────────────────────────────

output "lb_ip" {
  description = "Load balancer public IP. Point your DNS A record here."
  value       = utho_loadbalancer.main.ip
}

output "lb_dns" {
  description = "Load balancer DNS hostname. Use this for DNS CNAME records."
  value       = utho_loadbalancer.main.dns
}

output "lb_id" {
  description = "Load balancer ID."
  value       = utho_loadbalancer.main.id
}

output "backend_ips" {
  description = "Public IPs of backend servers."
  value       = utho_cloud.backend[*].ip
}

output "curl_test" {
  description = "Test the load balancer."
  value       = "curl http://${utho_loadbalancer.main.ip}"
}
