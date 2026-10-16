---
page_title: "Utho: utho_autoscaling"
subcategory: "Auto Scaling"
description: |-
  Create and manage Utho Auto Scaling groups.
---

# utho_autoscaling

Creates and manages an Auto Scaling group on Utho. Automatically scales cloud instances up or down based on CPU/RAM metrics or schedules. Instances are deployed from a snapshot or a stack image.

~> **Note:** Auto Scaling groups take 5-10 minutes to provision instances. Terraform polls for up to 6 minutes and shows a graceful warning if provisioning is not complete. Run `terraform apply` again once provisioning completes.

## Example Usage

### Basic ASG with snapshot + CPU policy

```hcl
resource "utho_autoscaling" "web" {
  name              = "web-asg"
  dcslug            = "inmumbaizone2"
  minsize           = "1"
  maxsize           = "5"
  desiredsize       = "2"
  planid            = "10404"
  planname          = "basic"
  snapshotid        = "203184"
  image_name        = "my-app-snapshot"
  os_disk_size      = 80
  public_ip_enabled = 1
  cpumodel          = "intel"

  policies = [
    {
      name     = "scale-out"
      type     = "cpu"
      compare  = "above"
      value    = "80"
      adjust   = 1
      period   = "5m"
      cooldown = "300"
    },
    {
      name     = "scale-in"
      type     = "cpu"
      compare  = "below"
      value    = "20"
      adjust   = -1
      period   = "5m"
      cooldown = "300"
    }
  ]
}
```

### ASG with stack image

```hcl
resource "utho_autoscaling" "web" {
  name              = "web-asg"
  dcslug            = "inmumbaizone2"
  minsize           = "1"
  maxsize           = "5"
  desiredsize       = "1"
  planid            = "10404"
  planname          = "basic"
  stack             = "6669436"
  stackid           = "6669436"
  stackimage        = "ubuntu-22.04-x86_64"
  image_name        = "ubuntu-22.04-x86_64"
  os_disk_size      = 80
  public_ip_enabled = 1
  cpumodel          = "intel"

  policies = [
    {
      name     = "scale-out-cpu"
      type     = "cpu"
      compare  = "above"
      value    = "80"
      adjust   = 1
      period   = "5m"
      cooldown = "300"
    }
  ]
}
```

### ASG with Load Balancer + VPC + security group

```hcl
resource "utho_autoscaling" "web" {
  name              = "web-asg"
  dcslug            = "inmumbaizone2"
  minsize           = "2"
  maxsize           = "10"
  desiredsize       = "2"
  planid            = "10404"
  planname          = "basic"
  snapshotid        = var.snapshot_id
  image_name        = var.snapshot_name
  os_disk_size      = 80
  public_ip_enabled = 1
  cpumodel          = "intel"
  vpc               = utho_subnet.public.id
  load_balancers    = utho_loadbalancer.main.id
  security_groups   = utho_firewall.web.id

  policies = [
    {
      name     = "scale-out-cpu"
      type     = "cpu"
      compare  = "above"
      value    = "75"
      adjust   = 2
      period   = "5m"
      cooldown = "300"
    },
    {
      name     = "scale-in-cpu"
      type     = "cpu"
      compare  = "below"
      value    = "25"
      adjust   = -1
      period   = "10m"
      cooldown = "600"
    }
  ]

  schedules = [
    {
      name                = "peak-hours"
      desiredsize         = "5"
      timezone            = "Asia/Kolkata"
      recurrence          = "Every day  09:00"
      recurrence_duration = "Every day"
      recurrence_week     = ""
      selected_time       = "09:00"
      selected_date       = "2026-10-01"
      start_date          = "2026-10-01T09:00:00.000+05:30"
    }
  ]
}
```

### Import existing ASG

```bash
terraform import utho_autoscaling.web <asg-id>
```

## Argument Reference

