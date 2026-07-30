## Why

Bangumi's live public collection API legitimately returns `comment: null` for
records without a comment. Version `v0.1.0` rejects that normal optional value
as a protocol error, so one null comment prevents callers from retrieving the
entire public collection and makes real personal queries fail.

## What Changes

- Treat an omitted or JSON-null optional `comment` as the existing public zero
  value `Subject.Comment == ""`.
- Continue rejecting every non-string, non-null present comment and preserve
  all existing required-field, identity, pagination, range, and size checks.
- Add focused decoder and aggregate tests using the observed nullable shape.
- Clarify the nullable optional-comment contract and publish the compatible
  patch as `v0.1.1` only after local and consumer acceptance.

## Capabilities

### New Capabilities

None.

### Modified Capabilities

- `collection-client-v0-1-0`: the optional comment contract accepts the
  upstream API's documented/live JSON-null representation as empty.

## Impact

| Boundary | Declaration |
|---|---|
| Status | Investigated and specified; implementation, verification, commit, push, tag, release, publication, consumer upgrade, and deployment are initially false. |
| Owner | Main agent owns specification, review, Git/release gates, and consumer acceptance; one delegated subagent owns implementation and local tests. |
| Writable/Owned paths | `types.go`; `types_test.go`; `README.md` only if required for the public nullable contract; `openspec/changes/accept-null-collection-comments/**`; lifecycle may synchronize `openspec/specs/collection-client-v0-1-0/spec.md` and archive this exact change. |
| Read-only protected inputs | Every other tracked path; the original worktree's user-owned `CLAUDE.md`, `note`, and ignored `.claude/settings.local.json`; the BangumiStaffStats repository until its separate dependency-integration change; live upstream response bodies and user data except bounded structural evidence; secrets and hosts. |
| Mutable refs | Only `codex/accept-null-comments`; later `v0.1.1` and its GitHub release may be created by the main agent after all gates pass. `main`, `codex/v0.1.0-hardening`, and `v0.1.0` remain fixed. |
| Consumes | Accepted `v0.1.0`, exact changes `harden-bangumi-collection-go-v0-1-0` and `pin-go-toolchain-patch`, official/live nullable-comment structure, and the existing public `Subject.Comment string` API. |
| Produces | A backward-compatible decoder fix, regression coverage, synchronized specification, and an accepted `v0.1.1` module revision. |
| Dependencies | `harden-bangumi-collection-go-v0-1-0`; `pin-go-toolchain-patch`. |
| Deliverables | Focused source/test/doc delta, strict OpenSpec evidence, full Go quality gates, reviewed commit, pushed branch, `v0.1.1`, and a separately verified BangumiStaffStats consumer upgrade. |
| Acceptance | `go test ./...`, `go test -race ./...`, `go vet ./...`, `go build ./...`, formatting/diff checks, strict OpenSpec validation, exact Git audit, consumer Actions, and one successful real personal query whose response is summarized without retaining personal payloads. |
| Non-goals | Relaxing required fields; accepting malformed comments; changing DTO signatures; changing pagination/retry/rate-limit behavior; logging or retaining upstream bodies; adding credentials; changing unrelated dependencies; modifying visual behavior. |
| Operations deferred | None inside this library. BangumiStaffStats build, production deployment, and live query verification are governed by its separate active change. |
| Stop/rollback conditions | Stop on any required-field regression, public API break, unexpected tracked/untracked overlap, failed race/vet/build/test/strict validation, inability to verify the exact module revision, or a consumer regression. Before tag, rollback is branch-local; after tag, do not move the tag and instead issue a new patch only with explicit review. |
