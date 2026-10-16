# ============================================================
# Example: Basic Auto Scaling Group with Snapshot
# ============================================================
# Creates an ASG that scales based on CPU usage.
# Uses a snapshot for instance deployment.
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

variable "snapshot_id" {
  type        = string
  description = "Snapshot ID to deploy instances from."
}

variable "snapshot_name" {
  type        = string
  description = "Snapshot name."
}

resource "utho_autoscaling" "main" {
  name              = "web-asg"
  dcslug            = "inmumbaizone2"
  minsize           = "1"
  maxsize           = "5"
  desiredsize       = "1"
  planid            = "10404"
  planname          = "basic"
  snapshotid        = var.snapshot_id
  image_name        = var.snapshot_name
  os_disk_size      = 80
  public_ip_enabled = 1
  cpumodel          = "intel"

  policies = [
    {
      name     = "scale-out-cpu"
      type     = "cpu"
      compare  = "above"
      value    = "80"
      adjust   = 1
      period   = "5m"
      cooldown = "300"
    },
    {
      name     = "scale-in-cpu"
      type     = "cpu"
      compare  = "below"
      value    = "20"
      adjust   = -1
      period   = "5m"
      cooldown = "300"
    }
  ]
}

output "asg_id"     { value = utho_autoscaling.main.id }
output "asg_status" { value = utho_autoscaling.main.status }
