# collection-client-v0-1-0 Specification

## Purpose
Define the first public anonymous Bangumi collection client contract, including complete DTOs, automatic pagination, shared limiting, bounded retries, sanitized errors, deterministic results, and local quality gates.
## Requirements
### Requirement: Hardening starts from accepted external governance and a committed plan
The capability SHALL operate only in canonical repository `/Users/luca/dev/bangumi-collection-go`, on branch `codex/v0.1.0-hardening`, from exact initial hardening checkpoint `59eba79b3e3621ce72b756a88b38aee970c00fcf`. Commit `8173f44911360150a5a5a7c6418021d1014fe85b` SHALL remain an ordinary ancestor. Exact dependency change `bootstrap-bangumi-collection-go-openspec` and root capability `external-openspec-bootstrap` SHALL exist and remain unchanged before apply begins.

After main-agent spec approval and before product apply, one delegated planning-checkpoint subagent SHALL stage exactly the five active-change files, require sorted path seal `a7da0916df9b59cdaff4d5a3b57d8b0231c8bad51921aedb308ee2d10262fe1f`, and create exact subject `docs(openspec): approve collection client hardening` with sole parent `59eba79b3e3621ce72b756a88b38aee970c00fcf`. Main-agent read-only acceptance SHALL record that exact commit as `HARDENING_PLANNING_HEAD`.

Apply SHALL begin only from exact `HARDENING_PLANNING_HEAD`, write only the 21 approved product/test/CI paths and task checkbox markers, and keep its index and every ref fixed. It SHALL NOT change `LICENSE`, generated OpenSpec skills/config, bootstrap archive/root spec, main-repository files, or protected user paths.

#### Scenario: Exact governed planning checkpoint is ready
- **WHEN** branch, `HARDENING_PLANNING_HEAD`, its exact subject/sole parent/five-path seal, ordinary ancestry, archived bootstrap change, root bootstrap capability, empty index, protected seals, and owned-path inventory all match
- **THEN** a delegated apply subagent may begin the approved implementation after explicit main-agent planning-checkpoint acceptance

#### Scenario: Governance or dependency differs
- **WHEN** branch, HEAD, ancestry, bootstrap content, protected input, index, writable-path complement, or dependency evidence differs
- **THEN** apply SHALL stop before product mutation and report the exact mismatch

### Requirement: The public DTO preserves the complete collection record
`Fetch` and `FetchPage` SHALL return `*Subject` values with public fields `ID`, `SubjectID`, `SubjectType`, `Type`, `Name`, `NameCn`, `Rate`, `Comment`, `Tags`, `UpdatedAt`, `VolStatus`, `EpStatus`, and `Private`. `SubjectID` SHALL contain upstream top-level `subject_id`; compatibility field `ID` SHALL equal `SubjectID`. `SubjectType` SHALL contain upstream `subject_type`; `Type` SHALL contain the collection state. Present comment, required tags, update timestamp, progress, rating, and private marker SHALL be preserved without semantic inference.

The decoder SHALL follow the current official Bangumi OAS: top-level
`subject_id`, `subject_type`, `rate`, `type`, `tags`, `ep_status`,
`vol_status`, `updated_at`, and `private` are required, while `comment` and
nested `subject` are optional. An absent comment SHALL map to `Comment == ""`.
An absent nested subject SHALL still set `ID == SubjectID` and leave
`Name`/`NameCn` empty; when present, its required identity/type/name tuple SHALL
be complete and consistent with the top level. A missing/null required tags
array SHALL fail, while an empty required array SHALL become a non-nil empty
slice and every returned tag slice SHALL be copied.

The decoder SHALL require a positive subject ID, valid subject type in `{1,2,3,4,6}`, valid collection type in `{1,2,3,4,5}`, rate in `0..10`, non-negative progress, RFC3339 `updated_at`, and exact normalized requested offset/limit metadata. Every page total SHALL be in `0..1_000_000` per collection type and `len(data)` SHALL NOT exceed the normalized requested limit. Any violation SHALL return `*ProtocolError` matching `ErrProtocol` before allocation or scheduling. Unknown additive upstream JSON fields MAY be ignored; missing/malformed required content, malformed present optional content, multiple JSON values, or trailing non-whitespace SHALL fail.

