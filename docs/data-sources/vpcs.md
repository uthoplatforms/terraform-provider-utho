---
page_title: "Utho: utho_vpcs"
subcategory: "VPC"
description: |-
  List all VPCs and their subnets in your Utho account.
---

# utho_vpcs

Fetches all VPCs in your account with their subnets. Use this to discover existing VPC and subnet IDs before attaching cloud instances, load balancers, databases, or other resources to a private network.

## Example Usage

### List all VPCs

```hcl
data "utho_vpcs" "all" {}

output "vpcs" {
  value = data.utho_vpcs.all.vpcs[*].name
}
```

### List VPCs in a specific DC

```hcl
data "utho_vpcs" "mumbai" {
  dcslug = "inmumbaizone2"
}

output "mumbai_vpcs" {
  value = data.utho_vpcs.mumbai.vpcs
}
```

### Find a subnet to attach an instance to

```hcl
data "utho_vpcs" "mumbai" {
  dcslug = "inmumbaizone2"
}

locals {
  # Find VPC by name
  my_vpc = one([
    for v in data.utho_vpcs.mumbai.vpcs :
    v if v.name == "production"
  ])

  # Get the public subnet from that VPC
  public_subnet = one([
    for s in local.my_vpc.subnets :
    s if s.subnet_type == "public"
  ])
}

resource "utho_cloud" "web" {
  hostname        = "web-01.mhc"
  dcslug          = "inmumbaizone2"
  planid          = "10308"
  billingcycle    = "hourly"
  auth            = "option1"
  root_password   = var.root_password
  image           = "ubuntu-22.04-x86_64"
  enable_publicip = "true"
  vpc             = local.public_subnet.id   # ← subnet numeric ID
}
```

### List all subnets across VPCs in a DC

```hcl
data "utho_vpcs" "mumbai" {
  dcslug = "inmumbaizone2"
}

output "all_subnets" {
  value = flatten([
    for vpc in data.utho_vpcs.mumbai.vpcs : [
      for subnet in vpc.subnets : {
        vpc_name    = vpc.name
        subnet_id   = subnet.id
        subnet_name = subnet.name
        type        = subnet.subnet_type
        network     = "${subnet.network}/${subnet.size}"
      }
    ]
  ])
}
```

### Find private subnets only

```hcl
data "utho_vpcs" "mumbai" {
  dcslug = "inmumbaizone2"
}

output "private_subnets" {
  value = flatten([
    for vpc in data.utho_vpcs.mumbai.vpcs : [
      for s in vpc.subnets :
      { id = s.id, name = s.name, vpc = vpc.name }
      if s.subnet_type == "private"
    ]
  ])
}
```

## Argument Reference

| Argument | Type   | Required | Description |
|----------|--------|----------|-------------|
| `dcslug` | String | No       | Filter VPCs by data center slug. Recommended when you know which DC you're deploying to. |

## Attribute Reference

### Top-level

| Attribute | Type | Description |
|-----------|------|-------------|
| `vpcs`    | List | List of all VPCs, each with their subnets. |

### vpcs

| Attribute   | Type   | Description |
|-------------|--------|-------------|
| `id`        | String | VPC UUID. |
| `name`      | String | VPC name. |
| `network`   | String | VPC network base address (e.g. `10.0.0.0`). |
| `size`      | String | Network prefix length (e.g. `16`). |
| `dcslug`    | String | Data center the VPC belongs to. |
| `total`     | Number | Total IP addresses in the VPC. |
| `available` | Number | Available IP addresses remaining. |
| `subnets`   | List   | Subnets within the VPC. See [subnets](#subnets). |

### subnets

| Attribute        | Type   | Description |
|------------------|--------|-------------|
| `id`             | String | Subnet numeric ID. **Use this as `vpc` in `utho_cloud`, `utho_loadbalancer`, `utho_database`, etc.** |
| `uuid`           | String | Subnet UUID. |
| `name`           | String | Subnet name. |
| `network`        | String | Subnet base network address. |
| `size`           | String | Subnet prefix length. |
| `subnet_type`    | String | `public` (instances get public IPs) or `private` (internal only). |
| `assign_publicip`| String | `"1"` if public IPs are auto-assigned, `"0"` if not. |
| `status`         | String | Subnet status (`Active`). |
| `gateway`        | String | Subnet gateway IP address. |
| `dcslug`         | String | Data center slug. |

## Notes

- Use `subnet.id` (the numeric ID) — not `subnet.uuid` — when attaching resources to a subnet.
- Subnets with `subnet_type = "public"` route traffic to the internet directly.
- Subnets with `subnet_type = "private"` require a NAT gateway (`utho_nat_gateway`) for outbound internet access.
- VPCs without subnets show an empty `subnets = []` list — create subnets with `utho_subnet`.
