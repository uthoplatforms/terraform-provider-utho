---
page_title: "Utho: utho_subnet"
subcategory: "VPC"
description: |-
  Create and manage subnets inside a Utho VPC.
---

# utho_subnet

Creates and manages a subnet inside a Utho VPC. Subnets divide a VPC
network into smaller segments and can be designated as public or private.

## Example Usage

### Public subnet

```hcl
resource "utho_subnet" "public" {
  name            = "public-subnet"
  vpc_id          = utho_vpc.main.id
  network         = "10.0.1.0"
  size            = 24
  type            = "public"
  assign_publicip = 1
  is_default      = 0
}

output "subnet_id" {
  value = utho_subnet.public.id
}

output "gateway" {
  value = utho_subnet.public.gateway
}
```

### Private subnet

```hcl
resource "utho_subnet" "private" {
  name            = "private-subnet"
  vpc_id          = utho_vpc.main.id
  network         = "10.0.2.0"
  size            = 24
  type            = "private"
  assign_publicip = 0
  is_default      = 0
}
```

### Deploy server into a subnet

```hcl
resource "utho_subnet" "public" {
  name            = "public-subnet"
  vpc_id          = utho_vpc.main.id
  network         = "10.0.1.0"
  size            = 24
  type            = "public"
  assign_publicip = 1
}

resource "utho_cloud" "app" {
  hostname = "app-01.mhc"
  dcslug   = "inmumbaizone2"
  planid   = "10360"
  billingcycle = "hourly"
  auth     = "option1"
  root_password = "StrongP@ssw0rd!"
  image    = "ubuntu-22.04-x86_64"
  vpc      = utho_subnet.public.id   # pass subnet ID here
}
```

## Argument Reference

| Argument         | Type   | Required | Description |
|------------------|--------|----------|-------------|
| `name`           | String | Yes      | Subnet name. |
| `vpc_id`         | String | Yes      | VPC ID to create the subnet in. Changing this forces a new resource. |
| `network`        | String | Yes      | Subnet network address (e.g. `10.0.1.0`). Changing this forces a new resource. |
| `size`           | Number | Yes      | CIDR prefix length (e.g. `24` for /24). |
| `type`           | String | Yes      | Subnet type. Accepted values: `public`, `private`. |
| `assign_publicip`| Number | No       | Auto-assign public IP to instances: `1` (yes) or `0` (no). |
| `is_default`     | Number | No       | Mark as default subnet: `1` (yes) or `0` (no). |

## Attribute Reference

| Attribute  | Type   | Description |
|------------|--------|-------------|
| `id`       | String | Unique subnet ID (UUID). |
| `gateway`  | String | Gateway IP address for the subnet. |
| `dcslug`   | String | Data center slug. |
| `status`   | String | Subnet status (e.g. `Active`). |

## Notes

- Changing `vpc_id` or `network` destroys and recreates the subnet.
- Delete all resources (instances, NAT gateways) inside the subnet before deleting it.
- When deploying `utho_cloud` into a subnet, pass the subnet `id` as the `vpc` argument.
