## Why

The current unversioned client omits collection fields required by the formal BangumiStaffStats backend, fetches pages with request-local rather than client-wide resource controls, exposes upstream response bodies through errors, and has no real test or CI gate. This change establishes the first publishable `v0.1.0` contract before the main repository is allowed to admit the client.

## What Changes

- Complete the public collection DTO with `subjectID`, `subjectType`, collection `type`, `rate`, `comment`, `tags`, `updatedAt`, `volStatus`, `epStatus`, and `private`, while retaining the existing subject identity/name fields used by current callers.
- Make `Fetch` retrieve every page automatically, normalize repeated requested collection types, deduplicate by the full collection identity, and return a deterministic total order independent of goroutine completion order.
- Make concurrency and QPS limits client-wide across concurrent `Fetch` and `FetchPage` calls. Every network attempt, including retries, passes the same context-aware limiter and in-flight gate.
- Retry only transport failures, HTTP 429, and HTTP 5xx. Honor a valid `Retry-After`, otherwise use bounded exponential full-jitter backoff; cancellation is terminal and all waits remain context-cancellable.
- Replace body-leaking failures with typed, bounded, sanitized errors that preserve `errors.Is`/`errors.As` classification for input, 401, 403, 404, 429, 5xx, transport, timeout, cancellation, decode, protocol, oversized response, and retry exhaustion.
- Add an injectable root endpoint and safe custom HTTP client handling, strict UID/User-Agent/option/rate validation, a per-attempt timeout that remains effective with a custom client, complete `httptest` coverage, race coverage, `go vet`, module-tidiness checks, and a test-only GitHub Actions workflow using Go 1.26. CI admits only GitHub-owned `actions/checkout@3d3c42e5aac5ba805825da76410c181273ba90b1` (`v7.0.1`) and `actions/setup-go@b7ad1dad31e06c5925ef5d2fc7ad053ef454303e` (`v7.0.0`).
- Pin the only runtime helper modules to `golang.org/x/sync v0.22.0` and `golang.org/x/time v0.15.0`; use only the standard `testing`/`httptest` stack for tests.
- **BREAKING**: remove `WithAccessToken` and all Authorization/Cookie behavior. The first public `v0.1.0` is anonymous-public-collection-only; the untagged baseline is evidence, not a compatibility promise for credentialed access.
- Keep `NewClient`, `Fetch`, `FetchPage`, existing enum values, existing non-auth options, existing sentinel compatibility, `Subject.ID`, `Subject.Name`, and `Subject.NameCn` source-compatible wherever this does not conflict with the anonymous and sanitized-error boundary.
- Commit the approved five-file OpenSpec planning checkpoint before product apply, then finish only with a separately accepted local development commit. Fetch, push, PR, tag, release, publication, deployment, and a public `v0.1.0` remain later, separately authorized actions.

## Capabilities

### New Capabilities

- `collection-client-v0-1-0`: Defines the public collection DTO, pagination, shared limiter, retry, error, deterministic output, compatibility, tests, CI, local-commit, and publication-gate contract for the first versioned client.

### Modified Capabilities

None. `external-openspec-bootstrap` remains unchanged and is consumed as the governance prerequisite.

## Impact

- Public package implementation and documentation: `README.md`, `client.go`, `collection.go`, `errors.go`, `limiter.go`, `options.go`, `request.go`, `retry.go`, `types.go`, and `example/example.go`.
- Module metadata: `go.mod`, `go.sum`.
- Tests: `client_test.go`, `collection_test.go`, `errors_test.go`, `limiter_test.go`, `options_test.go`, `request_test.go`, `retry_test.go`, `types_test.go`.
- Test-only CI: `.github/workflows/ci.yml`.
- No main-repository product path changes, no token/OAuth integration, no backend cache/statistics/domain behavior, and no remote mutation.

## Change Boundary

- **Status**
  - investigated: yes — the exact ten baseline product files, bootstrap root spec/archive, formal master plan external lane, and backend guide section 6.3 were read.
  - specified: yes — all five active-change files passed strict validation, independent semantic review, and main-agent approval; exact planning checkpoint commit pending.
  - implemented: no.
  - verified: specification validation and semantic review complete; no implementation verification.
  - committed: bootstrap governance only; no hardening planning or implementation commit.
  - pushed: no.
  - released: no.
  - deployed: no.
