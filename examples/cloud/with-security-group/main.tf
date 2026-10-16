# ============================================================
# Example: Cloud Instance with Security Group
# ============================================================
# Deploy a cloud instance with firewall rules controlling
# inbound and outbound traffic.
#
# What this creates:
#   - 1 SSH key
#   - 1 security group (firewall)
#   - 3 firewall rules (SSH, HTTP, HTTPS)
#   - 1 cloud instance with the security group attached
#
# Firewall rules:
#   - Port 22  (SSH)   — inbound from anywhere
#   - Port 80  (HTTP)  — inbound from anywhere
#   - Port 443 (HTTPS) — inbound from anywhere
#
# Usage:
#   terraform init
#   terraform plan
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

variable "hostname" {
  description = "Hostname for the instance."
  type        = string
  default     = "web-server.mhc"
}

variable "allowed_ssh_cidr" {
  description = "CIDR block allowed to SSH into the instance. Restrict to your IP for better security."
  type        = string
  default     = "0.0.0.0/0"   # change to your IP: "1.2.3.4/32"
}

# ── SSH Key ────────────────────────────────────────────────

resource "utho_ssh_key" "deploy" {
  name   = "web-server-key"
  sshkey = file("~/.ssh/id_ed25519.pub")
}

# ── Security Group ─────────────────────────────────────────
# Create a named security group and attach rules to it.
# The group can be reused across multiple instances.

resource "utho_firewall" "web" {
  name = "web-server-sg"
}

# Allow SSH access
resource "utho_firewall_rule" "ssh" {
  firewall_id  = utho_firewall.web.id
  type         = "incoming"
  service      = "SSH"
  protocol     = "tcp"
  port         = "22"
  port_range   = "22"
  addresses    = var.allowed_ssh_cidr
  source_range = var.allowed_ssh_cidr
}

# Allow HTTP traffic
resource "utho_firewall_rule" "http" {
  firewall_id  = utho_firewall.web.id
  type         = "incoming"
  service      = "HTTP"
  protocol     = "tcp"
  port         = "80"
  port_range   = "80"
  addresses    = "0.0.0.0/0"
  source_range = "0.0.0.0/0"
}

# Allow HTTPS traffic
resource "utho_firewall_rule" "https" {
  firewall_id  = utho_firewall.web.id
  type         = "incoming"
  service      = "HTTPS"
  protocol     = "tcp"
  port         = "443"
  port_range   = "443"
  addresses    = "0.0.0.0/0"
  source_range = "0.0.0.0/0"
}

# ── Instance ───────────────────────────────────────────────

resource "utho_cloud" "main" {
  hostname     = var.hostname
  dcslug       = "inmumbaizone2"
  planid       = "10308"
  billingcycle = "hourly"
  image        = "ubuntu-22.04-x86_64"

  auth    = "option2"
  sshkeys = utho_ssh_key.deploy.id

  enable_publicip = "true"
  cpumodel        = "amd"
  firewall        = utho_firewall.web.id   # attach security group at creation
}

# ── Outputs ────────────────────────────────────────────────

output "instance_id" {
  description = "Cloud instance ID."
  value       = utho_cloud.main.id
}

output "instance_ip" {
  description = "Public IP address."
  value       = utho_cloud.main.ip
}

output "firewall_id" {
  description = "Security group ID. Reuse this to attach to more instances."
  value       = utho_firewall.web.id
}

output "ssh_command" {
  description = "SSH command to connect."
  value       = "ssh root@${utho_cloud.main.ip}"
}

output "http_url" {
  description = "HTTP URL of the instance."
  value       = "http://${utho_cloud.main.ip}"
}
