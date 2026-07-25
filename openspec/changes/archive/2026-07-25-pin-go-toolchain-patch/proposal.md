## Why

The accepted client still lets a default local build select Go 1.26.0, whose
standard library has reachable vulnerabilities that are fixed by Go 1.26.5.
The first public tag must not silently build with that known-vulnerable patch.

## What Changes

- Keep the language compatibility floor at Go 1.26.0 while pinning the
  repository's preferred toolchain to Go 1.26.5.
- Document and verify the exact default local toolchain and zero reachable
  vulnerability result before the public v0.1.0 gate.

## Capabilities

### New Capabilities

None.

### Modified Capabilities

- `collection-client-v0-1-0`: require the Go 1.26.5 patch toolchain for local
  development and release acceptance while retaining Go 1.26.0 language
  compatibility.

## Impact

Only `go.mod`, the README toolchain guidance, the synchronized root capability,
and this change's artifacts are affected. The public Go API, runtime
dependencies, `go.sum`, source, tests, CI behavior, and supported data remain
unchanged.

| Boundary | Declaration |
|---|---|
| Status | investigated: complete; specified: approved; implemented: complete; verified: toolchain, tidy, module, vet, repeated/race tests, build, vulnerability, OpenSpec, and repository hygiene gates passed; committed/pushed/released/deployed: no |
| Owner | Main agent owns and reviews this minimal specification; one delegated subagent owns apply and verification. |
| Writable/Owned paths | `go.mod`, `README.md`, `openspec/specs/collection-client-v0-1-0/spec.md`, and `openspec/changes/pin-go-toolchain-patch/**` |
| Read-only protected inputs | All Go source/tests/examples, `go.sum`, `.github/**`, archived changes, `openspec/specs/external-openspec-bootstrap/**`, `CLAUDE.md`, `note`, `.claude/settings.local.json`, index, tags, remotes, other refs, and other repositories |
| Consumes | Accepted local client at `45fa2da` and the Go 1.26.5 vulnerability audit |
| Produces | Exact Go 1.26.5 toolchain selection and release-gate evidence |
| Dependencies | `harden-bangumi-collection-go-v0-1-0` |
| Deliverables | Minimal delta spec, toolchain pin, README update, tests/security verification, synchronized archive, and one local commit |
| Acceptance | `go env GOVERSION` under default auto selection reports `go1.26.5`; `go mod tidy -diff`, `go vet ./...`, `go test -count=20 ./...`, `go test -race ./...`, `go build ./...`, `go mod verify`, `govulncheck ./...`, strict OpenSpec validation/doctor, and `git diff --check` all pass; protected local paths and `go.sum` remain unchanged |
| Non-goals | API, DTO, pagination, limits, retries, errors, dependency or CI behavior changes |
| Operations deferred | Fetch, push, PR, tag, release, publication, deployment, and main-repository dependency admission |
| Mutable refs | After acceptance, only local `refs/heads/codex/v0.1.0-hardening` may advance once through one focused commit; no remote or tag ref may move |
| Stop/rollback conditions | Stop on protected-path/index drift, dependency or sum drift, nonzero reachable vulnerability, failed test/race/static/build/OpenSpec gate, or any requested remote mutation; reverse only exact unstaged owned changes |
