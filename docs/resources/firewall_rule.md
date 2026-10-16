---
page_title: "Utho: utho_firewall_rule"
subcategory: "Security Groups"
description: |-
  Add and manage rules in a Utho Security Group.
---

# utho_firewall_rule

Adds an inbound or outbound rule to a Utho Security Group. Rules control
what traffic is allowed to reach or leave the attached cloud instances.

## Example Usage

### Allow SSH from anywhere

```hcl
resource "utho_firewall_rule" "allow_ssh" {
  firewall_id  = utho_firewall.web.id
  type         = "incoming"
  service      = "SSH"
  protocol     = "tcp"
  port         = "22"
  port_range   = "22"
  addresses    = "0.0.0.0/0"
  source_range = "0.0.0.0/0"
}
```

### Allow HTTP and HTTPS

```hcl
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

resource "utho_firewall_rule" "allow_https" {
  firewall_id  = utho_firewall.web.id
  type         = "incoming"
  service      = "HTTPS"
  protocol     = "tcp"
  port         = "443"
  port_range   = "443"
  addresses    = "0.0.0.0/0"
  source_range = "0.0.0.0/0"
}
```

### Allow custom port range

```hcl
resource "utho_firewall_rule" "allow_app" {
  firewall_id  = utho_firewall.web.id
  type         = "incoming"
  service      = "CUSTOM"
  protocol     = "tcp"
  port         = "8000"
  port_range   = "8000-9000"
  addresses    = "10.0.0.0/8"
  source_range = "10.0.0.0/8"
}
```

### Allow all outbound traffic

```hcl
resource "utho_firewall_rule" "allow_all_out_tcp" {
  firewall_id  = utho_firewall.web.id
  type         = "outgoing"
  service      = "ALL TCP"
  protocol     = "tcp"
  port         = "ALL"
  port_range   = "ALL"
  addresses    = "0.0.0.0/0"
  source_range = "0.0.0.0/0"
}

resource "utho_firewall_rule" "allow_all_out_udp" {
  firewall_id  = utho_firewall.web.id
  type         = "outgoing"
  service      = "ALL UDP"
  protocol     = "udp"
  port         = "ALL"
  port_range   = "ALL"
  addresses    = "0.0.0.0/0"
  source_range = "0.0.0.0/0"
}
```

### Restrict SSH to private network only

```hcl
resource "utho_firewall_rule" "ssh_private" {
  firewall_id  = utho_firewall.web.id
  type         = "incoming"
  service      = "SSH"
  protocol     = "tcp"
  port         = "22"
  port_range   = "22"
  addresses    = "10.0.0.0/8"
  source_range = "10.0.0.0/8"
}
```

## Argument Reference

| Argument      | Type   | Required | Description |
|---------------|--------|----------|-------------|
| `firewall_id` | String | Yes      | Security Group ID to add this rule to. Changing this forces a new resource. |
| `type`        | String | Yes      | Rule direction. Accepted values: `incoming`, `outgoing`. Changing this forces a new resource. |
| `service`     | String | Yes      | Service name (e.g. `SSH`, `HTTP`, `HTTPS`, `CUSTOM`, `ALL TCP`, `ALL UDP`). Changing this forces a new resource. |
| `protocol`    | String | Yes      | Protocol. Accepted values: `tcp`, `udp`, `icmp`. Changing this forces a new resource. |
| `port`        | String | Yes      | Port number (e.g. `22`, `80`) or `ALL`. Changing this forces a new resource. |
| `port_range`  | String | Yes      | Port range (e.g. `8000-9000`) or same as `port` for a single port. Changing this forces a new resource. |
| `addresses`   | String | Yes      | Destination CIDR block (e.g. `0.0.0.0/0`, `10.0.0.0/8`). Changing this forces a new resource. |
| `source_range`| String | Yes      | Source CIDR block (e.g. `0.0.0.0/0`, `10.0.0.0/8`). Changing this forces a new resource. |

## Attribute Reference

| Attribute | Type   | Description |
|-----------|--------|-------------|
| `id`      | String | Unique rule ID assigned by Utho. |

## Notes

- All fields require destroy and recreate on change.
- To remove a rule, delete the resource block and run `terraform apply`.
- Multiple rules can be added to the same Security Group.