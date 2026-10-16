---
page_title: "Utho: utho_firewall_server"
subcategory: "Security Groups"
description: |-
  Attach or detach a cloud instance from a Utho Security Group.
---

# utho_firewall_server

Attaches a cloud instance to a Utho Security Group. Once attached, the
Security Group rules apply to the instance within a few seconds.

Destroying this resource detaches the instance from the Security Group.

## Example Usage

### Attach a server to a Security Group

```hcl
resource "utho_firewall_server" "attach" {
  firewall_id = utho_firewall.web.id
  cloud_id    = utho_cloud.app.id
}
```

### Attach multiple servers to the same Security Group

```hcl
resource "utho_cloud" "servers" {
  count    = 3
  hostname = "server-${count.index + 1}.mhc"
  # ...
}

resource "utho_firewall_server" "attach" {
  count       = 3
  firewall_id = utho_firewall.web.id
  cloud_id    = utho_cloud.servers[count.index].id
}
```

### Full example — Security Group with rules and servers

```hcl
resource "utho_firewall" "web" {
  name = "web-security-group"
}

resource "utho_firewall_rule" "allow_http" {
  firewall_id  = utho_firewall.web.id
  type         = "incoming"
  service      = "HTTP"
  protocol     = "tcp"
  port         = "80"
  port_range   = "80"
  addresses    = "0.0.0.0/0"
  source_range = "0.0.0.0/0"
}

resource "utho_cloud" "app" {
  count    = 2
  hostname = "app-${count.index + 1}.mhc"
  dcslug   = "inmumbaizone2"
  planid   = "10360"
  billingcycle = "hourly"
  auth     = "option1"
  root_password = "StrongP@ssw0rd!"
  image    = "ubuntu-22.04-x86_64"
  enable_publicip = "true"
}

resource "utho_firewall_server" "attach" {
  count       = 2
  firewall_id = utho_firewall.web.id
  cloud_id    = utho_cloud.app[count.index].id
}
```

## Argument Reference

| Argument      | Type   | Required | Description |
|---------------|--------|----------|-------------|
| `firewall_id` | String | Yes      | Security Group ID. Changing this forces a new resource. |
| `cloud_id`    | String | Yes      | Cloud instance ID to attach. Changing this forces a new resource. |

## Attribute Reference

| Attribute | Type   | Description |
|-----------|--------|-------------|
| `id`      | String | Identifier in the format `{firewall_id}:{cloud_id}`. |

## Notes

- Rules apply to the attached instance within a few seconds.
- Destroying this resource detaches the instance — the Security Group and its rules are not deleted.
- Changing `firewall_id` or `cloud_id` destroys and recreates the attachment.