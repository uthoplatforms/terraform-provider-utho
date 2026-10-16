---
page_title: "Utho: utho_route_table"
subcategory: "VPC"
description: |-
  Create and manage route tables for Utho VPCs.
---

# utho_route_table

Creates and manages a route table for a Utho VPC. Route tables control how network traffic is directed within and outside your VPC. Use `utho_route` to add individual routes to the table.

## Example Usage

```hcl
resource "utho_route_table" "main" {
  name   = "prod-route-table"
  vpc_id = utho_vpc.main.id
  dcslug = "inmumbaizone2"
}
```

## Argument Reference

| Argument | Type   | Required | Description |
|----------|--------|----------|-------------|
| `name`   | String | Yes      | Route table name. |
| `vpc_id` | String | Yes      | VPC ID to associate the route table with. Changing forces new resource. |
| `dcslug` | String | Yes      | Data center slug. Changing forces new resource. |

## Attribute Reference

| Attribute | Type   | Description |
|-----------|--------|-------------|
| `id`      | String | Route table ID (UUID). |

## Related Resources

- [utho_route](route) — Add routes to this route table
- [utho_vpc](vpc) — VPC that owns this route table
