# Changelog

All notable changes to this provider are documented in this file.
The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/), and this project follows [Semantic Versioning](https://semver.org/).

## [0.7.1] - 2026-10-07

### Added

- `README.md` with quick start, service overview, guides, examples, and development instructions.
- Guide: [Migrating from v0.6.x](docs/guides/migrating-from-v0-6.md), with a resource and argument mapping and a no-downtime upgrade path using `import` and `removed` blocks.
- Guide: [Best Practices for Production](docs/guides/best-practices.md).
- Guide: [Examples Catalog](docs/guides/examples.md), listing all 42 example configurations.
- Provider docs: Supported Services table, Requirements section, and `UTHO_API_KEY` authentication.
- `CONTRIBUTING.md`, `SECURITY.md`, and GitHub issue templates.

### Changed

- Version constraints in docs and examples updated to `~> 0.7`.

## [0.7.0] - 2026-10-07

Complete rewrite of the provider.

### Added

- 53 resources and 18 data sources across Cloud Instances, Auto Scaling, Kubernetes, Database (DBaaS), VPC and Networking, Load Balancing, Security Groups, Elastic IP, VPN (IPSec), DNS, SSL Certificates, Object Storage, Container Registry, Monitoring, IAM, Projects, and Billing.
- 10 guides and 42 runnable examples.
- `UTHO_API_KEY` environment variable support.
- Import support for `utho_cloud`, `utho_autoscaling`, and `utho_nat_gateway`.

### Changed (breaking)

- `utho_cloud_instance` → `utho_cloud`
- `utho_auto_scaling` → `utho_autoscaling`
- `utho_domain` → `utho_dns_zone`
- `utho_firewall` rules and attachments split into `utho_firewall_rule` and `utho_firewall_server`.
- `utho_loadbalancer` frontends and backends split into `utho_loadbalancer_frontend` and `utho_loadbalancer_backend`.
- `utho_images` data source → `utho_cloud_images`.
- `api_key` is now optional when `UTHO_API_KEY` is set.

### Removed

- `utho_target_group` resource (use the `utho_target_groups` data source).
- `utho_account` and `utho_object_storage_plan` data sources.

See the [migration guide](docs/guides/migrating-from-v0-6.md) for details.

## [0.6.4] - 2025-05-25

Last release of the v0.6 series. See the [v0.6.4 release](https://github.com/uthoplatforms/terraform-provider-utho/releases/tag/v0.6.4).

[0.7.1]: https://github.com/nitinuthocloud/terraform-provider-utho/compare/v0.7.0...v0.7.1
[0.7.0]: https://github.com/nitinuthocloud/terraform-provider-utho/releases/tag/v0.7.0
[0.6.4]: https://github.com/uthoplatforms/terraform-provider-utho/releases/tag/v0.6.4
