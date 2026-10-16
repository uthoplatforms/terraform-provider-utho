---
page_title: "Managed Kubernetes with the Kubernetes and Helm Providers"
subcategory: "Tutorials"
description: |-
  Create a Utho managed Kubernetes cluster and deploy workloads to it with the Kubernetes and Helm providers.
---

# Managed Kubernetes with the Kubernetes and Helm Providers

This tutorial creates a Utho managed Kubernetes cluster, adds a dedicated node pool, and then deploys workloads onto it using the official [Kubernetes](https://registry.terraform.io/providers/hashicorp/kubernetes/latest) and [Helm](https://registry.terraform.io/providers/hashicorp/helm/latest) providers.

## Prerequisites

* Complete [Getting Started](getting-started).
* [`kubectl`](https://kubernetes.io/docs/tasks/tools/) installed locally, if you want to inspect the cluster by hand.

## Recommended layout: two configurations

Terraform configures providers before it creates resources. If the Kubernetes provider's credentials come from a cluster created in the same configuration, the first `terraform plan` cannot connect to a cluster that does not exist yet, and later changes that replace the cluster can fail in confusing ways.

The reliable pattern, recommended by HashiCorp for every managed Kubernetes service, is to split the work:

1. **Cluster configuration**: creates the cluster and node pools with the Utho provider.
2. **Workloads configuration**: reads the cluster's kubeconfig with the [`utho_kubernetes_config`](../data-sources/kubernetes_config) data source and deploys to it.

## Part 1: Create the cluster

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

resource "utho_kubernetes" "main" {
  dcslug          = "inmumbaizone2"
  cluster_label   = "production"
  cluster_version = "1.32.10-utho"
  network_type    = "public"

  nodepools = [
    {
      label     = "general"
      size      = "10355"
      count     = "3"
      min_nodes = "2"
      max_nodes = "10"
      disk_size = "30"
      disk_type = "nvme"
    }
  ]
}

resource "utho_kubernetes_node_pool" "batch" {
  cluster_id = utho_kubernetes.main.id
  label      = "batch"
  size       = "10355"
  count      = "2"
  min_nodes  = "1"
  max_nodes  = "6"
}

output "cluster_id" {
  value = utho_kubernetes.main.id
}
```

Run `terraform apply` and note the `cluster_id` output.

For a cluster whose nodes live in a private network, set `network_type = "publicprivate"` and `vpc` to a private subnet ID. See [Private Networking with VPC and NAT](private-networking).

### Saving a kubeconfig for kubectl

To use `kubectl` against the new cluster, add:

```hcl
data "utho_kubernetes_config" "main" {
  cluster_id = utho_kubernetes.main.id
}

resource "local_file" "kubeconfig" {
  content         = data.utho_kubernetes_config.main.raw_config
  filename        = "${path.module}/kubeconfig.yaml"
  file_permission = "0600"
}
```

```bash
export KUBECONFIG=$PWD/kubeconfig.yaml
kubectl get nodes
```

!> **Warning:** The kubeconfig contains client credentials with full administrative access to the cluster. Add `kubeconfig.yaml` to `.gitignore`.

## Part 2: Deploy workloads

In a separate directory, read the cluster's kubeconfig and use it to configure the Kubernetes and Helm providers:

```hcl
terraform {
  required_providers {
    utho = {
      source  = "nitinuthocloud/utho"
      version = "~> 0.7"
    }
    kubernetes = {
      source  = "hashicorp/kubernetes"
      version = "~> 2.0"
    }
    helm = {
      source  = "hashicorp/helm"
      version = "~> 2.0"
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

variable "cluster_id" {
  type = string
}

data "utho_kubernetes_config" "main" {
  cluster_id = var.cluster_id
}

locals {
  kubeconfig = yamldecode(data.utho_kubernetes_config.main.raw_config)
  cluster    = local.kubeconfig.clusters[0].cluster
  user       = local.kubeconfig.users[0].user
}

provider "kubernetes" {
  host                   = local.cluster.server
  cluster_ca_certificate = base64decode(local.cluster["certificate-authority-data"])
  client_certificate     = base64decode(local.user["client-certificate-data"])
  client_key             = base64decode(local.user["client-key-data"])
}

provider "helm" {
  kubernetes {
    host                   = local.cluster.server
    cluster_ca_certificate = base64decode(local.cluster["certificate-authority-data"])
    client_certificate     = base64decode(local.user["client-certificate-data"])
    client_key             = base64decode(local.user["client-key-data"])
  }
}
```

-> **Note:** The `helm` block syntax above is for Helm provider 2.x. In Helm provider 3.x, `kubernetes` is an argument rather than a block: `kubernetes = { host = ..., ... }`.

Deploy an application and an ingress controller:

```hcl
resource "kubernetes_namespace" "app" {
  metadata {
    name = "app"
  }
}

resource "kubernetes_deployment" "web" {
  metadata {
    name      = "web"
    namespace = kubernetes_namespace.app.metadata[0].name
  }

  spec {
    replicas = 3

    selector {
      match_labels = { app = "web" }
    }

    template {
      metadata {
        labels = { app = "web" }
      }

      spec {
        container {
          name  = "web"
          image = "nginx:1.27"

          port {
            container_port = 80
          }
        }
      }
    }
  }
}

resource "helm_release" "ingress_nginx" {
  name             = "ingress-nginx"
  repository       = "https://kubernetes.github.io/ingress-nginx"
  chart            = "ingress-nginx"
  namespace        = "ingress-nginx"
  create_namespace = true
}
```

```bash
terraform init
terraform apply -var="cluster_id=<cluster-id-from-part-1>"
```

## Upgrading and scaling

* **Scaling**: change `count`, `min_nodes`, or `max_nodes` on a [`utho_kubernetes_node_pool`](../resources/kubernetes_node_pool) and apply. The pools defined inline in the `nodepools` argument of `utho_kubernetes` cannot be changed after the cluster is created, so manage any pool you expect to resize as a separate `utho_kubernetes_node_pool` resource.
* **Adding capacity for a new workload type**: add another `utho_kubernetes_node_pool`.
* **Kubernetes version upgrades**: `cluster_version` forces a new cluster. Plan upgrades as a blue/green migration: create a new cluster, move workloads with the Part 2 configuration, then destroy the old one.

See [`utho_kubernetes`](../resources/kubernetes) for the full list of available versions and arguments.
