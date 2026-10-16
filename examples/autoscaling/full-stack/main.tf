# ============================================================
# Example: Full Production Auto Scaling Stack
# ============================================================
# VPC + LB + SG + CPU policies + peak-hour schedule
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

variable "utho_api_key"   { type = string; sensitive = true }
variable "snapshot_id"    { type = string }
variable "snapshot_name"  { type = string }

resource "utho_vpc" "main" {
  name    = "production-vpc"
  network = "10.0.0.0"
  size    = "16"
  dcslug  = "inmumbaizone2"
  planid  = "1008"
}

resource "utho_subnet" "public" {
  name            = "public-subnet"
  vpc_id          = utho_vpc.main.id
  network         = "10.0.1.0"
  size            = 24
  type            = "public"
  assign_publicip = 1
}

resource "utho_firewall" "web" {
  name = "web-sg"
}

resource "utho_autoscaling" "web" {
  name              = "production-web-asg"
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
  vpc               = utho_subnet.public.id
  security_groups   = utho_firewall.web.id

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
    },
    {
      name     = "scale-out-ram"
      type     = "ram"
      compare  = "above"
      value    = "85"
      adjust   = 1
      period   = "5m"
      cooldown = "300"
    }
  ]

  schedules = [
    {
      name                = "peak-hours"
      desiredsize         = "5"
      timezone            = "Asia/Kolkata"
      recurrence          = "Every day  09:00"
      recurrence_duration = "Every day"
      recurrence_week     = ""
      selected_time       = "09:00"
      selected_date       = "2026-10-01"
      start_date          = "2026-10-01T09:00:00.000+05:30"
    }
  ]
}

output "asg_id"     { value = utho_autoscaling.web.id }
output "asg_status" { value = utho_autoscaling.web.status }
