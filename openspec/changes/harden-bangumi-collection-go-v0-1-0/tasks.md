## Pre-apply planning checkpoint protocol

After main-agent approval and before any product edit, one delegated checkpoint subagent stages exactly this active change's five files, requires a five-path count and sorted path seal `a7da0916df9b59cdaff4d5a3b57d8b0231c8bad51921aedb308ee2d10262fe1f`, and creates exact subject `docs(openspec): approve collection client hardening` with sole parent `59eba79b3e3621ce72b756a88b38aee970c00fcf`. It stages no product, protected, bootstrap, or main-repository path. Main-agent read-only acceptance records the resulting exact commit as `HARDENING_PLANNING_HEAD`; no apply checkbox may start before that acceptance.

## 1. Governed preflight and module boundary

- [ ] 1.1 Verify canonical root, branch `codex/v0.1.0-hardening`, `HEAD=HARDENING_PLANNING_HEAD`, exact planning subject/sole bootstrap parent/five-path seal, baseline ancestry, empty index, exact protected hashes/status, bootstrap dependency, and accepted planning-candidate status before any product edit.
- [ ] 1.2 Record the initial SHA-256 inventory for every approved existing product path and prove `LICENSE`, generated OpenSpec/config/bootstrap paths, protected user paths, and both main-repository guides are read-only.
- [ ] 1.3 Update `go.mod`/`go.sum` to Go `1.26.0`, direct `golang.org/x/sync v0.22.0`, and direct `golang.org/x/time v0.15.0`; run exact module/version/license admission and `go mod tidy -diff`.

## 2. Client configuration and complete DTO

- [ ] 2.1 Refactor `client.go` into immutable concurrently safe Client configuration with defaults for endpoint, request timeout, retries, retry bounds, QPS/burst, and in-flight concurrency; validate User-Agent as UTF-8/control-free/non-blank `1..256` bytes and record a nil Option as `ErrInvalidConfiguration` without panic.
- [ ] 2.2 Update `options.go` to retain non-auth baseline options, add endpoint/QPS/max-retry-delay options, accept only a root endpoint and finite positive QPS/positive burst, safely clone an injected HTTP client, clear its jar and Timeout, replace redirects with first-response refusal, make request timeout independent of option order/custom client, and remove `WithAccessToken`.
- [ ] 2.3 Update `types.go` with exact complete DTO fields, enum/range/page wire validation, compatibility aliases, deep tag copying, and non-nil empty collections.
- [ ] 2.4 Add `client_test.go`, `options_test.go`, and `types_test.go` covering defaults; nil Option; first-configuration-fault poisoning despite later valid options; NaN/positive infinity/negative infinity/zero/negative QPS; non-positive burst; root-path endpoint rules; User-Agent UTF-8/control/blank/byte limits; source-compatible non-auth call shapes; both custom-client/timeout option orders; caller HTTP-client immutability; enum/range/timestamp/identity validation; and caller mutation isolation.

## 3. Request and sanitized error pipeline

- [ ] 3.1 Implement `request.go` with root endpoint validation, Unicode-trim/case-preserving UID validation as UTF-8/control-free `1..256` bytes, one-segment escaping, anonymous headers, a fresh per-attempt timeout covering transport/read/validate/close even with a custom client, first-response redirect refusal, exact query parameters, bounded body reads, body close, strict single-value JSON decode, and DTO conversion.
- [ ] 3.2 Replace `errors.go` with the approved stable sentinels/codes and exact nil-context/HTTP/network/decode/protocol/retry matrix; preserve documented compatibility while keeping `HTTPError.Body` empty and every error string bounded/sanitized.
- [ ] 3.3 Add `request_test.go` and `errors_test.go` for safe request construction, UID UTF-8/control/byte boundaries, absence of credentials/cookies, custom transport/client timeout boundaries and option orders, fresh retry deadlines, first-response redirect refusal, 1xx/204/3xx/401/403/404/429/other-4xx/5xx mapping, nil context, transport/attempt-timeout/parent-cancel/parent-deadline races, cancellation during rate/semaphore/retry-wait/transport/body-read, `errors.Is/As`, secret/Location-marker sanitization, byte limits, trailing JSON, and body closure.

