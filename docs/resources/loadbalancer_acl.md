---
page_title: "Utho: utho_loadbalancer_acl"
subcategory: "Load Balancing"
description: |-
  Add and manage ACL rules on a Utho Load Balancer frontend.
---

# utho_loadbalancer_acl

Adds an ACL (Access Control List) rule to a Utho Load Balancer frontend. ACL rules allow you to route or filter traffic based on request attributes like URL path, user agent, host header, or referer.

## Example Usage

### URL path based routing

```hcl
resource "utho_loadbalancer_acl" "api_route" {
  loadbalancer_id = utho_loadbalancer.main.id
  frontend_id     = utho_loadbalancer_frontend.http.id
  name            = "api-route"
  condition_type  = "url_path"
  value           = "{\"path\":\"/api\"}"
}
```

### Host header match

```hcl
resource "utho_loadbalancer_acl" "host_match" {
  loadbalancer_id = utho_loadbalancer.main.id
  frontend_id     = utho_loadbalancer_frontend.http.id
  name            = "host-match"
  condition_type  = "host"
  value           = "{\"host\":\"api.example.com\"}"
}
```

### User agent filter

```hcl
resource "utho_loadbalancer_acl" "block_bot" {
  loadbalancer_id = utho_loadbalancer.main.id
  frontend_id     = utho_loadbalancer_frontend.http.id
  name            = "block-googlebot"
  condition_type  = "http_user_agent"
  value           = "{\"user_agent\":\"GoogleBot\"}"
}
```

## Argument Reference

| Argument          | Type   | Required | Description |
|-------------------|--------|----------|-------------|
| `loadbalancer_id` | String | Yes      | Load balancer ID. Changing this forces a new resource. |
| `frontend_id`     | String | Yes      | Frontend ID to attach this ACL rule to. |
| `name`            | String | Yes      | ACL rule name. |
| `condition_type`  | String | Yes      | Condition type. See [Condition Types](#condition-types) below. |
| `value`           | String | Yes      | ACL rule value as a JSON string. Format depends on condition_type — see table below. |

### Condition Types

| condition_type     | Description                          | value JSON format                          | Example value                        |
|--------------------|--------------------------------------|--------------------------------------------|--------------------------------------|
| `host`             | Match on Host header.                | `{"host": "<hostname>"}`                   | `{"host":"api.example.com"}`         |
| `url_path`         | Match on URL path prefix.            | `{"path": "<path>"}`                       | `{"path":"/api"}`                    |
| `url_path_regex`   | Match on URL path regex.             | `{"regex": "<pattern>"}`                   | `{"regex":"/api/.*"}`                |
| `http_user_agent`  | Match on User-Agent header.          | `{"user_agent": "<value>"}`                | `{"user_agent":"GoogleBot"}`         |
| `http_referer`     | Match on Referer header.             | `{"referer": "<value>"}`                   | `{"referer":"example.com"}`          |

~> **Note:** `http_method` is not currently supported by the Utho API and will return an `Invalid HTTP method` error. Use the condition types listed above.

## Attribute Reference

| Attribute | Type   | Description |
|-----------|--------|-------------|
| `id`      | String | Unique ACL rule ID. |