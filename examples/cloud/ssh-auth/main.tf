# ============================================================
# Example: Cloud Instance with SSH Key Authentication
# ============================================================
# SSH key authentication is the recommended way to access
# Utho Cloud instances. More secure than passwords and works
# seamlessly with automation tools like Ansible and CI/CD.
#
# What this creates:
#   - 1 SSH key (imported from your local machine)
#   - 1 cloud instance with SSH key auth
#   - Public IP address
#
# Prerequisites:
#   Generate an SSH key if you don't have one:
#   $ ssh-keygen -t ed25519 -C "your@email.com"
#
# Usage:
#   terraform init
#   terraform plan
#   terraform apply
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

variable "ssh_key_name" {
  description = "Name for your SSH key in Utho. Must be unique within your account."
  type        = string
  default     = "my-deploy-key"
}

variable "ssh_public_key_path" {
  description = "Path to your SSH public key file."
  type        = string
  default     = "~/.ssh/id_ed25519.pub"
}

variable "hostname" {
  description = "Hostname for the instance."
  type        = string
  default     = "my-server.mhc"
}

variable "dcslug" {
  description = "Data center slug. Run terraform apply to see available_dczones output."
  type        = string
  default     = "inmumbaizone2"
}

# ── SSH Key ────────────────────────────────────────────────
# Import your local SSH public key into Utho.
# Terraform will manage this key — destroying this resource
# removes the key from your Utho account.

resource "utho_ssh_key" "deploy" {
  name   = var.ssh_key_name
  sshkey = file(var.ssh_public_key_path)
}

# ── Instance ───────────────────────────────────────────────

resource "utho_cloud" "main" {
  hostname     = var.hostname
  dcslug       = var.dcslug
  planid       = "10308"                  # 2 vCPU / 4 GB / 80 GB
  billingcycle = "hourly"
  image        = "ubuntu-22.04-x86_64"

  auth    = "option2"                 # option2 = SSH key auth
  sshkeys = utho_ssh_key.deploy.id   # reference the key created above

  enable_publicip = "true"
  cpumodel        = "amd"
}

# ── Outputs ────────────────────────────────────────────────

output "instance_id" {
  description = "The unique ID of the cloud instance."
  value       = utho_cloud.main.id
}

output "instance_ip" {
  description = "The public IP address of the cloud instance."
  value       = utho_cloud.main.ip
}

output "ssh_key_id" {
  description = "The ID of the imported SSH key."
  value       = utho_ssh_key.deploy.id
}

output "ssh_command" {
  description = "SSH command to connect to your instance."
  value       = "ssh root@${utho_cloud.main.ip}"
}
