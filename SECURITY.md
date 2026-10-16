# Security Policy

## Supported Versions

| Version | Supported |
|---------|-----------|
| 0.7.x   | ✅ |
| < 0.7   | ❌ |

## Reporting a Vulnerability

Please **do not** report security vulnerabilities through public GitHub issues.

Report them privately through GitHub's [private vulnerability reporting](https://github.com/nitinuthocloud/terraform-provider-utho/security/advisories/new). Include:

* The affected provider version
* A description of the issue and its impact
* Steps to reproduce, with any API keys or passwords removed

We will acknowledge your report, investigate, and keep you informed until a fix is released.

## Handling Credentials

* Never commit Utho API keys. Supply them through `UTHO_API_KEY` or a `sensitive` Terraform variable.
* Terraform state can contain secrets in plain text. Store it in a backend with encryption and access control.
* If a key is exposed, revoke it immediately in the [Utho Console](https://console.utho.com) and create a new one.
