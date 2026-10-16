# ============================================================
# Example: Basic Monitoring Setup
# ============================================================
# Creates an alert contact and CPU/RAM alerts for a
# cloud instance.
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

variable "instance_id" {
  description = "Cloud instance ID to monitor."
  type        = string
}

resource "utho_alert_contact" "main" {
  name         = "ops-team"
  email        = "ops@mycompany.com"
  mobilenumber = "9999999999"
  status       = "1"
}

resource "utho_alert" "cpu" {
  name     = "high-cpu"
  ref_type = "cloud"
  type     = "cpu"
  compare  = "above"
  value    = "80"
  for      = "5m"
  contacts = utho_alert_contact.main.id
  status   = "1"
  ref_ids  = var.instance_id
}

resource "utho_alert" "ram" {
  name     = "high-ram"
  ref_type = "cloud"
  type     = "ram"
  compare  = "above"
  value    = "85"
  for      = "5m"
  contacts = utho_alert_contact.main.id
  status   = "1"
  ref_ids  = var.instance_id
}

output "contact_id" { value = utho_alert_contact.main.id }
output "cpu_alert"  { value = utho_alert.cpu.id }
output "ram_alert"  { value = utho_alert.ram.id }
