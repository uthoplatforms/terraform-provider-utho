# ============================================================
# Example: Basic Elastic IP Allocation
# ============================================================
# Allocates a static public IP without attaching it to an
# instance. Use this to reserve an IP before your instance
# is ready, or to keep an IP between deployments.
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

resource "utho_elastic_ip" "main" {
  dcslug       = var.dcslug
  billingcycle = "monthly"
}

output "static_ip" {
  description = "Your static public IP address."
  value       = utho_elastic_ip.main.ip
}