## 4. Shared limiting and bounded retry

- [ ] 4.1 Implement `limiter.go` so all Fetch/FetchPage attempts on one Client share one `x/time/rate` token bucket and one semaphore, rate-wait before acquiring concurrency, cancel every wait, and release exactly once.
- [ ] 4.2 Add `limiter_test.go` with timestamping/blocking loopback servers proving configured aggregate QPS/burst, maximum in-flight requests, cross-operation sharing, per-client isolation, cancellation, and balanced permits under race.
- [ ] 4.3 Implement `retry.go` with the exact retry matrix, `maxRetries + 1` attempt semantics, overflow-safe exponential full jitter, bounded delta/date Retry-After handling, context-cancellable waits, and typed retry exhaustion that preserves the last classification.
- [ ] 4.4 Add `retry_test.go` using injected clock/sleeper/random and safe transports for every retryable/terminal class, delta/date/malformed/past/excessive Retry-After, jitter bounds, maximum cap, successful early stop, exhaustion attempts, cancellation, and limiter participation.

## 5. Automatic pagination and deterministic aggregate

- [ ] 5.1 Refactor `collection.go` so Fetch validates and numerically normalizes requested types, requests each first 50-item page once, requires every total in `0..1_000_000`, exact returned offset/limit and `len(data)<=limit`, checks page-count/offset arithmetic before allocation, fixes the page plan from initial total, uses at most `min(remainingPages, configuredConcurrency)` cancellable workers, and returns no partial data on failure.
- [ ] 5.2 Implement deterministic source coordinates, dedupe by `(SubjectType, SubjectID, Type)` with smallest-coordinate winner, and final sort by `(SubjectType, SubjectID, Type)` while preserving FetchPage’s validated upstream page order and clamp compatibility.
- [ ] 5.3 Add `collection_test.go` for zero/one/multiple pages and types, no limit-1 probe, exact offsets, moving totals, total `1_000_000`, rejected `1_000_001`/near-`MaxInt`, `len(data)=limit+1`, mismatched metadata, checked arithmetic, worker peak, repeated requested types, out-of-order completion, conflicting duplicates, collection-type preservation, terminal-page cancellation, no partial result, empty success, and repeated-run determinism.

## 6. Documentation and test-only CI

- [ ] 6.1 Update `README.md`, package documentation in `client.go`, and `example/example.go` with the complete DTO, anonymous-only scope, exact compatibility/breaking changes, new limiter/retry/endpoint options, stable error use, and no publication claim.
- [ ] 6.2 Add `.github/workflows/ci.yml` with read-only contents permission, Go 1.26.x, exact `actions/checkout@3d3c42e5aac5ba805825da76410c181273ba90b1` (`v7.0.1`) and `actions/setup-go@b7ad1dad31e06c5925ef5d2fc7ad053ef454303e` (`v7.0.0`), no secrets/publication/deploy, and format/module/vet/test/race gates.
- [ ] 6.3 Mechanically verify the workflow has exactly those two `uses:` OIDs, no floating/other action ref, write permission, secret, artifact/release/registry/deploy/SSH step, and that tests contain a guard preventing non-loopback network access.

## 7. Unstaged implementation acceptance

