# ============================================================
# Example: Basic Cloud Instance (Password Authentication)
# ============================================================
# This is the simplest way to deploy a Utho Cloud instance.
# Use this as a starting point for development or testing.
#
# What this creates:
#   - 1 cloud instance (2 vCPU / 4 GB RAM / 80 GB disk)
#   - Public IP address
#   - Root password authentication
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
  description = "Your Utho API key. Get it from https://console.utho.com/api"
  type        = string
  sensitive   = true
}

variable "root_password" {
  description = "Root password for the instance. Must be at least 8 characters with uppercase, lowercase, numbers and special characters."
  type        = string
  sensitive   = true
}

variable "hostname" {
  description = "Hostname for the instance. Must be unique within your account."
  type        = string
  default     = "my-server.mhc"
}

# ── Discover valid values ──────────────────────────────────
# Use data sources to find valid DC slugs, plan IDs and image slugs.
# Run `terraform apply` on just the data sources to see available options.

data "utho_cloud_dczones" "all" {}

data "utho_cloud_plans" "mumbai" {
  dcslug = "inmumbaizone2"
}

data "utho_cloud_images" "ubuntu" {
  distro = "ubuntu"
}

# ── Instance ───────────────────────────────────────────────

resource "utho_cloud" "main" {
  hostname     = var.hostname
  dcslug       = "inmumbaizone2"          # Mumbai — change to your preferred DC
  planid       = "10308"                  # 2 vCPU / 4 GB / 80 GB — see data.utho_cloud_plans
  billingcycle = "hourly"                 # hourly | monthly | 3month | 6month | 12month | 24month | 36month
  image        = "ubuntu-22.04-x86_64"   # see data.utho_cloud_images

  auth          = "option1"       # option1 = password auth
  root_password = var.root_password

  enable_publicip = "true"
  cpumodel        = "amd"         # amd or intel — check utho_cloud_dczones for availability
}

# ── Outputs ────────────────────────────────────────────────

output "instance_id" {
  description = "The unique ID of the cloud instance."
  value       = utho_cloud.main.id
}

output "instance_ip" {
  description = "The public IP address of the cloud instance. Point your DNS A record here."
  value       = utho_cloud.main.ip
}

output "ssh_command" {
  description = "SSH command to connect to your instance."
  value       = "ssh root@${utho_cloud.main.ip}"
}

output "available_dczones" {
  description = "All available data center zones. Use slug value as dcslug."
  value = [
    for z in data.utho_cloud_dczones.all.zones :
    { slug = z.slug, city = z.city, country = z.country }
    if z.status == "active"
  ]
}

output "available_plans" {
  description = "Available plans in Mumbai with disk included."
  value = [
    for p in data.utho_cloud_plans.mumbai.plans :
    { id = p.id, cpu = p.cpu, ram_mb = p.ram, disk_gb = p.disk, price_inr = p.price }
    if p.disk != "0"
  ]
}

output "available_ubuntu_images" {
  description = "Available Ubuntu images. Use image value in resource."
  value       = data.utho_cloud_images.ubuntu.images[*].image
}
