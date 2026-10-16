# ============================================================
# Example: VPC Peering Connection
# ============================================================
# Creates a peering connection between two VPCs to allow
# private communication between their instances.
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

# App VPC
resource "utho_vpc" "app" {
  name    = "app-vpc"
  network = "10.0.0.0"
  size    = "20"
  dcslug  = "inmumbaizone2"
  planid  = "1008"
}

# DB VPC
resource "utho_vpc" "db" {
  name    = "db-vpc"
  network = "10.1.0.0"
  size    = "20"
  dcslug  = "inmumbaizone2"
  planid  = "1008"
}

# Peering connection
resource "utho_vpc_peering" "app_to_db" {
  name             = "app-to-db"
  dcslug           = "inmumbaizone2"
  requester_vpc_id = utho_vpc.app.id
  accepter_vpc_id  = utho_vpc.db.id
}

output "peering_id"     { value = utho_vpc_peering.app_to_db.id }
output "peering_status" { value = utho_vpc_peering.app_to_db.status }
