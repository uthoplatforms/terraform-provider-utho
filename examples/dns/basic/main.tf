# ============================================================
# Example: Basic DNS Zone with Common Records
# ============================================================
# Sets up a domain with A, CNAME, and TXT records.
# After applying, point your domain's nameservers to Utho
# at your domain registrar.
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
  description = "Your Utho API key."
  type        = string
  sensitive   = true
}

variable "domain" {
  description = "Your domain name (e.g. myapp.com)."
  type        = string
}

variable "server_ip" {
  description = "Your server IP address."
  type        = string
}

resource "utho_dns_zone" "main" {
  domain = var.domain
}

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

resource "utho_dns_record" "spf" {
  domain   = utho_dns_zone.main.domain
  type     = "TXT"
  hostname = "@"
  value    = "v=spf1 include:_spf.google.com ~all"
  ttl      = "3600"
}

output "zone_id"      { value = utho_dns_zone.main.id }
output "nspoint"      { value = utho_dns_zone.main.nspoint }
output "record_count" { value = utho_dns_zone.main.record_count }
