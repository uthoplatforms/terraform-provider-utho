---
page_title: "Utho: utho_billing_cost_by_project"
subcategory: "Billing"
description: |-
  Retrieve current month cost breakdown by Utho project.
---

# data.utho_billing_cost_by_project

Retrieves the current billing cycle cost breakdown per project, including attributed, unattributed, and total monthly spend.

## Example Usage

```hcl
data "utho_billing_cost_by_project" "current" {}

output "total_month"        { value = data.utho_billing_cost_by_project.current.total_month }
output "total_unattributed" { value = data.utho_billing_cost_by_project.current.total_unattributed }
```

## Attribute Reference

| Attribute            | Type   | Description |
|----------------------|--------|-------------|
| `id`                 | String | Static identifier `billing_cost_by_project`. |
| `total_attributed`   | Number | Total cost attributed to projects (INR). |
| `total_unattributed` | Number | Total cost not attributed to any project (INR). |
| `total_month`        | Number | Total month spend (INR). |
| `cycle_start`        | String | Billing cycle start date. |
| `cycle_end`          | String | Billing cycle end date. |
| `projects`           | List   | Cost breakdown per project. |

### projects

| Attribute       | Type    | Description |
|-----------------|---------|-------------|
| `project_id`    | Number  | Project ID. |
| `name`          | String  | Project name. |
| `is_default`    | Boolean | Whether this is the default project. |
| `subtotal`      | Number  | Project spend (INR). |
| `resource_count`| Number  | Number of resources in the project. |