- [ ] 7.1 Run empty `gofmt -d` over all 17 approved Go source/test files, `go mod tidy -diff`, exact module checks, `go vet ./...`, `go test -count=1 ./...`, and `go test -race -count=1 ./...` under Go 1.26.x; prove root tests are non-empty and no required case is skipped.
- [ ] 7.2 Run `openspec validate harden-bangumi-collection-go-v0-1-0 --strict`, `openspec validate --all --strict`, `openspec doctor --json`, `git diff --check`, symlink checks, and exact owned-path complement checks.
- [ ] 7.3 Recheck branch, exact `HEAD=HARDENING_PLANNING_HEAD` and its subject/parent/five-path seal, ancestry, empty index, exact protected hashes/status, unchanged governance and `LICENSE`, no real-network evidence, and a full physical/ignored/untracked status whose only non-protected product changes are the 21 approved paths plus this active change’s task markers.
- [ ] 7.4 Stop with the complete implementation unstaged and uncommitted for main-agent read-only acceptance; do not archive, stage, commit, fetch, push, open a PR, tag, release, publish, deploy, amend, rebase, or reset.

## Post-apply finalization protocol

These are acceptance-gated actions, not pre-archive checkboxes. Only after the main agent accepts the unstaged candidate may a delegated finalization subagent read `openspec-archive-change`, synchronize the capability with the exact approved Purpose, archive to the exact 2026-07-24 path, and stage the exact no-rename 32-path parent delta: 21 product/test/CI paths, five active-change deletions, five archive additions, and one root spec. The same tree must expose the exact 27-path cumulative inventory relative to bootstrap. It must stop for a second main-agent read-only acceptance. Only that acceptance authorizes the exact-subject local commit with sole parent `HARDENING_PLANNING_HEAD`. Post-commit proof is reported without amending or moving another ref. Fetch, push, PR, tag, release, publication, deployment, and main-repository admission remain separately unauthorized.

## Change Boundary

- **Status**
  - investigated: yes.
  - specified: approved with exact planning checkpoint commit pending.
  - implemented: no.
  - verified: specification strict validation and independent semantic review complete; no implementation verification.
  - committed: bootstrap governance only; no hardening planning or implementation commit.
  - pushed: no.
  - released: no.
  - deployed: no.
- **Owner**: delegated External collection client apply subagent; archive/finalization uses a separately delegated subagent; main agent only reviews/amends OpenSpec and performs read-only acceptance.
- **Writable/Owned paths**
  - Active planning: `openspec/changes/harden-bangumi-collection-go-v0-1-0/.openspec.yaml`, `openspec/changes/harden-bangumi-collection-go-v0-1-0/proposal.md`, `openspec/changes/harden-bangumi-collection-go-v0-1-0/design.md`, `openspec/changes/harden-bangumi-collection-go-v0-1-0/specs/collection-client-v0-1-0/spec.md`, `openspec/changes/harden-bangumi-collection-go-v0-1-0/tasks.md`.
  - Product/docs/module: `README.md`, `client.go`, `collection.go`, `errors.go`, `limiter.go`, `options.go`, `request.go`, `retry.go`, `types.go`, `example/example.go`, `go.mod`, `go.sum`.
  - Tests/CI: `client_test.go`, `collection_test.go`, `errors_test.go`, `limiter_test.go`, `options_test.go`, `request_test.go`, `retry_test.go`, `types_test.go`, `.github/workflows/ci.yml`.
  - Final archive/sync: `openspec/changes/archive/2026-07-24-harden-bangumi-collection-go-v0-1-0/.openspec.yaml`, `openspec/changes/archive/2026-07-24-harden-bangumi-collection-go-v0-1-0/proposal.md`, `openspec/changes/archive/2026-07-24-harden-bangumi-collection-go-v0-1-0/design.md`, `openspec/changes/archive/2026-07-24-harden-bangumi-collection-go-v0-1-0/specs/collection-client-v0-1-0/spec.md`, `openspec/changes/archive/2026-07-24-harden-bangumi-collection-go-v0-1-0/tasks.md`, `openspec/specs/collection-client-v0-1-0/spec.md`.
