---
page_title: "Best Practices for Production"
subcategory: ""
description: |-
  Recommendations for running Utho infrastructure with Terraform in production: state, secrets, environments, version pinning, and safe changes.
---

# Best Practices for Production

This guide collects recommendations for teams running Utho infrastructure with Terraform beyond a single laptop. None of it is required to use the provider, but following it avoids the most common production incidents.

## Pin Provider and Terraform Versions

Always pin the provider to a minor version so a new release never changes your infrastructure unexpectedly.

```hcl
terraform {
  required_version = ">= 1.5"

  required_providers {
    utho = {
      source  = "nitinuthocloud/utho"
      version = "~> 0.7" # allows 0.7.x, never 0.8.0
    }
  }
}
```

Commit the generated `.terraform.lock.hcl` file to version control. It records the exact provider build and its checksums, so every teammate and CI run uses the same binary.

Before upgrading, read the [changelog](https://github.com/nitinuthocloud/terraform-provider-utho/blob/test/new-provider/CHANGELOG.md), then run:

```bash
terraform init -upgrade
terraform plan
```

## Store State Remotely

Local `terraform.tfstate` files are easy to lose and cannot be shared safely. Use a remote backend with locking.

**HCP Terraform:**

```hcl
terraform {
  cloud {
    organization = "your-org"

    workspaces {
      name = "utho-production"
    }
  }
}
```

**Utho Object Storage (S3-compatible):**

```hcl
terraform {
  backend "s3" {
    bucket                      = "terraform-state"
    key                         = "production/terraform.tfstate"
    region                      = "us-east-1" # required by the backend, ignored by Utho
    endpoints                   = { s3 = "https://innoida.utho.io" }
    skip_credentials_validation = true
    skip_region_validation      = true
    skip_requesting_account_id  = true
    skip_metadata_api_check     = true
    use_path_style              = true
  }
}
```

Create the bucket and access keys with [`utho_object_storage`](../resources/object_storage) and [`utho_object_storage_key`](../resources/object_storage_key) in a separate bootstrap configuration, and enable versioning so you can recover an earlier state. Replace the endpoint with the one for your bucket's data center, as shown in the Utho Console.

~> **Note:** State files contain sensitive values such as passwords and API tokens in plain text. Restrict who can read the backend.

## Keep Secrets Out of Code

* Pass the API key through `UTHO_API_KEY` or a `sensitive` variable. Never commit it.
* Mark every secret input as `sensitive = true` so it is redacted from plan output.
* In CI/CD, read the key from your platform's secret store. See [Authenticating in CI/CD and HCP Terraform](authentication-ci-cd).
* Use a dedicated API key per environment and per pipeline, created with [`utho_api_token`](../resources/api_token) or in the console, so you can rotate or revoke one without affecting the others.
* Add these to `.gitignore`:

    ```gitignore
    .terraform/
    *.tfstate
    *.tfstate.*
    *.tfvars
    crash.log
    ```

## Separate Environments

Keep development, staging, and production in separate state files so a mistake in one cannot touch another. The simplest pattern is one directory per environment that calls shared modules:

```text
infra/
├── modules/
│   ├── network/      # utho_vpc, utho_subnet, utho_nat_gateway
│   ├── app/          # utho_cloud, utho_firewall, utho_loadbalancer
│   └── database/     # utho_database, utho_database_user
└── environments/
    ├── dev/
    │   └── main.tf
    ├── staging/
    │   └── main.tf
    └── production/
        └── main.tf
```

Use [`utho_project`](../resources/project) to group each environment's resources in the Utho Console and on your invoice. You can then track spend per environment with [`utho_billing_cost_by_project`](../data-sources/billing_cost_by_project).

## Look Up IDs Instead of Hard-Coding Them

Plan IDs, image slugs, and data center availability change over time. Use data sources so your configuration keeps working:

```hcl
data "utho_cloud_plans" "all" {}
data "utho_cloud_images" "all" {}
data "utho_cloud_dczones" "all" {}
```

See [Data Centers, Plans, and Images](data-centers-plans-images) for filtering examples.

## Protect Critical Resources

Use `prevent_destroy` on resources whose loss would cause an outage or data loss:

```hcl
resource "utho_database" "main" {
  # ...

  lifecycle {
    prevent_destroy = true
  }
}
```

Good candidates: `utho_database`, `utho_kubernetes`, `utho_object_storage`, `utho_vpc`, `utho_elastic_ip`, and `utho_dns_zone`.

## Review Every Plan

* Run `terraform plan -out=tfplan` and apply that exact file with `terraform apply tfplan`, so what you reviewed is what runs.
* Treat any `destroy` or `must be replaced` line in production as a stop sign until you understand why.
* In CI, post the plan to the pull request and require approval before apply.

## Back Up Data

Terraform manages infrastructure, not data. For stateful services:

* Enable `enablebackup = "true"` on [`utho_cloud`](../resources/cloud) instances that hold data.
* Take snapshots before risky changes with [`utho_cloud_snapshot`](../resources/cloud_snapshot).
* Use versioned buckets for backups. See the [`object-storage/with-versioning`](https://github.com/nitinuthocloud/terraform-provider-utho/tree/test/new-provider/examples/object-storage/with-versioning) example.

## Monitor What You Deploy

Create alerts alongside the infrastructure they watch, so new servers are never unmonitored:

```hcl
resource "utho_alert_contact" "oncall" {
  # ...
}

resource "utho_alert" "cpu" {
  # ...
}
```

See [`utho_alert`](../resources/alert) and [`utho_alert_contact`](../resources/alert_contact) for the full set of arguments, and the [`monitoring/full-stack`](https://github.com/nitinuthocloud/terraform-provider-utho/tree/test/new-provider/examples/monitoring/full-stack) example.

## Checklist

- [ ] Provider pinned with `~>` and `.terraform.lock.hcl` committed
- [ ] Remote state with locking and restricted access
- [ ] API key supplied by environment or secret store, never committed
- [ ] Separate state per environment
- [ ] Data sources instead of hard-coded plan and image IDs
- [ ] `prevent_destroy` on databases, clusters, and buckets
- [ ] Plans reviewed before every apply
- [ ] Backups and snapshots for stateful resources
- [ ] Alerts for every production server
