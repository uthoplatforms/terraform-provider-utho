---
page_title: "Utho: utho_ipsec_connection"
subcategory: "VPN (IPSec)"
description: |-
  Pair two Utho IPSec tunnels to establish a site-to-site VPN connection.
---

# utho_ipsec_connection

Pairs two `utho_ipsec` tunnels together to establish a site-to-site VPN connection. Both tunnels must be in `active` status before pairing.

## Example Usage

```hcl
resource "utho_ipsec" "site_a" {
  name         = "site-a-tunnel"
  dcslug       = "inmumbaizone2"
  vpc          = "4cea3ccb-3953-4cc9-a303-2d779bee8645"
  billingcycle = "monthly"
}

resource "utho_ipsec" "site_b" {
  name         = "site-b-tunnel"
  dcslug       = "inmumbaizone2"
  vpc          = "f2f749a5-6ea5-4251-a466-8d09704140ed"
  billingcycle = "monthly"
}

resource "utho_ipsec_connection" "vpn" {
  ipsec_id      = utho_ipsec.site_a.id
  peer_ipsec_id = utho_ipsec.site_b.id
  local_subnet  = "192.168.50.0/24"
  peer_subnet   = "192.168.60.0/24"
  name          = "site-a-to-site-b"
}
```

## Argument Reference

| Argument       | Type   | Required | Description |
|----------------|--------|----------|-------------|
| `ipsec_id`     | String | Yes      | Primary IPSec tunnel ID. Changing forces new resource. |
| `peer_ipsec_id`| String | Yes      | Peer IPSec tunnel ID. Changing forces new resource. |
| `local_subnet` | String | Yes      | Local network CIDR (e.g. `192.168.50.0/24`). Changing forces new resource. |
| `peer_subnet`  | String | Yes      | Peer network CIDR (e.g. `192.168.60.0/24`). Changing forces new resource. |
| `name`         | String | Yes      | Connection name. Changing forces new resource. |

## Attribute Reference

| Attribute    | Type   | Description |
|--------------|--------|-------------|
| `id`         | String | Connection ID on the primary tunnel side. |
| `status`     | String | Connection status (`active`, `inactive`). |
| `created_at` | String | Creation timestamp. |

## Notes

- Both tunnels must be `active` before pairing. Terraform waits for provisioning automatically.
- Destroy the connection before destroying the tunnels.
