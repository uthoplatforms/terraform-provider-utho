# ============================================================
# Example: Full Production Object Storage Stack
# ============================================================
# Complete setup with:
#   - 1 private bucket with versioning
#   - App service key (full access)
#   - Backup service key (write-only)
#   - Read-only reporting key (write access)
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

variable "project" {
  description = "Project name prefix."
  type        = string
  default     = "myapp"
}

# ── Bucket ─────────────────────────────────────────────────

resource "utho_object_storage" "main" {
  name            = "${var.project}-storage"
  dcslug          = "innoida"
  access          = "private"
  billingcycle    = "monthly"
  version_enabled = true
}

# ── Access Keys ────────────────────────────────────────────

resource "utho_object_storage_key" "app" {
  dcslug = "innoida"
  name   = "${var.project}-app-key"
}

resource "utho_object_storage_key" "backup" {
  dcslug = "innoida"
  name   = "${var.project}-backup-key"
}

# ── Permissions ────────────────────────────────────────────

resource "utho_object_storage_permission" "app" {
  dcslug      = "innoida"
  bucket_name = utho_object_storage.main.name
  access_key  = utho_object_storage_key.app.id
  permission  = "full"
}

resource "utho_object_storage_permission" "backup" {
  dcslug      = "innoida"
  bucket_name = utho_object_storage.main.name
  access_key  = utho_object_storage_key.backup.id
  permission  = "write"
}

# ── Outputs ────────────────────────────────────────────────

output "s3_endpoint" { value = "https://innoida.utho.io" }
output "bucket_name" { value = utho_object_storage.main.name }

output "app_access_key" {
  value     = utho_object_storage_key.app.access_key
  sensitive = true
}
output "app_secret_key" {
  value     = utho_object_storage_key.app.secret_key
  sensitive = true
}
output "backup_access_key" {
  value     = utho_object_storage_key.backup.access_key
  sensitive = true
}
output "backup_secret_key" {
  value     = utho_object_storage_key.backup.secret_key
  sensitive = true
}
