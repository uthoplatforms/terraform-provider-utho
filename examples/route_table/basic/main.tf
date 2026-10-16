# ============================================================
# Example: Basic Route Table with IGW Route
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

variable "vpc_id" {
  description = "VPC ID to associate the route table with."
  type        = string
}

resource "utho_route_table" "main" {
  name   = "main-route-table"
  vpc_id = var.vpc_id
  dcslug = "inmumbaizone2"
}

resource "utho_route" "internet" {
  route_table_id         = utho_route_table.main.id
  destination_cidr_block = "10.0.0.0/16"
  route_type             = "igw"
  target                 = "igw"
}

output "route_table_id" { value = utho_route_table.main.id }
output "route_id"       { value = utho_route.internet.id }
