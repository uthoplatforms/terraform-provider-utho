---
page_title: "Utho: utho_autoscaling_schedule"
subcategory: "Auto Scaling"
description: |-
  Manage scheduled scaling policies for a Utho Auto Scaling group.
---

# utho_autoscaling_schedule

Manages a scheduled scaling policy for a Utho Auto Scaling group. Scheduled scaling adjusts the desired instance count at a specific time, independent of metric-based policies.

Use this resource to update or delete schedules created inside `utho_autoscaling`. To create schedules at launch time use the `schedules` block inside `utho_autoscaling`.

## Example Usage

### Update a schedule

```hcl
resource "utho_autoscaling_schedule" "peak" {
  asg_id      = utho_autoscaling.app.id
  id          = "34235000"
  name        = "peak-hours"
  desiredsize = "5"
  timezone    = "Asia/Kolkata"
  recurrence  = "Every day 09:00"
  start_date  = "2026-09-15 09:00:00"
  status      = 1
}
```

### Disable a schedule

```hcl
resource "utho_autoscaling_schedule" "peak" {
  asg_id      = utho_autoscaling.app.id
  id          = "34235000"
  name        = "peak-hours"
  desiredsize = "5"
  timezone    = "Asia/Kolkata"
  recurrence  = "Every day 09:00"
  start_date  = "2026-09-15 09:00:00"
  status      = 0
}
```

## Argument Reference

| Argument      | Type   | Required | Description |
|---------------|--------|----------|-------------|
| `asg_id`      | String | Yes      | Auto scaling group ID. Changing this forces a new resource. |
| `name`        | String | Yes      | Schedule name. Updatable. |
| `desiredsize` | String | Yes      | Desired instance count at scheduled time. Updatable. |
| `timezone`    | String | Yes      | Timezone (e.g. `Asia/Kolkata`). Updatable. |
| `recurrence`  | String | Yes      | Recurrence expression (e.g. `Every day 09:00`). Updatable. |
| `start_date`  | String | Yes      | Start datetime (e.g. `2026-09-15 09:00:00`). Updatable. |
| `status`      | Number | No       | Schedule status: `1` (active) or `0` (inactive). Default: `1`. Updatable. |

## Attribute Reference

| Attribute | Type   | Description |
|-----------|--------|-------------|
| `id`      | String | Unique schedule ID. |

## Notes

- All fields except `asg_id` are updatable in place.
- Destroying this resource deletes the schedule from the auto scaling group.
- Schedules are also removed when the auto scaling group is deleted.