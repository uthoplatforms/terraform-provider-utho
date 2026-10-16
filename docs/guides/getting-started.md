---
page_title: "Getting Started with the Utho Provider"
subcategory: ""
description: |-
  Install the Utho provider, authenticate, and deploy your first cloud instance with Terraform.
---

# Getting Started with the Utho Provider

This guide walks you through installing the Utho provider, authenticating with the Utho API, and deploying your first cloud instance. By the end you will have a running virtual machine managed by Terraform, and you will know how to change and destroy it safely.

## Prerequisites

* [Terraform](https://developer.hashicorp.com/terraform/install) 1.0 or later.
* A [Utho Cloud](https://utho.com) account.
* A Utho API key. See [Authenticating with an API Key](authentication-api-key).
* An SSH key pair on your machine, for example `~/.ssh/id_ed25519.pub`.

## Step 1: Configure the provider

Create a new directory and add a file named `main.tf`:

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

provider "utho" {
  api_key = var.utho_api_key
}
```

The `required_providers` block tells Terraform where to download the provider from. The version constraint `~> 0.2` allows any `0.x` release from `0.2` onward, so you receive bug fixes and new features. Pin to an exact version, for example `version = "0.2.46"`, if you need fully reproducible runs.

The provider reads your API key from the `utho_api_key` variable. Marking it `sensitive` keeps the value out of plan and apply output.

## Step 2: Set your API key

Terraform reads any environment variable named `TF_VAR_<name>` as the value of input variable `<name>`:

```bash
export TF_VAR_utho_api_key="your-api-key"
```

## Step 3: Initialize the working directory

```bash
terraform init
```

Terraform downloads the Utho provider and records the exact version it selected in `.terraform.lock.hcl`. Commit this lock file to version control so everyone on your team uses the same provider version.

## Step 4: Define your first instance

Add the following to `main.tf`:

```hcl
resource "utho_ssh_key" "deploy" {
  name   = "deploy-key"
  sshkey = file("~/.ssh/id_ed25519.pub")
}

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

This configuration uploads your public SSH key and creates an Ubuntu 22.04 instance in Mumbai that you can log in to with that key. The reference `utho_ssh_key.deploy.id` tells Terraform to create the key before the instance.

-> **Tip:** The values for `dcslug`, `planid`, and `image` vary by account and data center. See [Data Centers, Plans, and Images](data-centers-plans-images) to look them up with data sources instead of hard-coding them.

## Step 5: Review and apply

Preview the changes Terraform will make:

```bash
terraform plan
```

The plan should show two resources to add. Apply it:

```bash
terraform apply
```

Type `yes` when prompted. When the apply finishes, Terraform prints the instance's public IP:

```text
Outputs:

web_ip = "203.0.113.25"
```

Connect to it:

```bash
ssh root@$(terraform output -raw web_ip)
```

## Step 6: Make a change

Terraform compares your configuration to the real infrastructure on every run. To attach a security group that allows HTTP traffic, add:

```hcl
resource "utho_firewall" "web" {
  name = "web-sg"
}

resource "utho_firewall_rule" "http" {
  firewall_id  = utho_firewall.web.id
  type         = "incoming"
  service      = "HTTP"
  protocol     = "tcp"
  port         = "80"
  port_range   = "80"
  addresses    = "0.0.0.0/0"
  source_range = "0.0.0.0/0"
}

resource "utho_firewall_server" "web" {
  firewall_id = utho_firewall.web.id
  cloud_id    = utho_cloud.web.id
}
```

Run `terraform plan` again. Terraform shows only the three new resources; the existing instance is left untouched.

## Step 7: Clean up

To delete everything this configuration created:

```bash
terraform destroy
```

## Next steps

* Store state remotely so your team can collaborate. See [Authenticating in CI/CD and HCP Terraform](authentication-ci-cd).
* Build a complete environment with the [Production Web Stack](production-web-stack) tutorial.
* Isolate workloads on a private network with [Private Networking with VPC and NAT](private-networking).
* Browse the resources and data sources for each service in the navigation on the left.
