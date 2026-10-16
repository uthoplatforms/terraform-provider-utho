# ============================================================
# Example: Email DNS Setup (Google Workspace / Gmail)
# ============================================================
# Sets up DNS records for email:
#   - MX records for mail routing
#   - SPF to authorize sending servers
#   - DMARC for email policy
#   - DKIM (add after getting key from email provider)
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

resource "utho_dns_zone" "main" {
  domain = var.domain
}

# MX record — route email to mail server
resource "utho_dns_record" "mx" {
  domain   = utho_dns_zone.main.domain
  type     = "MX"
  hostname = "@"
  value    = "mail.${var.domain}."
  ttl      = "3600"
}

# SPF — authorize which servers can send email
resource "utho_dns_record" "spf" {
  domain   = utho_dns_zone.main.domain
  type     = "TXT"
  hostname = "@"
  value    = "v=spf1 include:_spf.google.com ~all"
  ttl      = "3600"
}

# DMARC — email policy and reporting
resource "utho_dns_record" "dmarc" {
  domain   = utho_dns_zone.main.domain
  type     = "TXT"
  hostname = "_dmarc"
  value    = "v=DMARC1; p=quarantine; rua=mailto:dmarc@${var.domain}; pct=100"
  ttl      = "3600"
}

output "zone_id"   { value = utho_dns_zone.main.id }
output "next_step" { value = "Add DKIM record from your email provider after setup." }
