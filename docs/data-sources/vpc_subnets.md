---
page_title: "Utho: utho_vpc_subnets"
subcategory: "VPC"
description: |-
  List all subnets within a specific Utho VPC.
---

# utho_vpc_subnets

Fetches all subnets within a specific VPC. Use this when you know the VPC ID and need to discover its subnets before attaching resources. Unlike `utho_vpcs` which lists all VPCs, this data source gives you detailed subnet information for one specific VPC.

## Example Usage

### List subnets of a VPC

```hcl
data "utho_vpc_subnets" "prod" {
  vpc_id = "c9db1f6f-8018-406f-a495-a7d12ca2142f"
}

output "subnets" {
  value = data.utho_vpc_subnets.prod.subnets
}
```

### Find the public subnet and deploy instance

```hcl
data "utho_vpc_subnets" "prod" {
  vpc_id = "c9db1f6f-8018-406f-a495-a7d12ca2142f"
}

locals {
  public_subnet = one([
    for s in data.utho_vpc_subnets.prod.subnets :
    s if s.subnet_type == "public"
  ])
}

resource "utho_cloud" "web" {
  hostname        = "web-01.mhc"
  dcslug          = "inmumbaizone2"
  planid          = "10308"
  billingcycle    = "hourly"
  auth            = "option2"
  sshkeys         = utho_ssh_key.deploy.id
  image           = "ubuntu-22.04-x86_64"
  enable_publicip = "true"
  cpumodel        = "amd"
  vpc             = local.public_subnet.id
}
```

### Chain with utho_vpcs to find VPC by name first

```hcl
# Step 1 — find the VPC by name
data "utho_vpcs" "mumbai" {
  dcslug = "inmumbaizone2"
}

locals {
  prod_vpc = one([
    for v in data.utho_vpcs.mumbai.vpcs :
    v if v.name == "production"
  ])
}

# Step 2 — get its subnets
data "utho_vpc_subnets" "prod" {
  vpc_id = local.prod_vpc.id
}

# Step 3 — pick private subnet for backend
locals {
  private_subnet = one([
    for s in data.utho_vpc_subnets.prod.subnets :
    s if s.subnet_type == "private"
  ])
}

resource "utho_cloud" "backend" {
  hostname        = "backend-01.mhc"
  dcslug          = "inmumbaizone2"
  planid          = "10308"
  billingcycle    = "hourly"
  auth            = "option2"
  sshkeys         = utho_ssh_key.deploy.id
  image           = "ubuntu-22.04-x86_64"
  enable_publicip = "false"
  cpumodel        = "amd"
  vpc             = local.private_subnet.id
}
```

## Argument Reference

| Argument | Type   | Required | Description |
|----------|--------|----------|-------------|
| `vpc_id` | String | Yes      | VPC UUID. Get this from `data.utho_vpcs` or `utho_vpc.name.id`. |

## Attribute Reference

### Top-level

| Attribute | Type | Description |
|-----------|------|-------------|
| `subnets` | List | List of all subnets in the VPC. |

### subnets

| Attribute        | Type   | Description |
|------------------|--------|-------------|
| `id`             | String | Subnet numeric ID. **Use this as `vpc` in `utho_cloud`, `utho_loadbalancer`, `utho_database`, etc.** |
| `uuid`           | String | Subnet UUID. |
| `name`           | String | Subnet name. |
| `network`        | String | Subnet base network address. |
| `size`           | String | Subnet prefix length (e.g. `24` for /24). |
| `subnet_type`    | String | `public` (internet-routable) or `private` (internal only). |
| `assign_publicip`| String | `1` if public IPs are auto-assigned, `0` if not. |
| `status`         | String | Subnet status (`Active`). |
| `gateway`        | String | Subnet gateway IP address. |
| `dcslug`         | String | Data center the subnet belongs to. |

## Notes

- Use `subnet.id` (the numeric ID) — not `subnet.uuid` — when attaching resources.
- Public subnets (`subnet_type = "public"`) allow instances to have public IPs.
- Private subnets (`subnet_type = "private"`) require a NAT gateway for outbound internet access.
- Get the VPC ID from `data.utho_vpcs` or from the `id` attribute of an `utho_vpc` resource.
