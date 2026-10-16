---
page_title: "Utho: utho_object_storage"
subcategory: "Object Storage"
description: |-
  Create and manage Utho Object Storage buckets (S3-compatible).
---

# utho_object_storage

Creates and manages a Utho Object Storage bucket.

~> **Note:** Bucket names must be **lowercase only** and may contain letters, numbers, and hyphens. Uppercase letters and underscores are rejected by the API. Example: `my-app-uploads` ✅ `My_App_Uploads` ❌ Buckets are S3-compatible — use any AWS S3 SDK, CLI, or library to read and write objects. Terraform manages the bucket; your application manages the objects inside it.

~> **Note:** Object Storage is currently available only in `innoida` (Delhi/Noida) data center.

Common uses: application file uploads, static assets, database backups, Terraform remote state, logs.

## Example Usage

### Private bucket for application uploads

```hcl
resource "utho_object_storage" "uploads" {
  name         = "myapp-uploads"
  dcslug       = "innoida"
  access       = "private"
  billingcycle = "monthly"
}

output "s3_endpoint"   { value = "https://innoida.utho.io" }
output "s3_bucket"     { value = utho_object_storage.uploads.name }
output "s3_access_key" {
  value     = utho_object_storage.uploads.access_key
  sensitive = true
}
output "s3_secret_key" {
  value     = utho_object_storage.uploads.secret_key
  sensitive = true
}
```

### Public static assets bucket

```hcl
resource "utho_object_storage" "assets" {
  name         = "myapp-assets"
  dcslug       = "innoida"
  access       = "public"
  billingcycle = "monthly"
}
```

### Upload-only bucket

Allow users to upload files without being able to read others' files.

```hcl
resource "utho_object_storage" "submissions" {
  name         = "myapp-submissions"
  dcslug       = "innoida"
  access       = "upload"
  billingcycle = "monthly"
}
```

### Versioned bucket for backup storage

```hcl
resource "utho_object_storage" "backups" {
  name            = "myapp-backups"
  dcslug          = "innoida"
  access          = "private"
  billingcycle    = "monthly"
  version_enabled = true
}
```

### Update access policy in place

Change access policy without recreating the bucket.

```hcl
resource "utho_object_storage" "assets" {
  name         = "myapp-assets"
  dcslug       = "innoida"
  access       = "public"   # change to "private" and apply — updates in place
  billingcycle = "monthly"
}
```

### Terraform remote state backend

```hcl
resource "utho_object_storage" "tfstate" {
  name         = "mycompany-terraform-state"
  dcslug       = "innoida"
  access       = "private"
  billingcycle = "monthly"
}

# Add this to backend.tf and run terraform init
# terraform {
#   backend "s3" {
#     bucket                      = "mycompany-terraform-state"
#     key                         = "prod/terraform.tfstate"
#     endpoint                    = "https://innoida.utho.io"
#     access_key                  = "<access_key>"
#     secret_key                  = "<secret_key>"
#     region                      = "us-east-1"
#     skip_credentials_validation = true
#     skip_metadata_api_check     = true
#     force_path_style            = true
#   }
# }
```

## Argument Reference

| Argument          | Type   | Required | Description |
|-------------------|--------|----------|-------------|
| `name`            | String | Yes      | Bucket name. Must be globally unique, **lowercase only**, and may contain letters, numbers, and hyphens. Uppercase letters and underscores are not allowed. Changing this forces a new resource. |
| `dcslug`          | String | Yes      | Data center. Use `innoida` for object storage. Changing this forces a new resource. |
| `billingcycle`    | String | Yes      | Billing cycle: `monthly`. Changing this forces a new resource. |
| `access`          | String | No       | Bucket policy: `private` (default), `public` (public read), `upload` (write-only). Can be updated in place. |
| `version_enabled` | Bool   | No       | Enable object versioning. Protects against overwrites and deletions. Can be updated in place. Default: `false`. |

## Attribute Reference

| Attribute      | Type   | Description |
|----------------|--------|-------------|
| `id`           | String | Bucket name (used as unique identifier). |
| `access_key`   | String | S3 access key for this bucket. **Sensitive.** |
| `secret_key`   | String | S3 secret key for this bucket. **Sensitive.** |
| `status`       | String | Bucket status (`Active`). |
| `plan_gb`      | Number | Storage allocation in GB. |
| `used_gb`      | String | Current storage usage in GB. |
| `created_at`   | String | Creation timestamp. |

## Connecting with AWS SDK

```python
import boto3

s3 = boto3.client(
    "s3",
    endpoint_url="https://innoida.utho.io",
    aws_access_key_id=ACCESS_KEY,
    aws_secret_access_key=SECRET_KEY,
)

s3.upload_file("myfile.txt", "myapp-uploads", "myfile.txt")
```

```bash
# AWS CLI
aws s3 ls s3://myapp-uploads \
  --endpoint-url https://innoida.utho.io \
  --aws-access-key-id $ACCESS_KEY \
  --aws-secret-access-key $SECRET_KEY
```

## Related Resources

| Resource | Purpose |
|----------|---------|
| [utho_object_storage_key](object_storage_key) | Create dedicated access keys |
| [utho_object_storage_permission](object_storage_permission) | Grant key permissions on bucket |

## Notes

- S3 endpoint: `https://innoida.utho.io`
- Compatible with AWS SDK v2, boto3, s3cmd, rclone, MinIO client, and any S3-compatible tool
- Terraform manages the bucket — use the S3 SDK to manage objects inside it
- The bucket's own `access_key` and `secret_key` give full access — use `utho_object_storage_key` + `utho_object_storage_permission` for restricted per-service keys

## Billing

- ₹500/month minimum per bucket (100 GB included)
- Usage above 100 GB billed per GB