- **Owner**: delegated External collection client apply subagent; main agent may review and amend only OpenSpec artifacts and performs read-only acceptance.
- **Writable/Owned paths**
  - Planning phase only: `openspec/changes/harden-bangumi-collection-go-v0-1-0/.openspec.yaml`, `openspec/changes/harden-bangumi-collection-go-v0-1-0/proposal.md`, `openspec/changes/harden-bangumi-collection-go-v0-1-0/design.md`, `openspec/changes/harden-bangumi-collection-go-v0-1-0/specs/collection-client-v0-1-0/spec.md`, `openspec/changes/harden-bangumi-collection-go-v0-1-0/tasks.md`.
  - Approved apply product/docs/module paths: `README.md`, `client.go`, `collection.go`, `errors.go`, `limiter.go`, `options.go`, `request.go`, `retry.go`, `types.go`, `example/example.go`, `go.mod`, `go.sum`.
  - Approved apply test/CI paths: `client_test.go`, `collection_test.go`, `errors_test.go`, `limiter_test.go`, `options_test.go`, `request_test.go`, `retry_test.go`, `types_test.go`, `.github/workflows/ci.yml`.
  - Approved final archive/sync paths, only after implementation acceptance: `openspec/changes/archive/2026-07-24-harden-bangumi-collection-go-v0-1-0/.openspec.yaml`, `openspec/changes/archive/2026-07-24-harden-bangumi-collection-go-v0-1-0/proposal.md`, `openspec/changes/archive/2026-07-24-harden-bangumi-collection-go-v0-1-0/design.md`, `openspec/changes/archive/2026-07-24-harden-bangumi-collection-go-v0-1-0/specs/collection-client-v0-1-0/spec.md`, `openspec/changes/archive/2026-07-24-harden-bangumi-collection-go-v0-1-0/tasks.md`, `openspec/specs/collection-client-v0-1-0/spec.md`.
- **Read-only protected inputs**
  - User paths: untracked `CLAUDE.md` at SHA-256 `c69937414d643411e03f3d7982e4bf58e8ceb02802cf1338dada4ea6a630407d`; untracked `note` at `7b2376f605f5d9d9f677eaf1165a2ff911cec9092b6b2763b966a3d6f93e2e74`; ignored `.claude/settings.local.json` at `c67845ef6924b080344bc0010931336106d9df3a700f67f10a22763d07517e3e`.
  - Repository governance: `.codex/skills/openspec-apply-change/SKILL.md`, `.codex/skills/openspec-archive-change/SKILL.md`, `.codex/skills/openspec-explore/SKILL.md`, `.codex/skills/openspec-propose/SKILL.md`, `.codex/skills/openspec-sync-specs/SKILL.md`, `.codex/skills/openspec-update-change/SKILL.md`, `openspec/config.yaml`, `openspec/specs/external-openspec-bootstrap/spec.md`, and every path below `openspec/changes/archive/2026-07-24-bootstrap-bangumi-collection-go-openspec/`.
  - Repository product input not owned by apply: `LICENSE` at baseline blob `ebf24803f903b14347ecd6ee220df2d68b66c530`.
  - Cross-repository planning inputs: `/Users/luca/dev/BangumiStaffStats/tmp-formal-development/formal-development-master-plan.md` and `/Users/luca/dev/BangumiStaffStats/tmp-formal-development/backend-development-implementation-guide.md`.
  - Exact initial hardening checkpoint: branch `codex/v0.1.0-hardening`, `HEAD=59eba79b3e3621ce72b756a88b38aee970c00fcf`, with baseline `8173f44911360150a5a5a7c6418021d1014fe85b` as an ordinary ancestor.
- **Consumes**
  - Archived change `bootstrap-bangumi-collection-go-openspec` and root capability `external-openspec-bootstrap`.
  - The ten baseline paths `LICENSE`, `README.md`, `client.go`, `collection.go`, `errors.go`, `example/example.go`, `go.mod`, `go.sum`, `options.go`, and `types.go` as current behavior evidence.
  - The master-plan external-lane row and backend guide sections 6.2–6.3 as product admission authority.
- **Produces**
  - Capability `collection-client-v0-1-0`.
  - The exact approved product/test/CI paths above.
  - After planning approval, an exact five-file local planning commit with subject `docs(openspec): approve collection client hardening`; after two implementation/finalization read-only acceptance gates, a synchronized root spec, archived change, and one local development commit with subject `feat: harden collection client for v0.1.0`.
- **Dependencies**: exact change ID `bootstrap-bangumi-collection-go-openspec`; satisfied by local governance commit `59eba79b3e3621ce72b756a88b38aee970c00fcf`. No main-repository implementation change is an apply dependency.
- **Deliverables**
  - Complete anonymous collection DTO and compatibility documentation.
  - Correct automatic pagination, shared concurrency/QPS, bounded retry/`Retry-After`, cancellation, deterministic dedupe/sort, and typed sanitized errors.
  - Non-empty `httptest` suites, race/vet/module gates, and no-publication Go 1.26 CI.
  - Main-agent accepted five-file planning checkpoint, unstaged implementation candidate, staged archive candidate, and exact local implementation commit.
