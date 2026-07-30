## Context

The `v0.1.0` wire decoder already treats `comment` as optional when the member
is omitted, but rejects an explicit JSON `null`. Bangumi currently emits both
strings and nulls for that member in ordinary public collection pages. Because
one protocol failure rejects the whole page and aggregate, this representation
gap makes valid personal queries unavailable to downstream applications.

| Boundary | Declaration |
|---|---|
| Status | Investigated and specified; implemented, verified, committed, pushed, released, published, consumer-upgraded, and deployed remain false until their gates complete. |
| Owner | Main agent: decisions, spec, review, Git/release, and consumer acceptance. Delegated subagent: `types.go`, focused tests, optional README clarification, and task markers. |
| Writable/Owned paths | `types.go`; `types_test.go`; optional `README.md`; this active change; later its exact archive and synchronized `openspec/specs/collection-client-v0-1-0/spec.md`. |
| Read-only protected inputs | All other library paths and refs; original-worktree local files; consumer source before its own spec; full live collection payloads; hosts and secrets. |
| Mutable refs | `codex/accept-null-comments`, then immutable new tag `v0.1.1` only after acceptance. |
| Consumes | `v0.1.0`, exact changes `harden-bangumi-collection-go-v0-1-0` and `pin-go-toolchain-patch`, existing tests, and bounded structural evidence that optional comment is string or null. |
| Produces | One compatible decoder correction, tests/specs, and an accepted patch release. |
| Dependencies | `harden-bangumi-collection-go-v0-1-0`; `pin-go-toolchain-patch`. |
| Deliverables | Reviewed source/test/doc delta, strict validation, full Go gates, exact commit and patch release, then separate consumer integration. |
| Acceptance | Focused null/omitted/malformed cases plus `go test ./...`, `go test -race ./...`, `go vet ./...`, `go build ./...`, formatting/diff/strict OpenSpec checks, consumer Actions, and a bounded real-query summary. |
| Non-goals | Any other decode relaxation, DTO signature change, retry/pagination alteration, live payload storage, credential behavior, or UI change. |
| Operations deferred | Library has no operations. Consumer build/deploy remains under BangumiStaffStats governance. |
| Stop/rollback conditions | Stop for a required-field/API regression, scope overlap, failed quality gate, ambiguous module revision, or consumer failure. Never move a published tag. |

## Goals / Non-Goals

**Goals:**

- Accept both legal empty representations of optional comment.
- Preserve the current `string` DTO and exact empty-string behavior.
- Prove malformed non-null comments remain terminal protocol failures.

**Non-Goals:**

- Infer a comment value or distinguish omitted from null in the public DTO.
- Relax nested subject, tags, identity, enum, range, metadata, or body limits.
- Add an application-specific workaround outside the reusable client.

## Decisions

### Normalize optional null at the wire boundary

`decodeOptionalString` will return the empty string for either an absent raw
member or the exact JSON literal `null`. It will continue to JSON-decode every
other present value into a string and return `ProtocolError` for numbers,
objects, arrays, booleans, invalid JSON, or trailing content.

This is preferred over changing `Subject.Comment` to `*string`, because the
existing public contract intentionally collapses an absent optional comment to
empty and consumers do not need presence semantics. It is preferred over a
BangumiStaffStats adapter workaround because `v0.1.0` rejects the page before
the adapter receives it.

### Retain strict checks everywhere else

The change is limited to comment nullability. Tests will pair the accepted null
case with existing omitted/string cases and explicit malformed non-null cases.
No test will use or retain a real user's comment text.

### Publish a patch revision

The behavior is a backwards-compatible correction to `v0.1.0`, so the accepted
module version is `v0.1.1`. The consumer will pin that exact module revision
and prove its own build and production query separately.

## Risks / Trade-offs

- **Null can be confused with general permissiveness** → assert every other JSON
  type still returns the existing non-retryable protocol class.
- **A library-only test can miss the production path** → require the consumer's
  full Actions build and the same real personal query after deployment.
- **Live evidence may expose personal content** → record only field-type/count
  summaries and response/request identifiers, never response bodies or comment
  values.

## Migration Plan

1. Implement and run focused plus full library gates on the topic branch.
2. Review, archive/sync, commit, push, and create immutable `v0.1.1`.
3. In a separately specified BangumiStaffStats change, upgrade `go.mod` and
   `go.sum`, pass Actions, deploy the exact bundle, and rerun the bounded query.
4. If consumer acceptance fails before production, keep `v0.1.0` pinned. If an
   issue is discovered after publishing, retain `v0.1.1` and issue a separately
   reviewed later patch rather than moving the tag.

## Open Questions

None. The live wire shape and desired public zero-value mapping are explicit.
