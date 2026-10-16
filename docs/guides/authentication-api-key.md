---
page_title: "Authenticating with an API Key"
subcategory: "Authentication"
description: |-
  Create a Utho API key and supply it to the Utho provider securely.
---

# Authenticating with an API Key

The Utho provider authenticates every request with a Utho API key. This guide covers creating a key and the supported ways of passing it to Terraform, from most to least recommended.

## Creating an API key

1. Log in to the [Utho Console](https://console.utho.com).
2. Go to **Settings → API Tokens**.
3. Generate a new token and copy it. The full value is only shown once.

Choose the narrowest access that works for the job. A token used only to read data, for example by data sources in a reporting pipeline, does not need write access.

## Passing the key to the provider

The provider's `api_key` argument is required. Declare a sensitive input variable and pass it in:

```hcl
variable "utho_api_key" {
  description = "Utho API key"
  type        = string
  sensitive   = true
}

provider "utho" {
  api_key = var.utho_api_key
}
```

Marking the variable `sensitive` keeps its value out of plan and apply output. You then have two ways to give the variable a value.

## Option 1: Environment variable (recommended)

Terraform reads any environment variable named `TF_VAR_<name>` as the value of input variable `<name>`:

```bash
export TF_VAR_utho_api_key="your-api-key"
terraform plan
```

The key never appears in your configuration files, and the same configuration works unchanged on laptops, CI runners, and HCP Terraform.

## Option 2: A `.tfvars` file

Put the key in a `terraform.tfvars` file that is listed in `.gitignore`:

```hcl
utho_api_key = "your-api-key"
```

Terraform loads `terraform.tfvars` automatically.

## Option 3: Secrets manager

In larger teams, store the key in a secrets manager and read it at plan time. For example, with HashiCorp Vault:

```hcl
data "vault_kv_secret_v2" "utho" {
  mount = "secret"
  name  = "utho"
}

provider "utho" {
  api_key = data.vault_kv_secret_v2.utho.data["api_key"]
}
```

## Missing or empty keys

If `api_key` is omitted from the provider block, Terraform reports a missing required argument. If the variable is empty, the provider fails with a `Missing API Key` error.

## Using multiple accounts

Use provider aliases to manage resources in more than one Utho account from a single configuration:

```hcl
variable "utho_api_key_prod" {
  type      = string
  sensitive = true
}

variable "utho_api_key_staging" {
  type      = string
  sensitive = true
}

provider "utho" {
  alias   = "prod"
  api_key = var.utho_api_key_prod
}

provider "utho" {
  alias   = "staging"
  api_key = var.utho_api_key_staging
}

resource "utho_vpc" "staging" {
  provider = utho.staging

  name    = "staging"
  network = "10.1.0.0"
  size    = "16"
  dcslug  = "inmumbaizone2"
  planid  = "1008"
}
```

## Security recommendations

!> **Warning:** Never commit API keys to version control. If a key is exposed, revoke it in the Utho Console immediately and create a new one.

* Add `*.tfvars`, `.terraform/`, and `*.tfstate*` to `.gitignore`.
* Use separate keys for each environment and each pipeline, so one can be revoked without affecting the others.
* Rotate keys periodically. You can manage tokens themselves with the [`utho_api_token`](../resources/api_token) resource.
* Store state in an encrypted remote backend. State can contain sensitive attributes such as database passwords and kubeconfigs.
