---
page_title: "Utho: utho_loadbalancer_frontend"
subcategory: "Load Balancing"
description: |-
  Add and manage frontends on a Utho Load Balancer.
---

# utho_loadbalancer_frontend

Adds a frontend listener to a Utho Load Balancer. The frontend defines the protocol, port, and algorithm used to distribute traffic to backends.

~> **Note:** The load balancer must be fully ready before a frontend can be added. The provider automatically polls and waits up to 5 minutes. If it times out, run `terraform apply` again once the LB shows as ready in the Utho Console.

## Example Usage

### HTTP frontend with round-robin

```hcl
resource "utho_loadbalancer_frontend" "http" {
  loadbalancer_id = utho_loadbalancer.main.id
  name            = "http"
  algorithm       = "roundrobin"
  proto           = "http"
  port            = "80"
  cookie          = "0"
  redirecthttps   = "1"   # redirect to HTTPS
  certificate_id  = "0"
}
```

### HTTPS frontend with SSL certificate

```hcl
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

### Sticky sessions with cookie

```hcl
resource "utho_loadbalancer_frontend" "sticky" {
  loadbalancer_id = utho_loadbalancer.main.id
  name            = "sticky-http"
  algorithm       = "roundrobin"
  proto           = "http"
  port            = "80"
  cookie          = "1"           # enable sticky sessions
  cookiename      = "SERVERID"    # cookie name
  redirecthttps   = "0"
  certificate_id  = "0"
}
```

### TCP frontend with least connections

```hcl
resource "utho_loadbalancer_frontend" "tcp" {
  loadbalancer_id = utho_loadbalancer.main.id
  name            = "mysql"
  algorithm       = "leastconn"
  proto           = "tcp"
  port            = "3306"
  cookie          = "0"
  redirecthttps   = "0"
  certificate_id  = "0"
}
```

## Argument Reference

| Argument          | Type   | Required | Description |
|-------------------|--------|----------|-------------|
| `loadbalancer_id` | String | Yes      | Load balancer ID. Changing this forces a new resource. |
| `name`            | String | Yes      | Frontend name. Changing this forces a new resource. |
| `algorithm`       | String | Yes      | Load balancing algorithm: `roundrobin`, `leastconn`, `first`. Changing this forces a new resource. |
| `proto`           | String | Yes      | Protocol: `http`, `https`, `tcp`, `udp`. Changing this forces a new resource. |
| `port`            | String | Yes      | Port to listen on. Changing this forces a new resource. |
| `cookie`          | String | Yes      | Enable sticky sessions via cookie: `"1"` or `"0"`. Changing this forces a new resource. |
| `redirecthttps`   | String | Yes      | Redirect HTTP to HTTPS: `"1"` or `"0"`. Only applicable when `proto = "http"`. Changing this forces a new resource. |
| `certificate_id`  | String | Yes      | SSL certificate ID. Use `"0"` when no certificate is needed. Changing this forces a new resource. |
| `cookiename`      | String | No       | Cookie name when `cookie = "1"`. Changing this forces a new resource. |

## Attribute Reference

| Attribute | Type   | Description |
|-----------|--------|-------------|
| `id`      | String | Unique frontend ID. |

## Notes

- The LB must be fully ready before this resource can be created — the provider handles this automatically.
- You can add multiple frontends to a single LB — one per port (e.g. port 80 and port 443).
- Use `certificate_id = "0"` for HTTP frontends — never leave it empty.
