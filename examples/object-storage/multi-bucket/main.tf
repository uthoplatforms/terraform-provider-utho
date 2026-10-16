# ============================================================
# Example: Multi-Bucket Deployment
# ============================================================
# Deploy 3 buckets for different use cases:
#   - Private bucket for application data
#   - Public bucket for static assets
#   - Upload-only bucket for user submissions
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

variable "project" {
  type    = string
  default = "myapp"
}

resource "utho_object_storage" "data" {
  name         = "${var.project}-data"
  dcslug       = "innoida"
  access       = "private"
  billingcycle = "monthly"
}

resource "utho_object_storage" "assets" {
  name         = "${var.project}-assets"
  dcslug       = "innoida"
  access       = "public"
  billingcycle = "monthly"
}

resource "utho_object_storage" "uploads" {
  name         = "${var.project}-uploads"
  dcslug       = "innoida"
  access       = "upload"
  billingcycle = "monthly"
}

output "data_bucket"    { value = utho_object_storage.data.name }
output "assets_bucket"  { value = utho_object_storage.assets.name }
output "uploads_bucket" { value = utho_object_storage.uploads.name }
output "s3_endpoint"    { value = "https://innoida.utho.io" }
