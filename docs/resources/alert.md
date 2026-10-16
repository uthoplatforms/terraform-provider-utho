---
page_title: "Utho: utho_alert"
subcategory: "Monitoring"
description: |-
  Create and manage monitoring alerts for Utho cloud instances.
---

# utho_alert

Creates and manages a monitoring alert for Utho cloud instances. Alerts trigger when CPU or RAM usage crosses a threshold and notify the specified contacts.

## Example Usage

### CPU alert on a cloud instance

```hcl
resource "utho_alert_contact" "devops" {
  name         = "devops-team"
  email        = "devops@mycompany.com"
  mobilenumber = "9999999999"
  status       = "1"
}

resource "utho_alert" "cpu" {
  name     = "high-cpu-alert"
  ref_type = "cloud"
  type     = "cpu"
  compare  = "above"
  value    = "80"
  for      = "5m"
  contacts = utho_alert_contact.devops.id
  status   = "1"
  ref_ids  = utho_cloud.web.id
}
```

### RAM alert

```hcl
resource "utho_alert" "ram" {
  name     = "high-ram-alert"
  ref_type = "cloud"
  type     = "ram"
  compare  = "above"
  value    = "85"
  for      = "5m"
  contacts = utho_alert_contact.devops.id
  status   = "1"
  ref_ids  = utho_cloud.web.id
}
```

### Multiple instances + multiple contacts

```hcl
resource "utho_alert" "cpu" {
  name     = "prod-cpu-alert"
  ref_type = "cloud"
  type     = "cpu"
  compare  = "above"
  value    = "80"
  for      = "5m"
  contacts = "${utho_alert_contact.devops.id},${utho_alert_contact.oncall.id}"
  status   = "1"
  ref_ids  = "${utho_cloud.web1.id},${utho_cloud.web2.id},${utho_cloud.web3.id}"
}
```

## Argument Reference

| Argument   | Type   | Required | Description |
|------------|--------|----------|-------------|
| `name`     | String | Yes      | Alert name. |
| `ref_type` | String | Yes      | Resource type: `cloud`. |
| `type`     | String | Yes      | Metric: `cpu` or `ram`. |
| `compare`  | String | Yes      | Trigger direction: `above` or `below`. |
| `value`    | String | Yes      | Threshold percentage (e.g. `80`). Can be updated in place. |
| `for`      | String | Yes      | Evaluation window (e.g. `5m`, `10m`). |
| `contacts` | String | Yes      | Comma-separated alert contact IDs. |
| `status`   | String | Yes      | Enable alert: `1` or `0`. |
| `ref_ids`  | String | Yes      | Comma-separated cloud instance IDs to monitor. |

## Attribute Reference

| Attribute | Type   | Description |
|-----------|--------|-------------|
| `id`      | String | Alert ID. |

## Related Resources

| Resource | Purpose |
|----------|---------|
| [utho_alert_contact](alert_contact) | Contacts that receive alert notifications |
| [data.utho_clouds](../data-sources/clouds) | List cloud instances to use as ref_ids |

## Notes

- Multiple instances: `ref_ids = "id1,id2,id3"` — one alert watches all.
- Multiple contacts: `contacts = "id1,id2"` — all contacts notified.
- `value` and `contacts` can be updated in place without recreating the alert.