#### Scenario: Complete valid collection decodes
- **WHEN** an `httptest` response contains every required collection field and a consistent optional nested subject
- **THEN** both page and aggregate operations SHALL return every exact value
- **AND** `ID` SHALL equal `SubjectID`
- **AND** mutating a source decode buffer or another result SHALL NOT mutate the returned tags

#### Scenario: Required field or identity is invalid
- **WHEN** a success payload omits a required field, uses an invalid enum/range/timestamp, disagrees on subject identity, has invalid page metadata, or contains trailing JSON
- **THEN** the operation SHALL return a typed non-retryable decode or protocol error
- **AND** it SHALL return no partial aggregate result

#### Scenario: Optional comment and subject are absent
- **WHEN** a success payload contains every required top-level field but omits `comment` and nested `subject`
- **THEN** decoding SHALL succeed with empty `Comment`, `Name`, and `NameCn`
- **AND** `ID` SHALL still equal the required top-level `SubjectID`

#### Scenario: Empty public collection succeeds
- **WHEN** the first valid page reports total zero and contains no data
- **THEN** `Fetch` SHALL return a non-nil empty slice and nil error

### Requirement: The first public contract is anonymous and input-safe
The client SHALL retrieve only anonymous public collection data. `WithAccessToken` and all access-token state/behavior SHALL be removed. Library-built requests SHALL never set or forward Authorization, Cookie, or arbitrary caller request headers. They SHALL set only the configured `User-Agent` and `Accept: application/json`. The supplied User-Agent SHALL be valid UTF-8, `1..256` bytes, non-blank after a validation-only Unicode trim, and contain no Unicode control rune; once accepted, it SHALL be sent byte-for-byte as supplied.

The default endpoint SHALL be `https://api.bgm.tv`. `WithEndpoint` SHALL accept only an absolute HTTPS root endpoint or loopback HTTP root endpoint for tests, require a non-empty host and an empty or `/` path, reject opaque form, user-info, query, fragment, every other path, and every other insecure endpoint, and construct `/v0/users/{uid}/collections` with UID as one escaped path segment. After trimming both ends using Go Unicode whitespace semantics, UID SHALL be valid UTF-8, `1..256` bytes, contain no Unicode control rune, and preserve case. A supplied `http.Client` SHALL be shallow-cloned, have its cookie jar and `Timeout` cleared, and replace redirect handling so the first 3xx response is returned without following it. `WithRequestTimeout` SHALL govern every individual attempt regardless of option order or a supplied client. The caller-owned client SHALL remain unmodified. A caller-supplied `RoundTripper` is trusted implementation input but SHALL receive no credential header from this package.

Input classification SHALL be exact: nil context → `ErrNilContext`; blank trimmed UID → `ErrEmptyUserID`; malformed UTF-8/control-containing/over-256-byte UID → `ErrInvalidUserID`; invalid subject type → `ErrInvalidSubjectType`; invalid collection type → `ErrInvalidCollectionType`; empty `Fetch` type list → `ErrNoCollectionTypes`; and invalid User-Agent, endpoint, new-option value, or nil Option → `ErrInvalidConfiguration`. All SHALL fail before transport. Since `NewClient` cannot return an error, the first configuration fault SHALL poison the Client deterministically; a later valid Option SHALL NOT clear it.

#### Scenario: Anonymous request is constructed safely
- **WHEN** a UID containing path metacharacters and mixed case is requested through a loopback `httptest` endpoint
- **THEN** the server SHALL observe one escaped UID path segment, preserved case, expected query parameters, User-Agent, and JSON Accept header
- **AND** it SHALL observe no Authorization or Cookie header

#### Scenario: Dot-only UID remains opaque
- **WHEN** the valid normalized UID is `.` or `..`
- **THEN** its dots SHALL be percent-encoded in the one UID segment
- **AND** no client, proxy, or router path normalization SHALL change the route

#### Scenario: Custom client contains credentials or redirects
- **WHEN** `WithHTTPClient` receives a client with a cookie jar or redirect policy
- **THEN** the package-owned clone SHALL send no jar cookie and SHALL reject the first redirect without following it
- **AND** the caller-owned client object SHALL remain unmodified

#### Scenario: Input or endpoint is invalid
- **WHEN** context is nil, normalized UID or User-Agent violates its UTF-8/length/control/blank boundary, subject type is invalid, collection type is invalid, no collection types are supplied to `Fetch`, an Option is nil, QPS is non-finite/non-positive, burst is non-positive, or endpoint configuration is malformed/non-root/insecure
- **THEN** the operation SHALL fail before any transport call with the stable matching input/config sentinel

