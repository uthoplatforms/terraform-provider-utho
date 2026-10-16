---
page_title: "Utho: utho_container_registry_immutable_rule"
subcategory: "Container Registry"
description: |-
  Create immutable tag rules for a Utho Container Registry.
---

# utho_container_registry_immutable_rule

Creates an immutable tag rule to prevent specific image tags from being overwritten or deleted.

## Example Usage

```hcl
resource "utho_container_registry_immutable_rule" "prod_tags" {
  project_name = utho_container_registry.main.project_name
  tag_pattern  = "v*"
  repo_pattern = "**"
}
```

## Argument Reference

| Argument       | Type   | Required | Description |
|----------------|--------|----------|-------------|
| `project_name` | String | Yes      | Registry project name. Changing forces new resource. |
| `tag_pattern`  | String | Yes      | Tag pattern to protect (e.g. `v*` for all version tags). |
| `repo_pattern` | String | Yes      | Repository pattern (e.g. `**` for all repos). |

## Attribute Reference

| Attribute | Type   | Description |
|-----------|--------|-------------|
| `id`      | String | Rule ID. |
| `rule_id` | Number | Numeric rule ID assigned by the registry. |
