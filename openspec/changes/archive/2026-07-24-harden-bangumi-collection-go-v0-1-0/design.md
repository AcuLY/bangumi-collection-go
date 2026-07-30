## Context

The repository currently has one small package with an untagged API: `NewClient`, `Fetch`, `FetchPage`, functional options, flattened `Subject` values, an `errgroup` per `Fetch`, exponential retry, and errors that can print upstream bodies. `Fetch` probes with `limit=1`, refetches offset zero, and appends pages in goroutine completion order. Its concurrency limit is scoped to one call, it has no QPS gate, it drops comment/update/type fields, `WithAccessToken` emits credentials, and the package has no tests.

The formal backend will only consume anonymous public collections. It requires the exact collection fields named in backend guide section 6.3, stable output, classified errors, shared request limits, context propagation, and an independently published fixed `v0.1.0`. The OpenSpec bootstrap is already represented by commit `59eba79b3e3621ce72b756a88b38aee970c00fcf`; this design does not authorize publication.

## Goals / Non-Goals

**Goals:**

- Define a complete, race-safe, deterministic anonymous collection client contract suitable for the future `BangumiCollectionSource` adapter.
- Bound every attempt by client-wide QPS, client-wide in-flight concurrency, per-attempt timeout, retry count, retry delay, success-body size, and error-body size.
- Preserve the useful untagged call surface where compatible with the first explicit public contract, and document every intentional incompatibility before `v0.1.0`.
- Make all behavioral gates reproducible with local `httptest` servers, Go 1.26, race testing, standard tooling, and a test-only CI workflow.
- Preserve the exact repository governance, protected user paths, local-only commit gate, and downstream publication gate.

**Non-Goals:**

- Authenticated/private collections, OAuth/token storage, write APIs, caller header forwarding, or a general Bangumi SDK.
- Backend cache/stale/digest/domain/statistics behavior or any main-repository change.
- Real-upstream load tests, production rate tuning, monitoring, deployment, fetch, push, PR, tag, release, or publication.

## Decisions

### 1. Keep one immutable client configuration and split narrow concerns

`client.go` owns immutable configuration and shared state; `collection.go` owns public collection operations and page scheduling; `request.go` owns URL/request/response conversion; `limiter.go` owns shared QPS and in-flight gates; `retry.go` owns retry classification and delay; `errors.go` owns the public error taxonomy; `types.go` owns public and wire DTOs. The `Client` is safe for concurrent use after construction. Tests in the package may inject unexported clock/sleeper/random hooks before use; production callers cannot mutate them.

Alternative considered: keep all behavior in `collection.go`. Rejected because retry, limiting, HTTP decoding, and pagination have independent invariants and would remain difficult to race-test.

### 2. Preserve the call shape, establish a complete `Subject` collection DTO

`NewClient(userAgent string, options ...Option) *Client`, `Fetch`, and `FetchPage` retain their signatures. Existing `Subject.ID`, `Name`, `NameCn`, `Rate`, `VolStatus`, `EpStatus`, `Tags`, and `Private` remain. `Subject` adds:

- `SubjectID int` as the canonical top-level collection `subject_id`; `ID` is the compatibility alias and MUST equal it.
- `SubjectType SubjectType`.
- `Type CollectionType` for the collection state.
- `Comment string`.
- `UpdatedAt time.Time`.

Before the first public tag, the named subject constants are corrected to the
official OAS mapping: Book `1`, Anime `2`, Music `3`, Game `4`, Real `6`. The
untagged prototype had the `Game` and `Music` names reversed. Raw numeric
values and the valid set do not change, but callers using those two names
receive the corrected semantic request; the README lists this intentional
pre-`v0.1.0` break.

The wire decoder follows the current official Bangumi OAS checked on
2026-07-24. It requires the top-level collection fields `subject_id`,
`subject_type`, `rate`, `type`, `tags`, `ep_status`, `vol_status`,
`updated_at`, and `private`. `comment` and the nested `subject` projection are
optional upstream: an absent comment maps to the public zero-value string, and
an absent subject keeps `ID == SubjectID` while `Name`/`NameCn` remain empty.
When the nested subject is present its required identity/name tuple must be
complete and agree with the top level. Tags are required, copied, and returned
as a non-nil empty slice when the required array is empty. Unknown additive
upstream fields are ignored; missing or malformed required fields, a malformed
present optional field, identity disagreement, invalid page metadata, trailing
JSON, and a non-JSON success response are protocol/decode errors.