### Requirement: Existing non-authenticated call shapes have a documented compatibility boundary
The public signatures of `NewClient(userAgent string, options ...Option) *Client`, `Fetch(context.Context, string, SubjectType, ...CollectionType) ([]*Subject, error)`, and `FetchPage(context.Context, string, SubjectType, CollectionType, int, int) (*PageResult, error)` SHALL remain. Collection-type constants SHALL remain unchanged. Subject-type constants SHALL match the official OAS exactly: `SubjectTypeBook=1`, `SubjectTypeAnime=2`, `SubjectTypeMusic=3`, `SubjectTypeGame=4`, and `SubjectTypeReal=6`; this intentionally corrects the untagged prototype's reversed Game/Music names while preserving the valid raw numeric set. `WithHTTPClient(*http.Client) Option`, `WithConcurrencyLimit(int) Option`, `WithRequestTimeout(time.Duration) Option`, `WithMaxRetries(int) Option`, and `WithRetryInterval(time.Duration) Option` SHALL remain. New options SHALL be exactly `WithEndpoint(string) Option`, `WithRateLimit(float64, int) Option`, and `WithMaxRetryDelay(time.Duration) Option`. `FetchPage` SHALL continue clamping limit to `1..50` and negative offset to zero.

Invalid values supplied to an existing non-auth option SHALL retain the baseline fallback behavior: nil HTTP client, non-positive concurrency/timeout/retry interval, and negative maximum retries leave the corresponding default unchanged. Invalid new endpoint/rate/max-delay configuration or a nil `Option` SHALL retain the first client configuration fault and every operation SHALL return `ErrInvalidConfiguration` before transport without panic; no later valid Option SHALL clear it. QPS SHALL be finite and positive; NaN and either infinity SHALL be invalid. Burst SHALL be positive.

README/package/example documentation SHALL explicitly list the `v0.1.0` changes: corrected Music `3` / Game `4` named constants; `WithAccessToken` removal; anonymous-only requests; complete DTO fields and optional upstream projection behavior; canonical aggregate order; empty collection-type rejection; sanitized `HTTPError.Body`; and new endpoint, rate, and maximum-retry-delay options. Every error example SHALL prove a matching `errors.As` result before dereferencing a typed error.

#### Scenario: Named subject constants match upstream semantics
- **WHEN** Music and Game requests are constructed through their public named constants
- **THEN** the exact `subject_type` query values SHALL be `3` and `4` respectively
- **AND** docs SHALL identify the correction from the untagged prototype

#### Scenario: Baseline anonymous example migrates
- **WHEN** the baseline example adopts the corrected subject constants, removes authentication, and is formatted for the complete DTO
- **THEN** it SHALL compile against Go 1.26 with the retained constructor, Fetch signature, and non-auth options

#### Scenario: Credentialed untagged caller is evaluated
- **WHEN** a caller relies on `WithAccessToken` or Authorization behavior
- **THEN** compilation or migration documentation SHALL reveal the intentional pre-`v0.1.0` break
- **AND** no compatibility shim SHALL silently accept or transmit the token

### Requirement: Fetch performs complete bounded automatic pagination
The page size SHALL be 50. For each normalized collection type, `Fetch` SHALL request offset zero at limit 50 exactly once, require that first page's `len(data) <= total`, use its validated `total` in `0..1_000_000` as the fixed logical page plan, and request offsets `50, 100, ...` strictly below that total. Thus a first page with `total=0` and non-empty data SHALL fail as a protocol error. Page-count and offset arithmetic SHALL be checked for overflow before allocation or scheduling. It SHALL NOT make the old `limit=1` probe or refetch offset zero. Later-page total drift within `0..1_000_000` SHALL NOT expand the fixed plan; later metadata SHALL still satisfy protocol validity.

Page jobs SHALL run through bounded workers and the shared request pipeline. Remaining-page worker count SHALL never exceed `min(remainingPages, configuredConcurrency)`. A terminal page failure SHALL cancel sibling work, wait for started workers to exit, launch no further retry after cancellation, and return no partial result. No operation SHALL leave a goroutine, timer, body, or semaphore permit behind.

