---
page_title: "Authenticating in CI/CD and HCP Terraform"
subcategory: "Authentication"
description: |-
  Supply a Utho API key to Terraform running in GitHub Actions, GitLab CI, and HCP Terraform.
---

# Authenticating in CI/CD and HCP Terraform

When Terraform runs in automation, the Utho API key should come from the platform's secret store and reach Terraform as the `TF_VAR_utho_api_key` environment variable, which sets the `utho_api_key` input variable. Your configuration stays the same everywhere:

```hcl
variable "utho_api_key" {
  type      = string
  sensitive = true
}

provider "utho" {
  api_key = var.utho_api_key
}
```

-> **Tip:** Create a dedicated API key for each pipeline. You can manage these with the [`utho_api_token`](../resources/api_token) resource from an administrative workspace, and revoke a single pipeline's access without affecting anything else.

## GitHub Actions

1. In your repository, go to **Settings → Secrets and variables → Actions** and add a secret named `UTHO_API_KEY`.
2. Expose it to the Terraform steps:

```yaml
name: terraform

on:
  pull_request:
  push:
    branches: [main]

jobs:
  terraform:
    runs-on: ubuntu-latest
    env:
      TF_VAR_utho_api_key: ${{ secrets.UTHO_API_KEY }}
    steps:
      - uses: actions/checkout@v4

      - uses: hashicorp/setup-terraform@v3

      - run: terraform init
      - run: terraform plan -input=false

      - if: github.ref == 'refs/heads/main' && github.event_name == 'push'
        run: terraform apply -auto-approve -input=false
```

## GitLab CI/CD

1. Go to **Settings → CI/CD → Variables** and add `TF_VAR_utho_api_key`. Mark it **Masked** and **Protected**.
2. GitLab exposes it to jobs as an environment variable automatically:

```yaml
image:
  name: hashicorp/terraform:latest
  entrypoint: [""]

stages: [plan, apply]

plan:
  stage: plan
  script:
    - terraform init
    - terraform plan -input=false -out=tfplan
  artifacts:
    paths: [tfplan]

apply:
  stage: apply
  script:
    - terraform init
    - terraform apply -input=false tfplan
  when: manual
  only: [main]
```

## HCP Terraform and Terraform Enterprise

1. Open your workspace and go to **Variables**.
2. Add a variable with:
   * Key: `utho_api_key`
   * Category: **Terraform variable**
   * Sensitive: checked

To share one key across many workspaces, add it to a **variable set** instead and apply the set to those workspaces.

Connect your configuration to the workspace with a `cloud` block:

```hcl
terraform {
  cloud {
    organization = "my-org"

    workspaces {
      name = "utho-production"
    }
  }

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

## Remote state

Whatever runs Terraform, keep state in a shared, encrypted backend rather than on a single machine. Good options include HCP Terraform, Terraform Enterprise, or the `s3` backend. Remember that state can include sensitive values such as database passwords and kubeconfigs, so restrict who can read it.