Alternative considered: introduce a new `Collection` return type. Rejected for `v0.1.0` because extending the existing untagged `Subject` preserves current callers without sacrificing required fields.

### 3. Make the first public version anonymous-only

`WithAccessToken`, the `accessToken` field, and all Authorization behavior are removed. Requests set only the library-owned `User-Agent` and `Accept: application/json`; they never set Cookie or Authorization and never copy arbitrary caller request headers. The client shallow-clones an injected `http.Client`, clears its cookie jar and `Timeout`, and replaces redirect handling with first-response refusal so credentials or identity cannot cross an origin. `WithRequestTimeout` therefore governs every individual attempt regardless of option order or custom-client use; the caller-owned client remains unchanged. A custom `RoundTripper` remains a trusted caller dependency, but request-construction tests prove the library supplies no credential headers.

The default endpoint is `https://api.bgm.tv`. `WithEndpoint` accepts an absolute HTTPS endpoint or loopback HTTP endpoint for `httptest`, requires a non-empty host and an empty or `/` path, rejects opaque URLs, user-info, query, fragment, every other path, and every other insecure endpoint, and stores a parsed immutable root URL. UID normalization trims surrounding Unicode whitespace, then requires valid UTF-8, `1..256` bytes, and no Unicode control rune while preserving case; the normalized UID is escaped as exactly one opaque path segment. Dot-only values `.` and `..` remain valid UIDs but their dots are percent-encoded so a client, proxy, or router cannot normalize them as path traversal. The configured User-Agent must be valid UTF-8, `1..256` bytes, non-blank after a validation-only Unicode trim, and free of Unicode control runes; the accepted value is sent byte-for-byte as supplied.

Alternative considered: retain `WithAccessToken` as a no-op. Rejected because silently accepting a credential option would be misleading and would obscure the first version’s security boundary.

### 4. Use one client-wide limiter pipeline for every attempt

Every initial request and retry follows the same order:

1. fail immediately if the parent context is done;
2. wait for the shared `golang.org/x/time/rate.Limiter`;
3. acquire the shared in-flight semaphore;
4. create a fresh per-attempt timeout child context;
5. execute one HTTP round trip and complete the bounded body read, validation, and close under that child;
6. cancel the child and release the semaphore exactly once after those steps or after transport failure.

The rate default is 3 requests/second with burst 1; the existing concurrency default remains 10. New APIs are exactly `WithEndpoint(string) Option`, `WithRateLimit(float64, int) Option`, and `WithMaxRetryDelay(time.Duration) Option`. Existing invalid non-auth option values retain their default-preserving behavior; an invalid User-Agent/new option value, a non-finite/non-positive QPS, non-positive burst, or nil `Option` records the first `ErrInvalidConfiguration`, cannot be cleared by a later valid Option, and fails operations before transport without panic. All concurrent `Fetch`/`FetchPage` calls on the same client share both limiter objects; distinct clients do not. Waiting never holds an in-flight slot, and every wait/acquire is context-cancellable. If the rate limiter rejects a reservation because it would exceed the parent deadline before `ctx.Err()` becomes non-nil, that result is still the terminal parent-deadline classification and is never retried or wrapped as retry exhaustion. A parent cancellation/deadline always wins classification over an attempt timeout; a retry creates a new child deadline from the still-live parent.

`golang.org/x/time/rate v0.15.0` is admitted as a production dependency because the standard library has no concurrency-safe token-bucket with cancellable reservations. A home-grown ticker/bucket was rejected due to fairness, burst, cancellation, and race risk. Cost: one official `golang.org/x` module (BSD-3-Clause), no transitive runtime dependency, small maintenance/supply-chain surface, and no public type leakage. Removal gate: replace only in a later reviewed change if the standard library gains an equivalent and all limiter/race scenarios remain byte-for-byte behaviorally equivalent.

### 5. Page with bounded workers, then canonicalize deterministically

`Fetch` validates a non-empty collection-type list, rejects invalid enum values, removes repeats, and sorts the normalized types numerically. For each type it requests offset 0 at limit 50; the returned initial `total` fixes that type’s logical page plan. Every decoded page requires `total` in `0..1_000_000` per collection type, exact normalized requested offset/limit metadata, and `len(data) <= limit`; a first page additionally requires `len(data) <= total`, so `total=0` cannot carry records. A violation is a protocol error before allocation or scheduling. Page-count and offset arithmetic is checked for overflow. Remaining offsets `50, 100, ... < total` are handled through bounded `errgroup` workers whose count never exceeds `min(remainingPages, configuredConcurrency)`; no unbounded goroutine is created. A terminal error cancels sibling work and no background worker/timer survives return.

