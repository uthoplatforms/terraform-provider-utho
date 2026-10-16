---
page_title: "Utho: utho_api_token"
subcategory: "IAM"
description: |-
  Create and manage Utho API tokens for programmatic access.
---

# utho_api_token

Creates and manages a Utho API token. API tokens authenticate API requests from CI/CD pipelines, scripts, and third-party tools. You can create read-only tokens for monitoring tools and read-write tokens for deployment pipelines.

~> **Important:** The token value (`token` attribute) is only available at creation time. Save it to a secrets manager immediately — it cannot be retrieved later.

## Example Usage

### CI/CD deployment token

```hcl
resource "utho_api_token" "github_actions" {
  name  = "github-actions-deploy"
  write = "on"    # read + write access
}

# Store this in GitHub Actions secrets as UTHO_API_KEY
output "deploy_token" {
  value     = utho_api_token.github_actions.token
  sensitive = true
}
```

### Read-only monitoring token

```hcl
resource "utho_api_token" "monitoring" {
  name  = "datadog-monitoring"
  write = "off"   # read only
}
```

### Multiple tokens for different teams

```hcl
resource "utho_api_token" "infra_team" {
  name  = "infra-team"
  write = "on"
}

resource "utho_api_token" "security_audit" {
  name  = "security-audit-readonly"
  write = "off"
}
```

### Store token in a secret

```hcl
resource "utho_api_token" "deploy" {
  name  = "deploy-pipeline"
  write = "on"
}

# Example: store in AWS Secrets Manager
resource "aws_secretsmanager_secret_version" "utho_key" {
  secret_id     = aws_secretsmanager_secret.utho.id
  secret_string = utho_api_token.deploy.token
}
```

## Argument Reference

| Argument | Type   | Required | Description |
|----------|--------|----------|-------------|
| `name`   | String | Yes      | Token name for identification. Changing this forces a new resource. |
| `write`  | String | Yes      | Access level: `"on"` (read and write) or `"off"` (read only). Changing this forces a new resource. |

## Attribute Reference

| Attribute    | Type   | Description |
|--------------|--------|-------------|
| `id`         | String | Unique token ID. |
| `token`      | String | The API token value (starts with `live_`). **Sensitive — only shown once at creation.** |
| `created_at` | String | Creation timestamp. |

## Using the Token

```bash
# Set as environment variable for Terraform
export TF_VAR_utho_api_key="live_your_token_here"

# Or use directly in API calls
curl -H "Authorization: Bearer live_your_token_here" https://api.utho.com/v2/cloud
```

## Best Practices

- Create separate tokens for each service or team — don't share tokens
- Use `write = "off"` for anything that only needs to read data
- Rotate tokens periodically by creating a new one, updating the consumer, then deleting the old one
- Store tokens in a secrets manager (HashiCorp Vault, AWS Secrets Manager, GitHub Secrets) — never in `.tf` files or version control
