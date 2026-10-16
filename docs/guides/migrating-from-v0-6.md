---
page_title: "Migrating from v0.6.x to v0.7.0"
subcategory: ""
description: |-
  Upgrade an existing Utho provider configuration from v0.6.x to v0.7.0: renamed resources, changed arguments, and how to move state without recreating infrastructure.
---

# Migrating from v0.6.x to v0.7.0

Version 0.7.0 is a rewrite of the Utho provider. It grows coverage from 8 resources to 53 and from 3 data sources to 18, and it reorganizes several resources so that each one maps to a single Utho API object.

That reorganization means **v0.7.0 is not a drop-in upgrade**. Some resource types were renamed, some arguments were renamed or changed type, and a few v0.6 resources were split into smaller ones. This guide walks through every change and gives you a safe, step-by-step upgrade path.

~> **Important:** Changing a resource type in your configuration (for example `utho_cloud_instance` → `utho_cloud`) makes Terraform treat it as a *different* resource. If you change the code and run `terraform apply` without following the steps below, Terraform will plan to **destroy** the old infrastructure and create new infrastructure. Always read the plan carefully.

## Before You Begin

1. **Back up your state.**

    ```bash
    terraform state pull > terraform.tfstate.backup-v0.6
    ```

2. **Make sure your v0.6 configuration is clean.** Run `terraform plan` on the old version and confirm it reports `No changes`. Migrating from a drifted state makes the upgrade much harder to review.

