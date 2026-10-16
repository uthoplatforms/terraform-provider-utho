---
page_title: "Utho: utho_loadbalancer"
subcategory: "Load Balancing"
description: |-
  Create and manage Utho Load Balancers.
---

# utho_loadbalancer

Creates and manages a Utho Load Balancer. A load balancer distributes incoming traffic across multiple backend servers, improving availability and scalability. Choose `network` type for TCP/UDP traffic or `application` type for HTTP/HTTPS with advanced routing.

~> **Note:** A VPC subnet is required to create a load balancer. Use `utho_vpc` and `utho_subnet` to create one, or reference an existing subnet ID.

~> **Note:** After creation, the provider automatically waits up to 5 minutes for the load balancer to become fully ready before allowing frontends or backends to be added. If the LB takes longer, run `terraform apply` again.

## Load Balancer Types

| Type | Best For |
|------|----------|
| `network` | TCP/UDP traffic, non-HTTP protocols, lower latency |
| `application` | HTTP/HTTPS with path-based routing, SSL termination, sticky sessions |

## Example Usage

### Simple HTTP load balancer

```hcl
resource "utho_vpc" "main" {
  name    = "production-vpc"
  network = "10.0.0.0"
  size    = "24"
  dcslug  = "inmumbaizone2"
  planid  = "1008"
}

resource "utho_subnet" "public" {
  name            = "public-subnet"
  vpc_id          = utho_vpc.main.id
  network         = "10.0.0.0"
  size            = 24
  type            = "public"
  assign_publicip = 1
}

resource "utho_loadbalancer" "main" {
  name            = "web-lb"
  type            = "application"
  dcslug          = "inmumbaizone2"
  vpc             = utho_subnet.public.id   # required
  enable_publicip = "true"
}

resource "utho_loadbalancer_frontend" "http" {
  loadbalancer_id = utho_loadbalancer.main.id
  name            = "http"
  algorithm       = "roundrobin"
  proto           = "http"
  port            = "80"
  cookie          = "0"
  redirecthttps   = "0"
  certificate_id  = "0"
}

resource "utho_loadbalancer_backend" "web" {
  count           = 3
  loadbalancer_id = utho_loadbalancer.main.id
  frontend_id     = utho_loadbalancer_frontend.http.id
  backend_port    = "80"
  weight          = "1"
  type            = "cloud"
  cloudid         = utho_cloud.web[count.index].id
}

output "lb_ip"  { value = utho_loadbalancer.main.ip }
output "lb_dns" { value = utho_loadbalancer.main.dns }
```

### HTTPS load balancer with SSL termination

SSL is terminated at the load balancer — backends receive plain HTTP.

```hcl
resource "utho_ssl_certificate" "main" {
  name            = "myapp-cert"
  type            = "Custom"
  certificate_key = file("cert.pem")
  private_key     = file("key.pem")
}

resource "utho_loadbalancer" "main" {
  name            = "web-lb"
  type            = "application"
  dcslug          = "inmumbaizone2"
  vpc             = utho_subnet.public.id
  enable_publicip = "true"
}

# Redirect HTTP to HTTPS
resource "utho_loadbalancer_frontend" "http" {
  loadbalancer_id = utho_loadbalancer.main.id
  name            = "http"
  algorithm       = "roundrobin"
  proto           = "http"
  port            = "80"
  cookie          = "0"
  redirecthttps   = "1"
  certificate_id  = "0"
}

# HTTPS with SSL cert
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

### Private load balancer inside a VPC

Internal load balancer — only accessible from within the VPC.

```hcl
resource "utho_loadbalancer" "internal" {
  name            = "internal-lb"
  type            = "application"
  dcslug          = "inmumbaizone2"
  vpc             = utho_subnet.private.id
  enable_publicip = "false"
}
```

## Argument Reference

### Required

| Argument          | Type   | Description |
|-------------------|--------|-------------|
| `name`            | String | Load balancer name. |
| `type`            | String | `network` or `application`. Changing this forces a new resource. |
| `dcslug`          | String | Data center slug. Changing this forces a new resource. |
| `vpc`             | String | VPC subnet ID. Required. Use `utho_subnet.name.id`. Changing this forces a new resource. |
| `enable_publicip` | String | `"true"` to assign a public IP, `"false"` for internal only. |

### Optional

| Argument   | Type   | Description |
|------------|--------|-------------|
| `firewall` | String | Security group ID to attach at creation. Changing this forces a new resource. |

## Attribute Reference

| Attribute    | Type   | Description |
|--------------|--------|-------------|
| `id`         | String | Unique load balancer ID. |
| `ip`         | String | Public IP address. Point your DNS A record here. |
| `dns`        | String | DNS hostname provided by Utho. |
| `status`     | String | Load balancer status (`Active`). |
| `created_at` | String | Creation timestamp. |

## Related Resources

| Resource | Purpose |
|----------|---------|
| [utho_loadbalancer_frontend](loadbalancer_frontend) | Add a listener (port, protocol, algorithm) |
| [utho_loadbalancer_backend](loadbalancer_backend) | Add backend servers |
| [utho_loadbalancer_settings](loadbalancer_settings) | Configure timeouts, HTTP/2, compression |
| [utho_loadbalancer_acl](loadbalancer_acl) | Add ACL routing rules |
| [utho_ssl_certificate](ssl_certificate) | Upload SSL certificate for HTTPS |

## Notes

- A VPC subnet is required — create one with `utho_vpc` + `utho_subnet` or use an existing subnet ID.
- The `dns` hostname is stable even if the IP changes — use a DNS CNAME record pointing to it in production.
- For zero-downtime updates, add new backends before removing old ones.
- Advanced settings are available via `utho_loadbalancer_settings`.
