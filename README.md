<p align="center">
  <a href="https://utho.com">
    <img src="https://utho.com/favicon.ico" alt="Utho" width="64" height="64">
  </a>
</p>

<h1 align="center">Terraform Provider for Utho Cloud</h1>

<p align="center">
  Manage <a href="https://utho.com">Utho Cloud</a> infrastructure as code: cloud instances, Kubernetes, managed databases, VPC, load balancers, object storage, and more.
</p>

<p align="center">
  <a href="https://registry.terraform.io/providers/nitinuthocloud/utho/latest"><img src="https://img.shields.io/badge/terraform-registry-7B42BC?logo=terraform&logoColor=white" alt="Terraform Registry"></a>
  <a href="https://github.com/nitinuthocloud/terraform-provider-utho/releases/latest"><img src="https://img.shields.io/github/v/release/nitinuthocloud/terraform-provider-utho?label=version" alt="Latest release"></a>
  <a href="https://registry.terraform.io/providers/nitinuthocloud/utho/latest"><img src="https://img.shields.io/terraform/provider/dt/2438094?label=downloads" alt="Downloads"></a>
  <a href="https://github.com/nitinuthocloud/terraform-provider-utho/actions/workflows/release.yml"><img src="https://github.com/nitinuthocloud/terraform-provider-utho/actions/workflows/release.yml/badge.svg" alt="Release"></a>
  <img src="https://img.shields.io/badge/resources-53-blue" alt="53 resources">
  <img src="https://img.shields.io/badge/data%20sources-18-blue" alt="18 data sources">
</p>

<p align="center">
  <a href="https://registry.terraform.io/providers/nitinuthocloud/utho/latest/docs">Documentation</a> ·
  <a href="https://registry.terraform.io/providers/nitinuthocloud/utho/latest/docs/guides/getting-started">Getting started</a> ·
  <a href="examples">Examples</a> ·
  <a href="https://registry.terraform.io/providers/nitinuthocloud/utho/latest/docs/guides/migrating-from-v0-6">Migrating from v0.6</a> ·
  <a href="CHANGELOG.md">Changelog</a>
</p>

---

## Quick Start

