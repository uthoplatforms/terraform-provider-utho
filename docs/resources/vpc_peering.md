---
page_title: "Utho: utho_vpc_peering"
subcategory: "VPC"
description: |-
  Create and manage VPC peering connections between two Utho VPCs.
---

# utho_vpc_peering

Creates a VPC peering connection between two Utho VPCs in the same data center. Once active, instances in both VPCs can communicate using private IP addresses.

## Example Usage

```hcl
resource "utho_vpc_peering" "main" {
  name             = "app-to-db-peering"
  dcslug           = "inmumbaizone2"
  requester_vpc_id = utho_vpc.app.id
  accepter_vpc_id  = utho_vpc.db.id
}
```

### With existing VPCs

```hcl
resource "utho_vpc_peering" "main" {
  name             = "app-to-db-peering"
  dcslug           = "inmumbaizone2"
  requester_vpc_id = "eec77639-cfe7-4f01-87f9-659c05603230"
  accepter_vpc_id  = "1c689760-e9cf-48e3-b62d-8128b5e6edd7"
}

output "peering_id"     { value = utho_vpc_peering.main.id }
output "peering_status" { value = utho_vpc_peering.main.status }
```

## Argument Reference

| Argument          | Type   | Required | Description |
|-------------------|--------|----------|-------------|
| `name`            | String | Yes      | Peering connection name. Changing forces new resource. |
| `dcslug`          | String | Yes      | Data center slug. Both VPCs must be in the same DC. Changing forces new resource. |
| `requester_vpc_id`| String | Yes      | ID of the requester VPC. Changing forces new resource. |
| `accepter_vpc_id` | String | Yes      | ID of the accepter VPC. Changing forces new resource. |

## Attribute Reference

| Attribute | Type   | Description |
|-----------|--------|-------------|
| `id`      | String | Peering connection ID (UUID). |
| `status`  | String | Peering status (`Active`, `InActive`). |

## Notes

- Both VPCs must be in the same data center.
- VPC peering cannot be updated — destroy and recreate to change any configuration.
- Destroying the peering connection stops private communication between the VPCs.
