# ============================================================
# Example: Cloud Instance with EBS Block Storage
# ============================================================
# Deploy a cloud instance with additional EBS (Elastic Block
# Storage) volumes. EBS volumes are persistent — they survive
# instance reboots and can be used for databases, large
# datasets, or any data that must be retained.
#
# What this creates:
#   - 1 SSH key
#   - 1 cloud instance on an EBS plan (no included disk)
#   - 1 root EBS volume (80 GB NVMe)
#   - 1 data EBS volume (100 GB NVMe)
#
# EBS Plan vs General Plan:
#   - EBS plan: no disk included, all storage via EBS volumes
#   - General plan: disk included, add EBS for extra storage
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
  default     = "db-server.mhc"
}

variable "dcslug" {
  description = "Data center slug. EBS is available in: innoida, inmumbaizone2."
  type        = string
  default     = "innoida"   # EBS available in Delhi (Noida)
}

variable "root_disk_gb" {
  description = "Size of root OS volume in GB."
  type        = number
  default     = 80
}

variable "data_disk_gb" {
  description = "Size of data volume in GB. Used for application data."
  type        = number
  default     = 100
}

# ── SSH Key ────────────────────────────────────────────────

resource "utho_ssh_key" "deploy" {
  name   = "ebs-example-key"
  sshkey = file("~/.ssh/id_ed25519.pub")
}

# ── Instance with EBS ──────────────────────────────────────
# planid "10355" is an EBS-only plan (no included disk).
# All storage is via EBS volumes defined in the ebs block.
#
# ebs block volumes:
#   id   — sequential identifier ("1", "2", "3"...)
#   disk — size in GB
#   type — "nvme" (high performance) or "ssd" (standard)

resource "utho_cloud" "main" {
  hostname     = var.hostname
  dcslug       = var.dcslug
  planid       = "10355"                  # EBS plan — 2 vCPU / 4 GB / no disk
  billingcycle = "hourly"
  image        = "ubuntu-22.04-x86_64"

  auth    = "option2"
  sshkeys = utho_ssh_key.deploy.id

  enable_publicip = "true"
  cpumodel        = "intel"   # innoida uses intel CPUs
  delete_ebs      = true      # delete EBS volumes when instance is destroyed

  ebs = [
    {
      id   = "1"
      disk = var.root_disk_gb   # root OS volume
      type = "nvme"             # NVMe for best performance
    },
    {
      id   = "2"
      disk = var.data_disk_gb   # data volume for application storage
      type = "nvme"
    },
  ]
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

output "ssh_command" {
  description = "SSH command to connect."
  value       = "ssh root@${utho_cloud.main.ip}"
}

output "storage_info" {
  description = "EBS storage configuration."
  value = {
    root_disk_gb = var.root_disk_gb
    data_disk_gb = var.data_disk_gb
    total_gb     = var.root_disk_gb + var.data_disk_gb
  }
}