| Argument           | Type   | Required | Description |
|--------------------|--------|----------|-------------|
| `name`             | String | Yes      | Unique ASG name. Changing forces new resource. |
| `dcslug`           | String | Yes      | Data center slug. Changing forces new resource. |
| `minsize`          | String | Yes      | Minimum instances. Can be updated in place. |
| `maxsize`          | String | Yes      | Maximum instances. Can be updated in place. |
| `desiredsize`      | String | Yes      | Desired instances. Can be updated in place. |
| `planid`           | String | Yes      | Plan ID for the instance size (e.g. `10404`). Use the `utho_cloud_plans` data source to find available plan IDs for your `dcslug`. Changing forces new resource. |
| `planname`         | String | Yes      | Plan name label corresponding to `planid` (e.g. `basic`). This is a display label sent alongside `planid` — the API does not cross-validate them. Use the `name` field from `utho_cloud_plans` data source. Changing forces new resource. |
| `os_disk_size`     | Number | Yes      | Disk size in GB. Changing forces new resource. |
| `public_ip_enabled`| Number | Yes      | Enable public IP: `1` or `0`. |
| `policies`         | List   | Yes      | Scaling policies. At least 1 required. |
|        | String | No       | Snapshot ID. **Required** when using snapshot-based deployment. Mutually exclusive with stack deployment. When set, also set  to the snapshot name. |
|             | String | No       | Stack slug. **Required** when using stack-based deployment. Mutually exclusive with snapshot deployment. |
|           | String | No       | Stack ID. **Required** when using stack-based deployment. Set alongside , , and . |
|        | String | No       | Stack image slug. **Required** when using stack-based deployment. Also set  to the same value. |
|        | String | No       | Boot image name. **Required** for both snapshot and stack deployment. For snapshot: set to snapshot name. For stack: set to stack image slug (same as ). |
| `vpc`              | String | No       | VPC subnet ID. |
| `load_balancers`   | String | No       | Load balancer ID. Can be updated in place. |
| `security_groups`  | String | No       | Security group ID. Can be updated in place. |
| `target_groups`    | String | No       | Target group ID. Can be updated in place. |
| `cpumodel`         | String | No       | CPU model: `amd` or `intel`. |
| `schedules`        | List   | No       | Scaling schedules. 0 or more. |

### policies block

| Argument   | Type   | Required | Description |
|------------|--------|----------|-------------|
| `name`     | String | Yes      | Policy name. |
| `type`     | String | Yes      | Metric: `cpu` or `ram`. |
| `compare`  | String | Yes      | Trigger: `above` or `below`. |
| `value`    | String | Yes      | Threshold percentage (e.g. `80`). |
| `adjust`   | Number | Yes      | Instances to add (positive) or remove (negative). |
| `period`   | String | Yes      | Evaluation window (e.g. `5m`). |
| `cooldown` | String | Yes      | Cooldown in seconds (e.g. `300`). |

### schedules block

| Argument              | Type   | Required | Description |
|-----------------------|--------|----------|-------------|
| `name`                | String | Yes      | Schedule name. |
| `desiredsize`         | String | Yes      | Desired count at schedule time. |
| `timezone`            | String | Yes      | Timezone (e.g. `Asia/Kolkata`). |
| `recurrence`          | String | Yes      | Recurrence expression. |
| `recurrence_duration` | String | Yes      | Duration (e.g. `Every day`). |
| `selected_time`       | String | Yes      | Time (HH:MM). |
| `selected_date`       | String | Yes      | Date (YYYY-MM-DD). |
| `start_date`          | String | Yes      | Start datetime ISO 8601. |

## Attribute Reference

| Attribute    | Type   | Description |
|--------------|--------|-------------|
| `id`         | String | Auto Scaling group ID. |
| `status`     | String | Current status (`Active`, `Deploying`). |
| `created_at` | String | Creation timestamp. |

## Notes

- **Boot source is required** — either snapshot or stack deployment must be configured:
  - **Snapshot**: set  +  (snapshot name)
  - **Stack**: set  +  +  +  (stack image slug)
  - These two are mutually exclusive — do not mix snapshot and stack fields.
- `load_balancers`, `security_groups`, `target_groups` can be attached/detached after creation without replacing the ASG.
- `minsize`, `maxsize`, `desiredsize` can be updated in place.
- Use `lifecycle { ignore_changes = [...] }` for fields that differ between create and read API responses.