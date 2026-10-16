---
page_title: "Utho: utho_billing_invoices"
subcategory: "Billing"
description: |-
  Retrieve all billing invoices for your Utho account.
---

# data.utho_billing_invoices

Retrieves all billing invoices for the account including paid, unpaid, and overdue invoices.

## Example Usage

```hcl
data "utho_billing_invoices" "all" {}

output "invoice_count" { value = length(data.utho_billing_invoices.all.invoices) }
```

## Attribute Reference

| Attribute  | Type | Description |
|------------|------|-------------|
| `id`       | String | Static identifier `billing_invoices`. |
| `invoices` | List   | List of all invoices. |

### invoices

| Attribute    | Type   | Description |
|--------------|--------|-------------|
| `id`         | Number | Invoice ID. |
| `invoice_num`| String | Invoice number (e.g. `#227442`). |
| `date`       | String | Invoice date (YYYY-MM-DD). |
| `due_date`   | String | Payment due date (YYYY-MM-DD). |
| `amount`     | Number | Invoice amount (INR). |
| `status`     | String | Invoice status: `Paid`, `Unpaid`. |
