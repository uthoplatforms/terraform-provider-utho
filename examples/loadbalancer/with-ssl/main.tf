# ============================================================
# Example: Load Balancer with HTTPS + SSL Termination
# ============================================================
# HTTPS load balancer with SSL termination.
# - Port 80 redirects to HTTPS
# - Port 443 terminates SSL at the LB
# - Backends receive plain HTTP
#
# Prerequisites:
#   You need an SSL certificate. Either:
#   a) Upload one via utho_ssl_certificate (shown below)
#   b) Use an existing cert ID from data.utho_ssl_certificates
#
# Usage:
#   terraform init
#   terraform apply
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

variable "dcslug" {
  type    = string
  default = "inmumbaizone2"
}

variable "cert_pem_path" {
  description = "Path to your SSL certificate PEM file."
  type        = string
  default     = "cert.pem"
}

variable "key_pem_path" {
  description = "Path to your SSL private key PEM file."
  type        = string
  default     = "key.pem"
}

# ── SSL Certificate ────────────────────────────────────────
# Upload your SSL certificate to Utho.
# Alternatively, use data.utho_ssl_certificates to reference
# an existing certificate by name.

resource "utho_ssl_certificate" "main" {
  name            = "myapp-ssl-cert"
  type            = "Custom"
  certificate_key = file(var.cert_pem_path)
  private_key     = file(var.key_pem_path)
}

# To use an existing certificate instead:
#
# data "utho_ssl_certificates" "all" {}
# locals {
#   cert = one([
#     for c in data.utho_ssl_certificates.all.certificates :
#     c if c.name == "my-existing-cert" && c.state == "verified"
#   ])
# }
# Then use local.cert.id as certificate_id below.

# ── Network ────────────────────────────────────────────────

resource "utho_vpc" "main" {
  name    = "https-lb-vpc"
  network = "10.2.0.0"
  size    = "24"
  dcslug  = var.dcslug
  planid  = "1008"
}

resource "utho_subnet" "public" {
  name            = "https-lb-subnet"
  vpc_id          = utho_vpc.main.id
  network         = "10.2.0.0"
  size            = 24
  type            = "public"
  assign_publicip = 1
}

resource "utho_ssh_key" "deploy" {
  name   = "https-lb-key"
  sshkey = file("~/.ssh/id_ed25519.pub")
}

# ── Backend Instances ──────────────────────────────────────

resource "utho_cloud" "backend" {
  count = 2

  hostname        = "web-${count.index + 1}.mhc"
  dcslug          = var.dcslug
  planid          = "10308"
  billingcycle    = "hourly"
  auth            = "option2"
  sshkeys         = utho_ssh_key.deploy.id
  image           = "ubuntu-22.04-x86_64"
  enable_publicip = "true"
  cpumodel        = "amd"
}

# ── Load Balancer ──────────────────────────────────────────

resource "utho_loadbalancer" "main" {
  name            = "https-lb"
  type            = "application"
  dcslug          = var.dcslug
  vpc             = utho_subnet.public.id
  enable_publicip = "true"
}

# HTTP frontend — redirects to HTTPS
resource "utho_loadbalancer_frontend" "http" {
  loadbalancer_id = utho_loadbalancer.main.id
  name            = "http"
  algorithm       = "roundrobin"
  proto           = "http"
  port            = "80"
  cookie          = "0"
  redirecthttps   = "1"   # redirect all HTTP to HTTPS
  certificate_id  = "0"
}

# HTTPS frontend — SSL terminated here
resource "utho_loadbalancer_frontend" "https" {
  loadbalancer_id = utho_loadbalancer.main.id
  name            = "https"
  algorithm       = "roundrobin"
  proto           = "https"
  port            = "443"
  cookie          = "0"
  redirecthttps   = "0"
  certificate_id  = utho_ssl_certificate.main.id
}

# Backends on HTTPS frontend — receive plain HTTP from LB
resource "utho_loadbalancer_backend" "web" {
  count = 2

  loadbalancer_id = utho_loadbalancer.main.id
  frontend_id     = utho_loadbalancer_frontend.https.id
  backend_port    = "80"
  weight          = "1"
  type            = "cloud"
  cloudid         = utho_cloud.backend[count.index].id
}

# LB settings with HSTS enabled
resource "utho_loadbalancer_settings" "main" {
  loadbalancer_id        = utho_loadbalancer.main.id
  timeout_connect        = "5000"
  timeout_client         = "50000"
  timeout_server         = "50000"
  timeout_http_request   = "5000"
  timeout_http_keepalive = "60000"
  timeout_tunnel         = "3600000"
  max_connections        = "2000"
  http2                  = "1"
  compression            = "1"
  hsts                   = "1"   # enable HSTS for HTTPS
}

# ── Outputs ────────────────────────────────────────────────

output "lb_ip" {
  description = "Load balancer public IP."
  value       = utho_loadbalancer.main.ip
}

output "lb_dns" {
  description = "Load balancer DNS hostname. Set your DNS CNAME here."
  value       = utho_loadbalancer.main.dns
}

output "certificate_id" {
  description = "SSL certificate ID."
  value       = utho_ssl_certificate.main.id
}

output "https_url" {
  description = "Test HTTPS endpoint."
  value       = "https://${utho_loadbalancer.main.ip}"
}
