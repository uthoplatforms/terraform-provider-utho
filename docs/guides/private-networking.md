---
page_title: "Private Networking with VPC and NAT"
subcategory: "Tutorials"
description: |-
  Isolate workloads in a private subnet with outbound internet access through a NAT gateway.
---

# Private Networking with VPC and NAT

In production, application servers and databases usually should not be reachable from the internet. This tutorial builds a VPC with two subnets:

* A **public subnet** for internet-facing components, such as load balancers and bastion hosts, and for the NAT gateway.
* A **private subnet** for application servers that have no public IP but can still reach the internet for updates and API calls through the NAT gateway.

```text
               Internet
                  ▲
                  │ outbound only
         ┌────────┴─────────┐
         │ utho_nat_gateway │ ◄── utho_elastic_ip
         └────────┬─────────┘
 ┌────────────────┼─────────────────────────────┐
 │ VPC 10.0.0.0/16│                             │
 │  ┌─────────────┴──────┐   ┌────────────────┐ │
 │  │ public 10.0.1.0/24 │   │private 10.0.2.0│ │
 │  │  NAT, LB, bastion  │   │  app servers   │ │
 │  └────────────────────┘   └────────────────┘ │
 └──────────────────────────────────────────────┘
```

## Prerequisites

* Complete [Getting Started](getting-started).

## Configure the provider

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
```

## Create the VPC and subnets

Utho takes each address range as two separate arguments: the network base address (`network`) and the prefix length (`size`). The VPC below covers `10.0.0.0/16`, with a public subnet at `10.0.1.0/24` and a private subnet at `10.0.2.0/24`.

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

resource "utho_subnet" "private" {
  name            = "private"
  vpc_id          = utho_vpc.main.id
  network         = "10.0.2.0"
  size            = 24
  type            = "private"
  assign_publicip = 0
}
```

## Add a NAT gateway

The NAT gateway lives in the public subnet and uses an Elastic IP as its fixed outbound address. Because this address does not change, you can give it to third parties that need to allowlist your traffic.

```hcl
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
```

## Launch private instances

Instances in the private subnet get no public IP. They are reachable only from inside the VPC.

```hcl
resource "utho_ssh_key" "deploy" {
  name   = "deploy-key"
  sshkey = file("~/.ssh/id_ed25519.pub")
}

resource "utho_cloud" "app" {
  count = 2

  hostname        = "app-${count.index + 1}.mhc"
  dcslug          = var.dcslug
  planid          = "10308"
  billingcycle    = "hourly"
  auth            = "option2"
  sshkeys         = utho_ssh_key.deploy.id
  image           = "ubuntu-22.04-x86_64"
  enable_publicip = "false"
  vpc             = utho_subnet.private.id

  depends_on = [utho_nat_gateway.main]
}
```

The `depends_on` makes sure the NAT gateway exists before the instances boot, so their first-boot package updates can reach the internet.

## Allow traffic only from inside the VPC

Security group rules take a CIDR block. Build one from a subnet's `network` and `size` attributes with string interpolation:

```hcl
resource "utho_firewall" "app" {
  name = "app-sg"
}

resource "utho_firewall_rule" "app_from_public" {
  firewall_id  = utho_firewall.app.id
  type         = "incoming"
  service      = "CUSTOM"
  protocol     = "tcp"
  port         = "8080"
  port_range   = "8080"
  addresses    = "${utho_subnet.public.network}/${utho_subnet.public.size}"
  source_range = "${utho_subnet.public.network}/${utho_subnet.public.size}"
}

resource "utho_firewall_server" "app" {
  count = 2

  firewall_id = utho_firewall.app.id
  cloud_id    = utho_cloud.app[count.index].id
}
```

## Reaching private instances

Because private instances have no public IP, you need a path in for administration. Common options are:

* A small **bastion host** in the public subnet with SSH restricted to your office or VPN range, used as a jump host (`ssh -J root@<bastion-ip> root@<private-ip>`).
* A **site-to-site VPN** from your office network with [`utho_ipsec`](../resources/ipsec) and [`utho_ipsec_connection`](../resources/ipsec_connection).
* **VPC peering** with another VPC that already has access, using [`utho_vpc_peering`](../resources/vpc_peering).

## Outputs

```hcl
output "nat_public_ip" {
  description = "Outbound address of all private instances"
  value       = utho_elastic_ip.nat.ip
}

output "private_subnet_id" {
  value = utho_subnet.private.id
}
```

## Next steps

* Put a load balancer in the public subnet in front of the private instances. See [Deploying a Production Web Stack](production-web-stack).
* Run a [managed database](managed-database) in the same VPC.
* Use the private subnet for a [private Kubernetes cluster](kubernetes).
