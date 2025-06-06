---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "utho_dns_record Resource - utho"
subcategory: ""
description: |-
  
---

# utho_dns_record (Resource)



## Example Usage

```terraform
resource "utho_domain" "example" {
  domain = "example.com"
}

resource "utho_dns_record" "example" {
  domain   = utho_domain.example.domain
  type     = "A"
  hostname = "subdomain.${utho_domain.example.domain}"
  value    = "1.1.1.1"
  ttl      = "65444"
  porttype = "TCP"
  port     = "5060"
  priority = "10"
  weight   = "100"
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `domain` (String) Name of the domain
- `hostname` (String) Name (Hostname) The host name, alias, or service being defined by the record.
- `porttype` (String) This value is the time to live for the record, in seconds. This defines the time frame that
- `ttl` (String) The priority of the host (for SRV and MX records. null otherwise).
- `type` (String) The Record Type (A, AAAA, CAA, CNAME, MX, TXT, SRV, NS)
- `value` (String) Variable data depending on record type. For example, the value for an A record would be the IPv4 address to which the domain will be mapped. For a CAA record, it would contain the domain name of the CA being granted permission to issue certificates.

### Optional

- `port` (String) The port that the service is accessible on (for SRV records only. null otherwise).
- `priority` (String) priority
- `weight` (String) The weight of records with the same priority (for SRV records only. null otherwise).

### Read-Only

- `id` (String) id
