---
page_title: "Utho: utho_firewall"
subcategory: "Security Groups"
description: |-
  Create and manage Utho Security Groups (firewall rules for cloud instances).
---

# utho_firewall

Creates and manages a Utho Security Group — a stateful firewall that controls what traffic can reach your cloud instances. Security groups are attached to instances and define which ports and protocols are allowed in and out.

A security group on its own does nothing. Add rules with `utho_firewall_rule` and attach it to instances with `utho_firewall_server`.

## Example Usage

### Web server security group

Allow HTTP and HTTPS from anywhere, SSH only from a trusted IP range.

```hcl
resource "utho_firewall" "web" {
  name = "web-sg"
}

# Allow HTTP from anywhere
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

# Allow HTTPS from anywhere
resource "utho_firewall_rule" "https" {
  firewall_id  = utho_firewall.web.id
  type         = "incoming"
  service      = "HTTPS"
  protocol     = "tcp"
  port         = "443"
  port_range   = "443"
  addresses    = "0.0.0.0/0"
  source_range = "0.0.0.0/0"
}

# Allow SSH only from your office IP
resource "utho_firewall_rule" "ssh" {
  firewall_id  = utho_firewall.web.id
  type         = "incoming"
  service      = "SSH"
  protocol     = "tcp"
  port         = "22"
  port_range   = "22"
  addresses    = "203.0.113.0/24"
  source_range = "203.0.113.0/24"
}

# Allow all outbound traffic
resource "utho_firewall_rule" "outbound" {
  firewall_id  = utho_firewall.web.id
  type         = "outgoing"
  service      = "ALL TCP"
  protocol     = "tcp"
  port         = "ALL"
  port_range   = "ALL"
  addresses    = "0.0.0.0/0"
  source_range = "0.0.0.0/0"
}

# Attach to all web servers
resource "utho_firewall_server" "web" {
  count       = length(utho_cloud.web)
  firewall_id = utho_firewall.web.id
  cloud_id    = utho_cloud.web[count.index].id
}
```

### Database security group — private only

Only allow connections from within your VPC, block all public access.

```hcl
resource "utho_firewall" "db" {
  name = "database-sg"
}

resource "utho_firewall_rule" "postgres" {
  firewall_id  = utho_firewall.db.id
  type         = "incoming"
  service      = "CUSTOM"
  protocol     = "tcp"
  port         = "5432"
  port_range   = "5432"
  addresses    = "10.0.0.0/8"   # your VPC CIDR
  source_range = "10.0.0.0/8"
}
```

## Argument Reference

| Argument | Type   | Required | Description |
|----------|--------|----------|-------------|
| `name`   | String | Yes      | Security group name. Can be updated in place. |

## Attribute Reference

| Attribute       | Type   | Description |
|-----------------|--------|-------------|
| `id`            | String | Unique Security Group ID. Use this in `utho_firewall_rule` and `utho_firewall_server`. |
| `created_at`    | String | Creation timestamp. |
| `rule_count`    | String | Number of rules attached to this security group. |
| `servers_count` | String | Number of instances currently using this security group. |

## Notes

- Security groups are stateful — if you allow inbound port 80, the return traffic is automatically allowed.
- Rules apply within a few seconds after attaching to an instance.
- One instance can have multiple security groups via multiple `utho_firewall_server` resources.
- A security group can be attached to multiple instances — useful for shared rules across a fleet.
