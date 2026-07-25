## MODIFIED Requirements

### Requirement: Go 1.26 dependencies and CI have explicit value and cost gates
`go.mod` SHALL declare the Go language and module compatibility floor as
`1.26.0` and SHALL declare `toolchain go1.26.5` so default auto selection uses
the accepted patched standard library. Runtime dependencies SHALL be exactly
direct `golang.org/x/sync v0.22.0` and `golang.org/x/time v0.15.0` plus only
dependencies mechanically required by those modules. `x/sync/errgroup` SHALL
own bounded structured cancellation; `x/time/rate` SHALL own the
concurrency-safe token bucket. No runtime dependency type SHALL leak into the
public API.

Tests SHALL use standard-library `testing`, `httptest`, and helpers only; no assertion, mocking, HTTP-client, logging, or retry framework SHALL be introduced. `.github/workflows/ci.yml` SHALL be test/build-only, run Go 1.26.x format/module/vet/test/race gates, declare read-only contents permission, consume no secret, publish no artifact, deploy nothing, and use only `actions/checkout@3d3c42e5aac5ba805825da76410c181273ba90b1` with comment `v7.0.1` and `actions/setup-go@b7ad1dad31e06c5925ef5d2fc7ad053ef454303e` with comment `v7.0.0`.

The dependency admission record SHALL state: standard library lacks an equivalent cancellable token bucket; custom structured concurrency/limiting has higher cancellation/race risk; both runtime modules are official Go BSD-3-Clause modules with small module/supply-chain cost; CI actions are GitHub-owned and SHA-pinned. Removal SHALL require a later reviewed change, an equivalent standard/local replacement, and all behavior/race gates.

#### Scenario: Module and toolchain graph is exact
- **WHEN** default `GOTOOLCHAIN=auto` selection runs module and toolchain inspection
- **THEN** `go env GOVERSION` SHALL report `go1.26.5`, `go mod tidy -diff` SHALL have no diff, and `go.sum` SHALL remain unchanged
- **AND** direct versions SHALL be exactly `x/sync v0.22.0` and `x/time v0.15.0`
- **AND** no unreviewed production or test module SHALL appear

#### Scenario: Patched standard library has no reachable known vulnerability
- **WHEN** `govulncheck ./...` runs under the default selected Go 1.26.5 toolchain
- **THEN** it SHALL report no reachable vulnerability in the client or its standard-library call graph

#### Scenario: CI is least privilege and test-only
- **WHEN** the workflow is inspected mechanically
- **THEN** the two exact approved action OIDs SHALL be present, no other `uses:` target SHALL exist, permissions SHALL be read-only, Go SHALL be 1.26.x, and format/module/vet/test/race commands SHALL be present
- **AND** secret, credential, artifact publication, registry, release, deploy, SSH, and write-permission steps SHALL be absent
