---
page_title: "Utho: utho_dns_record"
subcategory: "DNS"
description: |-
  Create and manage DNS records inside a Utho DNS zone.
---

# utho_dns_record

Creates and manages a DNS record inside a Utho DNS zone. Records tell the DNS system where to route traffic for a hostname. All record types can be updated in place — no destroy/recreate needed when changing `value`, `ttl`, or `hostname`.

## Example Usage

### A record — map hostname to IP

```hcl
# Root domain → server IP
resource "utho_dns_record" "root" {
  domain   = utho_dns_zone.main.domain
  type     = "A"
  hostname = "@"
  value    = "203.0.113.10"
  ttl      = "300"
}

# Subdomain → different server
resource "utho_dns_record" "api" {
  domain   = utho_dns_zone.main.domain
  type     = "A"
  hostname = "api"
  value    = "203.0.113.20"
  ttl      = "300"
}
```

### AAAA record — IPv6

```hcl
resource "utho_dns_record" "ipv6" {
  domain   = utho_dns_zone.main.domain
  type     = "AAAA"
  hostname = "@"
  value    = "2001:db8::1"
  ttl      = "3600"
}
```

### CNAME record — alias one name to another

~> **Note:** CNAME values must end with a trailing dot (e.g. `myapp.com.`).

```hcl
resource "utho_dns_record" "www" {
  domain   = utho_dns_zone.main.domain
  type     = "CNAME"
  hostname = "www"
  value    = "myapp.com."
  ttl      = "3600"
}
```

### MX record — email routing

```hcl
resource "utho_dns_record" "mx" {
  domain   = utho_dns_zone.main.domain
  type     = "MX"
  hostname = "@"
  value    = "mail.myapp.com."
  ttl      = "3600"
}
```

### TXT record — SPF, DMARC, domain verification

```hcl
# SPF record
resource "utho_dns_record" "spf" {
  domain   = utho_dns_zone.main.domain
  type     = "TXT"
  hostname = "@"
  value    = "v=spf1 include:_spf.google.com ~all"
  ttl      = "3600"
}

# DMARC record
resource "utho_dns_record" "dmarc" {
  domain   = utho_dns_zone.main.domain
  type     = "TXT"
  hostname = "_dmarc"
  value    = "v=DMARC1; p=none; rua=mailto:dmarc@myapp.com"
  ttl      = "3600"
}

# Domain verification
resource "utho_dns_record" "verify" {
  domain   = utho_dns_zone.main.domain
  type     = "TXT"
  hostname = "_verification"
  value    = "verify=abc123"
  ttl      = "3600"
}
```

### Multiple records — complete domain setup

```hcl
resource "utho_dns_zone" "main" {
  domain = "myapp.com"
}

resource "utho_dns_record" "root" {
  domain   = utho_dns_zone.main.domain
  type     = "A"
  hostname = "@"
  value    = "203.0.113.10"
  ttl      = "300"
}

resource "utho_dns_record" "www" {
  domain   = utho_dns_zone.main.domain
  type     = "CNAME"
  hostname = "www"
  value    = "myapp.com."
  ttl      = "3600"
}

resource "utho_dns_record" "mx" {
  domain   = utho_dns_zone.main.domain
  type     = "MX"
  hostname = "@"
  value    = "mail.myapp.com."
  ttl      = "3600"
}

resource "utho_dns_record" "spf" {
  domain   = utho_dns_zone.main.domain
  type     = "TXT"
  hostname = "@"
  value    = "v=spf1 include:_spf.google.com ~all"
  ttl      = "3600"
}
```

## Argument Reference

| Argument   | Type   | Required | Description |
|------------|--------|----------|-------------|
| `domain`   | String | Yes      | Domain name the record belongs to. Must have a matching `utho_dns_zone`. Changing this forces a new resource. |
| `type`     | String | Yes      | Record type: `A`, `AAAA`, `CNAME`, `MX`, `TXT`, `SRV`, `NS`. Changing this forces a new resource. |
| `hostname` | String | Yes      | Subdomain or `@` for the root domain. Can be updated in place. |
| `value`    | String | Yes      | Record value — IP address, hostname, or text. Can be updated in place. |
| `ttl`      | String | Yes      | Time-to-live in seconds. Can be updated in place. |

## Attribute Reference

| Attribute | Type   | Description |
|-----------|--------|-------------|
| `id`      | String | Unique DNS record ID assigned by Utho. |

## TTL Guidelines

| TTL     | When to Use |
|---------|-------------|
| `60`    | During migrations or when you expect to change the value soon |
| `300`   | Default for most records |
| `3600`  | Stable records that rarely change |
| `86400` | Very stable records (nameservers, email) |

## Related Resources

| Resource | Purpose |
|----------|---------|
| [utho_dns_zone](dns_zone) | The DNS zone this record belongs to |

## Notes

- Use `@` as `hostname` for the root domain (e.g. `myapp.com` itself).
- CNAME and MX `value` fields must end with a trailing dot (e.g. `mail.myapp.com.`).
- All fields except `domain` and `type` can be updated in place — no resource replacement needed.
