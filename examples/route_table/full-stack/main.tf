# ============================================================
# Example: Full VPC Stack with Route Table and Routes
# ============================================================
# Creates a VPC with subnets, a route table, and routes
# for internet and internal traffic routing.
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

# ── VPC ────────────────────────────────────────────────────
resource "utho_vpc" "main" {
  name    = "prod-vpc"
  network = "10.0.0.0"
  size    = "20"
  dcslug  = "inmumbaizone2"
  planid  = "1008"
}

# ── Subnet ─────────────────────────────────────────────────
resource "utho_subnet" "public" {
  name            = "public-subnet"
  vpc_id          = utho_vpc.main.id
  network         = "10.0.1.0"
  size            = 24
  type            = "public"
  assign_publicip = 1
}

# ── Route Table ────────────────────────────────────────────
resource "utho_route_table" "main" {
  name   = "prod-route-table"
  vpc_id = utho_vpc.main.id
  dcslug = "inmumbaizone2"
}

# ── Routes ─────────────────────────────────────────────────

# Internet gateway route
resource "utho_route" "internet" {
  route_table_id         = utho_route_table.main.id
  destination_cidr_block = "10.0.0.0/16"
  route_type             = "igw"
  target                 = "igw"
}

# Local subnet route
resource "utho_route" "local" {
  route_table_id         = utho_route_table.main.id
  destination_cidr_block = "10.0.1.0/24"
  route_type             = "local"
  target                 = utho_subnet.public.id
}

# ── Outputs ────────────────────────────────────────────────
output "vpc_id"          { value = utho_vpc.main.id }
output "subnet_id"       { value = utho_subnet.public.id }
output "route_table_id"  { value = utho_route_table.main.id }
output "internet_route"  { value = utho_route.internet.id }
output "local_route"     { value = utho_route.local.id }
