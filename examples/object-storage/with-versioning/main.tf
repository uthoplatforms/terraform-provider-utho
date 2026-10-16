# ============================================================
# Example: Versioned Bucket for Backups
# ============================================================
# Versioning protects against accidental overwrites and
# deletions — every write creates a new version.
# Perfect for database backups, config files, critical data.
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

resource "utho_object_storage" "backups" {
  name            = "myapp-backups"
  dcslug          = "innoida"
  access          = "private"
  billingcycle    = "monthly"
  version_enabled = true
}

resource "utho_object_storage_key" "backup_writer" {
  dcslug = "innoida"
  name   = "backup-writer-key"
}

resource "utho_object_storage_permission" "backup_writer" {
  dcslug      = "innoida"
  bucket_name = utho_object_storage.backups.name
  access_key  = utho_object_storage_key.backup_writer.id
  permission  = "write"
}

output "s3_endpoint"   { value = "https://innoida.utho.io" }
output "bucket_name"   { value = utho_object_storage.backups.name }
output "versioning"    { value = utho_object_storage.backups.version_enabled }
output "access_key" {
  value     = utho_object_storage_key.backup_writer.access_key
  sensitive = true
}
output "secret_key" {
  value     = utho_object_storage_key.backup_writer.secret_key
  sensitive = true
}