#### Scenario: Pagination metadata exceeds a hard boundary
- **WHEN** a page reports total `1_000_001` or a near-`MaxInt` value, returns more items than the normalized limit, a first page returns more items than its total, or metadata disagrees with the requested offset/limit
- **THEN** the operation SHALL return `*ProtocolError` matching `ErrProtocol` before unbounded allocation, arithmetic, or scheduling
- **AND** total `1_000_000` SHALL remain valid with a peak remaining-page worker count no greater than configured concurrency

#### Scenario: Multi-page multi-type collection succeeds
- **WHEN** two normalized collection types require multiple pages and pages complete out of order
- **THEN** every planned offset SHALL be requested once through the shared client controls
- **AND** offset zero SHALL NOT be requested twice
- **AND** the final result SHALL contain all unique records in canonical order

#### Scenario: First page fixes a moving plan
- **WHEN** a later page reports a different in-range total because the upstream collection changed
- **THEN** the operation SHALL finish only the offsets derived from the first total
- **AND** it SHALL neither loop, append an unplanned offset, nor infer a snapshot guarantee

#### Scenario: One page fails
- **WHEN** one page reaches a terminal classified error while sibling page work is pending
- **THEN** siblings SHALL observe cancellation and finish
- **AND** Fetch SHALL return the error with no partial slice and no leaked worker

### Requirement: Fetch normalizes, deduplicates, and sorts deterministically
`Fetch` SHALL validate collection types, remove repeated requested values, and sort the normalized type set numerically before scheduling. Every decoded record SHALL carry a deterministic source coordinate `(collectionType, offset, itemIndex)`. Aggregate deduplication key SHALL be `(SubjectType, SubjectID, Type)` so collection state is never discarded. If the same key appears more than once, the record with the lexicographically smallest source coordinate SHALL win. The final total order SHALL be `(SubjectType ascending, SubjectID ascending, Type ascending)`.

`FetchPage` SHALL remain a page primitive and preserve the upstream order within that page after validation/conversion; cross-page dedupe/sort SHALL belong only to `Fetch`.

#### Scenario: Completion and input order vary
- **WHEN** identical fixture pages are delivered under different goroutine delays and collection types are requested in different/repeated orders
- **THEN** `Fetch` SHALL return byte-for-byte equivalent field values and key order on every run

#### Scenario: Duplicate record conflicts
- **WHEN** repeated/overlapping pages contain two records with the same full key but different non-key values
- **THEN** the smallest source coordinate SHALL win deterministically
- **AND** a record with a different collection `Type` SHALL remain a separate result

### Requirement: QPS and concurrency are client-wide and context-aware
Every HTTP attempt, including every retry, SHALL wait on the same per-Client `golang.org/x/time/rate.Limiter` and acquire the same per-Client in-flight semaphore. Distinct Client values SHALL not share limits. Default QPS SHALL be 3 requests/second with burst 1; default maximum in-flight requests SHALL remain 10. `WithRateLimit(qps, burst)` SHALL accept only finite positive QPS and positive burst; `WithConcurrencyLimit(limit)` SHALL accept positive values before the client is used.

An attempt SHALL rate-wait before acquiring concurrency, acquire no permit when its context is already done, and release any acquired permit exactly once. If the rate limiter refuses to wait because its reservation would exceed the parent deadline while `ctx.Err()` is not yet set, the call SHALL return the terminal parent-deadline classification directly, make no transport call, and SHALL NOT retry or match `ErrRetryExhausted`. After acquisition, every attempt SHALL create a fresh `WithRequestTimeout` child context covering transport, complete bounded body read, validation, and close, then cancel it. Waiting for QPS, concurrency, retry delay, transport, and body read SHALL all terminate when the parent operation context terminates. Parent cancellation/deadline SHALL take classification precedence over an attempt timeout.

#### Scenario: Custom HTTP client does not disable attempt timeout
- **WHEN** a custom client with zero, shorter, or longer `Client.Timeout` is supplied before or after `WithRequestTimeout`, and retries occur
- **THEN** the caller client SHALL remain unchanged, its clone timeout SHALL be cleared, and every attempt SHALL receive a fresh configured child deadline
- **AND** an earlier parent deadline SHALL terminate and classify the call before the attempt deadline

#### Scenario: Concurrent public calls share limits
- **WHEN** multiple `Fetch` and `FetchPage` calls execute concurrently on one Client against a timestamping/blocking test server
- **THEN** observed aggregate start rate and maximum in-flight count SHALL never exceed configured QPS/burst and concurrency

