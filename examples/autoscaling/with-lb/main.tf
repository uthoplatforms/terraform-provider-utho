# ============================================================
# Example: ASG with Load Balancer + Security Group
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

variable "snapshot_id"   { type = string }
variable "snapshot_name" { type = string }
variable "lb_id"         { type = string }
variable "firewall_id"   { type = string }

resource "utho_autoscaling" "main" {
  name              = "web-asg"
  dcslug            = "inmumbaizone2"
  minsize           = "2"
  maxsize           = "10"
  desiredsize       = "2"
  planid            = "10404"
  planname          = "basic"
  snapshotid        = var.snapshot_id
  image_name        = var.snapshot_name
  os_disk_size      = 80
  public_ip_enabled = 1
  cpumodel          = "intel"
  load_balancers    = var.lb_id
  security_groups   = var.firewall_id

  policies = [
    {
      name     = "scale-out-cpu"
      type     = "cpu"
      compare  = "above"
      value    = "75"
      adjust   = 2
      period   = "5m"
      cooldown = "300"
    },
    {
      name     = "scale-in-cpu"
      type     = "cpu"
      compare  = "below"
      value    = "25"
      adjust   = -1
      period   = "10m"
      cooldown = "600"
    }
  ]
}

output "asg_id"     { value = utho_autoscaling.main.id }
output "asg_status" { value = utho_autoscaling.main.status }
