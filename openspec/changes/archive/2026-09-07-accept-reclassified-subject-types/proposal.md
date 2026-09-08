## Why

The live public collection response for subject 631949 contains top-level `subject_type: 2` and nested `subject.type: 6` with the same ID. The v0.1.1 equality check rejects this valid-enum metadata difference and aborts the entire collection fetch.

## What Changes

- Accept differing top-level and nested subject types when each is a supported enum and the subject IDs match.
- Preserve top-level `subject_type` in the public DTO, as already required; preserve the nested names and all collection fields.
- Keep malformed/missing nested fields, null subjects, invalid enum values, different IDs and requested-filter mismatches invalid.
- Cover direct decoding, FetchPage and multi-page Fetch and document the compatibility rule.

## Capabilities

### New Capabilities

None.

### Modified Capabilities

- `collection-client-v0-1-0`: accept supported nested metadata type drift without changing the returned collection type.

## Impact

| Field | Decision |
| --- | --- |
| Status | Investigated/specified; implementation and verification pending; no commit/push/release/deployment. |
| Owner | Primary agent owns plan, types.go, README, acceptance and lifecycle; delegated regression-test owner owns types_test.go only after planning acceptance. |
| Writable/Owned paths | `types.go`, `types_test.go`, `README.md`, this change directory, `openspec/specs/collection-client-v0-1-0/spec.md`. |
| Read-only protected inputs | All other tracked paths, dependency/toolchain files, historical archives, any user notes/settings, BangumiStaffStats and production services. |
| Mutable refs | None. Current main at 748cfa4; no branch switch, commit, push or tag in this change. |
| Consumes | Current v0.1.1 decoder and the observed same-ID/different-valid-type response shape. |
| Produces | Compatible decoder repair, positive/negative regressions, docs and synchronized spec. |
| Dependencies | None. This is a maintenance change after the completed initial hardening and null-comment changes, not a replay of their historical checkpoint workflow. |
| Deliverables | Locally verified package patch ready for a later v0.1.2 publication decision. |
| Acceptance | Go 1.26.5 formatting, go mod tidy -diff, go test ./..., go test -race ./..., go vet ./..., go build ./..., strict OpenSpec, diff hygiene and bounded live public-collection reproduction. |
| Non-goals | Changing returned SubjectType authority, adding fields, dropping records, weakening ID/enum/required-field checks, changing limits/retries/dependencies, or editing the consumer before a published version exists. |
| Operations deferred | Push, PR, tag, release/publication, consumer version upgrade and production deployment remain separate actions. |
| Stop/rollback conditions | Stop on unrelated edits, public API regressions or new failure categories requiring wider handling. Preserve source and evidence; no reset, clean, stash or broad deletion. |

Apply begins only after complete artifacts pass strict validation and primary review. The user's current request authorizes this new local checkout at `D:/Luca/Code/MyProject/bangumi-collection-go`; historical bootstrap paths/checkpoints describe their completed initial-hardening lane.