#### Scenario: Waiter is canceled
- **WHEN** a call is canceled while rate-waiting or semaphore-waiting
- **THEN** it SHALL return the matching canceled/deadline classification promptly
- **AND** it SHALL make no transport call and leak no permit

#### Scenario: Reservation would outlive parent deadline
- **WHEN** the rate limiter rejects a wait because the next token lies after the parent deadline, before the context timer fires
- **THEN** the call SHALL return the terminal parent-deadline classification without retry exhaustion
- **AND** it SHALL make no transport call

#### Scenario: Separate clients run independently
- **WHEN** two Client values use independent test servers and limiters
- **THEN** one client’s budget SHALL NOT consume or release the other client’s budget

### Requirement: Retry behavior is closed, bounded, jittered, and Retry-After aware
`WithMaxRetries(n)` SHALL mean no more than `n` retries after one initial attempt; zero SHALL disable retries. Default retries SHALL be 3, default base interval 1 second, and default maximum retry delay 30 seconds. `WithRetryInterval` SHALL set a positive base; `WithMaxRetryDelay` SHALL set a positive cap.

Only live-parent-context transport failures, per-attempt timeouts, HTTP 429, and HTTP 500–599 SHALL retry. Caller cancellation/deadline, input/config, every other HTTP status including 1xx, non-200 2xx, 3xx, 401, 403, 404, and other 4xx, decode, protocol, oversize, and limiter/acquire cancellation SHALL be terminal. Before classification at every failure point, implementation SHALL check the parent context so caller termination cannot be misclassified as a retryable attempt timeout. Before a retry, its prior body SHALL be bounded/read as required and closed.

For retry number `n`, local full-jitter delay SHALL be in `[0, min(base*2^(n-1), maxDelay)]` using overflow-safe arithmetic. A valid non-negative `Retry-After` delta-seconds or HTTP-date on a retryable response SHALL contribute `min(retryAfter, maxDelay)`; an arbitrarily long all-ASCII-digit delta SHALL saturate to `maxDelay` without native-integer overflow. The actual wait SHALL be the greater of jitter and bounded Retry-After. Invalid, negative, or past values SHALL be ignored. Every wait SHALL be context-cancellable. Exhaustion SHALL return a typed retry error with total attempts, match `ErrRetryExhausted`, and unwrap the last sanitized error so its classification remains visible.

#### Scenario: Retryable matrix succeeds
- **WHEN** a table-driven server/transport returns a transport failure, attempt timeout, 429, or representative 5xx before a success
- **THEN** attempts SHALL pass through the shared limiter and stop immediately on success within the configured budget

#### Scenario: Retry-After formats are honored
- **WHEN** 429 or 503 supplies valid delta-seconds or HTTP-date Retry-After
- **THEN** the injected clock/sleeper evidence SHALL show the greater of bounded server delay and local full jitter

#### Scenario: Delay input is unusable
- **WHEN** Retry-After is malformed, negative, past, greater than the configured maximum, or an all-digit delta exceeds native integer range
- **THEN** malformed/negative/past input SHALL be ignored and an excessive valid value SHALL be capped
- **AND** no arithmetic overflow or unbounded sleep SHALL occur

#### Scenario: Terminal class is not retried
- **WHEN** the result is cancellation, 401, 403, 404, other 4xx, redirect, decode, protocol, or oversized response
- **THEN** exactly one attempt SHALL be observed

#### Scenario: Retry budget is exhausted
- **WHEN** every permitted attempt fails retryably
- **THEN** attempts SHALL equal `maxRetries + 1`
- **AND** the result SHALL match both retry exhaustion and the last underlying stable classification

### Requirement: Errors are typed, stable, bounded, and sanitized
The package SHALL expose these exact stable sentinel identifiers: existing `ErrInvalidUserID`, `ErrUnauthorized`, `ErrForbidden`, `ErrRateLimited`, `ErrServerError`, and `ErrEmptyUserID`; new `ErrNilContext`, `ErrNoCollectionTypes`, `ErrInvalidSubjectType`, `ErrInvalidCollectionType`, `ErrInvalidConfiguration`, `ErrNotFound`, `ErrHTTPStatus`, `ErrTransport`, `ErrTimeout`, `ErrCanceled`, `ErrDecode`, `ErrProtocol`, `ErrResponseTooLarge`, and `ErrRetryExhausted`. `errors.As` SHALL expose `*HTTPError`, `*NetworkError`, `*DecodeError`, `*ProtocolError`, and `*RetryError` as appropriate.

