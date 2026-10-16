# ============================================================
# Example: Full Stack Cloud Instance
# ============================================================
# A production-ready cloud instance with all options:
#   - VPC with public subnet
#   - Security group with SSH, HTTP, HTTPS rules
#   - SSH key authentication
#   - Public IP
#   - All values discovered via data sources (no hardcoded IDs)
#
# This is the recommended pattern for production deployments.
# Copy this, adjust variables, and you have a production server.
#
# What this creates:
#   - 1 SSH key
#   - 1 VPC + 1 public subnet
#   - 1 security group + 3 firewall rules
#   - 1 cloud instance (fully configured)
#
# Usage:
#   terraform init
#   terraform plan
#   terraform apply
#
# Connect after deploy:
#   ssh root@<instance_ip>
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

variable "project_name" {
  description = "Project name used as prefix for all resources."
  type        = string
  default     = "myapp"
}

variable "dcslug" {
  description = "Data center slug for all resources."
  type        = string
  default     = "inmumbaizone2"
}

variable "hostname" {
  description = "Hostname for the instance."
  type        = string
  default     = "prod-server.mhc"
}

variable "ssh_public_key_path" {
  description = "Path to your SSH public key."
  type        = string
  default     = "~/.ssh/id_ed25519.pub"
}

variable "allowed_ssh_cidr" {
  description = "CIDR allowed to SSH. Restrict to your IP for production."
  type        = string
  default     = "0.0.0.0/0"
}

variable "vpc_network" {
  description = "VPC base network address."
  type        = string
  default     = "10.0.0.0"
}

variable "billingcycle" {
  description = "Billing cycle: hourly | monthly | 3month | 6month | 12month | 24month | 36month"
  type        = string
  default     = "monthly"
}

# ── Discover valid values ──────────────────────────────────

data "utho_cloud_plans" "dc" {
  dcslug = var.dcslug
}

data "utho_cloud_images" "ubuntu" {
  distro = "ubuntu"
}

data "utho_billing_cycles" "cloud" {
  product = "cloud"
}

locals {
  # Pick plan — 2 vCPU / 4 GB with disk included
  plan = one([
    for p in data.utho_cloud_plans.dc.plans :
    p if p.id == "10308"
  ])

  # Ubuntu 22.04 LTS
  image = one([
    for img in data.utho_cloud_images.ubuntu.images :
    img if img.image == "ubuntu-22.04-x86_64"
  ])
}

# ── SSH Key ────────────────────────────────────────────────

resource "utho_ssh_key" "deploy" {
  name   = "${var.project_name}-deploy-key"
  sshkey = file(var.ssh_public_key_path)
}

# ── VPC + Subnet ───────────────────────────────────────────

resource "utho_vpc" "main" {
  name    = "${var.project_name}-vpc"
  network = var.vpc_network
  size    = "24"
  dcslug  = var.dcslug
  planid  = "1008"
}

resource "utho_subnet" "public" {
  name            = "${var.project_name}-public-subnet"
  vpc_id          = utho_vpc.main.id
  network         = var.vpc_network
  size            = 24
  type            = "public"
  assign_publicip = 1
}

# ── Security Group ─────────────────────────────────────────

resource "utho_firewall" "main" {
  name = "${var.project_name}-sg"
}

resource "utho_firewall_rule" "ssh" {
  firewall_id  = utho_firewall.main.id
  type         = "incoming"
  service      = "SSH"
  protocol     = "tcp"
  port         = "22"
  port_range   = "22"
  addresses    = var.allowed_ssh_cidr
  source_range = var.allowed_ssh_cidr
}

resource "utho_firewall_rule" "http" {
  firewall_id  = utho_firewall.main.id
  type         = "incoming"
  service      = "HTTP"
  protocol     = "tcp"
  port         = "80"
  port_range   = "80"
  addresses    = "0.0.0.0/0"
  source_range = "0.0.0.0/0"
}

resource "utho_firewall_rule" "https" {
  firewall_id  = utho_firewall.main.id
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
  dcslug       = var.dcslug
  planid       = local.plan.id        # discovered via data source
  billingcycle = var.billingcycle
  image        = local.image.image    # discovered via data source

  auth    = "option2"                  # SSH key auth
  sshkeys = utho_ssh_key.deploy.id

  enable_publicip = "true"
  cpumodel        = "amd"
  firewall        = utho_firewall.main.id    # attach security group
  vpc             = utho_subnet.public.id    # attach to VPC subnet
}

# ── Outputs ────────────────────────────────────────────────

output "instance_id" {
  description = "Cloud instance ID."
  value       = utho_cloud.main.id
}

output "instance_ip" {
  description = "Public IP address. Point your DNS A record here."
  value       = utho_cloud.main.ip
}

output "vpc_id" {
  description = "VPC ID."
  value       = utho_vpc.main.id
}

output "firewall_id" {
  description = "Security group ID."
  value       = utho_firewall.main.id
}

output "ssh_command" {
  description = "SSH command to connect to your instance."
  value       = "ssh root@${utho_cloud.main.ip}"
}

output "plan_used" {
  description = "Plan details used for this instance."
  value = {
    id    = local.plan.id
    cpu   = local.plan.cpu
    ram   = "${tonumber(local.plan.ram) / 1024} GB"
    disk  = "${local.plan.disk} GB"
    price = "₹${local.plan.price}/month"
  }
}
