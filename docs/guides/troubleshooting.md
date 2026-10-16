---
page_title: "Debugging and Troubleshooting"
subcategory: ""
description: |-
  Solutions to common Utho provider errors and how to collect debug logs.
---

# Debugging and Troubleshooting

This guide covers the most common errors when using the Utho provider and how to collect information for a bug report.

## Common errors

### Missing API Key

```text
Error: Missing API Key
Set api_key in provider block or export UTHO_API_KEY=your-key
```

The provider received an empty API key. Check that `api_key` is set in the provider block and that the variable it reads from has a value, for example by exporting `TF_VAR_utho_api_key` in the shell that runs Terraform. In CI/CD, check that the secret is exposed to the job as an environment variable. See [Authenticating with an API Key](authentication-api-key).

### API error (status 401) or (status 403)

The API key was found but rejected. The key may have been revoked, mistyped, or may not have write access for the operation. Create a new key in the Utho Console and try again.

### API error (status 4xx) with a validation message

The API rejected one of the argument values. The most common causes are:

* A `dcslug`, `planid`, or `image` that does not exist, or is not available in the chosen data center. Look up valid values with the data sources in [Data Centers, Plans, and Images](data-centers-plans-images).
* A plan that does not include a disk (`disk = "0"`) used without an `ebs` block on `utho_cloud`.
* A `cpumodel` that does not match the data center. Check `default_cpu` in [`utho_cloud_dczones`](../data-sources/cloud_dczones).
* A subnet range that falls outside its VPC or overlaps an existing subnet.

The response body included in the error message usually names the failing field.

### Update not supported

```text
Error: Update not supported
Utho Cloud instances cannot be updated in place. Destroy and recreate.
```

Many Utho resources cannot be modified after creation, including cloud instances, VPCs, subnets, and Kubernetes clusters. If you changed an argument on one of these resources, you have three options:

* Revert the change.
* Use a dedicated resource for the change, if one exists. For example, use [`utho_cloud_resize`](../resources/cloud_resize) to change an instance's plan, or [`utho_cloud_power`](../resources/cloud_power) to stop or start it.
* Replace the resource deliberately: `terraform apply -replace="utho_cloud.web"`. This destroys and recreates it, so any data on the instance is lost.

Each resource page lists which arguments force a new resource and which can be updated in place.

### Changes keep appearing after import

Arguments that are only used at creation time, such as `auth`, `planid`, `image`, and `root_password` on `utho_cloud`, are not returned by the API. Set them in your configuration to match the real resource. See [Importing Existing Resources](importing-resources).

### Timeouts on long-running operations

Database clusters, Kubernetes clusters, and auto scaling groups can take several minutes to provision. Terraform waits for them to become ready. If an apply is interrupted, run `terraform plan` to see what Terraform recorded before creating anything again, so you do not end up with duplicate resources.

## Enabling debug logs

Set `TF_LOG` to see detailed logs from Terraform and the provider:

```bash
export TF_LOG=DEBUG
export TF_LOG_PATH=./terraform-debug.log
terraform apply
```

To limit logging to providers only, use `TF_LOG_PROVIDER` instead:

```bash
export TF_LOG_PROVIDER=DEBUG
```

Unset these variables when you are done, because debug logs grow quickly.

!> **Warning:** Debug logs can contain API keys, passwords, and other secrets. Review and redact a log file before sharing it.

## Reporting an issue

If you think you have found a bug, open an issue on the provider's [GitHub repository](https://github.com/nitinuthocloud/terraform-provider-utho/issues) and include:

* The output of `terraform version`, including the Utho provider version.
* A minimal configuration that reproduces the problem, with secrets removed.
* The full error message.
* A redacted debug log, if possible.
