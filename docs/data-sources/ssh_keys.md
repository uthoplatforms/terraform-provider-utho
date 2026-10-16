---
page_title: "Utho: utho_ssh_keys"
subcategory: "Cloud Instances"
description: |-
  List all SSH keys in your Utho account.
---

# utho_ssh_keys

Fetches all SSH keys imported into your account. Use this to discover existing SSH key IDs before creating cloud instances with key-based authentication.

## Example Usage

### List all SSH keys

```hcl
data "utho_ssh_keys" "all" {}

output "key_names" {
  value = data.utho_ssh_keys.all.keys[*].name
}
```

### Find a key by name and use in instance

```hcl
data "utho_ssh_keys" "all" {}

locals {
  deploy_key = one([
    for k in data.utho_ssh_keys.all.keys :
    k if k.name == "deploy-key"
  ])
}

resource "utho_cloud" "web" {
  hostname        = "web-01.mhc"
  dcslug          = "inmumbaizone2"
  planid          = "10308"
  billingcycle    = "hourly"
  auth            = "option2"
  sshkeys         = local.deploy_key.id
  image           = "ubuntu-22.04-x86_64"
  enable_publicip = "true"
  cpumodel        = "amd"
}
```

### Use existing key or create new one

```hcl
data "utho_ssh_keys" "all" {}

locals {
  # Check if key already exists
  existing_key = one([
    for k in data.utho_ssh_keys.all.keys :
    k if k.name == "ci-key"
  ])
}

# Only create if it doesn't exist
resource "utho_ssh_key" "ci" {
  count  = local.existing_key == null ? 1 : 0
  name   = "ci-key"
  sshkey = file("~/.ssh/id_ed25519.pub")
}

locals {
  key_id = local.existing_key != null ? local.existing_key.id : utho_ssh_key.ci[0].id
}
```

## Attribute Reference

### Top-level

| Attribute | Type | Description |
|-----------|------|-------------|
| `keys`    | List | List of all SSH keys in your account. |

### keys

| Attribute    | Type   | Description |
|--------------|--------|-------------|
| `id`         | String | SSH key ID. Use this as `sshkeys` in `utho_cloud`. |
| `name`       | String | SSH key name. |
| `created_at` | String | Import timestamp. |

## Notes

- The public key value is not returned in the list response for security reasons.
- To import a new SSH key use `utho_ssh_key`.
- Always use `auth = "option2"` with `sshkeys` for production instances — password auth (`option1`) is less secure.
