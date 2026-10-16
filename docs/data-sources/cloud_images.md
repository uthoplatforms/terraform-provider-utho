---
page_title: "Utho: utho_cloud_images"
subcategory: "Cloud Instances"
description: |-
  List available OS images for Utho Cloud instances.
---

# utho_cloud_images

Fetches available OS images for cloud instances. Images are grouped by distribution. Use the `image` attribute value as the `image` argument in `utho_cloud`.

## Example Usage

### List all available images

```hcl
data "utho_cloud_images" "all" {}

output "all_images" {
  value = data.utho_cloud_images.all.images[*].image
}
```

### List Ubuntu images only

```hcl
data "utho_cloud_images" "ubuntu" {
  distro = "ubuntu"
}

output "ubuntu_images" {
  value = data.utho_cloud_images.ubuntu.images
}
```

### Find latest Ubuntu LTS

```hcl
data "utho_cloud_images" "ubuntu" {
  distro = "ubuntu"
}

locals {
  ubuntu_lts = one([
    for img in data.utho_cloud_images.ubuntu.images :
    img if img.image == "ubuntu-22.04-x86_64"
  ])
}

resource "utho_cloud" "web" {
  hostname        = "web-01.mhc"
  dcslug          = "inmumbaizone2"
  planid          = "10308"
  billingcycle    = "hourly"
  auth            = "option1"
  root_password   = var.root_password
  image           = local.ubuntu_lts.image
  enable_publicip = "true"
}
```

### Use in dynamic deployments

```hcl
data "utho_cloud_images" "all" {}

locals {
  # Build a map of image slug -> image object
  image_map = { for img in data.utho_cloud_images.all.images : img.image => img }
}

# Reference by slug
output "ubuntu_22_id" {
  value = local.image_map["ubuntu-22.04-x86_64"].id
}
```

## Argument Reference

| Argument | Type   | Required | Description |
|----------|--------|----------|-------------|
| `distro` | String | No       | Filter by distro name. See [Available Distros](#available-distros). If omitted, returns all images. |

## Attribute Reference

### Top-level

| Attribute | Type | Description |
|-----------|------|-------------|
| `images`  | List | List of available OS images. |

### images

| Attribute      | Type   | Description |
|----------------|--------|-------------|
| `id`           | String | Image ID. |
| `image`        | String | Image slug — use this as `image` in `utho_cloud`. |
| `version`      | String | OS version (e.g. `22.04`). |
| `distribution` | String | Full distribution name (e.g. `Ubuntu`). |
| `distro`       | String | Distro identifier (e.g. `ubuntu`). |
| `category`     | String | `distro` for OS images, `app` for application stacks. |

## Available Distros

| Distro Value  | Distribution  | Available Versions |
|---------------|---------------|--------------------|
| `ubuntu`      | Ubuntu        | 20.04, 22.04, 24.04, 25.04 |
| `debian`      | Debian        | 10, 11, 12 |
| `centos`      | CentOS        | 7.9, 8 |
| `almalinux`   | Alma Linux    | 8.4, 9.2 |
| `rockylinux`  | Rocky Linux   | 8.7, 8.8, 8.10, 9.1, 9.2 |
| `fedora`      | Fedora        | 32, 33, 34 |
| `windows`     | Windows       | Server 2012R2, 2016, 2019, 2022 |

## Common Image Slugs

| Image Slug | OS |
|------------|----|
| `ubuntu-22.04-x86_64` | Ubuntu 22.04 LTS (recommended) |
| `ubuntu-24.04-x86_64` | Ubuntu 24.04 LTS |
| `debian-12-x86_64` | Debian 12 |
| `centos-7.9-x86_64` | CentOS 7.9 |
| `almalinux-9.2-x86_64` | Alma Linux 9.2 |
| `rocky-9.2-x86_64` | Rocky Linux 9.2 |
