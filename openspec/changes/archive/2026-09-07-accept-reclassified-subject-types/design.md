## Context

The collection record owns the public SubjectType. The nested subject contributes names for the same subject ID. Upstream can return another valid type in that nested metadata; treating it as a different ID prevents otherwise usable complete collections.

## Goals / Non-Goals

Accept the observed metadata difference without changing DTO fields, dropping a record or relaxing malformed data handling. No dependency/default-limit changes, no consumer edits or publication in this phase.

## Boundary

| Field | Decision |
| --- | --- |
| Status | Specified; reviewed/strict-valid before implementation. |
| Owner | Primary: types.go, README and lifecycle; delegated test owner: types_test.go. |
| Writable/Owned paths | The exact paths in proposal.md Impact. |
| Read-only protected inputs | All other files, historical archives, user-owned notes/settings and consumer repository. |
| Mutable refs | None; main/748cfa4 baseline. |
| Consumes | Existing DTO authority and real structural evidence. |
| Produces | Supported-enum nested metadata compatibility with unchanged ID validation. |
| Dependencies | None. |
| Deliverables | Tests, decoder, documentation and accepted spec. |
| Acceptance | Go format/tidy/test/race/vet/build, strict spec and live reproduction as in tasks.md. |
| Non-goals | Defaults, new APIs, field inference, records filtering, publication and consumer upgrade. |
| Operations deferred | All remote mutations and deployment. |
| Stop/rollback conditions | Preserve unexpected work; stop on scope expansion or failed contract regressions. |

## Decisions

1. Keep the outer subject type equal to the requested type and preserve it in Subject.SubjectType. A valid different nested type cannot rewrite a collection record or silently remove it.
2. Replace nested-type equality with supported-enum validation. The nested object still needs matching positive ID, supported type and both name strings. Missing/null/invalid types remain errors.
3. Remove the now-unneeded expected-type argument from the private nested decoder; exported API stays identical.
4. Test all supported nested types for an anime collection, malformed nested values, wrong IDs and wrong outer filters, plus both page and aggregate calls. Include the observed 2/6 shape as a synthetic fixture without private comments.

## Risks / Trade-offs

The DTO intentionally represents the collection's type, which may differ from current subject metadata. This is already the documented field authority; document it explicitly. Statistical consumers must retain their own authoritative subject catalog rather than infer a reclassification from this optional metadata.

## Migration Plan

Complete local validation and archive the spec. A later authorized compatible tag (proposed v0.1.2) is required before the consumer can use it under its immutable-version contract.

## Open Questions

None for this repair. Publication authorization is separate from local implementation.
