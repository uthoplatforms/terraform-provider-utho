---
page_title: "Utho: utho_cloud_firewall"
subcategory: "Cloud Instances"
description: |-
  Attach or detach a security group from a Utho Cloud instance.
---

# utho_cloud_firewall

Attaches or detaches a security group (firewall) from an existing cloud instance. Use this to manage security group attachments post-deployment without recreating the instance.

~> **Note:** You can also attach a security group at creation time using the `firewall` argument in `utho_cloud`. Use `utho_cloud_firewall` to change the security group on already-running instances.

## Example Usage

### Attach a security group to an instance

```hcl
resource "utho_firewall" "web" {
  name = "web-sg"
}

resource "utho_firewall_rule" "http" {
  firewall_id  = utho_firewall.web.id
  type         = "incoming"
  service      = "HTTP"
  protocol     = "tcp"
  port         = "80"
  port_range   = "80"
  addresses    = "0.0.0.0/0"
  source_range = "0.0.0.0/0"
}

resource "utho_cloud_firewall" "web_attach" {
  cloud_id    = utho_cloud.web.id
  firewall_id = utho_firewall.web.id
}
```

### Use existing security group from data source

```hcl
data "utho_firewalls" "all" {}

locals {
  web_sg = one([
    for f in data.utho_firewalls.all.firewalls :
    f if f.name == "web-sg"
  ])
}

resource "utho_cloud_firewall" "attach" {
  cloud_id    = utho_cloud.web.id
  firewall_id = local.web_sg.id
}
```

### Attach same security group to multiple instances

```hcl
resource "utho_cloud_firewall" "web_sg" {
  count = length(utho_cloud.app)

  cloud_id    = utho_cloud.app[count.index].id
  firewall_id = utho_firewall.web.id
}
```

## Argument Reference

| Argument      | Type   | Required | Description |
|---------------|--------|----------|-------------|
| `cloud_id`    | String | Yes      | Cloud instance ID. Changing this forces a new resource. |
| `firewall_id` | String | Yes      | Security group ID. Use [utho_firewalls](../data-sources/firewalls) to discover existing IDs. |

## Attribute Reference

| Attribute | Type   | Description |
|-----------|--------|-------------|
| `id`      | String | Attachment ID. |

## Notes

- Security group rules take effect immediately after attachment.
- One instance can have multiple security groups — create multiple `utho_cloud_firewall` resources pointing to the same `cloud_id`.
- Destroying this resource detaches the security group from the instance but does NOT delete the security group itself.
- Use `data.utho_firewalls` to list existing security groups and find their IDs.
