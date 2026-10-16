# ============================================================
# Example: Full Private Network with NAT Gateway
# ============================================================
# Production-grade setup:
#   - VPC with public and private subnets
#   - NAT gateway for private instance outbound access
#   - Private cloud instance (no public IP)
#   - Load balancer in public subnet (optional)
#
# Architecture:
#   Internet → Load Balancer (public subnet)
#                  ↓
#             Private instances (private subnet)
#                  ↓ (outbound only)
#             NAT Gateway (public subnet)
#                  ↓
#              Internet
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

variable "ssh_public_key_path" {
  type    = string
  default = "~/.ssh/id_ed25519.pub"
}

# ── VPC ────────────────────────────────────────────────────

resource "utho_vpc" "main" {
  name    = "production-vpc"
  network = "10.0.0.0"
  size    = "16"
  dcslug  = var.dcslug
  planid  = "1008"
}

# ── Subnets ────────────────────────────────────────────────

resource "utho_subnet" "public" {
  name            = "public-subnet"
  vpc_id          = utho_vpc.main.id
  network         = "10.0.1.0"
  size            = 24
  type            = "public"
  assign_publicip = 1
}

resource "utho_subnet" "private" {
  name            = "private-subnet"
  vpc_id          = utho_vpc.main.id
  network         = "10.0.2.0"
  size            = 24
  type            = "private"
  assign_publicip = 0
}

# ── Elastic IP + NAT Gateway ───────────────────────────────

resource "utho_elastic_ip" "nat" {
  dcslug       = var.dcslug
  billingcycle = "monthly"
}

resource "utho_nat_gateway" "main" {
  name      = "production-nat"
  subnet_id = utho_subnet.public.id
  public_ip = utho_elastic_ip.nat.ip
  dcslug    = var.dcslug
}

# ── SSH Key ────────────────────────────────────────────────

resource "utho_ssh_key" "deploy" {
  name   = "deploy-key"
  sshkey = file(var.ssh_public_key_path)
}

# ── Private Cloud Instance ─────────────────────────────────
# No public IP — outbound traffic routes through NAT gateway

resource "utho_cloud" "worker" {
  hostname        = "worker-01.mhc"
  dcslug          = var.dcslug
  planid          = "10308"
  billingcycle    = "hourly"
  auth            = "option2"
  sshkeys         = utho_ssh_key.deploy.id
  image           = "ubuntu-22.04-x86_64"
  enable_publicip = "false"
  vpc             = utho_subnet.private.id
}

# ── Outputs ────────────────────────────────────────────────

output "vpc_id"          { value = utho_vpc.main.id }
output "public_subnet"   { value = utho_subnet.public.id }
output "private_subnet"  { value = utho_subnet.private.id }
output "nat_gateway_id"  { value = utho_nat_gateway.main.id }
output "nat_public_ip"   { value = utho_nat_gateway.main.public_ip }
output "worker_id"       { value = utho_cloud.worker.id }
output "worker_ip"       { value = utho_cloud.worker.ip }
