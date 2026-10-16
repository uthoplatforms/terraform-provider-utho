# ============================================================
# Example: EBS Volume Attached to Cloud Instance
# ============================================================
# Creates an EBS volume and attaches it to a cloud instance.
# Destroying the attachment detaches without deleting volume.
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

variable "cloud_id" {
  description = "Cloud instance ID to attach the EBS volume to."
  type        = string
}

# ── EBS Volume ─────────────────────────────────────────────
resource "utho_ebs" "data" {
  name       = "app-data-volume"
  dcslug     = "innoida"
  disk_type  = "SSD"
  disk       = "100"
  iops       = "500"
  throughput = "250"
}

# ── Attach to instance ─────────────────────────────────────
resource "utho_ebs_attachment" "data" {
  ebs_id   = utho_ebs.data.id
  cloud_id = var.cloud_id
}

# ── Outputs ────────────────────────────────────────────────
output "ebs_id"      { value = utho_ebs.data.id }
output "ebs_status"  { value = utho_ebs.data.status }
output "device"      { value = utho_ebs_attachment.data.device }
