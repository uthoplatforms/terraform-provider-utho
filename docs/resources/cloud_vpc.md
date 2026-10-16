---
page_title: "Utho: utho_cloud_vpc"
subcategory: "Cloud Instances"
description: |-
  Attach or detach a VPC subnet from a Utho Cloud instance.
---

# utho_cloud_vpc

Attaches or detaches a VPC subnet from an existing cloud instance. Use this to add an instance to a private network post-deployment without recreating it.

~> **Note:** You can also attach a VPC subnet at creation time using the `vpc` argument in `utho_cloud`. Use `utho_cloud_vpc` to change the VPC attachment on an already-running instance.

## Example Usage

### Attach instance to a VPC subnet

```hcl
resource "utho_cloud_vpc" "attach" {
  cloud_id  = utho_cloud.web.id
  subnet_id = utho_subnet.private.id
}
```

### Use existing subnet from data source

```hcl
data "utho_vpcs" "mumbai" {
  dcslug = "inmumbaizone2"
}

locals {
  private_subnet = one(flatten([
    for vpc in data.utho_vpcs.mumbai.vpcs : [
      for s in vpc.subnets :
      s if s.subnet_type == "private" && vpc.name == "production"
    ]
  ]))
}

resource "utho_cloud_vpc" "attach" {
  cloud_id  = utho_cloud.backend.id
  subnet_id = local.private_subnet.id
}
```

### Full example — create VPC, subnet, then attach

```hcl
resource "utho_vpc" "prod" {
  name    = "production"
  network = "10.0.0.0"
  size    = "24"
  dcslug  = "inmumbaizone2"
  planid  = "1008"
}

resource "utho_subnet" "private" {
  name            = "private-subnet"
  vpc_id          = utho_vpc.prod.id
  network         = "10.0.0.0"
  size            = 24
  type            = "private"
  assign_publicip = 0
}

resource "utho_cloud_vpc" "attach" {
  cloud_id  = utho_cloud.backend.id
  subnet_id = utho_subnet.private.id
}

output "private_ip" {
  value = utho_cloud_vpc.attach.private_ip
}
```

## Argument Reference

| Argument    | Type   | Required | Description |
|-------------|--------|----------|-------------|
| `cloud_id`  | String | Yes      | Cloud instance ID. Changing this forces a new resource. |
| `subnet_id` | String | Yes      | Subnet ID to attach. Use [utho_vpcs](../data-sources/vpcs) or [utho_vpc_subnets](../data-sources/vpc_subnets) to find valid subnet IDs. |

## Attribute Reference

| Attribute    | Type   | Description |
|--------------|--------|-------------|
| `id`         | String | Attachment ID (`cloud_id:subnet_id`). |
| `private_ip` | String | Private IP assigned to the instance within the subnet. |

## Notes

- Use the subnet's ID from `utho_subnet.name.id` or from `data.utho_vpcs`.
- The instance and subnet must be in the same data center.
- Destroying this resource detaches the instance from the VPC subnet.
- Use `data.utho_vpcs` or `data.utho_vpc_subnets` to discover existing subnet IDs.
