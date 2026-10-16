---
page_title: "Deploying a Production Web Stack"
subcategory: "Tutorials"
description: |-
  Build a complete production environment on Utho: VPC, security groups, load-balanced web servers, a managed PostgreSQL database, and DNS.
---

# Deploying a Production Web Stack

This tutorial builds a complete, production-style environment in a single configuration:

* A VPC with a public subnet
* A security group that allows HTTP and SSH
* Three web servers behind an application load balancer
* A managed PostgreSQL cluster with point-in-time recovery and a replica
* A DNS zone with the domain pointed at the load balancer

Terraform works out the dependencies between these resources from the references in the configuration and creates them in the right order.

```text
                 Internet
                    │
             ┌──────▼──────┐
             │ utho_dns_*  │  myapp.com → load balancer IP
             └──────┬──────┘
             ┌──────▼──────────────┐
             │ utho_loadbalancer   │  HTTP :80, round robin
             └──────┬──────────────┘
       ┌────────────┼────────────┐
  ┌────▼───┐   ┌────▼───┐   ┌────▼───┐
  │ app-1  │   │ app-2  │   │ app-3  │  utho_cloud × 3 (web-sg)
  └────┬───┘   └────┬───┘   └────┬───┘
       └────────────┼────────────┘
             ┌──────▼──────────────┐
             │ utho_database (pg)  │  private network
             └─────────────────────┘
```

## Prerequisites

* Complete [Getting Started](getting-started) so you have Terraform installed and `TF_VAR_utho_api_key` set.
* An SSH public key at `~/.ssh/id_ed25519.pub`.
* A domain you control, if you want to use the DNS step.

## Provider and variables

```hcl
terraform {
  required_providers {
    utho = {
      source  = "nitinuthocloud/utho"
      version = "~> 0.7"
    }
  }
}

variable "utho_api_key" {
  type      = string
  sensitive = true
}

provider "utho" {
  api_key = var.utho_api_key
}

variable "dcslug" {
  type    = string
  default = "inmumbaizone2"
}

variable "domain" {
  type    = string
  default = "myapp.com"
}

variable "app_count" {
  type    = number
  default = 3
}
```

## Network

```hcl
resource "utho_vpc" "main" {
  name    = "production"
  network = "10.0.0.0"
  size    = "16"
  dcslug  = var.dcslug
  planid  = "1008"
}

resource "utho_subnet" "public" {
  name            = "public"
  vpc_id          = utho_vpc.main.id
  network         = "10.0.1.0"
  size            = 24
  type            = "public"
  assign_publicip = 1
}
```

## Security group

```hcl
resource "utho_firewall" "web" {
  name = "web-sg"
}

resource "utho_firewall_rule" "http" {
  firewall_id  = utho_firewall.web.id
  type         = "incoming"
  service      = "HTTP"
  protocol     = "tcp"
  port         = "80"
  port_range   = "80"
  addresses    = "0.0.0.0/0"
  source_range = "0.0.0.0/0"
}

resource "utho_firewall_rule" "ssh" {
  firewall_id  = utho_firewall.web.id
  type         = "incoming"
  service      = "SSH"
  protocol     = "tcp"
  port         = "22"
  port_range   = "22"
  addresses    = "203.0.113.0/24" # replace with your office or VPN range
  source_range = "203.0.113.0/24"
}
```

~> **Note:** Restrict SSH to addresses you trust. Opening port 22 to `0.0.0.0/0` exposes every server to brute-force attempts.

## Web servers

```hcl
resource "utho_ssh_key" "deploy" {
  name   = "deploy-key"
  sshkey = file("~/.ssh/id_ed25519.pub")
}

resource "utho_cloud" "app" {
  count = var.app_count

  hostname        = "app-${count.index + 1}.mhc"
  dcslug          = var.dcslug
  planid          = "10360"
  billingcycle    = "hourly"
  auth            = "option2"
  sshkeys         = utho_ssh_key.deploy.id
  image           = "ubuntu-22.04-x86_64"
  enable_publicip = "true"
  vpc             = utho_subnet.public.id
}

resource "utho_firewall_server" "app" {
  count = var.app_count

  firewall_id = utho_firewall.web.id
  cloud_id    = utho_cloud.app[count.index].id
}
```

## Load balancer

```hcl
resource "utho_loadbalancer" "main" {
  name            = "main-lb"
  type            = "application"
  dcslug          = var.dcslug
  enable_publicip = "true"
  vpc             = utho_subnet.public.id
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

resource "utho_loadbalancer_backend" "app" {
  count = var.app_count

  loadbalancer_id = utho_loadbalancer.main.id
  frontend_id     = utho_loadbalancer_frontend.http.id
  backend_port    = "80"
  weight          = "1"
  type            = "cloud"
  cloudid         = utho_cloud.app[count.index].id
}
```

To serve HTTPS, upload a certificate with [`utho_ssl_certificate`](../resources/ssl_certificate), add a second frontend on port `443` with `proto = "https"` and its `certificate_id`, and set `redirecthttps = "1"` on the HTTP frontend.

## Database

```hcl
resource "utho_database" "pg" {
  cluster_name  = "production-db"
  dcslug        = var.dcslug
  engine        = "pg"
  version       = "17"
  size          = "10157"
  network_type  = "private"
  vpc           = utho_subnet.public.id
  billing       = "monthly"
  pitr_enabled  = "1"
  replica_count = "1"

  lifecycle {
    prevent_destroy = true
  }
}

resource "utho_database_db" "app" {
  cluster_id = utho_database.pg.id
  name       = "app"
}
```

`prevent_destroy` makes Terraform refuse any plan that would delete the cluster, which protects you from accidental data loss. Remove it deliberately if you ever need to tear the database down.

## DNS

```hcl
resource "utho_dns_zone" "main" {
  domain = var.domain
}

resource "utho_dns_record" "root" {
  domain   = utho_dns_zone.main.domain
  type     = "A"
  hostname = "@"
  value    = utho_loadbalancer.main.ip
  ttl      = "300"
}

resource "utho_dns_record" "www" {
  domain   = utho_dns_zone.main.domain
  type     = "CNAME"
  hostname = "www"
  value    = "${var.domain}."
  ttl      = "3600"
}
```

After the apply, point your domain's name servers at Utho at your registrar so the zone takes effect.

## Outputs

```hcl
output "lb_ip" {
  value = utho_loadbalancer.main.ip
}

output "app_ips" {
  value = utho_cloud.app[*].ip
}

output "db_host" {
  value = utho_database.pg.host_private
}

output "db_uri" {
  value     = utho_database.pg.uri_private
  sensitive = true
}
```

## Deploy

```bash
terraform init
terraform plan
terraform apply
```

Scale the web tier by changing one value, for example `terraform apply -var="app_count=5"`. Terraform adds the new instances, attaches them to the security group, and registers them with the load balancer in one run.

## Next steps

* Move the web servers into a private subnet behind a NAT gateway. See [Private Networking with VPC and NAT](private-networking).
* Replace the fixed web tier with an [auto scaling group](../resources/autoscaling).
* Add [monitoring alerts](../resources/alert) for CPU and disk usage.
