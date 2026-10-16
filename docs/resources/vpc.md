---
page_title: "Utho: utho_vpc"
subcategory: "VPC"
description: |-
  Create and manage Utho Virtual Private Cloud (VPC) networks.
---

# utho_vpc

Creates and manages a Utho VPC — an isolated private network in a single data center. Resources inside a VPC communicate over private IPs and are shielded from the public internet unless you explicitly attach a public IP or NAT gateway.

A complete network stack typically looks like: **VPC → Subnets → NAT Gateway → Route Table → Routes → Instances**.

## Example Usage

### Minimal VPC

```hcl
resource "utho_vpc" "main" {
  name    = "production"
  network = "10.0.0.0"
  size    = "16"
  dcslug  = "inmumbaizone2"
}
```

### VPC with public and private subnets

A common pattern: public subnet for load balancers and bastion hosts, private subnet for app servers and databases.

```hcl
resource "utho_vpc" "main" {
  name    = "production"
  network = "10.0.0.0"
  size    = "16"
  dcslug  = "inmumbaizone2"
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

### Full private network with NAT and routing

Private instances have no public IP but can reach the internet through the NAT gateway. Inbound connections from the internet are blocked by default.

```hcl
resource "utho_vpc" "main" {
  name    = "production"
  network = "10.0.0.0"
  size    = "16"
  dcslug  = "inmumbaizone2"
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

resource "utho_elastic_ip" "nat" {
  dcslug       = "inmumbaizone2"
  billingcycle = "monthly"
}

resource "utho_nat_gateway" "main" {
  name      = "main-nat"
  subnet_id = utho_subnet.public.id
  public_ip = utho_elastic_ip.nat.ip
  dcslug    = "inmumbaizone2"
}

resource "utho_route_table" "private" {
  name   = "private-routes"
  vpc_id = utho_vpc.main.id
  dcslug = "inmumbaizone2"
}

resource "utho_route" "default" {
  route_table_id         = utho_route_table.private.id
  destination_cidr_block = "0.0.0.0/0"
  route_type             = "nat"
  target                 = utho_nat_gateway.main.id
}

# Private instances — no public IP, internet via NAT
resource "utho_cloud" "app" {
  count           = 3
  hostname        = "app-${count.index + 1}.mhc"
  dcslug          = "inmumbaizone2"
  planid          = "10360"
  billingcycle    = "hourly"
  auth            = "option2"
  sshkeys         = utho_ssh_key.deploy.id
  image           = "ubuntu-22.04-x86_64"
  enable_publicip = "false"
  vpc             = utho_subnet.private.id
}
```

## Argument Reference

| Argument  | Type   | Required | Description |
|-----------|--------|----------|-------------|
| `name`    | String | Yes      | VPC name. |
| `network` | String | Yes      | Base network address in CIDR notation (e.g. `10.0.0.0`). Changing this forces a new resource. |
| `size`    | String | Yes      | CIDR prefix length (e.g. `16` for a /16 with 65,536 IPs). Changing this forces a new resource. |
| `dcslug`  | String | Yes      | Data center slug. Changing this forces a new resource. |

## Attribute Reference

| Attribute    | Type   | Description |
|--------------|--------|-------------|
| `id`         | String | Unique VPC ID (UUID). Reference this as `vpc_id` in subnet and other resources. |
| `is_default` | String | `1` if this is the account default VPC, `0` otherwise. |
| `total`      | Number | Total IP addresses in the VPC. |
| `available`  | Number | Available IP addresses remaining. |

## Notes

- VPCs are regional — all resources inside a VPC must be in the same data center.
- Changing `network`, `size`, or `dcslug` destroys and recreates the VPC and everything inside it. Plan carefully.
- Delete all subnets, NAT gateways, and instances inside a VPC before deleting it — Terraform handles this automatically if all resources are managed by the same config.
- The `/16` size (65,536 IPs) is a common choice for production. Use `/20` (4,096 IPs) for smaller environments.
