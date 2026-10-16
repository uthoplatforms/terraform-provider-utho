# ============================================================
# Example: Full Production Monitoring Stack
# ============================================================
# Two alert contacts (devops + on-call) and CPU/RAM alerts
# watching multiple instances. All contacts notified on alert.
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

variable "instance_ids" {
  description = "Comma-separated cloud instance IDs to monitor."
  type        = string
}

# ── Alert Contacts ─────────────────────────────────────────

resource "utho_alert_contact" "devops" {
  name         = "devops-team"
  email        = "devops@mycompany.com"
  mobilenumber = "9999999999"
  status       = "1"
}

resource "utho_alert_contact" "oncall" {
  name         = "oncall-engineer"
  email        = "oncall@mycompany.com"
  mobilenumber = "8888888888"
  status       = "1"
}

# ── Alerts ─────────────────────────────────────────────────

resource "utho_alert" "cpu" {
  name     = "prod-high-cpu"
  ref_type = "cloud"
  type     = "cpu"
  compare  = "above"
  value    = "80"
  for      = "5m"
  contacts = "${utho_alert_contact.devops.id},${utho_alert_contact.oncall.id}"
  status   = "1"
  ref_ids  = var.instance_ids
}

resource "utho_alert" "ram" {
  name     = "prod-high-ram"
  ref_type = "cloud"
  type     = "ram"
  compare  = "above"
  value    = "85"
  for      = "5m"
  contacts = "${utho_alert_contact.devops.id},${utho_alert_contact.oncall.id}"
  status   = "1"
  ref_ids  = var.instance_ids
}

# ── Outputs ────────────────────────────────────────────────

output "devops_contact" { value = utho_alert_contact.devops.id }
output "oncall_contact" { value = utho_alert_contact.oncall.id }
output "cpu_alert"      { value = utho_alert.cpu.id }
output "ram_alert"      { value = utho_alert.ram.id }