- **Read-only protected inputs**
  - `CLAUDE.md` SHA-256 `c69937414d643411e03f3d7982e4bf58e8ceb02802cf1338dada4ea6a630407d`, untracked.
  - `note` SHA-256 `7b2376f605f5d9d9f677eaf1165a2ff911cec9092b6b2763b966a3d6f93e2e74`, untracked.
  - `.claude/settings.local.json` SHA-256 `c67845ef6924b080344bc0010931336106d9df3a700f67f10a22763d07517e3e`, ignored.
  - `LICENSE`; `.codex/skills/openspec-apply-change/SKILL.md`; `.codex/skills/openspec-archive-change/SKILL.md`; `.codex/skills/openspec-explore/SKILL.md`; `.codex/skills/openspec-propose/SKILL.md`; `.codex/skills/openspec-sync-specs/SKILL.md`; `.codex/skills/openspec-update-change/SKILL.md`; `openspec/config.yaml`; `openspec/changes/archive/2026-07-24-bootstrap-bangumi-collection-go-openspec/.openspec.yaml`; `openspec/changes/archive/2026-07-24-bootstrap-bangumi-collection-go-openspec/proposal.md`; `openspec/changes/archive/2026-07-24-bootstrap-bangumi-collection-go-openspec/design.md`; `openspec/changes/archive/2026-07-24-bootstrap-bangumi-collection-go-openspec/specs/external-openspec-bootstrap/spec.md`; `openspec/changes/archive/2026-07-24-bootstrap-bangumi-collection-go-openspec/tasks.md`; `openspec/specs/external-openspec-bootstrap/spec.md`.
  - `/Users/luca/dev/BangumiStaffStats/tmp-formal-development/formal-development-master-plan.md`; `/Users/luca/dev/BangumiStaffStats/tmp-formal-development/backend-development-implementation-guide.md`.
  - Initial hardening `HEAD=59eba79b3e3621ce72b756a88b38aee970c00fcf`; baseline ancestor `8173f44911360150a5a5a7c6418021d1014fe85b`.
