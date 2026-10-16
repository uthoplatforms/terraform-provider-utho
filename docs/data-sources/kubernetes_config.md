---
page_title: "Utho: utho_kubernetes_config"
subcategory: "Kubernetes"
description: |-
  Fetch the kubeconfig file for a Utho Kubernetes cluster.
---

# utho_kubernetes_config

Fetches the kubeconfig YAML for a Utho Kubernetes cluster. The kubeconfig contains all the credentials and connection details needed for `kubectl`, Helm, and the Kubernetes Terraform provider to authenticate with your cluster.

~> **Security:** The kubeconfig contains cluster credentials equivalent to admin access. Always mark outputs as `sensitive = true` and never commit kubeconfig files to version control.

## Example Usage

### Save kubeconfig to a file

```hcl
data "utho_kubernetes_config" "main" {
  cluster_id = utho_kubernetes.main.id
}

resource "local_file" "kubeconfig" {
  content         = data.utho_kubernetes_config.main.raw_config
  filename        = "${path.module}/kubeconfig.yaml"
  file_permission = "0600"   # owner read/write only
}

output "kubectl_command" {
  value = "export KUBECONFIG=${abspath("${path.module}/kubeconfig.yaml")} && kubectl get nodes"
}
```

### Use with the Kubernetes Terraform provider

Deploy Kubernetes resources directly from Terraform by combining the Utho provider with the Kubernetes provider.

```hcl
data "utho_kubernetes_config" "main" {
  cluster_id = utho_kubernetes.main.id
}

# Decode the kubeconfig to extract connection details
locals {
  kubeconfig = yamldecode(data.utho_kubernetes_config.main.raw_config)
}

provider "kubernetes" {
  host                   = local.kubeconfig.clusters[0].cluster.server
  cluster_ca_certificate = base64decode(local.kubeconfig.clusters[0].cluster["certificate-authority-data"])
  client_certificate     = base64decode(local.kubeconfig.users[0].user["client-certificate-data"])
  client_key             = base64decode(local.kubeconfig.users[0].user["client-key-data"])
}

provider "helm" {
  kubernetes {
    host                   = local.kubeconfig.clusters[0].cluster.server
    cluster_ca_certificate = base64decode(local.kubeconfig.clusters[0].cluster["certificate-authority-data"])
    client_certificate     = base64decode(local.kubeconfig.users[0].user["client-certificate-data"])
    client_key             = base64decode(local.kubeconfig.users[0].user["client-key-data"])
  }
}
```

### Output for CI/CD use

```hcl
data "utho_kubernetes_config" "main" {
  cluster_id = utho_kubernetes.main.id
}

output "kubeconfig" {
  value     = data.utho_kubernetes_config.main.raw_config
  sensitive = true
}
```

Then in CI/CD:
```bash
terraform output -raw kubeconfig > kubeconfig.yaml
export KUBECONFIG=kubeconfig.yaml
kubectl get nodes
```

## Argument Reference

| Argument     | Type   | Required | Description |
|--------------|--------|----------|-------------|
| `cluster_id` | String | Yes      | Kubernetes cluster ID from `utho_kubernetes.main.id`. |

## Attribute Reference

| Attribute    | Type   | Description |
|--------------|--------|-------------|
| `id`         | String | Same as `cluster_id`. |
| `raw_config` | String | Complete kubeconfig YAML. **Sensitive.** Contains cluster CA certificate, client certificate, and client private key. |
