---
page_title: "Utho: utho_autoscaling_policy"
subcategory: "Auto Scaling"
description: |-
  Create and manage scaling policies for a Utho Auto Scaling group.
---

# utho_autoscaling_policy

Creates and manages a scaling policy for a Utho Auto Scaling group. Scaling policies define when and how the group scales based on CPU or RAM metrics.

Use this resource to add policies to an existing auto scaling group. To define policies at creation time, use the `policies` block inside `utho_autoscaling`.

## Example Usage

### Scale up on high CPU

```hcl
resource "utho_autoscaling_policy" "scale_up" {
  asg_id       = utho_autoscaling.app.id
  name         = "scale-up-cpu"
  type         = "cpu"
  compare      = "above"
  value        = "80"
  adjust       = 2
  period       = "5m"
  cooldown     = "300"
  scaling_type = "horizontal"
}
```

### Scale down on low CPU

```hcl
resource "utho_autoscaling_policy" "scale_down" {
  asg_id       = utho_autoscaling.app.id
  name         = "scale-down-cpu"
  type         = "cpu"
  compare      = "below"
  value        = "20"
  adjust       = -1
  period       = "5m"
  cooldown     = "300"
  scaling_type = "horizontal"
}
```

### Scale on RAM usage

```hcl
resource "utho_autoscaling_policy" "scale_ram" {
  asg_id       = utho_autoscaling.app.id
  name         = "scale-up-ram"
  type         = "ram"
  compare      = "above"
  value        = "75"
  adjust       = 1
  period       = "1m"
  cooldown     = "180"
  scaling_type = "horizontal"
}
```

## Argument Reference

| Argument       | Type   | Required | Description |
|----------------|--------|----------|-------------|
| `asg_id`       | String | Yes      | Auto scaling group ID. Changing this forces a new resource. |
| `name`         | String | Yes      | Policy name. Updatable. |
| `type`         | String | Yes      | Metric type: `cpu` or `ram`. Changing this forces a new resource. |
| `compare`      | String | Yes      | Trigger direction: `above` or `below`. Updatable. |
| `value`        | String | Yes      | Threshold percentage (e.g. `80` for 80%). Updatable. |
| `adjust`       | Number | Yes      | Number of instances to add (positive) or remove (negative). Updatable. |
| `period`       | String | Yes      | Evaluation period (e.g. `1m`, `5m`). Updatable. |
| `cooldown`     | String | Yes      | Cooldown in seconds after a scaling action. Updatable. |
| `scaling_type` | String | Yes      | Scaling type: `horizontal`. Changing this forces a new resource. |

## Attribute Reference

| Attribute | Type   | Description |
|-----------|--------|-------------|
| `id`      | String | Unique scaling policy ID. |

## Notes

- `name`, `compare`, `value`, `adjust`, `period`, and `cooldown` are all updatable in place.
- Changing `type` or `scaling_type` destroys and recreates the policy.
- Policies are automatically removed when the auto scaling group is deleted.