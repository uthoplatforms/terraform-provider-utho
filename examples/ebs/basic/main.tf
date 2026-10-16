# ============================================================
# Example: Basic EBS Volume
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

resource "utho_ebs" "data" {
  name       = "app-data-volume"
  dcslug     = "innoida"
  disk_type  = "SSD"
  disk       = "50"
  iops       = "250"
  throughput = "125"
}

output "ebs_id"     { value = utho_ebs.data.id }
output "ebs_status" { value = utho_ebs.data.status }