`HTTPError` SHALL retain exported `StatusCode int` and `Body string` and add exported `RetryAfter time.Duration`; `Body` SHALL always be empty. Only HTTP 200 SHALL be success. Every other final status, including 1xx, non-200 2xx, 3xx, and all 4xx/5xx, SHALL return `*HTTPError` matching `ErrHTTPStatus`. HTTP 401/403/404/429/5xx SHALL additionally match `ErrUnauthorized`/`ErrForbidden`/`ErrNotFound`/`ErrRateLimited`/`ErrServerError` respectively; 404 SHALL also match deprecated `ErrInvalidUserID`. Other statuses SHALL match no narrower HTTP sentinel. Redirects SHALL not be followed; the first 3xx SHALL use this terminal `HTTPError` mapping without retaining or rendering Location.

Nil context SHALL return `ErrNilContext` directly before transport and SHALL not be wrapped as a network error. `NetworkError` SHALL retain exported `Err error` and add exported `Timeout bool`; `Err` SHALL contain only a safe sentinel or `context.Canceled`/`context.DeadlineExceeded`, never raw URL/transport text. Parent cancellation SHALL return `*NetworkError{Timeout:false}` matching both `ErrCanceled` and `context.Canceled`; parent deadline SHALL return `*NetworkError{Timeout:true}` matching both `ErrTimeout` and `context.DeadlineExceeded`; both are terminal. A per-attempt timeout while the parent remains live SHALL return the same timeout/deadline matches and MAY retry. Another transport failure while the parent remains live SHALL return `*NetworkError` matching `ErrTransport` and MAY retry. Rate, semaphore, retry-wait, transport, and body-read cancellation SHALL all follow the parent classification.

`RetryError` SHALL expose exported `Attempts int` and `Err error`, with `Err` restricted to a sanitized typed error. `DecodeError` and `ProtocolError` SHALL expose no raw body, UID, URL, header, or transport field. Error strings SHALL be at most 256 bytes and SHALL contain none of UID, URL/path/query, response body, Authorization, Cookie, user agent, Location, arbitrary headers, or raw transport text.

The maximum success body SHALL be 16 MiB plus one detection byte. An oversized success body SHALL be a distinct non-retryable typed failure. At most 64 KiB plus one detection byte of an error body SHALL be read before it is discarded; an oversized error body SHALL NOT replace the already known HTTP status classification or suppress an otherwise permitted 429/5xx retry. Error bodies SHALL never be retained or rendered. Unreadable success body, malformed JSON, required-field failure, and trailing JSON SHALL remain distinct non-retryable typed failures.

#### Scenario: HTTP classes support Is and As
- **WHEN** a table server returns representative 1xx, 204, 3xx, 401, 403, 404, 429, another 4xx, and representative 5xx
- **THEN** each result SHALL `errors.As` to `*HTTPError`
- **AND** each SHALL match `ErrHTTPStatus` and only its documented narrower stable categories, with 404 also matching deprecated invalid-user-ID compatibility
- **AND** a 3xx SHALL not be followed and no Location value SHALL escape

#### Scenario: Transport, timeout, and cancellation are classified
- **WHEN** context is nil; a safe fake transport fails; a per-attempt timeout expires while the parent is live; or the parent is canceled/deadlined during rate, semaphore, retry wait, transport, or body read
- **THEN** nil context SHALL match only `ErrNilContext` before transport, while each runtime result SHALL `errors.As` to `*NetworkError` and `errors.Is` to the exact stable/context category above
- **AND** parent cancellation/deadline SHALL remain terminal while only live-parent transport/attempt-timeout failure MAY retry

#### Scenario: Sensitive material is present upstream
- **WHEN** UID, URL, headers, raw transport error, and a 64-KiB-plus secret response body contain unique marker strings
- **THEN** no returned error string, public Body field, unwrap string, or test log SHALL contain any marker
- **AND** the error string SHALL remain within 256 bytes

#### Scenario: Success or error body is oversized
- **WHEN** a success body exceeds 16 MiB or an HTTP error body exceeds 64 KiB
- **THEN** reading SHALL stop after the applicable limit plus detection byte and close the body
- **AND** oversized success SHALL return the documented non-retryable sanitized error
- **AND** oversized error content SHALL retain the status-based HTTP classification and retry eligibility without retaining its bytes