3. **Check your Terraform version.** This guide uses [`import` blocks](https://developer.hashicorp.com/terraform/language/import) (Terraform 1.5+) and [`removed` blocks](https://developer.hashicorp.com/terraform/language/resources/syntax#removing-resources) (Terraform 1.7+).

    ```bash
    terraform version
    ```

4. **Plan a maintenance window** for any resource in the [Resources That Must Be Recreated](#resources-that-must-be-recreated) section below.

## What's New in v0.7.0

| Area | v0.6.x | v0.7.0 |
|------|--------|--------|
| Resources | 8 | 53 |
| Data sources | 3 | 18 |
| Kubernetes, Database, Object Storage, Container Registry, IAM, Projects, Monitoring, IPSec, NAT, VPC Peering, Billing | — | ✅ |
| Authentication | `api_key` (required) | `api_key` **or** `UTHO_API_KEY` env var |
| Import support | Limited | `utho_cloud`, `utho_autoscaling`, `utho_nat_gateway` (more coming) |

## Resource Name Changes

| v0.6.x resource | v0.7.0 resource | Migration path |
|-----------------|-----------------|----------------|
| `utho_cloud_instance` | [`utho_cloud`](../resources/cloud) | ✅ In place, with an `import` block |
| `utho_auto_scaling` | [`utho_autoscaling`](../resources/autoscaling) (+ [`utho_autoscaling_policy`](../resources/autoscaling_policy), [`utho_autoscaling_schedule`](../resources/autoscaling_schedule)) | ✅ In place, with an `import` block |
| `utho_domain` | [`utho_dns_zone`](../resources/dns_zone) | ⚠️ Recreate, or release from state |
| `utho_dns_record` | [`utho_dns_record`](../resources/dns_record) | ⚠️ Same name, new schema. Recreate, or release from state |
| `utho_firewall` | [`utho_firewall`](../resources/firewall) + [`utho_firewall_rule`](../resources/firewall_rule) + [`utho_firewall_server`](../resources/firewall_server) | ⚠️ Recreate, or release from state |
| `utho_loadbalancer` | [`utho_loadbalancer`](../resources/loadbalancer) + [`utho_loadbalancer_frontend`](../resources/loadbalancer_frontend) + [`utho_loadbalancer_backend`](../resources/loadbalancer_backend) | ⚠️ Recreate, or release from state |
| `utho_vpc` | [`utho_vpc`](../resources/vpc) | ⚠️ Same name, new schema. Recreate, or release from state |
| `utho_target_group` | — (read with the [`utho_target_groups`](../data-sources/target_groups) data source) | ⚠️ Release from state and manage in the console |

## Data Source Changes

| v0.6.x data source | v0.7.0 replacement |
|--------------------|--------------------|
| `utho_images` | [`utho_cloud_images`](../data-sources/cloud_images) |
| `utho_account` | No direct replacement. For spend and usage, use [`utho_billing_usage`](../data-sources/billing_usage), [`utho_billing_invoices`](../data-sources/billing_invoices), and [`utho_billing_cost_by_project`](../data-sources/billing_cost_by_project). |
| `utho_object_storage_plan` | No direct replacement yet. Object storage buckets are now managed with [`utho_object_storage`](../resources/object_storage). |

## Argument Changes

### `utho_cloud_instance` → `utho_cloud`

| v0.6.x argument | v0.7.0 argument | Notes |
|-----------------|-----------------|-------|
| `name` | `hostname` | Renamed. |
| `vpc_id` | `vpc` | Renamed. Takes a VPC subnet ID. |
| `enablebackup` (Bool) | `enablebackup` (String) | Now `"true"` / `"false"`. |
| `management` | `support` | Use `support = "managed"` or `"unmanaged"`. |
| `subnetrequired` | — | Removed. Pass `vpc` instead. |
| `dcslug`, `planid`, `image`, `auth`, `root_password`, `sshkeys`, `billingcycle`, `cpumodel`, `enable_publicip`, `firewall`, `snapshotid`, `backupid` | Unchanged | `cpumodel` should be `amd` or `intel`. |
| — | `iso`, `stack`, `delete_ebs`, `ebs` | New in v0.7.0. |

**Before (v0.6.x):**

```hcl
resource "utho_cloud_instance" "web" {
  name            = "web-01.example.com"
  dcslug          = "inmumbaizone2"
  planid          = "10308"
  image           = "ubuntu-22.04-x86_64"
  auth            = "option2"
  sshkeys         = "4321"
  billingcycle    = "hourly"
  enable_publicip = "true"
  vpc_id          = "a1b2c3d4-..."
  enablebackup    = false
}
```

**After (v0.7.0):**

```hcl
resource "utho_cloud" "web" {
  hostname        = "web-01.example.com"
  dcslug          = "inmumbaizone2"
  planid          = "10308"
  image           = "ubuntu-22.04-x86_64"
  auth            = "option2"
  sshkeys         = "4321"
  billingcycle    = "hourly"
  enable_publicip = "true"
  vpc             = "a1b2c3d4-..."
  enablebackup    = "false"
}
```

### `utho_auto_scaling` → `utho_autoscaling`

| v0.6.x argument | v0.7.0 argument | Notes |
|-----------------|-----------------|-------|
| `vpc_id` | `vpc` | Renamed. |
| `loadbalancers_id` | `load_balancers` | Renamed. Can be updated in place. |
| `security_group_id` | `security_groups` | Renamed. Can be updated in place. |
| `target_groups_id` | `target_groups` | Renamed. Can be updated in place. |
| `public_ip_enabled` (Bool) | `public_ip_enabled` (Number) | Now `1` or `0`. |
| `instance_templateid` | — | Removed. Use `stack`, `stackid`, `stackimage`, `snapshotid`, or `image_name`. |
| `policies` | `policies` | Still inline. Standalone policies are also available as `utho_autoscaling_policy`. |
| — | `cpumodel`, `schedules` | New. Schedules are also available as `utho_autoscaling_schedule`. |

### `utho_domain` → `utho_dns_zone`

The only argument, `domain`, is unchanged.

### `utho_dns_record`

| v0.6.x argument | v0.7.0 argument | Notes |
|-----------------|-----------------|-------|
| `domain`, `type`, `hostname`, `value`, `ttl` | Unchanged | `ttl` is now required. |
| `porttype`, `port`, `priority`, `weight` | — | Not available yet. Keep SRV records and MX priorities in the Utho Console until support lands. |

### `utho_firewall`

In v0.6.x a firewall was a single resource. In v0.7.0 the security group, its rules, and its server attachments are separate resources, so each rule can be added or removed without touching the others.

```hcl
resource "utho_firewall" "web" {
  name = "web-sg"
}

resource "utho_firewall_rule" "https" {
  firewall_id  = utho_firewall.web.id
  type         = "incoming"
  service      = "HTTPS"
  protocol     = "tcp"
  port         = "443"
  port_range   = "443"
  addresses    = "0.0.0.0/0"
  source_range = "0.0.0.0/0"
}

resource "utho_firewall_server" "web" {
  firewall_id = utho_firewall.web.id
  cloud_id    = utho_cloud.web.id
}
```

See [`utho_firewall_rule`](../resources/firewall_rule) and [`utho_firewall_server`](../resources/firewall_server) for all arguments.

### `utho_loadbalancer`

| v0.6.x argument | v0.7.0 argument | Notes |
|-----------------|-----------------|-------|
| `vpc_id` | `vpc` | Renamed. |
| `cpu_model` | — | Removed. |
| `name`, `type`, `dcslug`, `enable_publicip`, `firewall` | Unchanged | |

Frontends, backends, ACLs, and settings are now separate resources: [`utho_loadbalancer_frontend`](../resources/loadbalancer_frontend), [`utho_loadbalancer_backend`](../resources/loadbalancer_backend), [`utho_loadbalancer_acl`](../resources/loadbalancer_acl), and [`utho_loadbalancer_settings`](../resources/loadbalancer_settings).

### `utho_vpc`

The arguments `name`, `network`, `size`, `dcslug`, and `planid` are unchanged, but the resource's internal state format changed. Subnets, route tables, NAT gateways, and peering are now managed with [`utho_subnet`](../resources/subnet), [`utho_route_table`](../resources/route_table), [`utho_nat_gateway`](../resources/nat_gateway), and [`utho_vpc_peering`](../resources/vpc_peering).

## Step-by-Step Upgrade

### Step 1: Pin the new version

```hcl
terraform {
  required_version = ">= 1.7"

  required_providers {
    utho = {
      source  = "nitinuthocloud/utho"
      version = "~> 0.7"
    }
  }
}
```

```bash
terraform init -upgrade
```

### Step 2: Move cloud instances and auto scaling groups (no downtime)

For every `utho_cloud_instance` and `utho_auto_scaling`, do three things in one change:

1. Rename the resource block and its arguments as shown above.
2. Add an `import` block that points the new resource at the existing object's ID.
3. Add a `removed` block so Terraform forgets the old address **without destroying it**.

```hcl
# 1. The new resource
resource "utho_cloud" "web" {
  hostname     = "web-01.example.com"
  dcslug       = "inmumbaizone2"
  planid       = "10308"
  image        = "ubuntu-22.04-x86_64"
  auth         = "option2"
  sshkeys      = "4321"
  billingcycle = "hourly"
}

# 2. Adopt the existing instance
import {
  to = utho_cloud.web
  id = "1671990" # instance ID from the console or `terraform state show`
}

# 3. Forget the old address without destroying the server
removed {
  from = utho_cloud_instance.web

  lifecycle {
    destroy = false
  }
}
```

Find the existing ID with:

```bash
terraform state show utho_cloud_instance.web | grep '^ *id '
```

Then review the plan:

```bash
terraform plan
```

You should see **1 to import** and **0 to destroy** for each instance. Some write-only arguments (`auth`, `planid`, `image`, `root_password`) may show as an in-place update after import because the API does not return them. That is expected and does not recreate the server.

The same pattern works for auto scaling groups:

```hcl
import {
  to = utho_autoscaling.web
  id = "12345"
}

removed {
  from = utho_auto_scaling.web

  lifecycle {
    destroy = false
  }
}
```

### Step 3: Release resources that cannot be imported yet

Domains, DNS records, firewalls, load balancers, VPCs, and target groups do not support import in v0.7.0. For each one, choose:

* **Keep the infrastructure, stop managing it with Terraform for now.** Add a `removed` block with `destroy = false`. The resource keeps running and stays untouched, and you can bring it back under Terraform once import support is added.

    ```hcl
    removed {
      from = utho_firewall.web

      lifecycle {
        destroy = false
      }
    }
    ```

* **Recreate it under the new schema.** Write the new resources, add a `removed` block with `destroy = false` for the old address, apply, then delete the old object in the Utho Console once traffic has moved. This is the cleanest option for firewalls and DNS records, which are cheap to recreate.

!> **Do not** simply delete the old block from your configuration. Without a `removed` block, Terraform will destroy the underlying infrastructure.

### Step 4: Replace data sources

```hcl
# v0.6.x
data "utho_images" "all" {}

# v0.7.0
data "utho_cloud_images" "all" {}
```

### Step 5: Apply and verify

```bash
terraform plan   # expect imports and in-place updates only, no destroys
terraform apply
terraform plan   # should now report: No changes
```

Once the apply succeeds you can delete the `import` and `removed` blocks, or keep them as a record of the migration.

## Resources That Must Be Recreated

| Resource | Impact of recreating | Recommendation |
|----------|---------------------|----------------|
| `utho_domain` / `utho_dns_record` | Brief DNS propagation delay | Release from state with `removed`, recreate in a low-traffic window. |
| `utho_firewall` | Instances briefly unprotected or unreachable | Create the new security group and rules first, attach them, then detach the old one. |
| `utho_loadbalancer` | New public IP | Create the new load balancer, switch DNS, then remove the old one. |
| `utho_vpc` | Instances must move subnets | Release from state with `removed` and keep the existing VPC. |
| `utho_target_group` | — | Release from state and read it with the `utho_target_groups` data source. |

## Rolling Back

Terraform does not change your state until you run `terraform apply`. If the plan shows anything you don't expect, especially a destroy, stop and revert your configuration:

```bash
git checkout -- .
terraform init -upgrade
terraform plan   # should report: No changes
```

Keep the `terraform.tfstate.backup-v0.6` file from [Before You Begin](#before-you-begin) until the migration is complete.

## Getting Help

* [Importing Existing Resources](importing-resources)
* [Debugging and Troubleshooting](troubleshooting)
* [Open an issue](https://github.com/nitinuthocloud/terraform-provider-utho/issues) with your `terraform plan` output (remove any secrets first).
