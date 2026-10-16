---
page_title: "Utho: utho_alert_contact"
subcategory: "Monitoring"
description: |-
  Create and manage alert contacts for Utho monitoring.
---

# utho_alert_contact

Creates and manages an alert contact. Alert contacts receive notifications when an alert is triggered via email or SMS.

## Example Usage

```hcl
resource "utho_alert_contact" "devops" {
  name         = "devops-team"
  email        = "devops@mycompany.com"
  mobilenumber = "9999999999"
  status       = "1"
}
```

## Argument Reference

| Argument       | Type   | Required | Description |
|----------------|--------|----------|-------------|
| `name`         | String | Yes      | Contact name. |
| `email`        | String | Yes      | Email address for alert notifications. |
| `mobilenumber` | String | Yes      | Mobile number for SMS notifications. |
| `status`       | String | Yes      | Enable contact: `1` (active) or `0` (inactive). |

## Attribute Reference

| Attribute | Type   | Description |
|-----------|--------|-------------|
| `id`      | String | Alert contact ID. |

## Related Resources

| Resource | Purpose |
|----------|---------|
| [utho_alert](alert) | Create alerts that use this contact |
