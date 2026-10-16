---
page_title: "Importing Existing Resources"
subcategory: ""
description: |-
  Bring infrastructure that was created outside Terraform under Terraform management.
---

# Importing Existing Resources

If you created infrastructure in the Utho Console or with the API, you can bring it under Terraform management without recreating it. After import, Terraform tracks the resource in state, and future changes go through `terraform plan` and `terraform apply`.

## Resources that support import

| Resource | Import ID | Example |
|----------|-----------|---------|
| [`utho_cloud`](../resources/cloud) | Instance ID | `1671990` |
| [`utho_autoscaling`](../resources/autoscaling) | Auto scaling group ID | `12345` |
| [`utho_nat_gateway`](../resources/nat_gateway) | NAT gateway ID (UUID) | `3f2b8c1e-...` |

Other resources do not support import yet. To find IDs for existing infrastructure, use the list data sources, such as [`utho_clouds`](../data-sources/clouds), [`utho_vpcs`](../data-sources/vpcs), and [`utho_firewalls`](../data-sources/firewalls).

## Import with an `import` block (Terraform 1.5+)

Configuration-driven import is the recommended approach. It is reviewed in `terraform plan` like any other change and can be committed to version control.

```hcl
import {
  to = utho_cloud.web
  id = "1671990"
}

resource "utho_cloud" "web" {
  hostname     = "web-01.mhc"
  dcslug       = "inmumbaizone2"
  planid       = "10308"
  billingcycle = "hourly"
  auth         = "option2"
  sshkeys      = utho_ssh_key.deploy.id
}
```

```bash
terraform plan
terraform apply
```

Once the apply finishes, you can remove the `import` block.

## Import with the CLI

On older Terraform versions, write the `resource` block first and then run:

```bash
terraform import utho_cloud.web 1671990
terraform import utho_autoscaling.web 12345
terraform import utho_nat_gateway.main 3f2b8c1e-0000-0000-0000-000000000000
```

## After importing

Run `terraform plan` and check the result. Some arguments, such as `auth`, `planid`, `image`, and `root_password` on `utho_cloud`, are only used when a resource is created, and the Utho API does not return them. Terraform cannot fill these in for you, so:

1. Set them in your configuration to the values the resource was created with.
2. Run `terraform plan` again.
3. Repeat until the plan shows no changes.

~> **Note:** Do not run `terraform apply` while the plan still shows unexpected changes to an imported resource. Resources that do not support in-place updates will fail with an "Update not supported" error rather than be modified.

## Removing a resource from Terraform without deleting it

To stop managing a resource without destroying it, use a `removed` block (Terraform 1.7+):

```hcl
removed {
  from = utho_cloud.web

  lifecycle {
    destroy = false
  }
}
```

On older versions, run `terraform state rm utho_cloud.web`.
