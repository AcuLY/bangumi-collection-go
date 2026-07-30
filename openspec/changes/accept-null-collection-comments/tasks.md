| Boundary | Declaration |
|---|---|
| Status | Investigated and specified. Implementation, verification, commit, push, tag/release/publication, consumer upgrade, and deployment are independently tracked below. |
| Owner | Main agent owns spec/audit/Git/release/consumer gates; delegated implementation agent owns only the source/test/doc paths and these task markers. |
| Writable/Owned paths | Apply: `types.go`, `types_test.go`, optional `README.md`, and checkbox markers here. Lifecycle: this exact change, its archive, and synchronized `openspec/specs/collection-client-v0-1-0/spec.md`. |
| Read-only protected inputs | Every other tracked path/ref; original-worktree user files; BangumiStaffStats until its separate change; full live payloads; hosts and secrets. |
| Mutable refs | Main agent only: `codex/accept-null-comments`, remote counterpart, immutable new `v0.1.1`, and release metadata. |
| Consumes | Accepted `v0.1.0`, changes `harden-bangumi-collection-go-v0-1-0` and `pin-go-toolchain-patch`, and bounded nullable-comment evidence. |
| Produces | Decoder correction, tests/specs, accepted commit/release, and an exact consumer dependency revision. |
| Dependencies | `harden-bangumi-collection-go-v0-1-0`; `pin-go-toolchain-patch`. |
| Deliverables | Focused code/tests/docs, strict/full Go evidence, exact Git lifecycle, patch release, consumer Actions/deployment/real-query evidence. |
| Acceptance | Commands below must exit zero; exact diffs/refs must match; the consumer query must return 200 without storing its body. |
| Non-goals | Other decoder relaxation, public API/pagination/retry changes, credentials, payload logging, UI changes, or unrelated dependency updates. |
| Operations deferred | Library operations are none; consumer operations stay in the BangumiStaffStats change. |
| Stop/rollback conditions | Stop on scope drift, protected-file change, API/required-field regression, failed gate, ambiguous revision, or consumer failure; never move a published tag. |

## 1. Planning acceptance

- [x] 1.1 Main agent audits proposal, design, delta spec, task boundaries, live
  structural evidence, exact worktree/ref state, and confirms zero P0/P1
  planning findings.
- [x] 1.2 Run OpenSpec 1.6.0 strict validation for
  `accept-null-collection-comments` before any product edit.

## 2. Delegated implementation

- [x] 2.1 Update only the optional-string wire decoder so omitted and exact
  JSON-null comments map to `""`, while every non-null non-string value keeps
  the existing typed protocol failure.
- [x] 2.2 Add focused page and aggregate regression tests for null comment,
  retain omitted/string coverage, and prove number/object/array/boolean
  comments remain rejected without partial output.
- [x] 2.3 Clarify `README.md` only if its current public nullable-comment text
  would otherwise contradict the corrected behavior; do not broaden docs.

## 3. Verification and main-agent audit

- [x] 3.1 Run `gofmt` on owned Go files, focused tests, `go test ./...`,
  `go test -race ./...`, `go vet ./...`, and `go build ./...`.
- [x] 3.2 Run strict change/all validation, `git diff --check`, exact
  owned/protected path audit, and hand off an unstaged/uncommitted candidate.
- [x] 3.3 Main agent reviews the implementation and evidence for zero P0/P1,
  verifies no live payload/user content entered the repository, and repeats
  proportionate checks before accepting.

## 4. Library and consumer lifecycle

- [ ] 4.1 Sync/archive the accepted OpenSpec, commit the exact library delta,
  push `codex/accept-null-comments`, and verify the pushed commit exactly.
- [ ] 4.2 Merge the accepted library change, create immutable `v0.1.1`, and
  verify `go list -m` can resolve its exact checksum; do not move `v0.1.0`.
- [ ] 4.3 Through a separate BangumiStaffStats OpenSpec, pin `v0.1.1`, pass
  Development Actions, deploy the exact admitted `linux/amd64` bundle, and
  rerun the same bounded `lucay126` personal ranking request.
- [ ] 4.4 Record investigated/specified/implemented/verified/committed/pushed/
  released/published/consumer-upgraded/deployed/live-query states separately.