Every collected item carries a deterministic source coordinate `(collectionType, offset, itemIndex)`. After all pages succeed, `Fetch` deduplicates by `(SubjectType, SubjectID, Type)`, selecting the smallest source coordinate on collision, and sorts by `(SubjectType ascending, SubjectID ascending, Type ascending)`. Thus completion timing and repeated requested types cannot change output. `FetchPage` remains a page primitive: it clamps limit to 1–50 and offset to at least zero for compatibility, validates/converts the page, and preserves that page’s upstream item order.

Alternative considered: append under a mutex as pages finish. Rejected because that makes public ordering nondeterministic and race-sensitive.

### 6. Use a closed retry matrix with bounded full jitter

`WithMaxRetries(n)` continues to mean retries after the first attempt; `n=0` disables retry. `WithRetryInterval` is the base delay, and new `WithMaxRetryDelay` bounds every retry wait (defaults: 1 second base, 30 seconds maximum).

Retryable classes are exactly:

- a transport failure while the parent context is still live, including a per-attempt timeout;
- HTTP 429;
- HTTP 500–599.

Caller cancellation/deadline, input/config error, 401/403/404/other 4xx, decode/protocol/oversize errors, and limiter/acquire cancellation are terminal. For attempt number `n`, local delay is full jitter in `[0, min(base*2^(n-1), maxDelay)]`. A syntactically valid non-negative `Retry-After` delta-seconds or HTTP-date on 429/5xx contributes `min(value, maxDelay)`; an all-digit delta too large for native integer parsing is valid and saturates directly to the configured cap. The actual wait is `max(localJitter, boundedRetryAfter)`. Invalid/past values are ignored. All arithmetic is overflow-safe. Retry exhaustion returns a typed wrapper with total attempts and the last sanitized typed error, preserving classification via `errors.Is/As`.

Alternative considered: always sleep the server value without a cap. Rejected because a context without a deadline could block a library call for an attacker-controlled/unbounded interval.

### 7. Sanitize errors while retaining compatibility classification

The implementation retains existing sentinels and public `NetworkError`/`HTTPError` names, adds stable sentinels/codes for nil context, no collection types, invalid enums/config, not found, generic HTTP status, timeout, canceled, transport, decode, protocol, oversized response, and retry exhaustion, and adds typed decode/protocol/retry errors. `errors.Is` classifies stable conditions; `errors.As` exposes typed metadata such as status, bounded `RetryAfter`, timeout flag, and attempt count.

For compatibility, `HTTPError.Body` remains present but is always the empty string; `NetworkError.Err` contains only a safe sentinel or `context.Canceled`/`context.DeadlineExceeded`, never a raw `url.Error` or request URL. Every `Error()` is bounded and excludes UID, URL, query, response body, Authorization/Cookie, user agent, Location, and arbitrary transport text. Only HTTP 200 is success. Every other final 1xx/2xx/3xx/4xx/5xx status returns `*HTTPError` matching `ErrHTTPStatus`; 401/403/404/429/5xx additionally match their existing stable category, and 404 also matches deprecated `ErrInvalidUserID`. Redirects are not followed and the first 3xx response is terminal. Nil context returns `ErrNilContext` directly before transport. Parent cancellation returns `*NetworkError{Timeout:false}` matching `ErrCanceled` and `context.Canceled`; parent deadline and a per-attempt timeout return `*NetworkError{Timeout:true}` matching `ErrTimeout` and `context.DeadlineExceeded`, but only the attempt timeout is retryable while the parent remains live. Other transport failure returns `*NetworkError` matching `ErrTransport`. Retry exhaustion unwraps the last exact sanitized typed classification.

Success bodies are limited to 16 MiB plus one detection byte; oversize/read/decode/trailing success failures are distinct and non-retryable. At most 64 KiB plus one detection byte of an error body is read before discard; its size never replaces the known HTTP status classification or suppresses a permitted 429/5xx retry. Error bodies are never stored or rendered.

Alternative considered: expose truncated response bodies for debugging. Rejected because upstream content can contain user-derived material and downstream logs cannot safely distinguish it.

### 8. Keep dependency and quality tooling narrow

Runtime modules are exactly `golang.org/x/sync v0.22.0` for structured cancellation/bounded workers and `golang.org/x/time v0.15.0` for the shared limiter. Both are official Go subrepositories under BSD-3-Clause and require no runtime service. Go source declares Go 1.26 in `go.mod`.

