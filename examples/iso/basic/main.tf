# ============================================================
# Example: Upload ISO and Mount on Cloud Instance
# ============================================================

terraform {
  required_providers {
    utho = {
      source  = "dharmajputho/utho-dev"
      version = ">= 0.2.64"
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

variable "cloud_id" {
  description = "Cloud instance ID to mount the ISO on."
  type        = string
}

# ── Upload ISO ─────────────────────────────────────────────
resource "utho_iso" "alpine" {
  name   = "alpine-linux"
  dcslug = "innoida"
  url    = "https://dl-cdn.alpinelinux.org/alpine/latest-stable/releases/x86_64/alpine-virt-3.24.3-x86_64.iso"
}

# ── Mount on instance ──────────────────────────────────────
resource "utho_cloud_iso" "mount" {
  cloud_id = var.cloud_id
  iso      = utho_iso.alpine.file
}

# ── Outputs ────────────────────────────────────────────────
output "iso_id"   { value = utho_iso.alpine.id }
output "iso_file" { value = utho_iso.alpine.file }