- **Consumes**: exact change `bootstrap-bangumi-collection-go-openspec`; root capability `external-openspec-bootstrap`; ten baseline product files; master-plan external lane; backend guide 6.2–6.3.
- **Produces**: exact five-file planning checkpoint; capability `collection-client-v0-1-0`; exact 21 product/test/CI paths; five archived change files; synchronized root spec; exact local implementation commit.
- **Dependencies**: `bootstrap-bangumi-collection-go-openspec`, satisfied by `59eba79b3e3621ce72b756a88b38aee970c00fcf`. No main-repository change, wave alias, or grouped dependency.
- **Deliverables**: committed approved plan; complete DTO; bounded automatic pagination; client-wide QPS/concurrency; bounded retry/Retry-After; typed sanitized errors; deterministic dedupe/sort; Go 1.26 httptest/race/vet/module/CI; compatibility docs; implementation/finalization acceptances; local implementation commit.
- **Acceptance**

  Planning candidate exact status:

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

  Planning commands:

  ```sh
  test "$(git branch --show-current)" = codex/v0.1.0-hardening
  test "$(git rev-parse HEAD)" = 59eba79b3e3621ce72b756a88b38aee970c00fcf
  git merge-base --is-ancestor 8173f44911360150a5a5a7c6418021d1014fe85b HEAD
  test -z "$(git diff --cached --name-only)"
  test "$(shasum -a 256 CLAUDE.md note .claude/settings.local.json | shasum -a 256 | awk '{print $1}')" = 44239e07874c868038a70e908e8c5fb42c4a48fadf8fa6ab6bdc6169cd9e3180
  test "$(git status --porcelain=v1 --ignored --untracked-files=all -- CLAUDE.md note .claude/settings.local.json | shasum -a 256 | awk '{print $1}')" = 221322afd4137d24875b5fbd369c9028e99f69f1cf155e3c49975996841a3938
  test "$(printf '%s\0' LICENSE README.md client.go collection.go errors.go example/example.go go.mod go.sum options.go types.go | xargs -0 shasum -a 256 | shasum -a 256 | awk '{print $1}')" = 7d4e55231e08c5d8286bffb921e53adbda9a5d82ced9c4688e4cadb129b1855e
  test "$(git status --porcelain=v1 --ignored --untracked-files=all | LC_ALL=C sort | shasum -a 256 | awk '{print $1}')" = e91f861f8e2d4e4398f6ca39ff39a257d71f3cf0fcf0a7382de82ca497796ef5
  test -z "$(find openspec/changes/harden-bangumi-collection-go-v0-1-0 -type l -print)"
  openspec status --change harden-bangumi-collection-go-v0-1-0
  openspec validate harden-bangumi-collection-go-v0-1-0 --strict
  openspec validate --all --strict
	  openspec doctor --json
	  git diff --check
	  ```

	  Planning-checkpoint staging and post-commit proof:

	  ```sh
	  test "$(git diff --cached --name-only --no-renames | wc -l | tr -d ' ')" = 5
	  test "$(git diff --cached --name-only --no-renames | LC_ALL=C sort | shasum -a 256 | awk '{print $1}')" = a7da0916df9b59cdaff4d5a3b57d8b0231c8bad51921aedb308ee2d10262fe1f
	  git diff --cached --check
	  test "$(git show -s --format=%P HEAD)" = 59eba79b3e3621ce72b756a88b38aee970c00fcf
	  test "$(git show -s --format=%s HEAD)" = "docs(openspec): approve collection client hardening"
	  test "$(git diff-tree --no-commit-id --name-only --no-renames -r HEAD | wc -l | tr -d ' ')" = 5
	  test "$(git diff-tree --no-commit-id --name-only --no-renames -r HEAD | LC_ALL=C sort | shasum -a 256 | awk '{print $1}')" = a7da0916df9b59cdaff4d5a3b57d8b0231c8bad51921aedb308ee2d10262fe1f
	  ```

	  Apply commands:

	  ```sh
	  test "$(git rev-parse HEAD)" = "$HARDENING_PLANNING_HEAD"
	  test "$(git show -s --format=%P HEAD)" = 59eba79b3e3621ce72b756a88b38aee970c00fcf
	  test "$(git show -s --format=%s HEAD)" = "docs(openspec): approve collection client hardening"
	  test "$(go env GOVERSION | sed -E 's/^go([0-9]+\.[0-9]+).*/\1/')" = 1.26
  test -z "$(gofmt -d -e client.go collection.go errors.go limiter.go options.go request.go retry.go types.go example/example.go client_test.go collection_test.go errors_test.go limiter_test.go options_test.go request_test.go retry_test.go types_test.go)"
  go mod tidy -diff
  test "$(go list -m -f '{{.Version}}' golang.org/x/sync)" = v0.22.0
  test "$(go list -m -f '{{.Version}}' golang.org/x/time)" = v0.15.0
  go vet ./...
  go test -count=1 ./...
  go test -race -count=1 ./...
  openspec validate harden-bangumi-collection-go-v0-1-0 --strict
  openspec validate --all --strict
  openspec doctor --json
  git diff --check
  test -z "$(git diff --cached --name-only)"
  ```

	  Finalization commands before the second acceptance:

	  ```sh
	  test "$(git rev-parse HEAD)" = "$HARDENING_PLANNING_HEAD"
	  test "$(git diff --cached --name-only --no-renames | wc -l | tr -d ' ')" = 32
	  test "$(git diff --cached --name-only --no-renames | LC_ALL=C sort | shasum -a 256 | awk '{print $1}')" = 740610fc014608d9294383d731e3a9df51b3521c25823824c6ef89dc582349b2
	  test "$(git diff --cached --name-only --no-renames 59eba79b3e3621ce72b756a88b38aee970c00fcf | wc -l | tr -d ' ')" = 27
	  test "$(git diff --cached --name-only --no-renames 59eba79b3e3621ce72b756a88b38aee970c00fcf | LC_ALL=C sort | shasum -a 256 | awk '{print $1}')" = b8ef7390c49b73f13061c3ca50f66a0ef260cf947c539d507ab1797de59c48c1
	  git diff --cached --check
  test "$(awk 'p && NF {print; exit} /^## Purpose$/{p=1}' openspec/specs/collection-client-v0-1-0/spec.md)" = "Define the first public anonymous Bangumi collection client contract, including complete DTOs, automatic pagination, shared limiting, bounded retries, sanitized errors, deterministic results, and local quality gates."
  test "$(awk '/^### Requirement:/{p=1} /^## Boundary Record$/{p=0} p' openspec/changes/archive/2026-07-24-harden-bangumi-collection-go-v0-1-0/specs/collection-client-v0-1-0/spec.md | shasum -a 256 | awk '{print $1}')" = "$(awk '/^### Requirement:/{p=1} /^## Boundary Record$/{p=0} p' openspec/specs/collection-client-v0-1-0/spec.md | shasum -a 256 | awk '{print $1}')"
  openspec validate --all --strict
  openspec doctor --json
  ```

	  Post-commit read-only proof:

	  ```sh
	  test "$(git show -s --format=%P HEAD)" = "$HARDENING_PLANNING_HEAD"
	  test "$(git show -s --format=%s HEAD)" = "feat: harden collection client for v0.1.0"
	  test "$(git diff-tree --no-commit-id --name-only --no-renames -r HEAD | wc -l | tr -d ' ')" = 32
	  test "$(git diff-tree --no-commit-id --name-only --no-renames -r HEAD | LC_ALL=C sort | shasum -a 256 | awk '{print $1}')" = 740610fc014608d9294383d731e3a9df51b3521c25823824c6ef89dc582349b2
	  test "$(git diff --name-only --no-renames 59eba79b3e3621ce72b756a88b38aee970c00fcf..HEAD | wc -l | tr -d ' ')" = 27
	  test "$(git diff --name-only --no-renames 59eba79b3e3621ce72b756a88b38aee970c00fcf..HEAD | LC_ALL=C sort | shasum -a 256 | awk '{print $1}')" = b8ef7390c49b73f13061c3ca50f66a0ef260cf947c539d507ab1797de59c48c1
  test "$(git rev-parse refs/heads/main)" = 8173f44911360150a5a5a7c6418021d1014fe85b
  test "$(git rev-parse refs/remotes/origin/main)" = 8173f44911360150a5a5a7c6418021d1014fe85b
  test -z "$(git tag --list)"
  test "$(git status --porcelain=v1 --ignored --untracked-files=all | shasum -a 256 | awk '{print $1}')" = 221322afd4137d24875b5fbd369c9028e99f69f1cf155e3c49975996841a3938
  ```
- **Non-goals**: token/OAuth/private/write/general-SDK behavior; main-repository cache/domain/statistics; real-network load tests; fetch/push/PR/tag/release/publication/deployment.
- **Operations deferred**: remote Git/GitHub, version tag/release/publication, consumer admission, production values, observability, rollout, deployment.
- **Stop/rollback conditions**: on any root/branch/ref/ancestry/index/path/hash/status/symlink/toolchain/dependency/format/static/test/race/strict/doctor mismatch, stop without clean/reset/stash/delete/overwrite/stage/commit and preserve evidence. Rollback/reversal requires a new reviewed change.
- **Mutable refs**: after spec approval, only `refs/heads/codex/v0.1.0-hardening` may advance once from bootstrap to the exact five-file planning commit. It remains fixed through apply and both implementation/finalization acceptances, then may advance once more to the exact implementation commit. `refs/heads/main`, `refs/remotes/origin/main`, tags, remote refs, and remote state remain immutable.
