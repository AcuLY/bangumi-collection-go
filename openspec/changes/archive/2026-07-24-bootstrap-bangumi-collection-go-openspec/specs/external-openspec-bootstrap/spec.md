## ADDED Requirements

### Requirement: Bootstrap is rooted in the exact external repository and baseline
The bootstrap SHALL operate only in canonical repository root `/Users/luca/dev/bangumi-collection-go`. Its local branch SHALL be `codex/v0.1.0-hardening`, and commit `8173f44911360150a5a5a7c6418021d1014fe85b` SHALL be its exact initial HEAD and an ordinary ancestor of all later authorized hardening work. Preflight SHALL verify that `refs/heads/main` and `refs/remotes/origin/main` resolve to that commit without fetching or moving either ref.

The branch creation and uncommitted `openspec init --tools codex --profile core` output SHALL be the only bootstrap exception to ordinary spec-before-apply, because the planning home cannot describe itself before it exists. That exception SHALL contain no product change, staging, or commit. The complete bootstrap change SHALL receive main-agent approval before validation-only apply, archive, or commit, and no later product change may reuse this exception.

#### Scenario: Exact baseline and branch match
- **WHEN** the canonical root, local branch, exact HEAD, local main, remote-tracking main, and ordinary ancestry all match
- **THEN** bootstrap validation may continue without changing a Git ref

#### Scenario: Repository or baseline differs
- **WHEN** any root, branch, commit, ref, or ancestry check differs from the approved values
- **THEN** bootstrap SHALL stop before further mutation and report the exact mismatch

#### Scenario: Bootstrap exception tries to include product work
- **WHEN** pre-spec initialization output contains a product/module/test change, staged path, commit, or any artifact beyond branch metadata and the generated OpenSpec framework
- **THEN** bootstrap acceptance SHALL fail
- **AND** ordinary hardening SHALL remain blocked behind its own approved change

### Requirement: The repository has one governed generated OpenSpec root
The repository SHALL contain one repo-local OpenSpec root with `openspec/config.yaml` declaring `schema: spec-driven`, repository-specific context, and proposal/design/tasks rules. Context SHALL identify the external public-collection client lane, exact baseline and branch, protected local paths, spec-before-apply rule, exact owner/path/dependency requirements, and the separation between local development commits and separately authorized push/PR/tag/release. Artifact rules SHALL require every proposal/design/tasks artifact to declare Owner, writable paths, protected inputs, consumes, produces, exact dependencies, deliverables, executable acceptance, non-goals, operations deferred, stop/rollback, and separate investigated/specified/implemented/verified/committed/pushed/released/deployed status. Default example comments alone SHALL be invalid.

Its generated Codex inventory SHALL contain exactly the six `SKILL.md` paths for `openspec-apply-change`, `openspec-archive-change`, `openspec-explore`, `openspec-propose`, `openspec-sync-specs`, and `openspec-update-change` below `.codex/skills/`, and each SHALL declare `generatedBy: "1.6.0"`. Before archive, the only additional candidate content SHALL be the active `bootstrap-bangumi-collection-go-openspec` change. After approved archive/sync, that active change SHALL be absent and the only additional OpenSpec content SHALL be its exact date-stamped archive plus `openspec/specs/external-openspec-bootstrap/spec.md`.

#### Scenario: Generated framework inventory matches
- **WHEN** config, skill paths, skill provenance, planning home, and active change enumeration equal the approved inventory
- **THEN** OpenSpec SHALL resolve the repository root locally and expose this change through the `spec-driven` workflow

#### Scenario: Framework inventory drifts
- **WHEN** a generated path is missing, extra, nested under another planning root, linked through a symlink, or has unapproved provenance
- **THEN** bootstrap acceptance SHALL fail without regenerating or deleting paths

### Requirement: User-owned local files remain byte- and state-exact
Bootstrap SHALL preserve `CLAUDE.md` as untracked with SHA-256 `c69937414d643411e03f3d7982e4bf58e8ceb02802cf1338dada4ea6a630407d`, `note` as untracked with SHA-256 `7b2376f605f5d9d9f677eaf1165a2ff911cec9092b6b2763b966a3d6f93e2e74`, and `.claude/settings.local.json` as ignored with SHA-256 `c67845ef6924b080344bc0010931336106d9df3a700f67f10a22763d07517e3e`. The exact three-path porcelain projection SHALL remain `?? CLAUDE.md`, `?? note`, and `!! .claude/settings.local.json`, hash to `221322afd4137d24875b5fbd369c9028e99f69f1cf155e3c49975996841a3938`, and have no index membership.

#### Scenario: Protected files remain unchanged
- **WHEN** their content hashes, porcelain projection, ignore classification, and index absence are rechecked after every bootstrap phase
- **THEN** all values SHALL equal the approved seals and their presence SHALL not be reported as bootstrap output

#### Scenario: Protected content or state drifts
- **WHEN** any protected byte, status code, ignore source, or index membership differs
- **THEN** work SHALL stop without clean, reset, stash, checkout restoration, deletion, staging, or overwrite

### Requirement: Baseline product content is read-only
The literal baseline paths `LICENSE`, `README.md`, `client.go`, `collection.go`, `errors.go`, `example/example.go`, `go.mod`, `go.sum`, `options.go`, and `types.go` SHALL remain byte-identical to commit `8173f44911360150a5a5a7c6418021d1014fe85b` in both worktree and index. Their ordered aggregate SHALL be computed with `printf '%s\0' LICENSE README.md client.go collection.go errors.go example/example.go go.mod go.sum options.go types.go | xargs -0 shasum -a 256 | shasum -a 256` and SHALL remain `7d4e55231e08c5d8286bffb921e53adbda9a5d82ced9c4688e4cadb129b1855e`. Current `git ls-files` SHALL NOT define this seal because accepted framework paths become tracked. No product, example, module, README, or license path SHALL enter the bootstrap candidate.

