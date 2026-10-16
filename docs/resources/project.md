---
page_title: "Utho: utho_project"
subcategory: "Projects"
description: |-
  Create and manage Utho Projects to organize resources and team members.
---

# utho_project

Creates and manages a Utho Project. Projects help you organize cloud resources and collaborate with team members across different environments.

## Example Usage

```hcl
resource "utho_project" "main" {
  name        = "my-app"
  description = "Production application project"
  environment = "production"
}
```

## Argument Reference

| Argument      | Type   | Required | Description |
|---------------|--------|----------|-------------|
| `name`        | String | Yes      | Project name. Can be updated in place. |
| `environment` | String | Yes      | Environment: `development`, `staging`, `production`, `testing`. Can be updated in place. |
| `description` | String | No       | Project description. Can be updated in place. |

## Attribute Reference

| Attribute       | Type    | Description |
|-----------------|---------|-------------|
| `id`            | String  | Project ID. |
| `status`        | String  | Project status (`active`). |
| `is_default`    | Boolean | Whether this is the default project. |
| `member_count`  | Number  | Number of project members. |
| `resource_count`| Number  | Number of resources in the project. |
| `created_at`    | String  | Creation timestamp. |

## Related Resources

- [utho_project_member](project_member) — Add IAM users to this project
