---
page_title: "Provider: Utho"
description: |-
  The Utho provider is used to manage infrastructure on Utho Cloud, including cloud instances, Kubernetes, managed databases, networking, storage, and more.
---

# Utho Provider

The Utho provider is used to configure infrastructure on [Utho Cloud](https://utho.com) using the Utho API. With it you can define cloud instances, managed Kubernetes clusters, managed databases, VPCs, load balancers, DNS, object storage, container registries, monitoring, and IAM as code, and manage their full lifecycle with Terraform.

Documentation for the resources and data sources supported by the Utho provider can be found in the navigation to the left, grouped by service.

To learn the basics of Terraform using this provider, follow the [Getting Started](guides/getting-started) guide. For end-to-end examples, see the tutorials under **Guides**, such as [Deploying a Production Web Stack](guides/production-web-stack) and [Private Networking with VPC and NAT](guides/private-networking).

-> **Upgrading from v0.6.x?** Resource names and several arguments changed in v0.7.0. Read the [Migrating from v0.6.x](guides/migrating-from-v0-6) guide before you upgrade an existing configuration.

## Supported Services

| Service | Resources | Data Sources |
|---------|-----------|--------------|
| **Cloud Instances** | [`utho_cloud`](resources/cloud), [`utho_cloud_ebs`](resources/cloud_ebs), [`utho_cloud_storage`](resources/cloud_storage), [`utho_cloud_snapshot`](resources/cloud_snapshot), [`utho_cloud_iso`](resources/cloud_iso), [`utho_cloud_power`](resources/cloud_power), [`utho_cloud_resize`](resources/cloud_resize), [`utho_cloud_public_ip`](resources/cloud_public_ip), [`utho_cloud_vpc`](resources/cloud_vpc), [`utho_cloud_firewall`](resources/cloud_firewall), [`utho_ssh_key`](resources/ssh_key) | [`utho_clouds`](data-sources/clouds), [`utho_cloud_plans`](data-sources/cloud_plans), [`utho_cloud_images`](data-sources/cloud_images), [`utho_cloud_dczones`](data-sources/cloud_dczones), [`utho_cloud_isos`](data-sources/cloud_isos), [`utho_cloud_snapshots`](data-sources/cloud_snapshots), [`utho_ssh_keys`](data-sources/ssh_keys) |
| **Auto Scaling** | [`utho_autoscaling`](resources/autoscaling), [`utho_autoscaling_policy`](resources/autoscaling_policy), [`utho_autoscaling_schedule`](resources/autoscaling_schedule) | [`utho_autoscalings`](data-sources/autoscalings), [`utho_target_groups`](data-sources/target_groups) |
| **Kubernetes** | [`utho_kubernetes`](resources/kubernetes), [`utho_kubernetes_node_pool`](resources/kubernetes_node_pool) | [`utho_kubernetes_config`](data-sources/kubernetes_config) |
| **Database (DBaaS)** | [`utho_database`](resources/database), [`utho_database_db`](resources/database_db), [`utho_database_user`](resources/database_user), [`utho_database_pool`](resources/database_pool) | [`utho_database_plans`](data-sources/database_plans) |
| **VPC and Networking** | [`utho_vpc`](resources/vpc), [`utho_subnet`](resources/subnet), [`utho_nat_gateway`](resources/nat_gateway), [`utho_route_table`](resources/route_table), [`utho_route`](resources/route), [`utho_vpc_peering`](resources/vpc_peering) | [`utho_vpcs`](data-sources/vpcs), [`utho_vpc_subnets`](data-sources/vpc_subnets) |
| **Load Balancing** | [`utho_loadbalancer`](resources/loadbalancer), [`utho_loadbalancer_frontend`](resources/loadbalancer_frontend), [`utho_loadbalancer_backend`](resources/loadbalancer_backend), [`utho_loadbalancer_acl`](resources/loadbalancer_acl), [`utho_loadbalancer_settings`](resources/loadbalancer_settings) | |
| **Security Groups** | [`utho_firewall`](resources/firewall), [`utho_firewall_rule`](resources/firewall_rule), [`utho_firewall_server`](resources/firewall_server) | [`utho_firewalls`](data-sources/firewalls) |
| **Elastic IP** | [`utho_elastic_ip`](resources/elastic_ip) | |
| **Elastic Block Storage** | [`utho_ebs`](resources/ebs), [`utho_ebs_attachment`](resources/ebs_attachment) | [`utho_ebs_dczones`](data-sources/ebs_dczones) |
| **VPN (IPSec)** | [`utho_ipsec`](resources/ipsec), [`utho_ipsec_connection`](resources/ipsec_connection) | |
| **DNS** | [`utho_dns_zone`](resources/dns_zone), [`utho_dns_record`](resources/dns_record) | |
| **SSL Certificates** | [`utho_ssl_certificate`](resources/ssl_certificate) | |
| **Object Storage** | [`utho_object_storage`](resources/object_storage), [`utho_object_storage_key`](resources/object_storage_key), [`utho_object_storage_permission`](resources/object_storage_permission) | |
| **Container Registry** | [`utho_container_registry`](resources/container_registry), [`utho_container_registry_robot`](resources/container_registry_robot), [`utho_container_registry_webhook`](resources/container_registry_webhook), [`utho_container_registry_immutable_rule`](resources/container_registry_immutable_rule) | |
| **Monitoring** | [`utho_alert`](resources/alert), [`utho_alert_contact`](resources/alert_contact) | |
| **IAM** | [`utho_iam_user`](resources/iam_user), [`utho_api_token`](resources/api_token) | |
| **Projects** | [`utho_project`](resources/project), [`utho_project_member`](resources/project_member) | |
| **Billing** | | [`utho_billing_usage`](data-sources/billing_usage), [`utho_billing_invoices`](data-sources/billing_invoices), [`utho_billing_cost_by_project`](data-sources/billing_cost_by_project), [`utho_billing_cycles`](data-sources/billing_cycles) |

## Requirements

* [Terraform](https://developer.hashicorp.com/terraform/install) 1.0 or later. [`import` blocks](guides/importing-resources) need Terraform 1.5 or later.
* A Utho account and an API key.

## Example Usage

```hcl
terraform {
  required_providers {
    utho = {
      source  = "nitinuthocloud/utho"
      version = "~> 0.7"
    }
  }
}

variable "utho_api_key" {
  type      = string
  sensitive = true
}

# Configure the Utho provider
provider "utho" {
  api_key = var.utho_api_key
}

# Create an SSH key
resource "utho_ssh_key" "deploy" {
  name   = "deploy-key"
  sshkey = file("~/.ssh/id_ed25519.pub")
}

# Create a cloud instance
resource "utho_cloud" "web" {
  hostname        = "web-01.mhc"
  dcslug          = "inmumbaizone2"
  planid          = "10308"
  billingcycle    = "hourly"
  auth            = "option2"
  sshkeys         = utho_ssh_key.deploy.id
  image           = "ubuntu-22.04-x86_64"
  enable_publicip = "true"
}

output "web_ip" {
  value = utho_cloud.web.ip
}
```

## Authentication and Configuration

The Utho provider authenticates with a Utho API key. You can create one in the [Utho Console](https://console.utho.com) under **Settings → API Tokens**.

The provider reads the key from the `api_key` argument or, if that is not set, from the `UTHO_API_KEY` environment variable. Always supply it through a sensitive input variable or the environment rather than writing the key into your configuration.

!> **Warning:** Hard-coded credentials are not recommended in any Terraform configuration and risk secret leakage should the file ever be committed to a public version control system.

```hcl
variable "utho_api_key" {
  type      = string
  sensitive = true
}

provider "utho" {
  api_key = var.utho_api_key
}
```

Set the variable with the `TF_VAR_utho_api_key` environment variable, so the key never appears in your files:

```bash
export TF_VAR_utho_api_key="your-api-key"
terraform plan
```

Alternatively, leave `api_key` out of the provider block entirely and export `UTHO_API_KEY`:

```hcl
provider "utho" {}
```

```bash
export UTHO_API_KEY="your-api-key"
terraform plan
```

For more detail, including secrets managers, multiple accounts, and authenticating from CI/CD pipelines and HCP Terraform, see the [Authenticating with an API Key](guides/authentication-api-key) and [Authenticating in CI/CD and HCP Terraform](guides/authentication-ci-cd) guides.

## Argument Reference

The following arguments are supported in the `utho` provider block:

* `api_key` - (Optional) Utho API key. Supply it through a sensitive variable, as shown above. If omitted, the provider reads the `UTHO_API_KEY` environment variable. One of the two must be set.
* `base_url` - (Optional) Override the Utho API base URL. Defaults to `https://api.utho.com/v2`. This is only needed for testing against a non-production API endpoint.

## Data Centers

Most Utho resources take a `dcslug` argument that selects the data center they are created in.

| Slug | Location |
|------|----------|
| `innoida` | Delhi (Noida), India |
| `inmumbaizone2` | Mumbai, India |
| `inbangalore` | Bangalore, India |
| `defra1` | Frankfurt, Germany |
| `uslosangeles` | Los Angeles, United States |

Not every service, plan, or feature is available in every data center. Use the [`utho_cloud_dczones`](data-sources/cloud_dczones) data source to look up current availability, and see the [Data Centers, Plans, and Images](guides/data-centers-plans-images) guide for discovering valid `dcslug`, `planid`, and `image` values without hard-coding them.

## Getting Help

* Browse the [Guides](guides/getting-started) for tutorials and common workflows.
* Upgrading an existing configuration? See [Migrating from v0.6.x](guides/migrating-from-v0-6).
* Ready-to-run configurations for every service live in the [`examples`](https://github.com/nitinuthocloud/terraform-provider-utho/tree/test/new-provider/examples) directory, catalogued in the [Examples](guides/examples) guide.
* See [Debugging and Troubleshooting](guides/troubleshooting) for common errors and how to enable debug logging.
* Report bugs or request features on the provider's [GitHub repository](https://github.com/nitinuthocloud/terraform-provider-utho/issues).