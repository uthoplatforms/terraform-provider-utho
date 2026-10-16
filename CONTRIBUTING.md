# Contributing

Thanks for helping improve the Utho Terraform provider.

## Reporting Issues

* **Bugs:** open a [bug report](https://github.com/nitinuthocloud/terraform-provider-utho/issues/new?template=bug_report.yml) with your Terraform and provider versions, a minimal configuration, and the error output. Remove API keys and passwords first.
* **Feature requests:** open a [feature request](https://github.com/nitinuthocloud/terraform-provider-utho/issues/new?template=feature_request.yml) describing the Utho service or argument you need.
* **Security issues:** do not open a public issue. See [SECURITY.md](SECURITY.md).

## Development Setup

Requirements: [Go](https://go.dev/dl/) (version in `go.mod`) and [Terraform](https://developer.hashicorp.com/terraform/install) 1.0+.

```bash
git clone https://github.com/nitinuthocloud/terraform-provider-utho.git
cd terraform-provider-utho
go build ./...
go vet ./...
```

To try your build against real infrastructure, use a [development override](README.md#developing-the-provider).

## Making Changes

1. Create a branch from the default branch.
2. Keep each pull request focused on one resource or fix.
3. Run `gofmt -w .`, `go vet ./...`, and `go build ./...` before pushing.
4. Update documentation in `docs/` for any schema change. Every resource needs a page in `docs/resources/` with a `subcategory`, an example, an argument reference, and an attribute reference.
5. Add or update an example in `examples/` when you add a resource.
6. Add an entry to the `Unreleased` section of [CHANGELOG.md](CHANGELOG.md).

## Adding a Resource

1. Add the API calls to `internal/client/`.
2. Add the resource in `internal/resources/`, implementing Create, Read, Update, Delete and, where the API allows it, `ImportState`.
3. Register it in `internal/provider/provider.go`.
4. Document it in `docs/resources/<name>.md` and add an example under `examples/`.

## Commit Messages

Use a short imperative subject, for example `Add utho_kubernetes_node_pool autoscaling` or `Fix utho_cloud Read when instance is deleted`.
