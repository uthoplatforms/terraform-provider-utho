---
page_title: "Utho: utho_billing_cycles"
subcategory: "Billing"
description: |-
  List available billing cycles for a Utho product.
---

# utho_billing_cycles

Fetches the available billing cycles for a specific Utho product. Use this to discover valid `billingcycle` values before creating resources — available cycles vary by product.

## Example Usage

### List billing cycles for cloud instances

```hcl
data "utho_billing_cycles" "cloud" {
  product = "cloud"
}

output "cloud_billing_cycles" {
  value = data.utho_billing_cycles.cloud.cycles
}

# Output:
# ["hourly", "monthly", "3month", "6month", "12month", "24month", "36month"]
```

### List billing cycles for Kubernetes

```hcl
data "utho_billing_cycles" "kubernetes" {
  product = "kubernetes"
}

output "k8s_cycles" {
  value = data.utho_billing_cycles.kubernetes.cycles
}
```

### Use in resource

```hcl
data "utho_billing_cycles" "cloud" {
  product = "cloud"
}

resource "utho_cloud" "web" {
  hostname        = "web-01.mhc"
  dcslug          = "inmumbaizone2"
  planid          = "10308"
  billingcycle    = "monthly"   # must be one of data.utho_billing_cycles.cloud.cycles
  auth            = "option2"
  sshkeys         = utho_ssh_key.deploy.id
  image           = "ubuntu-22.04-x86_64"
  enable_publicip = "true"
  cpumodel        = "amd"
}
```

## Argument Reference

| Argument  | Type   | Required | Description |
|-----------|--------|----------|-------------|
| `product` | String | Yes      | Product name. See [Supported Products](#supported-products). |

## Attribute Reference

### Top-level

| Attribute | Type         | Description |
|-----------|--------------|-------------|
| `cycles`  | List(String) | List of valid billing cycle values for this product. |

## Supported Products

| Product         | Description |
|-----------------|-------------|
| `cloud`         | Cloud instances |
| `kubernetes`    | Managed Kubernetes clusters |
| `database`      | Managed databases |
| `objectstorage` | Object storage buckets |

## Billing Cycle Reference

| Value     | Duration | Console Label |
|-----------|----------|---------------|
| `hourly`  | Pay per hour | Hourly |
| `monthly` | 1 month | Monthly |
| `3month`  | 3 months | — |
| `6month`  | 6 months | — |
| `12month` | 1 year | Annual |
| `24month` | 2 years | — |
| `36month` | 3 years | 3 Years (20% OFF) |

~> **Tip:** Longer billing cycles offer significant discounts. `36month` gives 20% off compared to monthly billing.
