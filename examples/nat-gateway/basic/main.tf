# ============================================================
# Example: Basic NAT Gateway
# ============================================================
# Creates a VPC with a public subnet and a NAT Gateway.
# Use this as the foundation before adding private subnets
# and instances.
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

variable "dcslug" {
  type    = string
  default = "inmumbaizone2"
}

resource "utho_vpc" "main" {
  name    = "main-vpc"
  network = "10.0.0.0"
  size    = "16"
  dcslug  = var.dcslug
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

resource "utho_elastic_ip" "nat" {
  dcslug       = var.dcslug
  billingcycle = "monthly"
}

resource "utho_nat_gateway" "main" {
  name      = "main-nat-gateway"
  subnet_id = utho_subnet.public.id
  public_ip = utho_elastic_ip.nat.ip
  dcslug    = var.dcslug
}

output "vpc_id"         { value = utho_vpc.main.id }
output "nat_gateway_id" { value = utho_nat_gateway.main.id }
output "nat_public_ip"  { value = utho_nat_gateway.main.public_ip }
output "nat_status"     { value = utho_nat_gateway.main.status }
