---
page_title: "Utho: utho_cloud_resize"
subcategory: "Cloud Instances"
description: |-
  Resize a Utho Cloud instance to a different plan.
---

# utho_cloud_resize

Resizes an existing Utho Cloud instance to a different plan (more or less CPU, RAM, and disk). The instance must be powered off before resizing.

## Example Usage

### Resize to a larger plan

```hcl
# Step 1 — power off first
resource "utho_cloud_power" "web" {
  cloud_id = utho_cloud.web.id
  action   = "poweroff"
}

# Step 2 — resize
resource "utho_cloud_resize" "web" {
  cloud_id    = utho_cloud.web.id
  plan_id     = "10313"   # 4 vCPU / 8 GB / 160 GB
  resize_type = "full"    # full = CPU/RAM + disk, ramcpu = CPU/RAM only

  depends_on = [utho_cloud_power.web]
}

# Step 3 — power back on
# Change utho_cloud_power.web action to "poweron" and apply
```

### Find a bigger plan before resizing

```hcl
data "utho_cloud_plans" "mumbai" {
  dcslug = "inmumbaizone2"
}

locals {
  bigger_plan = one([
    for p in data.utho_cloud_plans.mumbai.plans :
    p if p.cpu == "4" && p.disk != "0" && p.slug == "basic"
  ])
}

resource "utho_cloud_resize" "web" {
  cloud_id    = utho_cloud.web.id
  plan_id     = local.bigger_plan.id
  resize_type = "full"
}
```

## Argument Reference

| Argument      | Type   | Required | Description |
|---------------|--------|----------|-------------|
| `cloud_id`    | String | Yes      | Cloud instance ID. Changing this forces a new resource. |
| `plan_id`     | String | Yes      | New plan ID. Use [utho_cloud_plans](../data-sources/cloud_plans) to find valid plan IDs. |
| `resize_type` | String | Yes      | `full` — resize CPU, RAM and disk. `ramcpu` — resize CPU and RAM only (no disk change). |

## Attribute Reference

| Attribute | Type   | Description |
|-----------|--------|-------------|
| `id`      | String | Same as `cloud_id`. |

## Resize Types

| Value    | What Changes |
|----------|--------------|
| `full`   | CPU, RAM, and disk — use when upgrading to a plan with larger disk |
| `ramcpu` | CPU and RAM only — use when you want to keep existing disk size |

## Notes

- **Power off the instance before resizing** using `utho_cloud_power`. Resizing a running instance may cause data corruption.
- After resize, power the instance back on with `utho_cloud_power`.
- You can only resize to plans available in the same data center.
- Use `data.utho_cloud_plans` to discover available plan IDs and their specs.
- Disk resizes are permanent — you cannot shrink disk after expanding.
- Destroying this resource does NOT destroy the cloud instance.
