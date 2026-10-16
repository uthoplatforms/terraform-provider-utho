---
page_title: "Utho: utho_firewalls"
subcategory: "Security Groups"
description: |-
  List all security groups in your Utho account.
---

# utho_firewalls

Fetches all security groups (firewalls) in your account. Use this to discover existing security group IDs before attaching them to cloud instances, load balancers, or other resources.

## Example Usage

### List all security groups

```hcl
data "utho_firewalls" "all" {}

output "firewalls" {
  value = data.utho_firewalls.all.firewalls[*].name
}
```

### Find a security group by name and attach to instance

```hcl
data "utho_firewalls" "all" {}

locals {
  web_sg = one([
    for f in data.utho_firewalls.all.firewalls :
    f if f.name == "web-sg"
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
  firewall        = local.web_sg.id
}
```

### List security groups with server counts

```hcl
data "utho_firewalls" "all" {}

output "sg_summary" {
  value = [
    for f in data.utho_firewalls.all.firewalls : {
      name    = f.name
      rules   = f.rule_count
      servers = f.servers_count
    }
  ]
}
```

## Attribute Reference

### Top-level

| Attribute   | Type | Description |
|-------------|------|-------------|
| `firewalls` | List | List of all security groups in your account. |

### firewalls

| Attribute      | Type   | Description |
|----------------|--------|-------------|
| `id`           | String | Security group ID. Use this as `firewall` in `utho_cloud`, `utho_loadbalancer`, etc. |
| `name`         | String | Security group name. |
| `created_at`   | String | Creation timestamp. |
| `rule_count`   | String | Number of inbound and outbound rules. |
| `servers_count`| String | Number of instances currently using this security group. |

## Notes

- To create a new security group use `utho_firewall` and add rules with `utho_firewall_rule`.
- A security group can be attached to multiple instances — changes to rules apply to all attached instances immediately.
- Use `servers_count = "0"` to identify unused security groups.