Tests use only `testing`, `httptest`, and standard-library helpers; no assertion/mock framework is admitted because table tests and custom round trippers cover the small public surface without another API, transitive graph, or maintenance owner. CI uses only GitHub-owned `actions/checkout@3d3c42e5aac5ba805825da76410c181273ba90b1` (`v7.0.1`) and `actions/setup-go@b7ad1dad31e06c5925ef5d2fc7ad053ef454303e` (`v7.0.0`), discovered from their official repository tag refs and pinned to immutable commits. It declares `permissions: contents: read`, uses no secrets, publishes no artifact, and runs only format/module/vet/test/race jobs. Removal gates are explicit: either runtime module or CI action can be removed only in a reviewed change after its owned behavior has a standard-library/local replacement and all quality gates remain green.

### 9. Treat `v0.1.0` publication as a later compatibility gate

This change creates a local implementation commit, not the tag. README and package docs list the preserved call surface and the intentional removals: credential option/behavior is gone; HTTP error body is empty; output order is canonical rather than completion order; `Fetch` rejects an empty type list; DTOs are extended. The future public tag may only point to the accepted local commit or a separately reviewed descendant. The main repository may admit only a separately authorized, publicly available fixed `v0.1.0`, with no `replace`.

## Risks / Trade-offs

- **[Risk] Collection contents can mutate between concurrent page requests** → Fix the page plan from the first response, attach deterministic source coordinates, deduplicate, and document that the upstream API offers no snapshot transaction. Never invent missing items or loop indefinitely on a moving `total`.
- **[Risk] Shared QPS creates head-of-line delay across callers** → Make waits cancellable, avoid holding concurrency while rate-waiting, expose explicit pre-construction configuration, and test two simultaneous public operations.
- **[Risk] Retry can amplify upstream load** → Use a closed retry matrix, client-wide limiter, small bounded count, full jitter, and a maximum wait.
- **[Risk] Removing authentication breaks untagged callers** → Mark it breaking before the first public tag, retain all compatible call shapes, and document anonymous-only scope prominently.
- **[Risk] Sanitization loses debugging detail** → Preserve stable code/status/retry/attempt metadata without storing body, URL, UID, headers, or raw transport strings.
- **[Risk] Go 1.26 is newer than the currently installed 1.25.4 toolchain** → Apply must obtain/use a Go 1.26.x toolchain outside repository product paths, record `go env GOVERSION`, and stop rather than weakening `go.mod` or CI.
- **[Risk] CI action/module supply chain expands** → Pin exact modules, SHA-pin only two GitHub-owned actions, keep read-only permissions, and enforce no publish/deploy step.

## Migration Plan

1. Main agent approves all artifacts while `HEAD` remains `59eba79b3e3621ce72b756a88b38aee970c00fcf` and the index is empty.
2. A delegated planning-checkpoint subagent stages exactly the five active-change files, proves their sorted path seal `a7da0916df9b59cdaff4d5a3b57d8b0231c8bad51921aedb308ee2d10262fe1f`, and creates subject `docs(openspec): approve collection client hardening` with sole parent `59eba79b3e3621ce72b756a88b38aee970c00fcf`. Main-agent read-only acceptance records that exact commit as `HARDENING_PLANNING_HEAD`.
3. A delegated apply subagent, pinned to `HARDENING_PLANNING_HEAD`, edits only the 21 approved product/test/CI paths, checks task boxes, runs every local gate, and stops with an unstaged candidate.
4. Main agent performs read-only behavior/diff/seal acceptance. Failure returns to a new OpenSpec amendment; no product fix is made by the main agent.
5. A delegated archive/finalization subagent reads the archive skill, synchronizes `collection-client-v0-1-0`, moves the active change to the exact 2026-07-24 archive path, and stages a no-rename 32-path parent delta: 21 product/test/CI paths, five active-change deletions, five archive additions, and one root spec. The same tree has a 27-path cumulative inventory relative to the bootstrap checkpoint.
6. Main agent performs a second read-only index/tree/seal acceptance.
7. Only then may finalization create one local commit with sole parent `HARDENING_PLANNING_HEAD` and subject `feat: harden collection client for v0.1.0`; it then reruns read-only proof and stops.
8. Fetch, push, PR, tag, release, publication, main-repository admission, and deployment await separate authorization.

