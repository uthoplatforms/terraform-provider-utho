# ============================================================
# Example: Site-to-Site VPN with Two IPSec Tunnels
# ============================================================
# Creates two IPSec tunnels and pairs them together to
# establish a site-to-site VPN connection.
#
# Site A network: 192.168.50.0/24
# Site B network: 192.168.60.0/24
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

variable "site_a_subnet" {
  description = "Subnet ID for site A tunnel."
  type        = string
}

variable "site_b_subnet" {
  description = "Subnet ID for site B tunnel."
  type        = string
}

# ── Site A Tunnel ──────────────────────────────────────────
resource "utho_ipsec" "site_a" {
  name         = "site-a-vpn"
  dcslug       = "inmumbaizone2"
  vpc          = var.site_a_subnet
  billingcycle = "monthly"
}

# ── Site B Tunnel ──────────────────────────────────────────
resource "utho_ipsec" "site_b" {
  name         = "site-b-vpn"
  dcslug       = "inmumbaizone2"
  vpc          = var.site_b_subnet
  billingcycle = "monthly"
}

# ── Pair the tunnels ────────────────────────────────────────
resource "utho_ipsec_connection" "vpn" {
  ipsec_id      = utho_ipsec.site_a.id
  peer_ipsec_id = utho_ipsec.site_b.id
  local_subnet  = "192.168.50.0/24"
  peer_subnet   = "192.168.60.0/24"
  name          = "site-a-to-site-b"

  depends_on = [utho_ipsec.site_a, utho_ipsec.site_b]
}

# ── Outputs ─────────────────────────────────────────────────
output "site_a_id"      { value = utho_ipsec.site_a.id }
output "site_b_id"      { value = utho_ipsec.site_b.id }
output "connection_id"  { value = utho_ipsec_connection.vpn.id }
output "conn_status"    { value = utho_ipsec_connection.vpn.status }
output "site_a_psk" {
  value     = utho_ipsec.site_a.psk
  sensitive = true
}
output "site_b_psk" {
  value     = utho_ipsec.site_b.psk
  sensitive = true
}
