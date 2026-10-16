---
page_title: "Utho: utho_elastic_ip"
subcategory: "Elastic IP"
description: |-
  Allocate and manage Elastic IPs on Utho Cloud.
---

# utho_elastic_ip

Allocates and manages a static public IP address (Elastic IP) on Utho Cloud. Unlike instance IPs that change when an instance is destroyed and recreated, an Elastic IP is permanent — you keep it until you explicitly release it. Attach it to any cloud instance and move it between instances without changing your DNS records.

## Example Usage

### Allocate an unattached Elastic IP

```hcl
resource "utho_elastic_ip" "main" {
  dcslug       = "inmumbaizone2"
  billingcycle = "monthly"
}

output "static_ip" { value = utho_elastic_ip.main.ip }
```

### Allocate and attach to a cloud instance

```hcl
resource "utho_cloud" "web" {
  hostname        = "web-server.mhc"
  dcslug          = "inmumbaizone2"
  planid          = "10308"
  billingcycle    = "hourly"
  image           = "ubuntu-22.04-x86_64"
  enable_publicip = "true"
  auth            = "option2"
  sshkeys         = utho_ssh_key.deploy.id
}

resource "utho_elastic_ip" "web" {
  dcslug       = "inmumbaizone2"
  billingcycle = "monthly"
  cloud_id     = utho_cloud.web.id
}

output "static_ip"   { value = utho_elastic_ip.web.ip }
output "instance_ip" { value = utho_cloud.web.ip }
```

### Move Elastic IP between instances

Change `cloud_id` to move the IP to a different instance — updates in place without reallocating the IP.

```hcl
# Change cloud_id from old instance to new instance
resource "utho_elastic_ip" "web" {
  dcslug       = "inmumbaizone2"
  billingcycle = "monthly"
  cloud_id     = utho_cloud.web_v2.id  # ← change this and apply
}
```

### Zero-downtime deployment pattern

```hcl
# 1. Keep the elastic IP pointing to v1 while deploying v2
resource "utho_cloud" "web_v1" {
  # ... instance configuration ...
}

resource "utho_cloud" "web_v2" {
  # ... instance configuration ...
}

resource "utho_elastic_ip" "web" {
  dcslug       = "inmumbaizone2"
  billingcycle = "monthly"
  cloud_id     = utho_cloud.web_v1.id  # switch to web_v2.id when ready
}
```

## Argument Reference

| Argument       | Type   | Required | Description |
|----------------|--------|----------|-------------|
| `dcslug`       | String | Yes      | Data center slug. Changing this forces a new resource. |
| `billingcycle` | String | Yes      | Billing cycle: `monthly` or `hourly`. Changing this forces a new resource. |
| `cloud_id`     | String | No       | Cloud instance ID to attach this IP to. Leave empty to keep unattached. Can be updated in place. |

## Attribute Reference

| Attribute     | Type   | Description |
|---------------|--------|-------------|
| `id`          | String | The Elastic IP address (used as identifier). |
| `ip`          | String | The allocated static IP address. |
| `ipid`        | String | Numeric ID of the Elastic IP. |
| `cloud_id`    | String | ID of the attached cloud instance (`"0"` if unattached). |
| `assigned_at` | String | Timestamp when the IP was assigned. |

## Notes

- Elastic IPs are billed even when unattached — release them when not in use.
- One Elastic IP can only be attached to one instance at a time.
- Changing `cloud_id` detaches from the old instance and attaches to the new one in a single apply.
- The IP address remains the same when moved between instances — your DNS records don't need to change.
- Data center must match the instance's data center.