- **Acceptance**
  - The sorted planning-candidate full status SHALL be exactly:

    ```text
    !! .claude/settings.local.json
    ?? CLAUDE.md
    ?? note
    ?? openspec/changes/harden-bangumi-collection-go-v0-1-0/.openspec.yaml
    ?? openspec/changes/harden-bangumi-collection-go-v0-1-0/design.md
    ?? openspec/changes/harden-bangumi-collection-go-v0-1-0/proposal.md
    ?? openspec/changes/harden-bangumi-collection-go-v0-1-0/specs/collection-client-v0-1-0/spec.md
    ?? openspec/changes/harden-bangumi-collection-go-v0-1-0/tasks.md
    ```

  - Spec-candidate commands:

    ```sh
    test "$(git branch --show-current)" = codex/v0.1.0-hardening
    test "$(git rev-parse HEAD)" = 59eba79b3e3621ce72b756a88b38aee970c00fcf
    git merge-base --is-ancestor 8173f44911360150a5a5a7c6418021d1014fe85b HEAD
    test -z "$(git diff --cached --name-only)"
    test "$(shasum -a 256 CLAUDE.md | awk '{print $1}')" = c69937414d643411e03f3d7982e4bf58e8ceb02802cf1338dada4ea6a630407d
    test "$(shasum -a 256 note | awk '{print $1}')" = 7b2376f605f5d9d9f677eaf1165a2ff911cec9092b6b2763b966a3d6f93e2e74
    test "$(shasum -a 256 .claude/settings.local.json | awk '{print $1}')" = c67845ef6924b080344bc0010931336106d9df3a700f67f10a22763d07517e3e
    test "$(git status --porcelain=v1 --ignored --untracked-files=all -- CLAUDE.md note .claude/settings.local.json | shasum -a 256 | awk '{print $1}')" = 221322afd4137d24875b5fbd369c9028e99f69f1cf155e3c49975996841a3938
    test "$(printf '%s\0' LICENSE README.md client.go collection.go errors.go example/example.go go.mod go.sum options.go types.go | xargs -0 shasum -a 256 | shasum -a 256 | awk '{print $1}')" = 7d4e55231e08c5d8286bffb921e53adbda9a5d82ced9c4688e4cadb129b1855e
    test "$(git status --porcelain=v1 --ignored --untracked-files=all | LC_ALL=C sort | shasum -a 256 | awk '{print $1}')" = e91f861f8e2d4e4398f6ca39ff39a257d71f3cf0fcf0a7382de82ca497796ef5
    openspec status --change harden-bangumi-collection-go-v0-1-0
    openspec validate harden-bangumi-collection-go-v0-1-0 --strict
    openspec validate --all --strict
    openspec doctor --json
    git diff --check
    test -z "$(find openspec/changes/harden-bangumi-collection-go-v0-1-0 -type l -print)"
    ```

  - Planning checkpoint acceptance stages exactly the five active-change files, requires sorted path seal `a7da0916df9b59cdaff4d5a3b57d8b0231c8bad51921aedb308ee2d10262fe1f`, and creates exact subject `docs(openspec): approve collection client hardening` with sole parent `59eba79b3e3621ce72b756a88b38aee970c00fcf`. The main agent records its accepted commit as `HARDENING_PLANNING_HEAD`.
  - Apply acceptance additionally runs from exact `HARDENING_PLANNING_HEAD`, under a Go 1.26.x toolchain and with no real Bangumi network request:

    ```sh
    test "$(go env GOVERSION | sed -E 's/^go([0-9]+\.[0-9]+).*/\1/')" = 1.26
    test -z "$(gofmt -d -e client.go collection.go errors.go limiter.go options.go request.go retry.go types.go example/example.go client_test.go collection_test.go errors_test.go limiter_test.go options_test.go request_test.go retry_test.go types_test.go)"
    go mod tidy -diff
    test "$(go list -m -f '{{.Version}}' golang.org/x/sync)" = v0.22.0
    test "$(go list -m -f '{{.Version}}' golang.org/x/time)" = v0.15.0
    go vet ./...
    go test -count=1 ./...
    go test -race -count=1 ./...
    git diff --check
    ```

  - Exact expected evidence: tests exist in the root package and example package is buildable; bounded pagination, limiter, retry, errors, dedupe/sort, cancellation, credential absence, response bounds, input boundaries, custom-client timeout behavior, and compatibility scenarios pass; the index is empty at the first implementation acceptance; protected seals/status are unchanged; every changed path is in the approved list; `HEAD` remains `HARDENING_PLANNING_HEAD`.
  - Finalization acceptance: after archive/sync and staging, the cached no-rename parent delta is exactly 32 paths: the 21 product/test/CI paths, five active-change deletions, five date-stamped archive additions, and `openspec/specs/collection-client-v0-1-0/spec.md`. Its sorted path-only SHA-256 is `740610fc014608d9294383d731e3a9df51b3521c25823824c6ef89dc582349b2`. The cumulative final tree relative to bootstrap contains exactly the 27 product/archive/root paths with sorted path-only SHA-256 `b8ef7390c49b73f13061c3ca50f66a0ef260cf947c539d507ab1797de59c48c1`. `git diff --cached --check`, strict/all/doctor and all Go gates pass; main agent accepts that exact index before the local implementation commit. The commit has sole parent `HARDENING_PLANNING_HEAD`, exact subject `feat: harden collection client for v0.1.0`, no other paths, and leaves only the three protected user paths in tolerated status, verified with:

    ```sh
    test "$(git diff --cached --name-only --no-renames | LC_ALL=C sort | shasum -a 256 | awk '{print $1}')" = 740610fc014608d9294383d731e3a9df51b3521c25823824c6ef89dc582349b2
    test "$(git diff --cached --name-only --no-renames | wc -l | tr -d ' ')" = 32
    test "$(git diff --cached --name-only --no-renames 59eba79b3e3621ce72b756a88b38aee970c00fcf | LC_ALL=C sort | shasum -a 256 | awk '{print $1}')" = b8ef7390c49b73f13061c3ca50f66a0ef260cf947c539d507ab1797de59c48c1
    test "$(git diff --cached --name-only --no-renames 59eba79b3e3621ce72b756a88b38aee970c00fcf | wc -l | tr -d ' ')" = 27
    git diff --cached --check
    test "$(awk 'p && NF {print; exit} /^## Purpose$/{p=1}' openspec/specs/collection-client-v0-1-0/spec.md)" = "Define the first public anonymous Bangumi collection client contract, including complete DTOs, automatic pagination, shared limiting, bounded retries, sanitized errors, deterministic results, and local quality gates."
    test "$(awk '/^### Requirement:/{p=1} /^## Boundary Record$/{p=0} p' openspec/changes/archive/2026-07-24-harden-bangumi-collection-go-v0-1-0/specs/collection-client-v0-1-0/spec.md | shasum -a 256 | awk '{print $1}')" = "$(awk '/^### Requirement:/{p=1} /^## Boundary Record$/{p=0} p' openspec/specs/collection-client-v0-1-0/spec.md | shasum -a 256 | awk '{print $1}')"
    ```
