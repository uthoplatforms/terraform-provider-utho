# ============================================================
# Example: Multiple Cloud Instances (Horizontal Scaling)
# ============================================================
# Deploy multiple identical cloud instances in parallel using
# Terraform's count feature. All instances share the same
# SSH key and security group but get unique hostnames and IPs.
#
# Use cases:
#   - Web server farm behind a load balancer
#   - Worker nodes for a queue-based system
#   - Identical staging/production environments
#
# What this creates:
#   - 1 SSH key (shared across all instances)
#   - 1 security group + rules (shared)
#   - N cloud instances (default: 3)
#
# Scaling:
#   Change instance_count variable to add/remove instances.
#   Terraform will only create/destroy the delta.
#
# Usage:
#   terraform init
#   terraform plan
#   terraform apply
#   terraform apply -var="instance_count=5"  # scale up to 5
#   terraform apply -var="instance_count=2"  # scale down to 2
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

variable "instance_count" {
  description = "Number of instances to deploy. Change to scale up or down."
  type        = number
  default     = 3
}

variable "hostname_prefix" {
  description = "Prefix for instance hostnames. Instances will be named prefix-1.mhc, prefix-2.mhc, etc."
  type        = string
  default     = "worker"
}

variable "dcslug" {
  description = "Data center slug."
  type        = string
  default     = "inmumbaizone2"
}

variable "planid" {
  description = "Plan ID for all instances. Use data.utho_cloud_plans to find IDs."
  type        = string
  default     = "10308"   # 2 vCPU / 4 GB / 80 GB
}

# ── SSH Key ────────────────────────────────────────────────

resource "utho_ssh_key" "deploy" {
  name   = "multi-instance-key"
  sshkey = file("~/.ssh/id_ed25519.pub")
}

# ── Security Group ─────────────────────────────────────────
# One security group shared across all instances.

resource "utho_firewall" "workers" {
  name = "worker-sg"
}

resource "utho_firewall_rule" "ssh" {
  firewall_id  = utho_firewall.workers.id
  type         = "incoming"
  service      = "SSH"
  protocol     = "tcp"
  port         = "22"
  port_range   = "22"
  addresses    = "0.0.0.0/0"
  source_range = "0.0.0.0/0"
}

resource "utho_firewall_rule" "app" {
  firewall_id  = utho_firewall.workers.id
  type         = "incoming"
  service      = "Custom"
  protocol     = "tcp"
  port         = "8080"
  port_range   = "8080"
  addresses    = "0.0.0.0/0"
  source_range = "0.0.0.0/0"
}

# ── Instances ──────────────────────────────────────────────
# count creates N identical instances in parallel.
# Each gets a unique hostname using count.index.

resource "utho_cloud" "workers" {
  count = var.instance_count

  hostname     = "${var.hostname_prefix}-${count.index + 1}.mhc"
  dcslug       = var.dcslug
  planid       = var.planid
  billingcycle = "hourly"
  image        = "ubuntu-22.04-x86_64"

  auth    = "option2"
  sshkeys = utho_ssh_key.deploy.id

  enable_publicip = "true"
  cpumodel        = "amd"
  firewall        = utho_firewall.workers.id
}

# ── Outputs ────────────────────────────────────────────────

output "instance_ids" {
  description = "IDs of all deployed instances."
  value       = utho_cloud.workers[*].id
}

output "instance_ips" {
  description = "Public IPs of all deployed instances. Add these to your load balancer."
  value       = utho_cloud.workers[*].ip
}

output "instance_hostnames" {
  description = "Hostnames of all deployed instances."
  value       = utho_cloud.workers[*].hostname
}

output "ssh_commands" {
  description = "SSH commands to connect to each instance."
  value = [
    for instance in utho_cloud.workers :
    "ssh root@${instance.ip} # ${instance.hostname}"
  ]
}

output "instance_count" {
  description = "Number of instances deployed."
  value       = var.instance_count
}
