---
page_title: "Utho: utho_nat_gateway"
subcategory: "VPC"
description: |-
  Create and manage NAT Gateways for Utho VPC subnets.
---

# utho_nat_gateway

Creates and manages a NAT Gateway for a Utho VPC subnet. A NAT Gateway gives instances in a private subnet outbound internet access (software updates, API calls, external services) without exposing them to inbound traffic from the internet.

When this resource is destroyed, Terraform automatically detaches the NAT Gateway from the subnet before deleting it.

## Example Usage

### Basic NAT Gateway with Elastic IP

```hcl
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

resource "utho_elastic_ip" "nat" {
  dcslug       = "inmumbaizone2"
  billingcycle = "monthly"
}

resource "utho_nat_gateway" "main" {
  name      = "main-nat-gateway"
  subnet_id = utho_subnet.public.id
  public_ip = utho_elastic_ip.nat.ip
  dcslug    = "inmumbaizone2"
}

output "nat_gateway_id" { value = utho_nat_gateway.main.id }
output "nat_public_ip"  { value = utho_nat_gateway.main.public_ip }
```

### Full private subnet architecture

```hcl
resource "utho_vpc" "main" {
  name    = "production-vpc"
  network = "10.0.0.0"
  size    = "16"
  dcslug  = "inmumbaizone2"
  planid  = "1008"
}

# Public subnet — NAT gateway and load balancer live here
resource "utho_subnet" "public" {
  name            = "public-subnet"
  vpc_id          = utho_vpc.main.id
  network         = "10.0.1.0"
  size            = 24
  type            = "public"
  assign_publicip = 1
}

# Private subnet — instances live here, no direct internet access
resource "utho_subnet" "private" {
  name            = "private-subnet"
  vpc_id          = utho_vpc.main.id
  network         = "10.0.2.0"
  size            = 24
  type            = "private"
  assign_publicip = 0
}

# Static IP for NAT gateway
resource "utho_elastic_ip" "nat" {
  dcslug       = "inmumbaizone2"
  billingcycle = "monthly"
}

# NAT gateway on public subnet
resource "utho_nat_gateway" "main" {
  name      = "main-nat"
  subnet_id = utho_subnet.public.id
  public_ip = utho_elastic_ip.nat.ip
  dcslug    = "inmumbaizone2"
}

# Private instances route outbound traffic through the NAT gateway
resource "utho_cloud" "worker" {
  hostname        = "worker-01.mhc"
  dcslug          = "inmumbaizone2"
  planid          = "10308"
  billingcycle    = "hourly"
  auth            = "option2"
  sshkeys         = utho_ssh_key.deploy.id
  image           = "ubuntu-22.04-x86_64"
  enable_publicip = "false"
  vpc             = utho_subnet.private.id
}
```

### Import existing NAT gateway

```bash
terraform import utho_nat_gateway.main <nat-gateway-id>
```

## Argument Reference

| Argument    | Type   | Required | Description |
|-------------|--------|----------|-------------|
| `name`      | String | Yes      | NAT Gateway name. Changing this forces a new resource. |
| `subnet_id` | String | Yes      | Public subnet ID to place the NAT Gateway in. Changing this forces a new resource. |
| `public_ip` | String | Yes      | Public Elastic IP address for outbound traffic. Use `utho_elastic_ip.name.ip`. Changing this forces a new resource. |
| `dcslug`    | String | Yes      | Data center slug. Must match the subnet's DC. Changing this forces a new resource. |

## Attribute Reference

| Attribute | Type   | Description |
|-----------|--------|-------------|
| `id`      | String | Unique NAT Gateway ID (UUID). |
| `status`  | String | NAT Gateway status (`active`). |

## Related Resources

| Resource | Purpose |
|----------|---------|
| [utho_vpc](vpc) | The VPC containing the subnets |
| [utho_subnet](subnet) | Public subnet for NAT gateway placement |
| [utho_elastic_ip](elastic_ip) | Static public IP for the NAT gateway |

## Notes

- Place the NAT gateway in a **public subnet** — private instances route through it.
- Use `utho_elastic_ip` for `public_ip` — static IP ensures your outbound address never changes.
- All fields require destroy and recreate on change — NAT gateways are immutable.
- NAT gateway creation takes ~30 seconds.
- The NAT gateway is automatically detached from the subnet before deletion.
- One NAT gateway per subnet — a subnet cannot have two NAT gateways.
