# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

Terraform provider for CMC Cloud (`cmccloudv2_*` resources and data sources), built on the legacy `hashicorp/terraform-plugin-sdk` v1 (`helper/schema`, non-context CRUD functions returning `error`).

## Commands

The `GNUmakefile` targets depend on `sh`, `find` and `xargs`; neither `make` nor `sh` is on PATH on this Windows machine, so call Go directly here. CI (Linux) runs the make targets.

| Task | Command |
| --- | --- |
| Build everything | `go build ./...` |
| Build the provider binary | `go build -o terraform-provider-cmccloudv2.exe` |
| Vet (CI: `make vet`) | `go vet ./...` |
| Lint (CI: golangci-lint v1.61.0) | `golangci-lint run ./...` |
| Unit tests (CI: `make test`) | `go test ./...` |
| Single test | `go test ./cmccloudv2 -run TestName -v` |
| Acceptance tests | `TF_ACC=1 go test ./... -v -timeout 120m` |
| Regenerate `docs/` | `docs-gen.bat` (needs `tfplugindocs` installed) |

- There are currently no `*_test.go` files; build, vet and lint are the only automated checks.
- `gofmt -l` reports every file on a Windows checkout because `core.autocrlf` gives CRLF working copies (the index is LF). This is not a real formatting failure; do not mass-reformat.
- CI uses Go 1.22 (`go.mod`); `.go-version` (1.15) is stale.
- Releases: pushing a `v*` tag runs GoReleaser, which builds and GPG-signs a draft GitHub release.

## Architecture

Two packages, both in this module:

- `gocmcapiv2/` — the HTTP API client. It is vendored in-tree and imported as `github.com/cmc-cloud/terraform-provider-cmccloudv2/gocmcapiv2`; the commented-out external `github.com/cmc-cloud/gocmcapiv2` dependency in `go.mod` is no longer used. API changes are made here directly.
- `cmccloudv2/` — the Terraform provider.

### API client (`gocmcapiv2`)

- `client.go` defines `Client`, which holds one service interface per product (`client.EFS`, `client.Server`, `client.Kubernetesv2`, ...). Each service lives in its own file as: an exported `XxxService` interface, exported JSON structs, and an unexported struct holding `client *Client`. A new service must be added to the `Client` struct **and** instantiated in `NewClient`.
- Services call `client.Get/Post/Put/Patch/Delete` (returning the raw JSON string to unmarshal) or the `PerformCreate/PerformUpdate/PerformAction/PerformDelete/PerformPatch` wrappers that decode `CreateResponse`/`ActionResponse`. Paths are relative to `api_endpoint`.
- Every request sends `api_key` as a query parameter and `Project-Id` / `Region-Id` headers.
- `parseResponse` normalises the different error shapes of the backend APIs (core, devops, security, DNS, CDN). 404s are wrapped with `ErrNotFound`; check with `errors.Is(err, gocmcapiv2.ErrNotFound)`.
- Debug logging: set `DEBUG_CMCCLOUD_TERRAFORM=DEBUG` to append every API request and response to `log.txt` (or the path in `DEBUG_CMCCLOUD_LOGFILE`) via `Logs`/`Logo`.

### Provider (`cmccloudv2`)

- `provider.go` holds the provider schema (`api_endpoint`, `api_key`, `project_id`, `region_id`) and the `ResourcesMap` / `DataSourcesMap` registries. `meta` is a `*CombinedConfig`; get the client with `getClient(meta)`.
- File naming per resource:
  - `resource_cmccloud_<name>.go` — `resource<Name>()` plus CRUD, import and wait functions.
  - `schema_cmccloud_<name>.go` — `<name>Schema()` returning the schema map, kept separate from the CRUD code.
  - `datasource_cmccloud_<name>.go` — data sources (schema inline).
- Several resource types share one parameterised implementation: `resourceDBv2Configuration("Mysql")`, `resourceDBv2Backup("mysql")`, `resourceDBv2AutoBackup("Mongo")` are registered once per database engine with the engine name passed in.
- `v2` resources (`kubernetesv2`, `autoscalingv2_*`) coexist with the older versions; both are registered and maintained.
- Shared helpers in `util.go`:
  - Async operations: `waitUntilResourceStatusChanged` / `waitUntilResourceIdStatusChanged` (poll until a target status, fail on listed error statuses) and `waitUntilResourceDeleted` (poll until not-found), all configured with `WaitConf`. Each resource wraps these in its own `waitUntilXxx...` function.
  - Read functions: computed/required fields use `d.Set` directly, while optional fields use `setString` / `setInt` / `setBool` / `setTypeSet`, which only write to state when the user set the field, to avoid spurious diffs against defaults.
  - Block and set helpers: `getFirstBlock`, `isSubBlockFieldChanged`, `getDiffSet`, `getStringArrayFromTypeSet`.
- `validator.go` holds shared `ValidateFunc`s (`validateName`, `validateBillingMode`, ...); `schema_cmccloud_common.go` holds `tagSchema()`.
- Billing mode changes go through the separate `client.BillingMode` service (`SetXxxBilingMode`), not the resource's own update endpoint.

### Adding a resource

1. Add the service file in `gocmcapiv2/` and register it in `Client` and `NewClient`.
2. Add `schema_cmccloud_<name>.go` and `resource_cmccloud_<name>.go` in `cmccloudv2/`.
3. Register it in `ResourcesMap` (or `DataSourcesMap`) in `provider.go`.
4. Add `examples/resources/cmccloudv2_<name>/resource.tf` (or `examples/data-sources/cmccloudv2_<name>/data-source.tf`), then run `docs-gen.bat`.

### Docs

`docs/` is generated by `tfplugindocs` from schema `Description` fields and `examples/`; edit those, not the markdown. The example directory name must equal the registered type name exactly for the example to be picked up. Some existing directories do not match (e.g. `cmccloudv2_api_gateway` vs registered `cmccloudv2_apigateway`, `cmccloudv2_container_registry` vs `cmccloudv2_container_registry_repo`).

## Conventions

- Code comments are frequently in Vietnamese; keep existing ones intact.
- Errors are returned as `fmt.Errorf("error <doing what> [%s]: %v", id, err)`.
