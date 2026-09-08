## Boundary

| Field | Decision |
| --- | --- |
| Status | Investigated/specified; implemented, verified, committed, pushed, released and deployed initially false. |
| Owner | Primary plans/reviews/implements decoder/docs; one delegated owner implements regression tests only. |
| Writable/Owned paths | types.go, types_test.go, README.md, this change, and its accepted main capability spec. |
| Read-only protected inputs | All other files and repositories; historical bootstrap artifacts and user-local files. |
| Mutable refs | None. |
| Consumes | Reviewed proposal/design/delta and current package implementation. |
| Produces | Verified local package patch and evidence. |
| Dependencies | None. |
| Deliverables | Fix, regression tests, docs, spec sync/archive. |
| Acceptance | Commands listed below. |
| Non-goals | Any dependency/limit/API shape change, consumer upgrade or remote publication. |
| Operations deferred | All remote mutations and production. |
| Stop/rollback conditions | Check main/748cfa4 and non-overlapping writable files; preserve all unrelated work; no destructive Git cleanup. |

## 1. Plan

- [x] 1.1 Finish and review artifacts and run strict change validation before apply.

## 2. Implement

- [x] 2.1 Primary changes only nested valid-type compatibility and corresponding README rule.
- [x] 2.2 Delegated owner adds positive and negative decoder/page/aggregate regressions in types_test.go.

## 3. Verify and close

- [x] 3.1 Primary reviews actual diff and runs pinned Go formatting, tidy, tests, race, vet and build.
- [x] 3.2 Reproduce the real 631949 page and complete zhong_mo collection with the local package; record only structural/count evidence.
- [x] 3.3 Sync accepted specification, archive this change and validate all specs/diff.

## Evidence

Primary review accepted the same-ID/valid-enum compatibility rule, with outer SubjectType and requested-filter checks retained. Strict change validation passed before apply. Delegated implementation owner for types_test.go: collection_regression_tests; primary owns types.go and README.

Live probe: local checkout's module, Go 1.26.5 Windows build, shared rate 5/s burst 10, per-attempt 10s and overall 90s. zhong_mo anime completed+in_progress returned 7,179 items in 28.3499s with 145 HTTP 200 pages, no transport errors. Subject 631949 was retained as SubjectType=2, name=銀河少年隊, nameCN=银河少年队 despite nested type=6. No upstream body/comment/tag payload was retained.

Primary diff review confirmed only supported nested-type equality is relaxed; outer requested-type, nested ID/name/enum and all other existing checks remain. Test-owner diff was audited and rerun by primary.

Clean LF Linux export acceptance (HEAD 748cfa4 plus exact types.go/types_test.go/README patch; Go 1.26.5): gofmt clean, go mod tidy -diff clean, go test -count=1 ./... passed (0.650s package), go test -race -count=1 ./... passed (2.173s package), go vet ./... and go build ./... passed. Windows vet/build also passed; Windows tidy comparison only differed due checkout CRLF and was validated on the clean LF export. Test-only tools and live probe were outside the repository. No dependency versions or CI configuration changed.

Status: investigated, specified, implemented and verified; not committed, pushed, tagged, released, published or deployed. Consumer remains on v0.1.1 pending publication authorization.
