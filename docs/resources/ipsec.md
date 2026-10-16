---
page_title: "Utho: utho_ipsec"
subcategory: "VPN (IPSec)"
description: |-
  Create and manage Utho IPSec site-to-site VPN tunnels.
---

# utho_ipsec

Creates and manages an IPSec VPN tunnel on Utho. Each tunnel provisions a dedicated gateway instance in your VPC subnet. Once active, two tunnels can be paired using `utho_ipsec_connection` to establish a site-to-site VPN.

~> **Note:** IPSec tunnel provisioning takes 5-10 minutes. Terraform polls until the tunnel is active before proceeding.

## Example Usage

```hcl
resource "utho_ipsec" "tunnel" {
  name         = "prod-vpn-tunnel"
  dcslug       = "inmumbaizone2"
  vpc          = "4cea3ccb-3953-4cc9-a303-2d779bee8645"
  billingcycle = "monthly"
}

output "tunnel_id" { value = utho_ipsec.tunnel.id }
output "tunnel_ip" { value = utho_ipsec.tunnel.cloudid }
output "tunnel_psk" {
  value     = utho_ipsec.tunnel.psk
  sensitive = true
}
```

## Argument Reference

| Argument       | Type   | Required | Description |
|----------------|--------|----------|-------------|
| `name`         | String | Yes      | Tunnel name. Changing forces new resource. |
| `dcslug`       | String | Yes      | Data center slug (e.g. `inmumbaizone2`). Changing forces new resource. |
| `vpc`          | String | Yes      | Subnet ID to deploy the tunnel gateway in. If the VPC has subnets, provide the subnet ID; otherwise provide the VPC ID. Changing forces new resource. |
| `billingcycle` | String | Yes      | Billing cycle: `monthly`, `3month`, `6month`, `12month`. Changing forces new resource. |

## Attribute Reference

| Attribute    | Type   | Description |
|--------------|--------|-------------|
| `id`         | String | IPSec tunnel ID. |
| `psk`        | String | Auto-generated pre-shared key (sensitive). |
| `status`     | String | Tunnel status (`pending`, `active`). |
| `cloudid`    | String | Cloud instance ID of the gateway. |
| `created_at` | String | Creation timestamp. |

## Related Resources

- [utho_ipsec_connection](ipsec_connection) — Pair two tunnels to establish site-to-site VPN