Rollback is not an in-place reset. Before commit, preserve the unstaged or staged candidate and report the failed gate. After the accepted local commit, any reversal requires a new reviewed change and an additive/revert commit; amend, rebase, reset, or checkout restoration is not authorized.

## Open Questions

None. Default QPS/burst, retry bounds, DTO compatibility, dedupe key/order, error taxonomy, dependency versions, owned paths, and local/publication gates are fixed by this change. Future production tuning or authenticated scope requires a new OpenSpec change.

## Change Boundary

- **Status**: investigated=yes; specified=approved and committed; implemented=yes; verified=independent code review plus format, module, vet, build, full-test, race, randomized-order, coverage, CI, strict OpenSpec, and boundary gates passed; committed=yes at local `63036ca990eaade8cbf691daf8b1db31ca39ac83`; pushed=no; released=no; deployed=no.
- **Owner**: delegated External collection client apply subagent; main agent only reviews/amends OpenSpec and performs read-only acceptance.
- **Writable/Owned paths**
  - Planning: `openspec/changes/harden-bangumi-collection-go-v0-1-0/.openspec.yaml`, `proposal.md`, `design.md`, `specs/collection-client-v0-1-0/spec.md`, and `tasks.md` below that exact active root.
  - Product/docs/module: `README.md`, `client.go`, `collection.go`, `errors.go`, `limiter.go`, `options.go`, `request.go`, `retry.go`, `types.go`, `example/example.go`, `go.mod`, `go.sum`.
  - Tests/CI: `client_test.go`, `collection_test.go`, `errors_test.go`, `limiter_test.go`, `options_test.go`, `request_test.go`, `retry_test.go`, `types_test.go`, `.github/workflows/ci.yml`.
  - Final archive/sync: `openspec/changes/archive/2026-07-24-harden-bangumi-collection-go-v0-1-0/.openspec.yaml`, `openspec/changes/archive/2026-07-24-harden-bangumi-collection-go-v0-1-0/proposal.md`, `openspec/changes/archive/2026-07-24-harden-bangumi-collection-go-v0-1-0/design.md`, `openspec/changes/archive/2026-07-24-harden-bangumi-collection-go-v0-1-0/specs/collection-client-v0-1-0/spec.md`, `openspec/changes/archive/2026-07-24-harden-bangumi-collection-go-v0-1-0/tasks.md`, `openspec/specs/collection-client-v0-1-0/spec.md`.
