# ============================================================
# Example: Elastic IP Attached to Cloud Instance
# ============================================================
# Allocates a static IP and attaches it to a cloud instance.
# The IP stays the same even if you destroy and recreate
# the instance — your DNS records never need to change.
#
# To move the IP to a new instance:
#   1. Change cloud_id to the new instance ID
#   2. Run terraform apply
#   The IP moves without any DNS changes needed.
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

variable "utho_api_key" {
  type      = string
  sensitive = true
}

variable "dcslug" {
  type    = string
  default = "inmumbaizone2"
}

# ── SSH Key ────────────────────────────────────────────────

resource "utho_ssh_key" "deploy" {
  name   = "deploy-key"
  sshkey = file("~/.ssh/id_ed25519.pub")
}

# ── Cloud Instance ─────────────────────────────────────────

resource "utho_cloud" "web" {
  hostname        = "web-server.mhc"
  dcslug          = var.dcslug
  planid          = "10308"
  billingcycle    = "hourly"
  auth            = "option2"
  sshkeys         = utho_ssh_key.deploy.id
  image           = "ubuntu-22.04-x86_64"
  enable_publicip = "true"
  cpumodel        = "amd"
}

# ── Elastic IP ─────────────────────────────────────────────

resource "utho_elastic_ip" "web" {
  dcslug       = var.dcslug
  billingcycle = "monthly"
  cloud_id     = utho_cloud.web.id
}

# ── Outputs ────────────────────────────────────────────────

output "static_ip" {
  description = "Static IP — use this in your DNS A record."
  value       = utho_elastic_ip.web.ip
}

output "instance_id" {
  description = "Cloud instance ID."
  value       = utho_cloud.web.id
}

output "ssh_command" {
  description = "SSH into your server."
  value       = "ssh root@${utho_elastic_ip.web.ip}"
}
