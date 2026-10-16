---
page_title: "Utho: utho_dns_zone"
subcategory: "DNS"
description: |-
  Create and manage public DNS zones on Utho.
---

# utho_dns_zone

Creates and manages a public DNS zone on Utho. A DNS zone holds all DNS records for a domain. Once created, add records with `utho_dns_record`, then point your domain's nameservers to Utho at your registrar to activate resolution.

## Example Usage

### Create a zone and add records

```hcl
resource "utho_dns_zone" "main" {
  domain = "myapp.com"
}

# Root domain → load balancer IP
resource "utho_dns_record" "root" {
  domain   = utho_dns_zone.main.domain
  type     = "A"
  hostname = "@"
  value    = utho_loadbalancer.main.ip
  ttl      = "300"
}

# www subdomain
resource "utho_dns_record" "www" {
  domain   = utho_dns_zone.main.domain
  type     = "CNAME"
  hostname = "www"
  value    = "myapp.com."
  ttl      = "3600"
}

# API subdomain → different server
resource "utho_dns_record" "api" {
  domain   = utho_dns_zone.main.domain
  type     = "A"
  hostname = "api"
  value    = utho_cloud.api.ip
  ttl      = "60"
}
```

### After creating — update your registrar

Point your domain's nameservers to Utho at your domain registrar (GoDaddy, Namecheap, Cloudflare, etc.). The nameservers are shown in the Utho Console under DNS after the zone is created.

## Argument Reference

| Argument | Type   | Required | Description |
|----------|--------|----------|-------------|
| `domain` | String | Yes      | Domain name (e.g. `myapp.com`). Do not include a trailing dot. Changing this forces a new resource. |

## Attribute Reference

| Attribute      | Type   | Description |
|----------------|--------|-------------|
| `id`           | String | The domain name (used as identifier). |
| `nspoint`      | String | `YES` if nameservers are pointed to Utho, `NO` otherwise. |
| `record_count` | String | Number of DNS records in this zone. |
| `created_at`   | String | Creation timestamp. |

## Related Resources

| Resource | Purpose |
|----------|---------|
| [utho_dns_record](dns_record) | Add DNS records to this zone |

## Notes

- Creating a zone does not make the domain resolve — you must update your registrar's nameservers.
- DNS propagation after changing nameservers takes 1–48 hours (typically 1–4 hours).
- Use a low TTL (`60`) before changing a record's value — reduces cache time.
- Use a high TTL (`3600`) for stable records — reduces DNS query load.