### Requirement: Client lifecycle is race-safe and cancellation-complete
A Client SHALL be safe for concurrent use after `NewClient` returns. Public operations SHALL not mutate configuration or return shared mutable DTO/tag backing storage. First terminal failure SHALL cancel related work. Every started goroutine SHALL be joined; every timer/wait SHALL be stoppable; every response body SHALL close; every limiter/semaphore resource SHALL be balanced.

#### Scenario: Race suite exercises shared client
- **WHEN** the same Client performs concurrent multi-page Fetch and FetchPage operations, cancellations, retries, and duplicate normalization under `go test -race`
- **THEN** the suite SHALL report no race
- **AND** server counters SHALL return to zero after each test

#### Scenario: Caller mutates returned data
- **WHEN** a caller mutates a returned Subject or Tags slice after one operation
- **THEN** no later result or internal client state SHALL change

### Requirement: Go 1.26 dependencies and CI have explicit value and cost gates
`go.mod` SHALL declare the Go language and module compatibility floor as
`1.26.0` and SHALL declare `toolchain go1.26.5` so default auto selection uses
the accepted patched standard library. Runtime dependencies SHALL be exactly
direct `golang.org/x/sync v0.22.0` and `golang.org/x/time v0.15.0` plus only
dependencies mechanically required by those modules. `x/sync/errgroup` SHALL
own bounded structured cancellation; `x/time/rate` SHALL own the
concurrency-safe token bucket. No runtime dependency type SHALL leak into the
public API.

Tests SHALL use standard-library `testing`, `httptest`, and helpers only; no assertion, mocking, HTTP-client, logging, or retry framework SHALL be introduced. `.github/workflows/ci.yml` SHALL be test/build-only, run Go 1.26.x format/module/vet/test/race gates, declare read-only contents permission, consume no secret, publish no artifact, deploy nothing, and use only `actions/checkout@3d3c42e5aac5ba805825da76410c181273ba90b1` with comment `v7.0.1` and `actions/setup-go@b7ad1dad31e06c5925ef5d2fc7ad053ef454303e` with comment `v7.0.0`.

The dependency admission record SHALL state: standard library lacks an equivalent cancellable token bucket; custom structured concurrency/limiting has higher cancellation/race risk; both runtime modules are official Go BSD-3-Clause modules with small module/supply-chain cost; CI actions are GitHub-owned and SHA-pinned. Removal SHALL require a later reviewed change, an equivalent standard/local replacement, and all behavior/race gates.

#### Scenario: Module and toolchain graph is exact
- **WHEN** default `GOTOOLCHAIN=auto` selection runs module and toolchain inspection
- **THEN** `go env GOVERSION` SHALL report `go1.26.5`, `go mod tidy -diff` SHALL have no diff, and `go.sum` SHALL remain unchanged
- **AND** direct versions SHALL be exactly `x/sync v0.22.0` and `x/time v0.15.0`
- **AND** no unreviewed production or test module SHALL appear

#### Scenario: Patched standard library has no reachable known vulnerability
- **WHEN** `govulncheck ./...` runs under the default selected Go 1.26.5 toolchain
- **THEN** it SHALL report no reachable vulnerability in the client or its standard-library call graph

#### Scenario: CI is least privilege and test-only
- **WHEN** the workflow is inspected mechanically
- **THEN** the two exact approved action OIDs SHALL be present, no other `uses:` target SHALL exist, permissions SHALL be read-only, Go SHALL be 1.26.x, and format/module/vet/test/race commands SHALL be present
- **AND** secret, credential, artifact publication, registry, release, deploy, SSH, and write-permission steps SHALL be absent

### Requirement: Automated evidence covers every admission behavior
The root package SHALL contain non-empty tests for DTO validation, anonymous request construction, endpoint/custom-client safety, automatic pagination, fixed moving-total behavior, QPS, in-flight concurrency, Retry-After delta/date, full jitter bounds, retry matrix/exhaustion, typed error `Is/As`, sanitization, size limits, deterministic dedupe/sort, cancellation/leak balance, compatibility, and caller mutation isolation. Tests SHALL use loopback `httptest` or safe fake transports only and SHALL make no request to `api.bgm.tv` or another public host.

Local acceptance SHALL run under Go 1.26.x: `gofmt`, `go mod tidy -diff`, exact module checks, `go vet ./...`, `go test -count=1 ./...`, and `go test -race -count=1 ./...`. A root package result of `[no test files]`, a skipped required scenario, a real-network dependency, or a flaky timing assertion SHALL fail acceptance.

