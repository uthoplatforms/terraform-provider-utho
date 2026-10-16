---
page_title: "Utho: utho_billing_usage"
subcategory: "Billing"
description: |-
  Retrieve current billing cycle usage breakdown by resource category.
---

# data.utho_billing_usage

Retrieves the current billing cycle usage summary including spend by resource category, top resources by cost, and wallet/credit balance.

## Example Usage

```hcl
data "utho_billing_usage" "current" {}

output "total_due"     { value = data.utho_billing_usage.current.due_amount }
output "cycle_start"   { value = data.utho_billing_usage.current.cycle_start }
output "days_left"     { value = data.utho_billing_usage.current.days_left }
```

### Conditional provisioning based on balance

```hcl
data "utho_billing_usage" "current" {}

locals {
  has_sufficient_balance = data.utho_billing_usage.current.wallet_balance > 0
}
```

## Attribute Reference

| Attribute            | Type   | Description |
|----------------------|--------|-------------|
| `id`                 | String | Static identifier `billing_usage`. |
| `cycle_start`        | String | Billing cycle start date (YYYY-MM-DD). |
| `cycle_end`          | String | Billing cycle end date (YYYY-MM-DD). |
| `days_in`            | Number | Days elapsed in current cycle. |
| `days_left`          | Number | Days remaining in current cycle. |
| `usage_excl_tax`     | Number | Total usage excluding tax (INR). |
| `tax`                | Number | Tax amount at 18% GST (INR). |
| `total_before_wallet`| Number | Total before wallet deduction (INR). |
| `due_amount`         | Number | Amount currently due (INR). |
| `wallet_balance`     | Number | Current wallet balance (INR). |
| `available_credit`   | Number | Available free credits (INR). |
| `by_category`        | List   | Usage breakdown by resource category. |
| `top_resources`      | List   | Top 5 resources by spend this cycle. |

### by_category

| Attribute       | Type   | Description |
|-----------------|--------|-------------|
| `key`           | String | Category key (e.g. `cloud`, `kubernetes`, `dbaas`). |
| `label`         | String | Human-readable category label. |
| `subtotal`      | Number | Category spend (INR). |
| `resource_count`| Number | Number of resources in this category. |
| `share_pct`     | Number | Percentage of total spend. |

### top_resources

| Attribute       | Type   | Description |
|-----------------|--------|-------------|
| `name`          | String | Resource name. |
| `category`      | String | Resource category key. |
| `category_label`| String | Human-readable category. |
| `amount`        | Number | Spend amount (INR). |
| `share_pct`     | Number | Percentage of total spend. |
