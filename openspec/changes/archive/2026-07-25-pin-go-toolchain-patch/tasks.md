## 1. Preflight

- [x] 1.1 Verify branch `codex/v0.1.0-hardening`, `HEAD=45fa2da`, empty index, accepted dependency, exact owned/protected status, unchanged `go.sum`, and strict-valid change before product mutation.

## 2. Toolchain Patch

- [x] 2.1 Add `toolchain go1.26.5` while retaining `go 1.26.0`; update README to distinguish the compatibility floor from the required local/release patch toolchain.
- [x] 2.2 Prove default auto selection reports `go1.26.5`, `go mod tidy -diff` and `go mod verify` are clean, and `go.sum` is byte-unchanged.

## 3. Acceptance

- [x] 3.1 Run `go vet ./...`, `go test -count=20 ./...`, `go test -race ./...`, `go build ./...`, and `govulncheck ./...`; require zero failures and zero reachable vulnerabilities.
- [x] 3.2 Run exact-path `git diff --check`, strict change/all validation, and `openspec doctor`; hand off only the unstaged owned diff with protected user paths/index unchanged and every delivery state reported separately.

| Boundary | Declaration |
|---|---|
| Status | investigated: complete; specified: approved; implemented: complete; verified: all tasks and acceptance gates passed; committed/pushed/released/deployed: no |
| Owner | Main agent owns specification decisions; one delegated subagent owns tasks 1.1–3.2. |
| Writable/Owned paths | `go.mod`, `README.md`, `openspec/specs/collection-client-v0-1-0/spec.md`, and `openspec/changes/pin-go-toolchain-patch/**` |
| Read-only protected inputs | All Go source/tests/examples, `go.sum`, `.github/**`, archived changes, bootstrap spec, `CLAUDE.md`, `note`, `.claude/settings.local.json`, index, tags, remotes, other refs, and other repositories |
| Consumes | Accepted `harden-bangumi-collection-go-v0-1-0` at `45fa2da` and Go 1.26.5 audit evidence |
| Produces | Go 1.26.5 default toolchain pin, documentation, and executable release evidence |
| Dependencies | `harden-bangumi-collection-go-v0-1-0` |
| Deliverables | Reviewed delta; two product edits; test/security/OpenSpec evidence; later archive/sync and local commit |
| Acceptance | Exact commands in tasks 2.2–3.2 pass; `go.sum`, protected paths, index, remote refs, and tags remain unchanged |
| Non-goals | Product API/behavior, dependencies, CI, source/test changes |
| Operations deferred | Fetch, push, PR, tag, release, publication, deployment, and main-repository admission |
| Mutable refs | None during apply; only the local hardening branch may advance once after main acceptance |
| Stop/rollback conditions | Stop on branch/HEAD/dependency/protected/index/sum drift, failed gate, reachable vulnerability, or remote mutation; reverse only exact unstaged owned edits |
