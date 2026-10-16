---
page_title: "Utho: utho_project_member"
subcategory: "Projects"
description: |-
  Add IAM sub-users as members to a Utho Project.
---

# utho_project_member

Adds an IAM sub-user as a member of a Utho Project.

## Example Usage

```hcl
resource "utho_project" "main" {
  name        = "my-app"
  environment = "production"
}

resource "utho_project_member" "dev" {
  project_id = utho_project.main.id
  user_id    = 269944
}
```

## Argument Reference

| Argument     | Type   | Required | Description |
|--------------|--------|----------|-------------|
| `project_id` | String | Yes      | Project ID. Changing forces new resource. |
| `user_id`    | Number | Yes      | IAM sub-user ID to add. Changing forces new resource. |

## Attribute Reference

| Attribute | Type   | Description |
|-----------|--------|-------------|
| `id`      | String | Composite ID: `{project_id}:{user_id}`. |

## Notes

- The user must be an existing IAM sub-user in your account.
- Destroying this resource removes the member from the project.
