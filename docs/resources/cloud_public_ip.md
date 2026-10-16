---
page_title: "Utho: utho_cloud_public_ip"
subcategory: "Cloud Instances"
description: |-
  Assign or release additional public IPs on a Utho Cloud instance.
---

# utho_cloud_public_ip

Assigns an additional public IP address to an existing cloud instance. Use this when an instance needs multiple public IPs — for example, hosting multiple SSL certificates, running multiple services on separate IPs, or IP-based routing.

~> **Note:** This is for assigning *additional* public IPs. The primary public IP is configured via `enable_publicip` in `utho_cloud` at creation time.

## Example Usage

### Assign an additional public IP

```hcl
resource "utho_cloud_public_ip" "extra" {
  cloud_id = utho_cloud.web.id
}

output "extra_ip" {
  value = utho_cloud_public_ip.extra.ip
}
```

### Multiple public IPs for multi-service hosting

```hcl
resource "utho_cloud_public_ip" "service_a" {
  cloud_id = utho_cloud.web.id
}

resource "utho_cloud_public_ip" "service_b" {
  cloud_id = utho_cloud.web.id
}

output "service_ips" {
  value = {
    service_a = utho_cloud_public_ip.service_a.ip
    service_b = utho_cloud_public_ip.service_b.ip
  }
}
```

## Argument Reference

| Argument   | Type   | Required | Description |
|------------|--------|----------|-------------|
| `cloud_id` | String | Yes      | Cloud instance ID. Changing this forces a new resource. |

## Attribute Reference

| Attribute | Type   | Description |
|-----------|--------|-------------|
| `id`      | String | Unique IP assignment ID. |
| `ip`      | String | Assigned public IP address. |

## Notes

- Additional public IPs are billed separately.
- Destroying this resource releases the IP — it will no longer be reachable.
- Configure the additional IP inside the OS using standard networking tools (`ip addr`, `netplan`, etc.).
- For static IPs that persist across instance recreations, use `utho_elastic_ip` instead.
