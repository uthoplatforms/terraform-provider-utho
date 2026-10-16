---
page_title: "Utho: utho_container_registry_robot"
subcategory: "Container Registry"
description: |-
  Create robot accounts for automated access to a Utho Container Registry.
---

# utho_container_registry_robot

Creates a robot account for automated CI/CD access to a container registry. Robot accounts have scoped permissions and optional expiry.

## Example Usage

```hcl
resource "utho_container_registry_robot" "ci" {
  project_name = utho_container_registry.main.project_name
  name         = "ci-robot"
  description  = "CI/CD pipeline robot"
  duration     = 90

  access = [
    {
      resource = "repository"
      action   = "push"
    },
    {
      resource = "repository"
      action   = "pull"
    }
  ]
}

output "robot_secret" {
  value     = utho_container_registry_robot.ci.secret
  sensitive = true
}
```

## Argument Reference

| Argument       | Type   | Required | Description |
|----------------|--------|----------|-------------|
| `project_name` | String | Yes      | Registry project name. Changing forces new resource. |
| `name`         | String | Yes      | Robot account name. Changing forces new resource. |
| `description`  | String | No       | Description of the robot account. |
| `duration`     | Number | Yes      | Token validity in days. Changing forces new resource. |
| `access`       | List   | Yes      | List of access permissions. |

### access block

| Argument   | Type   | Required | Description |
|------------|--------|----------|-------------|
| `resource` | String | Yes      | Resource type: `repository`. |
| `action`   | String | Yes      | Action: `push`, `pull`, `delete`. |

## Attribute Reference

| Attribute    | Type   | Description |
|--------------|--------|-------------|
| `id`         | String | Robot account ID. |
| `secret`     | String | Robot secret token (sensitive). |
| `expires_at` | Number | Expiry timestamp. |