- **Read-only protected inputs**: untracked `CLAUDE.md` SHA-256 `c69937414d643411e03f3d7982e4bf58e8ceb02802cf1338dada4ea6a630407d`; untracked `note` SHA-256 `7b2376f605f5d9d9f677eaf1165a2ff911cec9092b6b2763b966a3d6f93e2e74`; ignored `.claude/settings.local.json` SHA-256 `c67845ef6924b080344bc0010931336106d9df3a700f67f10a22763d07517e3e`; `LICENSE`; `.codex/skills/openspec-apply-change/SKILL.md`; `.codex/skills/openspec-archive-change/SKILL.md`; `.codex/skills/openspec-explore/SKILL.md`; `.codex/skills/openspec-propose/SKILL.md`; `.codex/skills/openspec-sync-specs/SKILL.md`; `.codex/skills/openspec-update-change/SKILL.md`; `openspec/config.yaml`; `openspec/specs/external-openspec-bootstrap/spec.md`; `openspec/changes/archive/2026-07-24-bootstrap-bangumi-collection-go-openspec/.openspec.yaml`; `openspec/changes/archive/2026-07-24-bootstrap-bangumi-collection-go-openspec/proposal.md`; `openspec/changes/archive/2026-07-24-bootstrap-bangumi-collection-go-openspec/design.md`; `openspec/changes/archive/2026-07-24-bootstrap-bangumi-collection-go-openspec/specs/external-openspec-bootstrap/spec.md`; `openspec/changes/archive/2026-07-24-bootstrap-bangumi-collection-go-openspec/tasks.md`; `/Users/luca/dev/BangumiStaffStats/tmp-formal-development/formal-development-master-plan.md`; `/Users/luca/dev/BangumiStaffStats/tmp-formal-development/backend-development-implementation-guide.md`; initial checkpoint `59eba79b3e3621ce72b756a88b38aee970c00fcf`; baseline ancestor `8173f44911360150a5a5a7c6418021d1014fe85b`.
- **Consumes**: exact change `bootstrap-bangumi-collection-go-openspec`, capability `external-openspec-bootstrap`, ten baseline product files, master-plan external lane, backend guide sections 6.2–6.3.
- **Produces**: one accepted five-file planning checkpoint; capability `collection-client-v0-1-0`; 21 exact product/test/CI paths; synchronized root spec; archived change; and one accepted local implementation commit.
- **Dependencies**: `bootstrap-bangumi-collection-go-openspec`, satisfied by `59eba79b3e3621ce72b756a88b38aee970c00fcf`; no grouped/wave alias and no main-repository apply dependency.
- **Deliverables**: complete DTO; automatic paging; client-wide QPS/concurrency; retry/Retry-After; typed sanitized errors; deterministic dedupe/sort; Go 1.26 tests/race/vet/module/CI; compatibility docs; local commit proof.
- **Acceptance**

  The sorted planning-candidate full status SHALL be exactly:

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

  ```sh
  test "$(git branch --show-current)" = codex/v0.1.0-hardening
  test "$(git rev-parse HEAD)" = 59eba79b3e3621ce72b756a88b38aee970c00fcf
  git merge-base --is-ancestor 8173f44911360150a5a5a7c6418021d1014fe85b HEAD
  test -z "$(git diff --cached --name-only)"
  test "$(shasum -a 256 CLAUDE.md note .claude/settings.local.json | shasum -a 256 | awk '{print $1}')" = 44239e07874c868038a70e908e8c5fb42c4a48fadf8fa6ab6bdc6169cd9e3180
  test "$(git status --porcelain=v1 --ignored --untracked-files=all -- CLAUDE.md note .claude/settings.local.json | shasum -a 256 | awk '{print $1}')" = 221322afd4137d24875b5fbd369c9028e99f69f1cf155e3c49975996841a3938
  test "$(printf '%s\0' LICENSE README.md client.go collection.go errors.go example/example.go go.mod go.sum options.go types.go | xargs -0 shasum -a 256 | shasum -a 256 | awk '{print $1}')" = 7d4e55231e08c5d8286bffb921e53adbda9a5d82ced9c4688e4cadb129b1855e
  test "$(git status --porcelain=v1 --ignored --untracked-files=all | LC_ALL=C sort | shasum -a 256 | awk '{print $1}')" = e91f861f8e2d4e4398f6ca39ff39a257d71f3cf0fcf0a7382de82ca497796ef5
  openspec validate harden-bangumi-collection-go-v0-1-0 --strict
  openspec validate --all --strict
  openspec doctor --json
  git diff --check
  test -z "$(find openspec/changes/harden-bangumi-collection-go-v0-1-0 -type l -print)"
  ```

	  Apply additionally requires:

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
  ```

  Planning checkpoint acceptance requires exactly five staged active-change paths, sorted path seal `a7da0916df9b59cdaff4d5a3b57d8b0231c8bad51921aedb308ee2d10262fe1f`, sole parent `59eba79b3e3621ce72b756a88b38aee970c00fcf`, and exact subject `docs(openspec): approve collection client hardening`. Finalization additionally requires a `--no-renames` 32-path parent delta with sorted path seal `740610fc014608d9294383d731e3a9df51b3521c25823824c6ef89dc582349b2`, a cumulative 27-path inventory relative to `59eba79b3e3621ce72b756a88b38aee970c00fcf` with seal `b8ef7390c49b73f13061c3ca50f66a0ef260cf947c539d507ab1797de59c48c1`, exact root Purpose, archived/root requirement-body hash equality, cached diff/strict/all/doctor/Go gates, unchanged protection seals, and the two main-agent implementation/finalization acceptances described above.
- **Non-goals**: credentials/OAuth/private/write APIs; main-repository behavior; generic SDK scope; real upstream/load/production work; fetch/push/PR/tag/release/publication/deployment.
- **Operations deferred**: public version/tag/release, GitHub/remote mutation, consumer admission, monitoring, rollout, production tuning, deployment.
- **Stop/rollback conditions**: stop on any branch/ancestry/index/path/hash/status/toolchain/dependency/format/static/test/race/strict/doctor mismatch without clean/reset/stash/delete/overwrite/stage/commit; preserve evidence. Rollback requires a new reviewed change.
- **Mutable refs**: after spec approval, only `refs/heads/codex/v0.1.0-hardening` may advance from the bootstrap checkpoint once to the exact five-file planning commit. It then remains fixed through apply and both candidate acceptances, after which it may advance once more to the exact implementation commit. Main, origin/main, tags, remote refs, and remote state remain fixed.
