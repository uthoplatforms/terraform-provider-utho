# ============================================================
# Example: Cloud Instance inside a VPC
# ============================================================
# Deploy a cloud instance inside a private network (VPC).
# Ideal for production workloads that should not be directly
# exposed to the internet.
#
# What this creates:
#   - 1 VPC (private network)
#   - 1 public subnet inside the VPC
#   - 1 SSH key
#   - 1 cloud instance inside the VPC subnet
#
# Network layout:
#   VPC: 10.0.0.0/24
#   └── Public subnet: 10.0.0.0/24
#       └── Instance (gets both public + private IP)
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

variable "dcslug" {
  description = "Data center slug for all resources. All must be in the same DC."
  type        = string
  default     = "inmumbaizone2"
}

variable "vpc_name" {
  description = "Name for the VPC."
  type        = string
  default     = "production-vpc"
}

variable "hostname" {
  description = "Hostname for the instance."
  type        = string
  default     = "app-server.mhc"
}

# ── SSH Key ────────────────────────────────────────────────

resource "utho_ssh_key" "deploy" {
  name   = "vpc-example-key"
  sshkey = file("~/.ssh/id_ed25519.pub")
}

# ── VPC ────────────────────────────────────────────────────
# Create a private network for your infrastructure.
# All instances in this VPC can communicate privately.

resource "utho_vpc" "main" {
  name    = var.vpc_name
  network = "10.0.0.0"      # base network address
  size    = "24"             # /24 = 256 IPs
  dcslug  = var.dcslug
  planid  = "1008"           # standard VPC plan
}

# ── Subnet ─────────────────────────────────────────────────
# Create a public subnet inside the VPC.
# Instances in this subnet get a public IP + private IP.

resource "utho_subnet" "public" {
  name            = "public-subnet"
  vpc_id          = utho_vpc.main.id
  network         = "10.0.0.0"
  size            = 24
  type            = "public"
  assign_publicip = 1       # auto-assign public IPs to instances
}

# ── Instance ───────────────────────────────────────────────
# Deploy inside the VPC subnet.
# The instance gets both a public IP and a private IP.

resource "utho_cloud" "main" {
  hostname     = var.hostname
  dcslug       = var.dcslug    # must match VPC datacenter
  planid       = "10308"
  billingcycle = "hourly"
  image        = "ubuntu-22.04-x86_64"

  auth    = "option2"
  sshkeys = utho_ssh_key.deploy.id

  enable_publicip = "true"
  cpumodel        = "amd"
  vpc             = utho_subnet.public.id  # attach to VPC subnet at creation
}

# ── Outputs ────────────────────────────────────────────────

output "instance_id" {
  description = "Cloud instance ID."
  value       = utho_cloud.main.id
}

output "instance_ip" {
  description = "Public IP of the instance."
  value       = utho_cloud.main.ip
}

output "vpc_id" {
  description = "VPC ID."
  value       = utho_vpc.main.id
}

output "subnet_id" {
  description = "Subnet ID. Use this to attach more instances to the same subnet."
  value       = utho_subnet.public.id
}

output "ssh_command" {
  description = "SSH command to connect."
  value       = "ssh root@${utho_cloud.main.ip}"
}
