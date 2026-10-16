# ============================================================
# Example: Basic IPSec Tunnel
# ============================================================
# Creates a single IPSec VPN tunnel in a VPC subnet.
# Use utho_ipsec_connection to pair with another tunnel.
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

variable "subnet_id" {
  description = "Subnet ID to deploy the IPSec gateway in."
  type        = string
}

resource "utho_ipsec" "main" {
  name         = "prod-vpn"
  dcslug       = "inmumbaizone2"
  vpc          = var.subnet_id
  billingcycle = "monthly"
}

output "tunnel_id" { value = utho_ipsec.main.id }
output "tunnel_status" { value = utho_ipsec.main.status }
output "tunnel_cloudid" { value = utho_ipsec.main.cloudid }
output "tunnel_psk" {
  value     = utho_ipsec.main.psk
  sensitive = true
}