- **Non-goals**
  - Main-repository cache, digest, stale policy, statistics, Archive, handler, or domain implementation.
  - Token/OAuth/private collection support, caller Cookie/Authorization forwarding, write endpoints, collection persistence, or a generic Bangumi API SDK.
  - Fetch, push, PR, tag, release, publication, deployment, live upstream load tests, or changes outside this external repository.
- **Operations deferred**: all remote Git/GitHub actions, `v0.1.0` tag/release/publication, package-consumer admission, production rate values, monitoring, rollout, and deployment.
- **Stop/rollback conditions**
  - Stop without cleaning, resetting, stashing, deleting, overwriting, staging, or committing if branch/ancestry, protected hashes/status, root governance files, owned-path complement, dependency versions, Go version, strict validation, tests, race, vet, module tidiness, or full status differs.
  - A failed apply remains an unstaged candidate. Rollback is a new reviewed OpenSpec decision; no destructive restoration is authorized.
  - A failed staged acceptance remains uncommitted. A failed local commit proof does not authorize amend/rebase/reset; report the exact mismatch.
- **Mutable refs**
  - After spec approval, one delegated planning-checkpoint subagent may advance `refs/heads/codex/v0.1.0-hardening` once from `59eba79b3e3621ce72b756a88b38aee970c00fcf` to the exact accepted five-file planning commit; no other ref may move.
  - During apply and both implementation/finalization candidate reviews, `HEAD=HARDENING_PLANNING_HEAD` and every ref remain fixed.
  - Only after unstaged implementation acceptance and staged archive acceptance may a delegated finalization subagent advance the same branch once more with the exact local implementation commit. `refs/heads/main`, `refs/remotes/origin/main`, tags, remote refs, and all remote state remain immutable.