#### Scenario: Product complement is unchanged
- **WHEN** tracked inventory, aggregate seal, worktree diff, cached diff, and physical status paths are audited
- **THEN** every non-OpenSpec tracked path SHALL remain identical to the baseline

#### Scenario: Product or dependency file changes
- **WHEN** any tracked product byte changes or a generated Go/module/product path appears
- **THEN** bootstrap acceptance SHALL fail and the candidate SHALL remain unstaged

### Requirement: Bootstrap apply is validation-only and path-exact
Apply SHALL accept the already generated config and skills without rerunning `openspec init`. Its only mutable content SHALL be checkbox markers in `openspec/changes/bootstrap-bangumi-collection-go-openspec/tasks.md`; all other framework, product, protected, index, and Git-ref state SHALL remain immutable. Apply SHALL NOT implement or test v0.1.0 hardening behavior.

#### Scenario: Validation-only apply runs
- **WHEN** the exact approved framework bytes and preflight seals match
- **THEN** a delegated apply subagent MAY complete this change's checkbox markers and run its read-only validation commands

#### Scenario: Apply attempts regeneration or hardening
- **WHEN** a command would rerun initialization, rewrite generated skills/config, edit Go code, add a hardening artifact, or mutate any path outside task checkboxes
- **THEN** apply SHALL stop before that command and report the boundary violation

### Requirement: Acceptance is strict, unstaged, and local
Acceptance SHALL require artifact completeness, `openspec validate --all --strict`, `openspec doctor --json`, `git diff --check`, byte-format and symlink checks, exact generated-path enumeration, unchanged protection/product seals, and an empty Git index. The implementation subagent SHALL stop with the candidate unstaged and uncommitted for main-agent read-only review.

#### Scenario: All bootstrap gates pass
- **WHEN** every strict, doctor, path, byte, seal, branch, and index check passes with no unresolved P0 or P1 finding
- **THEN** status MAY be reported as investigated, specified, framework-generated, and locally verified
- **AND** committed, pushed, released, and deployed SHALL remain false

#### Scenario: A validation gate fails
- **WHEN** any artifact, strict validation, doctor, path, byte, seal, branch, or index gate fails
- **THEN** the candidate SHALL remain unstaged and no commit, archive, push, PR, tag, release, or deployment SHALL be authorized

### Requirement: Archive and commit require a second acceptance
Only after the main agent accepts the unstaged validation-only apply candidate MAY one delegated finalization subagent archive/synchronize `bootstrap-bangumi-collection-go-openspec`. It SHALL read and follow `openspec-archive-change`, produce exactly `openspec/changes/archive/2026-07-24-bootstrap-bangumi-collection-go-openspec/**`, and synchronize `openspec/specs/external-openspec-bootstrap/spec.md`. It SHALL replace any generated `TBD` Purpose in that root spec with exactly `Define the independent OpenSpec governance, baseline ancestry, protected local-file boundary, and validation gates for the external bangumi-collection-go hardening lane.` while keeping the root requirements byte-identical to the archived delta requirements. It SHALL stage only the six accepted skill files, `openspec/config.yaml`, the five archived change files, and that root spec: exactly 13 added files. It SHALL keep every baseline product path and protected user path out of the index and stop with the exact archive candidate staged for a second main-agent read-only acceptance.

Only that second acceptance SHALL authorize one local commit on `codex/v0.1.0-hardening` with exact subject `chore: establish external openspec governance` and sole parent `8173f44911360150a5a5a7c6418021d1014fe85b`. The commit SHALL have exactly those 13 accepted added paths, contain no protected/product path, leave the active change absent, pass strict validation/doctor/diff checks, and leave only the three protected user paths as tolerated local status. It SHALL NOT fetch, push, open a PR, tag, release, publish, amend, rebase, or begin hardening apply.

#### Scenario: Staged archive candidate awaits acceptance
- **WHEN** the unstaged apply candidate was accepted and finalization has archived/synchronized and staged the exact approved paths
- **THEN** `HEAD` SHALL still equal the baseline
- **AND** no commit SHALL exist until the main agent accepts that exact index

#### Scenario: Accepted local governance commit is created
- **WHEN** the second acceptance names the exact staged path/hash inventory
- **THEN** finalization MAY create the exact-subject single-parent local commit
- **AND** committed SHALL become true while pushed, released, and deployed remain false

### Requirement: Later hardening has a separate approval gate
The future `harden-bangumi-collection-go-v0-1-0` change SHALL depend on accepted bootstrap governance and SHALL separately declare its exact product paths, behavior, tests, Go version, commit, and publication gates. This bootstrap capability SHALL grant no product-write or remote-mutation authority.

#### Scenario: Hardening is proposed after bootstrap
- **WHEN** bootstrap has independent main-agent acceptance and a later hardening proposal is requested
- **THEN** that proposal SHALL be created and reviewed as a separate OpenSpec change before product implementation

#### Scenario: Product or remote work is attempted early
- **WHEN** bootstrap work attempts hardening code, fetch/push, PR, tag, release, publication, or modification of the `BangumiStaffStats` repository
- **THEN** the action SHALL be rejected as outside this capability
