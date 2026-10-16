---
page_title: "Utho: utho_ssl_certificate"
subcategory: "SSL Certificates"
description: |-
  Upload and manage SSL certificates on Utho Cloud.
---

# utho_ssl_certificate

Uploads and manages an SSL/TLS certificate on Utho Cloud. Once uploaded, the certificate ID can be referenced in load balancer frontends for HTTPS termination.

## Example Usage

### Upload a custom SSL certificate

```hcl
resource "utho_ssl_certificate" "main" {
  name            = "my-cert"
  type            = "Custom"
  certificate_key = file("cert.pem")
  private_key     = file("key.pem")
}

output "certificate_id" {
  value = utho_ssl_certificate.main.id
}
```

### Use with a load balancer HTTPS frontend

```hcl
resource "utho_ssl_certificate" "main" {
  name            = "my-cert"
  type            = "Custom"
  certificate_key = file("cert.pem")
  private_key     = file("key.pem")
}

resource "utho_loadbalancer_frontend" "https" {
  loadbalancer_id = utho_loadbalancer.main.id
  name            = "https"
  algorithm       = "roundrobin"
  proto           = "https"
  port            = "443"
  cookie          = "0"
  redirecthttps   = "0"
  certificate_id  = utho_ssl_certificate.main.id
}
```

### Upload with certificate chain

```hcl
resource "utho_ssl_certificate" "main" {
  name              = "my-cert"
  type              = "Custom"
  certificate_key   = file("cert.pem")
  private_key       = file("key.pem")
  certificate_chain = file("chain.pem")
}
```

## Argument Reference

| Argument             | Type   | Required | Description |
|----------------------|--------|----------|-------------|
| `name`               | String | Yes      | Certificate name. Changing this forces a new resource. |
| `type`               | String | Yes      | Certificate type: `Custom`. Changing this forces a new resource. |
| `certificate_key`    | String | Yes      | PEM-encoded certificate content. **Sensitive.** Changing this forces a new resource. |
| `private_key`        | String | Yes      | PEM-encoded private key. **Sensitive.** Changing this forces a new resource. |
| `certificate_chain`  | String | No       | PEM-encoded intermediate certificate chain. **Sensitive.** Changing this forces a new resource. |

## Attribute Reference

| Attribute      | Type   | Description |
|----------------|--------|-------------|
| `id`           | String | Unique certificate ID (e.g. `managed-19`). |
| `status`       | String | Certificate status. |
| `primary_domain` | String | Primary domain on the certificate. |
| `is_wildcard`  | String | Whether this is a wildcard certificate. |
| `auto_renew`   | String | Whether auto-renewal is enabled. |
| `issuer`       | String | Certificate issuer. |
| `created_at`   | String | Creation timestamp. |
| `expire_at`    | String | Expiry timestamp. |
