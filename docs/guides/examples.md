---
page_title: "Examples Catalog"
subcategory: ""
description: |-
  Ready-to-run Terraform configurations for every Utho service, from a single instance to full production stacks.
---

# Examples Catalog

The provider repository ships 42 ready-to-run configurations across 16 services. Each one is a self-contained Terraform root module: clone the repository, `cd` into the folder, set your API key, and run `terraform apply`.

```bash
git clone https://github.com/nitinuthocloud/terraform-provider-utho.git
cd terraform-provider-utho/examples/cloud/basic

export UTHO_API_KEY="your-api-key"
terraform init
terraform plan
terraform apply
```

~> **Note:** Examples create real, billable infrastructure. Run `terraform destroy` when you are done experimenting.

Most examples come in three sizes:

* **basic**: the smallest working configuration for a service. Start here.
* **with-…**: one service combined with another, such as an instance inside a VPC.
* **full-stack**: a production-style layout with networking, security groups, and monitoring.

## Auto Scaling

| Example | What it creates | Resources used |
|---------|-----------------|----------------|
| [Basic Auto Scaling Group with Snapshot](https://github.com/nitinuthocloud/terraform-provider-utho/tree/test/new-provider/examples/autoscaling/basic) | — | `utho_autoscaling` |
| [Full Production Auto Scaling Stack](https://github.com/nitinuthocloud/terraform-provider-utho/tree/test/new-provider/examples/autoscaling/full-stack) | — | `utho_autoscaling`, `utho_firewall`, `utho_subnet`, `utho_vpc` |
| [ASG with Load Balancer + Security Group](https://github.com/nitinuthocloud/terraform-provider-utho/tree/test/new-provider/examples/autoscaling/with-lb) | — | `utho_autoscaling` |

## Billing

| Example | What it creates | Resources used |
|---------|-----------------|----------------|
| [Billing Module](https://github.com/nitinuthocloud/terraform-provider-utho/tree/test/new-provider/examples/billing) | — |  |

## Cloud Instances

| Example | What it creates | Resources used |
|---------|-----------------|----------------|
| [Basic Cloud Instance (Password Authentication)](https://github.com/nitinuthocloud/terraform-provider-utho/tree/test/new-provider/examples/cloud/basic) | 1 cloud instance (2 vCPU / 4 GB RAM / 80 GB disk)<br>Public IP address<br>Root password authentication | `utho_cloud` |
| [Full Stack Cloud Instance](https://github.com/nitinuthocloud/terraform-provider-utho/tree/test/new-provider/examples/cloud/full-stack) | 1 SSH key<br>1 VPC + 1 public subnet<br>1 security group + 3 firewall rules<br>1 cloud instance (fully configured) | `utho_cloud`, `utho_firewall`, `utho_firewall_rule`, `utho_ssh_key`, `utho_subnet`, `utho_vpc` |
| [Multiple Cloud Instances (Horizontal Scaling)](https://github.com/nitinuthocloud/terraform-provider-utho/tree/test/new-provider/examples/cloud/multi-instance) | 1 SSH key (shared across all instances)<br>1 security group + rules (shared)<br>N cloud instances (default: 3) | `utho_cloud`, `utho_firewall`, `utho_firewall_rule`, `utho_ssh_key` |
| [Cloud Instance with SSH Key Authentication](https://github.com/nitinuthocloud/terraform-provider-utho/tree/test/new-provider/examples/cloud/ssh-auth) | 1 SSH key (imported from your local machine)<br>1 cloud instance with SSH key auth<br>Public IP address | `utho_cloud`, `utho_ssh_key` |
| [Cloud Instance with EBS Block Storage](https://github.com/nitinuthocloud/terraform-provider-utho/tree/test/new-provider/examples/cloud/with-ebs) | 1 SSH key<br>1 cloud instance on an EBS plan (no included disk)<br>1 root EBS volume (80 GB NVMe)<br>1 data EBS volume (100 GB NVMe) | `utho_cloud`, `utho_ssh_key` |
| [Cloud Instance with Security Group](https://github.com/nitinuthocloud/terraform-provider-utho/tree/test/new-provider/examples/cloud/with-security-group) | 1 SSH key<br>1 security group (firewall)<br>3 firewall rules (SSH, HTTP, HTTPS)<br>1 cloud instance with the security group attached | `utho_cloud`, `utho_firewall`, `utho_firewall_rule`, `utho_ssh_key` |
| [Cloud Instance inside a VPC](https://github.com/nitinuthocloud/terraform-provider-utho/tree/test/new-provider/examples/cloud/with-vpc) | 1 VPC (private network)<br>1 public subnet inside the VPC<br>1 SSH key<br>1 cloud instance inside the VPC subnet | `utho_cloud`, `utho_ssh_key`, `utho_subnet`, `utho_vpc` |

## Container Registry

| Example | What it creates | Resources used |
|---------|-----------------|----------------|
| [Basic Container Registry](https://github.com/nitinuthocloud/terraform-provider-utho/tree/test/new-provider/examples/container_registry/basic) | — | `utho_container_registry` |
| [Full Container Registry Stack](https://github.com/nitinuthocloud/terraform-provider-utho/tree/test/new-provider/examples/container_registry/full-stack) | — | `utho_container_registry`, `utho_container_registry_immutable_rule`, `utho_container_registry_robot`, `utho_container_registry_webhook` |

## DNS

| Example | What it creates | Resources used |
|---------|-----------------|----------------|
| [Basic DNS Zone with Common Records](https://github.com/nitinuthocloud/terraform-provider-utho/tree/test/new-provider/examples/dns/basic) | — | `utho_dns_record`, `utho_dns_zone` |
| [Email DNS Setup (Google Workspace / Gmail)](https://github.com/nitinuthocloud/terraform-provider-utho/tree/test/new-provider/examples/dns/email-setup) | — | `utho_dns_record`, `utho_dns_zone` |
| [Full Production DNS Setup](https://github.com/nitinuthocloud/terraform-provider-utho/tree/test/new-provider/examples/dns/full-stack) | — | `utho_dns_record`, `utho_dns_zone` |

## Database (DBaaS)

| Example | What it creates | Resources used |
|---------|-----------------|----------------|
| [Basic MySQL Database Cluster](https://github.com/nitinuthocloud/terraform-provider-utho/tree/test/new-provider/examples/database/basic) | — | `utho_database` |
| [Full Production Database Stack](https://github.com/nitinuthocloud/terraform-provider-utho/tree/test/new-provider/examples/database/full-stack) | — | `utho_database`, `utho_database_db`, `utho_database_pool`, `utho_database_user` |
| [App Server + Database](https://github.com/nitinuthocloud/terraform-provider-utho/tree/test/new-provider/examples/database/with-app) | — | `utho_cloud`, `utho_database`, `utho_ssh_key` |

## Elastic IP

| Example | What it creates | Resources used |
|---------|-----------------|----------------|
| [Basic Elastic IP Allocation](https://github.com/nitinuthocloud/terraform-provider-utho/tree/test/new-provider/examples/elastic-ip/basic) | — | `utho_elastic_ip` |
| [Elastic IP Attached to Cloud Instance](https://github.com/nitinuthocloud/terraform-provider-utho/tree/test/new-provider/examples/elastic-ip/with-instance) | — | `utho_cloud`, `utho_elastic_ip`, `utho_ssh_key` |

## IAM

| Example | What it creates | Resources used |
|---------|-----------------|----------------|
| [Basic IAM User](https://github.com/nitinuthocloud/terraform-provider-utho/tree/test/new-provider/examples/iam/basic) | — | `utho_iam_user` |
| [Multiple IAM Users with Different Roles](https://github.com/nitinuthocloud/terraform-provider-utho/tree/test/new-provider/examples/iam/full-stack) | — | `utho_iam_user` |

## Load Balancing

| Example | What it creates | Resources used |
|---------|-----------------|----------------|
| [Basic Load Balancer (HTTP)](https://github.com/nitinuthocloud/terraform-provider-utho/tree/test/new-provider/examples/loadbalancer/basic) | 1 VPC + 1 subnet (required for LB)<br>1 SSH key<br>2 cloud instances (backends)<br>1 application load balancer<br>1 HTTP frontend (port 80)<br>2 backends | `utho_cloud`, `utho_loadbalancer`, `utho_loadbalancer_backend`, `utho_loadbalancer_frontend`, `utho_ssh_key`, `utho_subnet`, `utho_vpc` |
| [Full Production Load Balancer Stack](https://github.com/nitinuthocloud/terraform-provider-utho/tree/test/new-provider/examples/loadbalancer/full-stack) | — | `utho_cloud`, `utho_firewall`, `utho_firewall_rule`, `utho_loadbalancer`, `utho_loadbalancer_backend`, `utho_loadbalancer_frontend`, `utho_loadbalancer_settings`, `utho_ssh_key`, `utho_subnet`, `utho_vpc` |
| [Load Balancer with Weighted Backends](https://github.com/nitinuthocloud/terraform-provider-utho/tree/test/new-provider/examples/loadbalancer/multi-backend) | — | `utho_cloud`, `utho_loadbalancer`, `utho_loadbalancer_backend`, `utho_loadbalancer_frontend`, `utho_ssh_key`, `utho_subnet`, `utho_vpc` |
| [Load Balancer with HTTPS + SSL Termination](https://github.com/nitinuthocloud/terraform-provider-utho/tree/test/new-provider/examples/loadbalancer/with-ssl) | — | `utho_cloud`, `utho_loadbalancer`, `utho_loadbalancer_backend`, `utho_loadbalancer_frontend`, `utho_loadbalancer_settings`, `utho_ssh_key`, `utho_ssl_certificate`, `utho_subnet`, `utho_vpc` |

## Monitoring

| Example | What it creates | Resources used |
|---------|-----------------|----------------|
| [Basic Monitoring Setup](https://github.com/nitinuthocloud/terraform-provider-utho/tree/test/new-provider/examples/monitoring/basic) | — | `utho_alert`, `utho_alert_contact` |
| [Full Production Monitoring Stack](https://github.com/nitinuthocloud/terraform-provider-utho/tree/test/new-provider/examples/monitoring/full-stack) | — | `utho_alert`, `utho_alert_contact` |

## NAT Gateway

| Example | What it creates | Resources used |
|---------|-----------------|----------------|
| [Basic NAT Gateway](https://github.com/nitinuthocloud/terraform-provider-utho/tree/test/new-provider/examples/nat-gateway/basic) | — | `utho_elastic_ip`, `utho_nat_gateway`, `utho_subnet`, `utho_vpc` |
| [Full Private Network with NAT Gateway](https://github.com/nitinuthocloud/terraform-provider-utho/tree/test/new-provider/examples/nat-gateway/full-stack) | — | `utho_cloud`, `utho_elastic_ip`, `utho_nat_gateway`, `utho_ssh_key`, `utho_subnet`, `utho_vpc` |

## Object Storage

| Example | What it creates | Resources used |
|---------|-----------------|----------------|
| [Basic Object Storage Bucket](https://github.com/nitinuthocloud/terraform-provider-utho/tree/test/new-provider/examples/object-storage/basic) | — | `utho_object_storage`, `utho_object_storage_key`, `utho_object_storage_permission` |
| [Full Production Object Storage Stack](https://github.com/nitinuthocloud/terraform-provider-utho/tree/test/new-provider/examples/object-storage/full-stack) | — | `utho_object_storage`, `utho_object_storage_key`, `utho_object_storage_permission` |
| [Multi-Bucket Deployment](https://github.com/nitinuthocloud/terraform-provider-utho/tree/test/new-provider/examples/object-storage/multi-bucket) | — | `utho_object_storage` |
| [Versioned Bucket for Backups](https://github.com/nitinuthocloud/terraform-provider-utho/tree/test/new-provider/examples/object-storage/with-versioning) | — | `utho_object_storage`, `utho_object_storage_key`, `utho_object_storage_permission` |

## Projects

| Example | What it creates | Resources used |
|---------|-----------------|----------------|
| [Basic Project with Member](https://github.com/nitinuthocloud/terraform-provider-utho/tree/test/new-provider/examples/project/basic) | — | `utho_project`, `utho_project_member` |
| [Multi-Environment Projects](https://github.com/nitinuthocloud/terraform-provider-utho/tree/test/new-provider/examples/project/multi-env) | — | `utho_project`, `utho_project_member` |

## Route Tables

| Example | What it creates | Resources used |
|---------|-----------------|----------------|
| [Basic Route Table with IGW Route](https://github.com/nitinuthocloud/terraform-provider-utho/tree/test/new-provider/examples/route_table/basic) | — | `utho_route`, `utho_route_table` |
| [Full VPC Stack with Route Table and Routes](https://github.com/nitinuthocloud/terraform-provider-utho/tree/test/new-provider/examples/route_table/full-stack) | — | `utho_route`, `utho_route_table`, `utho_subnet`, `utho_vpc` |

## VPC Peering

| Example | What it creates | Resources used |
|---------|-----------------|----------------|
| [VPC Peering Connection](https://github.com/nitinuthocloud/terraform-provider-utho/tree/test/new-provider/examples/vpc_peering/basic) | — | `utho_vpc`, `utho_vpc_peering` |

## VPN (IPSec)

| Example | What it creates | Resources used |
|---------|-----------------|----------------|
| [Basic IPSec Tunnel](https://github.com/nitinuthocloud/terraform-provider-utho/tree/test/new-provider/examples/ipsec/basic) | — | `utho_ipsec` |
| [Site-to-Site VPN with Two IPSec Tunnels](https://github.com/nitinuthocloud/terraform-provider-utho/tree/test/new-provider/examples/ipsec/site-to-site) | — | `utho_ipsec`, `utho_ipsec_connection` |

## Contributing an Example

Have a configuration that others would find useful? Open a pull request that adds a folder under `examples/<service>/<name>/` with a `main.tf` that starts with an `# Example:` header and a `What this creates:` list, like the existing examples.
