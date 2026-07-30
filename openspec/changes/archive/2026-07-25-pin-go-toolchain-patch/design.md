## Context

`go 1.26.0` is the compatibility floor, not a safe release-toolchain pin.
With `GOTOOLCHAIN=auto`, the current module can select the initial 1.26.0
standard library; the release audit found reachable vulnerabilities that are
absent under 1.26.5. The product API and dependency graph are otherwise ready.

## Goals / Non-Goals

**Goals:**

- Make default repository development and release verification select Go
  1.26.5 while retaining 1.26.0 language/module compatibility.
- Prove the unchanged client under the patched standard library with all
  existing gates and `govulncheck`.

**Non-Goals:**

- No API, behavior, runtime dependency, CI action, or remote publication
  change.

## Decisions

Add `toolchain go1.26.5` beside `go 1.26.0`. Raising only the `go` directive
would unnecessarily change the module compatibility floor; relying only on CI
would leave default local and release-candidate commands ambiguous. README
states the distinction explicitly.

Acceptance runs with default auto toolchain selection and also records
`go env GOVERSION`. `govulncheck ./...` is an acceptance tool, not a runtime or
module dependency.

## Risks / Trade-offs

- [Risk] Go may download the pinned patch when absent locally. → Keep the
  standard `GOTOOLCHAIN=auto` mechanism and fail closed if 1.26.5 cannot be
  selected.
- [Risk] A future patch supersedes 1.26.5. → A later reviewed change may raise
  only the toolchain line after repeating all gates.

| Boundary | Declaration |
|---|---|
| Status | investigated: complete; specified: in progress; implemented/verified/committed/pushed/released/deployed: no |
| Owner | Main agent owns specification decisions; one delegated subagent owns apply and verification. |
| Writable/Owned paths | `go.mod`, `README.md`, `openspec/specs/collection-client-v0-1-0/spec.md`, and `openspec/changes/pin-go-toolchain-patch/**` |
| Read-only protected inputs | All Go source/tests/examples, `go.sum`, `.github/**`, archived changes, bootstrap spec, protected user paths, index, tags, remotes, other refs, and other repositories |
| Consumes | Accepted `45fa2da` client and Go 1.26.5 audit evidence |
| Produces | Deterministic patched toolchain selection and release verification |
| Dependencies | `harden-bangumi-collection-go-v0-1-0` |
| Deliverables | Delta spec, two product-file edits, verification, archive/sync, and one local commit |
| Acceptance | Default `go env GOVERSION` is `go1.26.5`; tidy/vet/test×20/race/build/mod-verify/govulncheck and strict OpenSpec/diff gates pass; `go.sum` and protected paths are unchanged |
| Non-goals | Product behavior, dependencies, CI, remote or tag mutation |
| Operations deferred | Fetch, push, PR, tag, release, publication, deployment, and main-repository admission |
| Mutable refs | Only the local hardening branch may advance once after acceptance |
| Stop/rollback conditions | Stop on protected/index/sum drift, any failed gate, reachable vulnerability, or remote mutation request; reverse only exact unstaged owned edits |
