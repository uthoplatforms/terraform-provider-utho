---
page_title: "Utho: utho_loadbalancer_settings"
subcategory: "Load Balancing"
description: |-
  Configure advanced settings for a Utho Load Balancer.
---

# utho_loadbalancer_settings

Configures advanced timeout, connection, and performance settings for a Utho Load Balancer.

~> **Note:** All timeout values are in **milliseconds**, not seconds. For example, a 5-second connect timeout is `"5000"`.

## Example Usage

### Recommended production settings

```hcl
resource "utho_loadbalancer_settings" "main" {
  loadbalancer_id        = utho_loadbalancer.main.id
  timeout_connect        = "5000"     # 5 seconds
  timeout_client         = "50000"    # 50 seconds
  timeout_server         = "50000"    # 50 seconds
  timeout_http_request   = "5000"     # 5 seconds
  timeout_http_keepalive = "60000"    # 60 seconds
  timeout_tunnel         = "3600000"  # 1 hour (for WebSocket)
  max_connections        = "2000"
  http2                  = "1"
  compression            = "1"
  hsts                   = "0"
}
```

### High performance API gateway

```hcl
resource "utho_loadbalancer_settings" "api" {
  loadbalancer_id        = utho_loadbalancer.api.id
  timeout_connect        = "3000"     # 3 seconds
  timeout_client         = "30000"    # 30 seconds
  timeout_server         = "30000"    # 30 seconds
  timeout_http_request   = "3000"     # 3 seconds
  timeout_http_keepalive = "30000"    # 30 seconds
  timeout_tunnel         = "3600000"  # 1 hour
  max_connections        = "10000"
  http2                  = "1"
  compression            = "1"
  hsts                   = "1"
}
```

### WebSocket / real-time app

```hcl
resource "utho_loadbalancer_settings" "ws" {
  loadbalancer_id        = utho_loadbalancer.ws.id
  timeout_connect        = "5000"
  timeout_client         = "3600000"   # 1 hour — keep WebSocket alive
  timeout_server         = "3600000"   # 1 hour
  timeout_http_request   = "5000"
  timeout_http_keepalive = "3600000"   # 1 hour
  timeout_tunnel         = "86400000"  # 24 hours
  max_connections        = "5000"
  http2                  = "0"         # disable for WebSocket
  compression            = "0"
  hsts                   = "0"
}
```

## Argument Reference

| Argument                 | Type   | Required | Description |
|--------------------------|--------|----------|-------------|
| `loadbalancer_id`        | String | Yes      | Load balancer ID. |
| `timeout_connect`        | String | No       | Connection timeout in **milliseconds**. Default: `5000` (5s). |
| `timeout_client`         | String | No       | Client inactivity timeout in **milliseconds**. Default: `50000` (50s). |
| `timeout_server`         | String | No       | Server response timeout in **milliseconds**. Default: `50000` (50s). |
| `timeout_http_request`   | String | No       | HTTP request timeout in **milliseconds**. Default: `5000` (5s). |
| `timeout_http_keepalive` | String | No       | HTTP keep-alive timeout in **milliseconds**. Default: `60000` (60s). |
| `timeout_tunnel`         | String | No       | Tunnel timeout in **milliseconds** for WebSocket/TCP. Default: `3600000` (1h). |
| `max_connections`        | String | No       | Maximum concurrent connections. Default: `2000`. |
| `http2`                  | String | No       | Enable HTTP/2: `"1"` or `"0"`. |
| `compression`            | String | No       | Enable response compression: `"1"` or `"0"`. |
| `hsts`                   | String | No       | Enable HTTP Strict Transport Security: `"1"` or `"0"`. |

## Attribute Reference

| Attribute | Type   | Description |
|-----------|--------|-------------|
| `id`      | String | Same as `loadbalancer_id`. |

## Notes

- All timeout values are in **milliseconds**. `5000` = 5 seconds, `60000` = 60 seconds.
- For WebSocket applications set `timeout_tunnel` to a high value (e.g. `86400000` for 24 hours).
- Enable `http2 = "1"` for better performance with modern browsers.
- Enable `hsts = "1"` only when you have HTTPS configured — it tells browsers to always use HTTPS.
