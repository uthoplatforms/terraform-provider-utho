---
page_title: "Utho: utho_iam_user"
subcategory: "IAM"
description: |-
  Create and manage IAM sub-users in your Utho account.
---

# utho_iam_user

Creates and manages an IAM sub-user in your Utho account. IAM users can be granted granular permissions to access specific Utho services.

## Example Usage

### Basic IAM user

```hcl
resource "utho_iam_user" "dev" {
  fullname    = "John Doe"
  email       = "john.doe@mycompany.com"
  mobilecc    = "91"
  mobile      = "9999999999"
  permissions = "compute_read,compute_write,database_read,database_write"
}
```

### IAM user with full permissions

```hcl
resource "utho_iam_user" "admin" {
  fullname    = "Admin User"
  email       = "admin@mycompany.com"
  mobilecc    = "91"
  mobile      = "8888888888"
  permissions = "compute_read,compute_write,compute_delete,database_read,database_write,database_delete,objectstorage_read,objectstorage_write,objectstorage_delete,monitoring_read,monitoring_write,monitoring_delete,loadbalancer_read,loadbalancer_write,loadbalancer_delete,dns_read,dns_write,dns_delete,vpc_read,vpc_write,vpc_delete,firewall_read,firewall_write,firewall_delete"
}
```

## Argument Reference

| Argument      | Type   | Required | Description |
|---------------|--------|----------|-------------|
| `fullname`    | String | Yes      | Full name of the IAM user. Changing forces new resource. |
| `email`       | String | Yes      | Email address. Must be unique. Changing forces new resource. |
| `mobilecc`    | String | Yes      | Mobile country code without `+` (e.g. `91`). Changing forces new resource. |
| `mobile`      | String | Yes      | Mobile number. Changing forces new resource. |
| `permissions` | String | Yes      | Comma-separated permissions string. Can be updated in place. |

## Attribute Reference

| Attribute    | Type   | Description |
|--------------|--------|-------------|
| `id`         | String | IAM user account access ID. |
| `status`     | String | User status (`Pending`, `Active`). |
| `date_added` | String | Date the user was added. |

## Available Permissions

Permissions follow the format `{service}_{action}` where action is `read`, `write`, or `delete`.

| Service | Permissions |
|---------|-------------|
| Compute | `compute_read`, `compute_write`, `compute_delete` |
| Database | `database_read`, `database_write`, `database_delete` |
| Object Storage | `objectstorage_read`, `objectstorage_write`, `objectstorage_delete` |
| Kubernetes | `kubernetes_read`, `kubernetes_write`, `kubernetes_delete` |
| Load Balancer | `loadbalancer_read`, `loadbalancer_write`, `loadbalancer_delete` |
| DNS | `dns_read`, `dns_write`, `dns_delete` |
| VPC | `vpc_read`, `vpc_write`, `vpc_delete` |
| Firewall | `firewall_read`, `firewall_write`, `firewall_delete` |
| Monitoring | `monitoring_read`, `monitoring_write`, `monitoring_delete` |
| Auto Scaling | `autoscaling_read`, `autoscaling_write`, `autoscaling_delete` |
| VPN | `vpn_read`, `vpn_write`, `vpn_delete` |
| Snapshots | `snapshot_read`, `snapshot_write`, `snapshot_delete` |
| SSH Keys | `sshkey_read`, `sshkey_write`, `sshkey_delete` |
| Billing | `billing_read`, `billing_write`, `billing_delete` |
| API | `api_read`, `api_write`, `api_delete` |
| Container Registry | `container_registry_read`, `container_registry_write`, `container_registry_delete` |

## Notes

- The invited user receives an email to set their password and activate their account.
- `fullname`, `email`, `mobilecc`, and `mobile` cannot be changed after creation — changing any of these forces a new resource.
- `permissions` can be updated in place at any time.
