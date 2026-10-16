# ============================================================
# Example: Load Balancer with Weighted Backends
# ============================================================
# Distribute traffic unevenly across backends using weights.
# Useful for A/B testing, canary deployments, or gradually
# shifting traffic to new servers.
#
# Traffic distribution with weights 3, 3, 1:
#   - backend-1: ~43% of traffic
#   - backend-2: ~43% of traffic
#   - backend-3: ~14% of traffic (canary)
#
# Usage:
#   terraform init
#   terraform apply
#
# Scale up:
#   terraform apply -var="backend_count=5"
#
# Scale down:
#   terraform apply -var="backend_count=2"
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

# ── Variables ─────────────────────────────────────────────

variable "utho_api_key" {
  description = "Your Utho API key."
  type        = string
  sensitive   = true
}

variable "dcslug" {
  type    = string
  default = "inmumbaizone2"
}

variable "backend_count" {
  description = "Number of primary backend servers."
  type        = number
  default     = 3
}

# ── Network ────────────────────────────────────────────────

resource "utho_vpc" "main" {
  name    = "weighted-lb-vpc"
  network = "10.1.0.0"
  size    = "24"
  dcslug  = var.dcslug
  planid  = "1008"
}

resource "utho_subnet" "public" {
  name            = "weighted-lb-subnet"
  vpc_id          = utho_vpc.main.id
  network         = "10.1.0.0"
  size            = 24
  type            = "public"
  assign_publicip = 1
}

resource "utho_ssh_key" "deploy" {
  name   = "weighted-lb-key"
  sshkey = file("~/.ssh/id_ed25519.pub")
}

# ── Primary backends (equal weight) ───────────────────────

resource "utho_cloud" "primary" {
  count = var.backend_count

  hostname        = "primary-${count.index + 1}.mhc"
  dcslug          = var.dcslug
  planid          = "10308"
  billingcycle    = "hourly"
  auth            = "option2"
  sshkeys         = utho_ssh_key.deploy.id
  image           = "ubuntu-22.04-x86_64"
  enable_publicip = "true"
  cpumodel        = "amd"
}

# ── Canary backend (lower weight) ─────────────────────────
# Gets ~14% of traffic when primary weight=3

resource "utho_cloud" "canary" {
  hostname        = "canary-1.mhc"
  dcslug          = var.dcslug
  planid          = "10308"
  billingcycle    = "hourly"
  auth            = "option2"
  sshkeys         = utho_ssh_key.deploy.id
  image           = "ubuntu-22.04-x86_64"
  enable_publicip = "true"
  cpumodel        = "amd"
}

# ── Load Balancer ──────────────────────────────────────────

resource "utho_loadbalancer" "main" {
  name            = "weighted-lb"
  type            = "application"
  dcslug          = var.dcslug
  vpc             = utho_subnet.public.id
  enable_publicip = "true"
}

resource "utho_loadbalancer_frontend" "http" {
  loadbalancer_id = utho_loadbalancer.main.id
  name            = "http"
  algorithm       = "roundrobin"
  proto           = "http"
  port            = "80"
  cookie          = "0"
  redirecthttps   = "0"
  certificate_id  = "0"
}

# Primary backends — weight 3 each
resource "utho_loadbalancer_backend" "primary" {
  count = var.backend_count

  loadbalancer_id = utho_loadbalancer.main.id
  frontend_id     = utho_loadbalancer_frontend.http.id
  backend_port    = "80"
  weight          = "3"
  type            = "cloud"
  cloudid         = utho_cloud.primary[count.index].id
}

# Canary backend — weight 1 (lower traffic)
resource "utho_loadbalancer_backend" "canary" {
  loadbalancer_id = utho_loadbalancer.main.id
  frontend_id     = utho_loadbalancer_frontend.http.id
  backend_port    = "80"
  weight          = "1"
  type            = "cloud"
  cloudid         = utho_cloud.canary.id
}

# ── Outputs ────────────────────────────────────────────────

output "lb_ip"        { value = utho_loadbalancer.main.ip }
output "lb_dns"       { value = utho_loadbalancer.main.dns }
output "primary_ips"  { value = utho_cloud.primary[*].ip }
output "canary_ip"    { value = utho_cloud.canary.ip }