**1. Get an API key** from the [Utho Console](https://console.utho.com) under **Settings → API Tokens**, and export it:

```bash
export UTHO_API_KEY="your-api-key"
```

**2. Write your first configuration** in `main.tf`:

```hcl
terraform {
  required_providers {
    utho = {
      source  = "nitinuthocloud/utho"
      version = "~> 0.7"
    }
  }
}

provider "utho" {} # reads UTHO_API_KEY

resource "utho_ssh_key" "me" {
  name   = "my-laptop"
  sshkey = file("~/.ssh/id_ed25519.pub")
}

resource "utho_cloud" "web" {
  hostname        = "web-01"
  dcslug          = "inmumbaizone2"
  planid          = "10308"
  billingcycle    = "hourly"
  image           = "ubuntu-22.04-x86_64"
  auth            = "option2"
  sshkeys         = utho_ssh_key.me.id
  enable_publicip = "true"
}

output "ip" {
  value = utho_cloud.web.ip
}
```

**3. Deploy it:**

```bash
terraform init
terraform apply
ssh root@$(terraform output -raw ip)
```

When you're done, run `terraform destroy`.

## Supported Services

| Service | Resources | Data sources |
|---------|:---------:|:------------:|
| ☁️ **Cloud Instances**: VMs, block storage, snapshots, ISOs, power, resize, public IPs | 11 | 7 |
| 📈 **Auto Scaling**: groups, policies, schedules | 3 | 2 |
| ☸️ **Kubernetes**: managed clusters and node pools | 2 | 1 |
| 🗄️ **Database (DBaaS)**: clusters, databases, users, connection pools | 4 | 1 |
| 🌐 **VPC and Networking**: VPCs, subnets, NAT gateways, route tables, peering | 6 | 2 |
| ⚖️ **Load Balancing**: load balancers, frontends, backends, ACLs, settings | 5 | — |
| 🛡️ **Security Groups**: firewalls, rules, server attachments | 3 | 1 |
| 📍 **Elastic IP** | 1 | — |
| 🔐 **VPN (IPSec)**: tunnels and connections | 2 | — |
| 🔤 **DNS**: zones and records | 2 | — |
| 🔒 **SSL Certificates** | 1 | — |
| 🪣 **Object Storage**: S3-compatible buckets, keys, permissions | 3 | — |
| 📦 **Container Registry**: registries, robot accounts, webhooks, immutable rules | 4 | — |
| 🔔 **Monitoring**: alerts and alert contacts | 2 | — |
| 👤 **IAM**: users and API tokens | 2 | — |
| 🗂️ **Projects**: projects and members | 2 | — |
| 💳 **Billing**: usage, invoices, cost by project, billing cycles | — | 4 |
| **Total** | **53** | **18** |

The full reference for every resource and data source is on the [Terraform Registry](https://registry.terraform.io/providers/nitinuthocloud/utho/latest/docs).

## Guides

| Guide | What you'll learn |
|-------|-------------------|
| [Getting Started](https://registry.terraform.io/providers/nitinuthocloud/utho/latest/docs/guides/getting-started) | Install the provider and deploy your first server |
| [Authenticating with an API Key](https://registry.terraform.io/providers/nitinuthocloud/utho/latest/docs/guides/authentication-api-key) | Supply credentials safely, including from secret managers |
| [Authenticating in CI/CD and HCP Terraform](https://registry.terraform.io/providers/nitinuthocloud/utho/latest/docs/guides/authentication-ci-cd) | GitHub Actions, GitLab CI, and HCP Terraform workspaces |
| [Data Centers, Plans, and Images](https://registry.terraform.io/providers/nitinuthocloud/utho/latest/docs/guides/data-centers-plans-images) | Discover valid `dcslug`, `planid`, and `image` values |
| [Deploying a Production Web Stack](https://registry.terraform.io/providers/nitinuthocloud/utho/latest/docs/guides/production-web-stack) | VPC, load balancer, app servers, and a managed database |
| [Private Networking with VPC and NAT](https://registry.terraform.io/providers/nitinuthocloud/utho/latest/docs/guides/private-networking) | Private subnets with outbound internet through NAT |
| [Running a Managed Database](https://registry.terraform.io/providers/nitinuthocloud/utho/latest/docs/guides/managed-database) | DBaaS clusters, users, and connection pools |
| [Managed Kubernetes](https://registry.terraform.io/providers/nitinuthocloud/utho/latest/docs/guides/kubernetes) | Create a cluster and deploy to it with the Kubernetes and Helm providers |
| [Importing Existing Resources](https://registry.terraform.io/providers/nitinuthocloud/utho/latest/docs/guides/importing-resources) | Bring console-created infrastructure under Terraform |
| [Best Practices for Production](https://registry.terraform.io/providers/nitinuthocloud/utho/latest/docs/guides/best-practices) | State, secrets, environments, and safe changes |
| [Migrating from v0.6.x](https://registry.terraform.io/providers/nitinuthocloud/utho/latest/docs/guides/migrating-from-v0-6) | Upgrade existing configurations without downtime |
| [Debugging and Troubleshooting](https://registry.terraform.io/providers/nitinuthocloud/utho/latest/docs/guides/troubleshooting) | Common errors and debug logging |

## Examples

The [`examples`](examples) directory has 42 ready-to-run configurations across 16 services, from a single VM to complete production stacks. Each folder is a self-contained Terraform root module.

```bash
cd examples/loadbalancer/basic
export UTHO_API_KEY="your-api-key"
terraform init && terraform apply
```

| Service | Examples |
|---------|----------|
| Cloud Instances | [basic](examples/cloud/basic) · [ssh-auth](examples/cloud/ssh-auth) · [with-vpc](examples/cloud/with-vpc) · [with-security-group](examples/cloud/with-security-group) · [with-ebs](examples/cloud/with-ebs) · [multi-instance](examples/cloud/multi-instance) · [full-stack](examples/cloud/full-stack) |
| Load Balancing | [basic](examples/loadbalancer/basic) · [multi-backend](examples/loadbalancer/multi-backend) · [with-ssl](examples/loadbalancer/with-ssl) · [full-stack](examples/loadbalancer/full-stack) |
| Auto Scaling | [basic](examples/autoscaling/basic) · [with-lb](examples/autoscaling/with-lb) · [full-stack](examples/autoscaling/full-stack) |
| Database | [basic](examples/database/basic) · [with-app](examples/database/with-app) · [full-stack](examples/database/full-stack) |
| Object Storage | [basic](examples/object-storage/basic) · [multi-bucket](examples/object-storage/multi-bucket) · [with-versioning](examples/object-storage/with-versioning) · [full-stack](examples/object-storage/full-stack) |
| DNS | [basic](examples/dns/basic) · [email-setup](examples/dns/email-setup) · [full-stack](examples/dns/full-stack) |
| Networking | [nat-gateway](examples/nat-gateway) · [route_table](examples/route_table) · [vpc_peering](examples/vpc_peering) · [elastic-ip](examples/elastic-ip) · [ipsec](examples/ipsec) |
| Platform | [container_registry](examples/container_registry) · [monitoring](examples/monitoring) · [iam](examples/iam) · [project](examples/project) · [billing](examples/billing) |

See the [Examples Catalog](https://registry.terraform.io/providers/nitinuthocloud/utho/latest/docs/guides/examples) for what each one creates.

## Authentication

The provider looks for an API key in this order:

1. The `api_key` argument in the `provider "utho"` block
2. The `UTHO_API_KEY` environment variable

```hcl
variable "utho_api_key" {
  type      = string
  sensitive = true
}

provider "utho" {
  api_key = var.utho_api_key
}
```

> [!WARNING]
> Never commit API keys to version control. Use environment variables, a `sensitive` variable, or your CI platform's secret store.

| Argument | Environment variable | Required | Description |
|----------|---------------------|----------|-------------|
| `api_key` | `UTHO_API_KEY` | One of the two | Utho API key |
| `base_url` | — | No | API endpoint override. Defaults to `https://api.utho.com/v2` |

## Requirements

| Tool | Version |
|------|---------|
| [Terraform](https://developer.hashicorp.com/terraform/install) | 1.0 or later (1.5+ for `import` blocks, 1.7+ for `removed` blocks) |
| [Go](https://go.dev/dl/) (only to build from source) | See `go.mod` |

## Upgrading from v0.6.x

v0.7.0 is a major rewrite. Several resources were renamed (for example `utho_cloud_instance` → `utho_cloud`), so **do not** upgrade an existing configuration without reading the [migration guide](https://registry.terraform.io/providers/nitinuthocloud/utho/latest/docs/guides/migrating-from-v0-6). It shows how to move cloud instances and auto scaling groups in place with `import` and `removed` blocks, with no downtime.

## Developing the Provider

```bash
git clone https://github.com/nitinuthocloud/terraform-provider-utho.git
cd terraform-provider-utho
go build -o terraform-provider-utho
```

To test a local build without publishing it, point Terraform at your binary with a [development override](https://developer.hashicorp.com/terraform/cli/config/config-file#development-overrides-for-provider-developers) in `~/.terraformrc`:

```hcl
provider_installation {
  dev_overrides {
    "nitinuthocloud/utho" = "/path/to/terraform-provider-utho"
  }
  direct {}
}
```

Then run `terraform plan` in any configuration. You don't need `terraform init` when using overrides.

### Project layout

```text
.
├── main.go                    # provider entry point
├── internal/
│   ├── provider/              # provider schema and configuration
│   ├── client/                # Utho API client
│   ├── resources/             # one file per resource family
│   └── datasources/           # one file per data source family
├── docs/                      # registry documentation (resources, data sources, guides)
├── examples/                  # runnable example configurations
└── .goreleaser.yml            # release build configuration
```

### Releasing

Releases are fully automated. Pushing a `v*` tag runs the [release workflow](.github/workflows/release.yml), which builds binaries for Linux, macOS, Windows, and FreeBSD, signs the checksums with GPG, and publishes a GitHub release. The Terraform Registry picks up the new version automatically within a few minutes.

```bash
git tag -a v0.7.1 -m "v0.7.1"
git push origin v0.7.1
```

## Contributing

Bug reports, feature requests, and pull requests are welcome. See [CONTRIBUTING.md](CONTRIBUTING.md) to get started, and [SECURITY.md](SECURITY.md) to report a vulnerability privately.

## Support

* 📖 [Provider documentation](https://registry.terraform.io/providers/nitinuthocloud/utho/latest/docs)
* 🐛 [Report a bug or request a feature](https://github.com/nitinuthocloud/terraform-provider-utho/issues/new/choose)
* 🌐 [Utho Cloud](https://utho.com) · [Utho Console](https://console.utho.com) · [Utho API](https://utho.com/api-docs/)
