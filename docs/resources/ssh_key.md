---
page_title: "Utho: utho_ssh_key"
subcategory: "Cloud Instances"
description: |-
  Import and manage SSH public keys on Utho Cloud.
---

# utho_ssh_key

Imports and manages an SSH public key on Utho Cloud. When you create a cloud instance with `auth = "option2"`, Utho installs this key so you can connect without a password.

SSH key authentication is strongly recommended over password authentication — it's more secure and works seamlessly with automation tools.

## Example Usage

### Import your local SSH key

```hcl
resource "utho_ssh_key" "deploy" {
  name   = "deploy-key"
  sshkey = file("~/.ssh/id_ed25519.pub")
}
```

### Use the key when creating instances

```hcl
resource "utho_ssh_key" "deploy" {
  name   = "deploy-key"
  sshkey = file("~/.ssh/id_ed25519.pub")
}

resource "utho_cloud" "web" {
  hostname        = "web-01.mhc"
  dcslug          = "inmumbaizone2"
  planid          = "10360"
  billingcycle    = "hourly"
  auth            = "option2"
  sshkeys         = utho_ssh_key.deploy.id   # ← reference the key ID
  image           = "ubuntu-22.04-x86_64"
  enable_publicip = "true"
}
```

### Share one key across many servers

```hcl
resource "utho_ssh_key" "deploy" {
  name   = "deploy-key"
  sshkey = file("~/.ssh/id_ed25519.pub")
}

resource "utho_cloud" "worker" {
  count = 10

  hostname        = "worker-${count.index + 1}.mhc"
  dcslug          = "inmumbaizone2"
  planid          = "10360"
  billingcycle    = "hourly"
  auth            = "option2"
  sshkeys         = utho_ssh_key.deploy.id
  image           = "ubuntu-22.04-x86_64"
  enable_publicip = "true"
}
```

### Connect after deployment

```bash
terraform output -raw   # get the IP
ssh ubuntu@<ip>         # or root@<ip> depending on the image
```

## Argument Reference

| Argument | Type   | Required | Description |
|----------|--------|----------|-------------|
| `name`   | String | Yes      | Key name for identification. Changing this forces a new resource. |
| `sshkey` | String | Yes      | The public key content. Accepted formats: `ssh-rsa`, `ssh-ed25519`, `ecdsa-sha2-nistp256`. **Sensitive.** Changing this forces a new resource. |

## Attribute Reference

| Attribute    | Type   | Description |
|--------------|--------|-------------|
| `id`         | String | Unique key ID. Use this as `sshkeys` in `utho_cloud`. |
| `created_at` | String | Import timestamp. |

## Notes

- Only the **public key** is imported — never share or import your private key.
- SSH keys are immutable — changing `name` or `sshkey` destroys and recreates the key.
- Deleting the key resource does not revoke access to existing instances (the key was already installed at creation time).
- To generate an ed25519 key pair: `ssh-keygen -t ed25519 -C "deploy@mycompany.com"`
