---
page_title: "Utho: utho_clouds"
subcategory: "Cloud Instances"
description: |-
  List all cloud instances in your Utho account.
---

# utho_clouds

Fetches all cloud instances in your account. Use this to discover existing instance IDs, IPs, and status without importing them into Terraform state.

## Example Usage

### List all instances

```hcl
data "utho_clouds" "all" {}

output "instances" {
  value = [
    for c in data.utho_clouds.all.clouds :
    { id = c.id, hostname = c.hostname, ip = c.ip, status = c.status }
  ]
}
```

### Find instance by hostname

```hcl
data "utho_clouds" "all" {}

locals {
  web = one([
    for c in data.utho_clouds.all.clouds :
    c if c.hostname == "web-01.mhc"
  ])
}

output "web_ip" { value = local.web.ip }
```

## Attribute Reference

| Attribute      | Type   | Description |
|----------------|--------|-------------|
| `id`           | String | Instance ID. |
| `hostname`     | String | Instance hostname. |
| `ip`           | String | Primary public IP address. |
| `status`       | String | Instance status (e.g. Active). |
| `power_status` | String | Power state (Running, Shutdown). |
| `dcslug`       | String | Data center slug. |
| `created_at`   | String | Creation timestamp. |
