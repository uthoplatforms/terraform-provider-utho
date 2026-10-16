---
page_title: "Utho: utho_container_registry_webhook"
subcategory: "Container Registry"
description: |-
  Create webhooks for Utho Container Registry events.
---

# utho_container_registry_webhook

Creates a webhook that triggers on registry events such as image pushes.

## Example Usage

```hcl
resource "utho_container_registry_webhook" "deploy" {
  project_name    = utho_container_registry.main.project_name
  name            = "deploy-trigger"
  address         = "https://ci.mycompany.com/webhook/registry"
  event_types     = ["PUSH_ARTIFACT"]
  payload_format  = "Default"
  skip_cert_verify = false
}
```

## Argument Reference

| Argument          | Type   | Required | Description |
|-------------------|--------|----------|-------------|
| `project_name`    | String | Yes      | Registry project name. Changing forces new resource. |
| `name`            | String | Yes      | Webhook name. Changing forces new resource. |
| `address`         | String | Yes      | Webhook URL endpoint. Changing forces new resource. |
| `event_types`     | List   | Yes      | Events to trigger on: `PUSH_ARTIFACT`, `DELETE_ARTIFACT`. |
| `payload_format`  | String | Yes      | Payload format: `Default` or `Slack`. Changing forces new resource. |
| `auth_header`     | String | No       | Authorization header value. |
| `skip_cert_verify`| Bool   | No       | Skip TLS certificate verification. Default: `false`. |

## Attribute Reference

| Attribute | Type   | Description |
|-----------|--------|-------------|
| `id`      | String | Webhook ID. |