#### Scenario: Complete local matrix passes
- **WHEN** all acceptance commands run in the exact governed candidate
- **THEN** every named behavior SHALL have executable passing evidence
- **AND** root package tests SHALL not report `[no test files]`

#### Scenario: Test attempts public network
- **WHEN** a test would resolve/connect to a non-loopback endpoint
- **THEN** the guard transport SHALL fail the test before the request leaves the process

### Requirement: Development uses an accepted planning commit, one implementation commit, and no publication
Before product apply, the exact approved five-file active change SHALL be committed with subject `docs(openspec): approve collection client hardening`, sole parent `59eba79b3e3621ce72b756a88b38aee970c00fcf`, five paths, and sorted path seal `a7da0916df9b59cdaff4d5a3b57d8b0231c8bad51921aedb308ee2d10262fe1f`. Main-agent read-only acceptance SHALL name that commit `HARDENING_PLANNING_HEAD`.

The apply subagent SHALL stop first with an unstaged implementation candidate, empty index, exact `HEAD=HARDENING_PLANNING_HEAD`, exact owned-path diff, complete tasks, and passing acceptance. Only after main-agent read-only acceptance MAY a delegated archive/finalization subagent synchronize this capability, replace generated root Purpose with exactly `Define the first public anonymous Bangumi collection client contract, including complete DTOs, automatic pagination, shared limiting, bounded retries, sanitized errors, deterministic results, and local quality gates.`, archive to `openspec/changes/archive/2026-07-24-harden-bangumi-collection-go-v0-1-0/`, and stage the exact final candidate.

The staged candidate SHALL have a no-rename parent delta of exactly 32 paths relative to `HARDENING_PLANNING_HEAD`: the 21 approved product/test/CI paths; five active-change deletions; `.openspec.yaml`, `proposal.md`, `design.md`, `specs/collection-client-v0-1-0/spec.md`, and `tasks.md` below the exact date-stamped archive; and `openspec/specs/collection-client-v0-1-0/spec.md`. Its sorted path seal SHALL be `740610fc014608d9294383d731e3a9df51b3521c25823824c6ef89dc582349b2`. Rename detection SHALL be disabled for acceptance and post-commit proof.

The same final tree SHALL have exactly 27 cumulative paths relative to bootstrap checkpoint `59eba79b3e3621ce72b756a88b38aee970c00fcf`, with sorted path seal `b8ef7390c49b73f13061c3ca50f66a0ef260cf947c539d507ab1797de59c48c1`. The synchronized root requirements SHALL be byte-identical to the archived delta requirements. After a second main-agent read-only acceptance, finalization MAY create one local commit with sole parent `HARDENING_PLANNING_HEAD` and exact subject `feat: harden collection client for v0.1.0`.

The commit SHALL contain no protected/governance/complement path, leave only `CLAUDE.md`, `note`, and ignored `.claude/settings.local.json` in tolerated local status, and pass every OpenSpec/Go/diff/seal gate. It SHALL NOT fetch, push, create a PR, tag, release, publish, deploy, amend, rebase, reset, or move any other ref. Public `v0.1.0` and main-repository admission SHALL require later explicit authorization and a fixed public tag.

#### Scenario: Unstaged implementation awaits first acceptance
- **WHEN** implementation and tests pass in approved paths
- **THEN** the index SHALL remain empty and HEAD SHALL remain `HARDENING_PLANNING_HEAD`
- **AND** archive/staging/commit SHALL remain blocked until main-agent acceptance

#### Scenario: Exact archive candidate awaits second acceptance
- **WHEN** the accepted implementation is archived/synchronized and the exact no-rename 32-path parent delta / 27-path cumulative inventory is staged
- **THEN** HEAD SHALL still equal `HARDENING_PLANNING_HEAD`
- **AND** no commit SHALL be created before the main agent accepts that exact index/tree/seal

#### Scenario: Accepted local commit is created
- **WHEN** both acceptances pass and finalization creates the exact-subject sole-parent commit
- **THEN** investigated, specified, implemented, verified, and committed MAY become true
- **AND** pushed, released, and deployed SHALL remain false

#### Scenario: Remote or publication action is attempted
- **WHEN** any step would fetch, push, open a PR, create/move a tag, release, publish, deploy, or admit the package in the main repository
- **THEN** it SHALL stop as unauthorized even if the local commit passes
