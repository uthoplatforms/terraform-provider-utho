---
page_title: "Utho: utho_loadbalancer_backend"
subcategory: "Load Balancing"
description: |-
  Add and manage backends on a Utho Load Balancer frontend.
---

# utho_loadbalancer_backend

Adds a backend server to a Utho Load Balancer frontend. Backends can be Utho cloud instances or custom IP addresses.

~> **Note:** The load balancer must be fully ready before backends can be added. The provider automatically polls and waits. Add backends after frontends are created.

## Example Usage

### Add cloud instances as backends

```hcl
resource "utho_loadbalancer_backend" "app" {
  count           = 3
  loadbalancer_id = utho_loadbalancer.main.id
  frontend_id     = utho_loadbalancer_frontend.http.id
  backend_port    = "80"
  weight          = "1"
  type            = "cloud"
  cloudid         = utho_cloud.app[count.index].id
}
```

### Weighted backends (70/30 split)

```hcl
resource "utho_loadbalancer_backend" "primary" {
  loadbalancer_id = utho_loadbalancer.main.id
  frontend_id     = utho_loadbalancer_frontend.http.id
  backend_port    = "80"
  weight          = "7"   # 70% of traffic
  type            = "cloud"
  cloudid         = utho_cloud.primary.id
}

resource "utho_loadbalancer_backend" "secondary" {
  loadbalancer_id = utho_loadbalancer.main.id
  frontend_id     = utho_loadbalancer_frontend.http.id
  backend_port    = "80"
  weight          = "3"   # 30% of traffic
  type            = "cloud"
  cloudid         = utho_cloud.secondary.id
}
```

### Add a custom IP backend

```hcl
resource "utho_loadbalancer_backend" "external" {
  loadbalancer_id = utho_loadbalancer.main.id
  frontend_id     = utho_loadbalancer_frontend.http.id
  backend_port    = "80"
  weight          = "1"
  type            = "custom"
  ip              = "10.0.0.42"
}
```

### Scale backends up and down

```hcl
variable "backend_count" {
  description = "Number of backend servers."
  type        = number
  default     = 3
}

resource "utho_loadbalancer_backend" "app" {
  count           = var.backend_count

  loadbalancer_id = utho_loadbalancer.main.id
  frontend_id     = utho_loadbalancer_frontend.http.id
  backend_port    = "80"
  weight          = "1"
  type            = "cloud"
  cloudid         = utho_cloud.app[count.index].id
}
```

Change `backend_count` and run `terraform apply` to scale up or down.

## Argument Reference

| Argument          | Type   | Required | Description |
|-------------------|--------|----------|-------------|
| `loadbalancer_id` | String | Yes      | Load balancer ID. Changing this forces a new resource. |
| `frontend_id`     | String | Yes      | Frontend ID to attach this backend to. Changing this forces a new resource. |
| `backend_port`    | String | Yes      | Port the backend server listens on. Changing this forces a new resource. |
| `weight`          | String | Yes      | Backend weight for traffic distribution. Higher weight = more traffic. Changing this forces a new resource. |
| `type`            | String | Yes      | Backend type: `cloud` or `custom`. Changing this forces a new resource. |
| `cloudid`         | String | No       | Cloud instance ID. Required when `type = "cloud"`. Changing this forces a new resource. |
| `ip`              | String | No       | Custom IP address. Required when `type = "custom"`. Changing this forces a new resource. |

## Attribute Reference

| Attribute | Type   | Description |
|-----------|--------|-------------|
| `id`      | String | Unique backend ID. |

## Notes

- Use `count` to add multiple backends from a list of cloud instances.
- Weight is relative — if backends have weights `1`, `1`, `1` they share traffic equally. If `3`, `1` they get 75%/25%.
- To drain a backend before removal, set its weight to `0` in a separate step.
- Removing a backend from the list removes it from the LB and destroys the resource.
