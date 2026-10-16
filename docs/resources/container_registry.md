---
page_title: "Utho: utho_container_registry"
subcategory: "Container Registry"
description: |-
  Create and manage Utho Container Registries for storing Docker images.
---

# utho_container_registry

Creates and manages a Utho Container Registry. Use it to store, manage and deploy Docker container images.

## Example Usage

```hcl
resource "utho_container_registry" "main" {
  project_name = "my-registry"
  dcslug       = "innoida"
  planid       = 10276
  billingcycle = "monthly"
  public       = false
}
```

## Argument Reference

| Argument       | Type   | Required | Description |
|----------------|--------|----------|-------------|
| `project_name` | String | Yes      | Registry project name. Unique identifier. Changing forces new resource. |
| `dcslug`       | String | Yes      | Data center slug. Currently `innoida` is supported. Changing forces new resource. |
| `planid`       | Number | Yes      | Plan ID. See available plans below. Changing forces new resource. |
| `billingcycle` | String | Yes      | Billing cycle: `monthly`. Changing forces new resource. |
| `public`       | Bool   | Yes      | Set `true` for public registry, `false` for private. |

## Available Plans

| Plan ID | Name         | Storage | Price/month |
|---------|--------------|---------|-------------|
| `10276` | Starter      | 250 GB  | ₹1,500      |
| `10275` | Growth       | 500 GB  | ₹4,000      |
| `10274` | Professional | 1000 GB | ₹7,500      |
| `10273` | Enterprise   | 2000 GB | ₹15,000     |

## Attribute Reference

| Attribute    | Type   | Description |
|--------------|--------|-------------|
| `id`         | String | Registry project name (used as ID). |
| `created_at` | String | Creation timestamp. |

## Related Resources

- [utho_container_registry_robot](container_registry_robot) — Robot accounts for CI/CD
- [utho_container_registry_webhook](container_registry_webhook) — Webhooks for image push events
- [utho_container_registry_immutable_rule](container_registry_immutable_rule) — Immutable tag rules
