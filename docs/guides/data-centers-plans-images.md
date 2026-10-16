---
page_title: "Data Centers, Plans, and Images"
subcategory: ""
description: |-
  Look up valid data center slugs, plan IDs, and OS images with data sources instead of hard-coding them.
---

# Data Centers, Plans, and Images

Most Utho resources need three values that depend on your account and location:

| Argument | Meaning | Data source |
|----------|---------|-------------|
| `dcslug` | Data center the resource is created in | [`utho_cloud_dczones`](../data-sources/cloud_dczones) |
| `planid` | Instance size (vCPU, RAM, disk) | [`utho_cloud_plans`](../data-sources/cloud_plans) |
| `image`  | Operating system image | [`utho_cloud_images`](../data-sources/cloud_images) |

You can hard-code these, as most examples in this documentation do for readability. In shared modules and long-lived configurations it is better to look them up, so your code keeps working when plans or images change.

## Data centers

| Slug | Location |
|------|----------|
| `innoida` | Delhi (Noida), India |
| `inmumbaizone2` | Mumbai, India |
| `inbangalore` | Bangalore, India |
| `defra1` | Frankfurt, Germany |
| `uslosangeles` | Los Angeles, United States |

Feature availability differs between data centers. For example, EBS block volumes are not available everywhere. Check before you deploy:

```hcl
data "utho_cloud_dczones" "all" {}

output "ebs_capable_zones" {
  value = [
    for z in data.utho_cloud_dczones.all.zones :
    z.slug if z.ebs_available && z.status == "active"
  ]
}
```

Each zone also reports a `default_cpu` (`amd` or `intel`), which is the value to use for `cpumodel` on resources that accept it.

## Plans

Plans are specific to a data center, so always filter by `dcslug`:

```hcl
data "utho_cloud_plans" "mumbai" {
  dcslug = "inmumbaizone2"
}

locals {
  # Smallest plan with at least 2 vCPU and 4 GB RAM that includes a disk
  web_plan = [
    for p in data.utho_cloud_plans.mumbai.plans :
    p if tonumber(p.cpu) >= 2 && tonumber(p.ram) >= 4096 && p.disk != "0"
  ][0]
}

output "web_plan" {
  value = "${local.web_plan.id}: ${local.web_plan.name}"
}
```

Plans with `disk = "0"` do not include an OS disk. Use the `ebs` block of [`utho_cloud`](../resources/cloud) to attach one.

Managed databases use their own plans. See [`utho_database_plans`](../data-sources/database_plans).

## Images

```hcl
data "utho_cloud_images" "ubuntu" {
  distro = "ubuntu"
}

output "ubuntu_images" {
  value = data.utho_cloud_images.ubuntu.images[*].image
}
```

Use the `image` attribute, such as `ubuntu-22.04-x86_64`, as the `image` argument of `utho_cloud`. To boot from your own snapshot or ISO instead, see [`utho_cloud_snapshots`](../data-sources/cloud_snapshots) and [`utho_cloud_isos`](../data-sources/cloud_isos).

## Putting it together

```hcl
variable "dcslug" {
  type    = string
  default = "inmumbaizone2"
}

data "utho_cloud_plans" "this" {
  dcslug = var.dcslug
}

data "utho_cloud_images" "ubuntu" {
  distro = "ubuntu"
}

locals {
  plan = [
    for p in data.utho_cloud_plans.this.plans :
    p if tonumber(p.cpu) >= 2 && p.disk != "0"
  ][0]

  image = one([
    for i in data.utho_cloud_images.ubuntu.images :
    i if i.image == "ubuntu-22.04-x86_64"
  ])
}

resource "utho_cloud" "web" {
  hostname        = "web-01.mhc"
  dcslug          = var.dcslug
  planid          = local.plan.id
  billingcycle    = "hourly"
  auth            = "option2"
  sshkeys         = utho_ssh_key.deploy.id
  image           = local.image.image
  enable_publicip = "true"
}
```

~> **Note:** `utho_cloud` instances cannot be updated in place. If a lookup starts returning a different plan or image, the next `terraform apply` fails with an "Update not supported" error. For production instances, pin the selected values with variables once chosen. To change an instance's size, use [`utho_cloud_resize`](../resources/cloud_resize); to rebuild it deliberately, run `terraform apply -replace="utho_cloud.web"`.
