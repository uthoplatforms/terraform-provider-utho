# ============================================================
# Example: Full Production DNS Setup
# ============================================================
# Complete DNS setup for a production domain:
#   - Root A record
#   - IPv6 AAAA record
#   - www CNAME
#   - API subdomain
#   - MX for email
#   - SPF + DMARC TXT records
#
# Single terraform apply creates everything.
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

variable "domain" {
  description = "Your domain (e.g. myapp.com)."
  type        = string
}

variable "server_ip" {
  description = "Primary server IP."
  type        = string
}

variable "server_ipv6" {
  description = "Primary server IPv6 address."
  type        = string
  default     = ""
}

variable "api_ip" {
  description = "API server IP."
  type        = string
  default     = ""
}

# ── DNS Zone ───────────────────────────────────────────────

resource "utho_dns_zone" "main" {
  domain = var.domain
}

# ── Web Records ────────────────────────────────────────────

resource "utho_dns_record" "root" {
  domain   = utho_dns_zone.main.domain
  type     = "A"
  hostname = "@"
  value    = var.server_ip
  ttl      = "300"
}

resource "utho_dns_record" "www" {
  domain   = utho_dns_zone.main.domain
  type     = "CNAME"
  hostname = "www"
  value    = "${var.domain}."
  ttl      = "3600"
}

resource "utho_dns_record" "ipv6" {
  count    = var.server_ipv6 != "" ? 1 : 0
  domain   = utho_dns_zone.main.domain
  type     = "AAAA"
  hostname = "@"
  value    = var.server_ipv6
  ttl      = "3600"
}

resource "utho_dns_record" "api" {
  count    = var.api_ip != "" ? 1 : 0
  domain   = utho_dns_zone.main.domain
  type     = "A"
  hostname = "api"
  value    = var.api_ip
  ttl      = "60"
}

# ── Email Records ──────────────────────────────────────────

resource "utho_dns_record" "mx" {
  domain   = utho_dns_zone.main.domain
  type     = "MX"
  hostname = "@"
  value    = "mail.${var.domain}."
  ttl      = "3600"
}

resource "utho_dns_record" "spf" {
  domain   = utho_dns_zone.main.domain
  type     = "TXT"
  hostname = "@"
  value    = "v=spf1 include:_spf.google.com ~all"
  ttl      = "3600"
}

resource "utho_dns_record" "dmarc" {
  domain   = utho_dns_zone.main.domain
  type     = "TXT"
  hostname = "_dmarc"
  value    = "v=DMARC1; p=none; rua=mailto:dmarc@${var.domain}"
  ttl      = "3600"
}

# ── Outputs ────────────────────────────────────────────────

output "zone_id"      { value = utho_dns_zone.main.id }
output "nspoint"      { value = utho_dns_zone.main.nspoint }
output "next_step"    { value = "Point your domain nameservers to Utho at your registrar." }
