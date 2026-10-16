# ============================================================
# Example: Basic Object Storage Bucket
# ============================================================
# Creates a private S3-compatible bucket with a dedicated
# access key. Use the outputs to configure your application.
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

variable "bucket_name" {
  description = "Unique bucket name."
  type        = string
  default     = "myapp-storage"
}

resource "utho_object_storage" "main" {
  name         = var.bucket_name
  dcslug       = "innoida"
  access       = "private"
  billingcycle = "monthly"
}

resource "utho_object_storage_key" "app" {
  dcslug = "innoida"
  name   = "${var.bucket_name}-key"
}

resource "utho_object_storage_permission" "app" {
  dcslug      = "innoida"
  bucket_name = utho_object_storage.main.name
  access_key  = utho_object_storage_key.app.id
  permission  = "full"
}

output "s3_endpoint"   { value = "https://innoida.utho.io" }
output "s3_bucket"     { value = utho_object_storage.main.name }
output "s3_access_key" {
  value     = utho_object_storage_key.app.access_key
  sensitive = true
}
output "s3_secret_key" {
  value     = utho_object_storage_key.app.secret_key
  sensitive = true
}
