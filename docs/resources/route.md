---
page_title: "Utho: utho_route"
subcategory: "VPC"
description: |-
  Create and manage individual routes inside a Utho route table.
---

# utho_route

Creates and manages an individual route inside a Utho route table. Routes direct traffic based on destination CIDR blocks to various targets such as internet gateways, subnets, IPSec tunnels, or peering connections.

## Example Usage

### Internet Gateway route

```hcl
resource "utho_route" "internet" {
  route_table_id         = utho_route_table.main.id
  destination_cidr_block = "10.0.0.0/16"
  route_type             = "igw"
  target                 = "igw"
}
```

### Local subnet route

```hcl
resource "utho_route" "internal" {
  route_table_id         = utho_route_table.main.id
  destination_cidr_block = "172.20.0.0/24"
  route_type             = "local"
  target                 = "<subnet-id>"
}
```

## Argument Reference

| Argument                | Type   | Required | Description |
|-------------------------|--------|----------|-------------|
| `route_table_id`        | String | Yes      | Route table ID. Changing forces new resource. |
| `destination_cidr_block`| String | Yes      | Destination CIDR block. Can be updated in place. |
| `route_type`            | String | Yes      | Route type: `igw`, `local`, `ipsec`, `peering`. Can be updated in place. |
| `target`                | String | Yes      | Route target. For `igw` use `igw`, for `local` use subnet ID, for others use resource ID. Can be updated in place. |

## Attribute Reference

| Attribute | Type   | Description |
|-----------|--------|-------------|
| `id`      | String | Route ID (UUID). |

## Notes

- Destroy routes before destroying the route table.
- `destination_cidr_block`, `route_type`, and `target` can all be updated in place.
