## Context

`bangumi-collection-go` is an external dependency lane, not one of the `BangumiStaffStats` repository's internal changes. Its fixed baseline is commit `8173f44911360150a5a5a7c6418021d1014fe85b`; at bootstrap preflight, `HEAD`, `refs/heads/main`, and `refs/remotes/origin/main` all resolved to that commit and the tracked tree/index were clean. The local branch `codex/v0.1.0-hardening` was then created directly from that commit without cleaning, stashing, resetting, or overwriting the checkout.

The checkout intentionally contains three user-owned paths that are not repository inputs:

| Path | Required Git state | SHA-256 |
|---|---|---|
| `CLAUDE.md` | untracked (`??`) | `c69937414d643411e03f3d7982e4bf58e8ceb02802cf1338dada4ea6a630407d` |
| `note` | untracked (`??`) | `7b2376f605f5d9d9f677eaf1165a2ff911cec9092b6b2763b966a3d6f93e2e74` |
| `.claude/settings.local.json` | ignored (`!!`) | `c67845ef6924b080344bc0010931336106d9df3a700f67f10a22763d07517e3e` |

Before initialization, `git status --porcelain=v1 --untracked-files=all --ignored=matching` contained exactly those three entries and had SHA-256 `221322afd4137d24875b5fbd369c9028e99f69f1cf155e3c49975996841a3938`. After framework files exist, that same value is used only for the three-path status projection; the full status necessarily also contains the approved OpenSpec candidate.

OpenSpec 1.6.0 has already generated a root `spec-driven` config and the Codex core skill set using `openspec init --tools codex --profile core`. The existing tracked product inventory remains exactly `LICENSE`, `README.md`, `client.go`, `collection.go`, `errors.go`, `example/example.go`, `go.mod`, `go.sum`, `options.go`, and `types.go`. Its deterministic aggregate seal is computed only over that literal ordered list with `printf '%s\0' LICENSE README.md client.go collection.go errors.go example/example.go go.mod go.sum options.go types.go | xargs -0 shasum -a 256 | shasum -a 256`; the approved value is `7d4e55231e08c5d8286bffb921e53adbda9a5d82ced9c4688e4cadb129b1855e`. Current `git ls-files` is not the authority because finalization intentionally adds framework/archive paths.

This is the lane's one unavoidable control-plane bootstrap: `openspec init` must create the planning home before an active OpenSpec change can exist. The exception authorized only branch creation plus generated framework bytes, never product behavior, staging, or a commit. Retaining, validating, archiving, and committing those bytes is governed by this now-complete change and still occurs only after main-agent spec approval. Every later external-lane product change follows ordinary spec-before-apply without an exception.

### Change boundary

| Field | Declaration |
|---|---|
| Status | investigated: complete; specified: pending approval; framework-generated: planning candidate only; implemented: no; verified: planning-only; committed: no; pushed: no; released: no; deployed: no |
| Owner | External OpenSpec bootstrap subagent for framework generation, validation-only apply, archive, and local commit; main agent amends OpenSpec only and performs both read-only acceptances |
| Writable / owned paths | Planning: `openspec/config.yaml`, six exact generated `.codex/skills/openspec-*/SKILL.md` files, and `openspec/changes/bootstrap-bangumi-collection-go-openspec/**`. Apply: only task-checkbox markers. Finalization: accepted framework paths, removal of that active change through archive, `openspec/changes/archive/2026-07-24-bootstrap-bangumi-collection-go-openspec/**`, and `openspec/specs/external-openspec-bootstrap/spec.md` |
| Read-only protected inputs | The ten exact baseline product paths at `8173f44911360150a5a5a7c6418021d1014fe85b`; local/remote-tracking main refs; three sealed user paths; main-repository master-plan external lane |
| Consumes | Exact baseline, protected-file/status seals, master-plan bootstrap and hardening boundary |
| Produces | Independent OpenSpec root/config/skills, archived bootstrap change, synchronized root capability, exact local governance commit |
| Dependencies | None |
| Deliverables | Config with repository-specific context/rules, six generator-provenance skills, strict-valid accepted capability/archive, protected/product evidence, clean local commit |
| Acceptance | Exact inventory/seals, strict/doctor/format/path/index checks, unstaged-candidate acceptance, staged archive-candidate acceptance, single-parent/subject/delta/clean-branch post-commit proof |
| Non-goals | Product behavior or tests, Go 1.26 migration, hardening proposal/apply, remote/publication activity, main-repository write |
| Operations deferred | Push, PR, tag, release/module publication, deploy, and any host/service work |
| Mutable refs | The local branch was created once by the delegated bootstrap preflight from the fixed baseline. Refs freeze during planning/apply. After staged archive acceptance, finalization may advance only `refs/heads/codex/v0.1.0-hardening` once with subject `chore: establish external openspec governance` |
| Stop / rollback conditions | Any root/branch/HEAD/ancestry/seal/path/approval/validation/index/archive/staged-delta mismatch stops in place; no reset/clean/stash/checkout restoration/broad delete/history rewrite |

### Executable acceptance

The delegated apply/finalization owner runs each phase from `/Users/luca/dev/bangumi-collection-go` with `zsh`. Before archive, all commands below must exit zero; the exact inventory also makes the six owned skill paths explicit:

```zsh
test "$(git rev-parse --show-toplevel)" = "/Users/luca/dev/bangumi-collection-go"
test "$(git branch --show-current)" = "codex/v0.1.0-hardening"
test "$(git rev-parse HEAD)" = "8173f44911360150a5a5a7c6418021d1014fe85b"
test "$(git rev-parse refs/heads/main)" = "8173f44911360150a5a5a7c6418021d1014fe85b"
test "$(git rev-parse refs/remotes/origin/main)" = "8173f44911360150a5a5a7c6418021d1014fe85b"
git merge-base --is-ancestor 8173f44911360150a5a5a7c6418021d1014fe85b HEAD
test -z "$(git diff --cached --name-only)"
test "$(find .codex openspec -type f -print | LC_ALL=C sort)" = "$(printf '%s\n' \
  .codex/skills/openspec-apply-change/SKILL.md \
  .codex/skills/openspec-archive-change/SKILL.md \
  .codex/skills/openspec-explore/SKILL.md \
  .codex/skills/openspec-propose/SKILL.md \
  .codex/skills/openspec-sync-specs/SKILL.md \
  .codex/skills/openspec-update-change/SKILL.md \
  openspec/changes/bootstrap-bangumi-collection-go-openspec/.openspec.yaml \
  openspec/changes/bootstrap-bangumi-collection-go-openspec/design.md \
  openspec/changes/bootstrap-bangumi-collection-go-openspec/proposal.md \
  openspec/changes/bootstrap-bangumi-collection-go-openspec/specs/external-openspec-bootstrap/spec.md \
  openspec/changes/bootstrap-bangumi-collection-go-openspec/tasks.md \
  openspec/config.yaml | LC_ALL=C sort)"
bootstrap_expected_pre_status="$({
  find .codex openspec -type f -print | while IFS= read -r bootstrap_path; do printf '?? %s\n' "$bootstrap_path"; done
  printf '%s\n' '?? CLAUDE.md' '?? note' '!! .claude/settings.local.json'
} | LC_ALL=C sort)"
test "$(git status --porcelain=v1 --untracked-files=all --ignored=matching | LC_ALL=C sort)" = "$bootstrap_expected_pre_status"
for bootstrap_skill in .codex/skills/openspec-*/SKILL.md; do
  test "$(grep -Fxc '  generatedBy: "1.6.0"' "$bootstrap_skill")" = "1"
done
while IFS= read -r bootstrap_path; do
  test -f "$bootstrap_path" && test ! -L "$bootstrap_path"
  iconv -f UTF-8 -t UTF-8 "$bootstrap_path" >/dev/null
  test "$(head -c 3 "$bootstrap_path" | od -An -tx1 | tr -d ' \n')" != "efbbbf"
  if LC_ALL=C grep -n $'\r' "$bootstrap_path"; then exit 1; fi
  if LC_ALL=C grep -n '[[:blank:]]$' "$bootstrap_path"; then exit 1; fi
  test "$(tail -c 1 "$bootstrap_path" | od -An -tx1 | tr -d ' \n')" = "0a"
done < <(find .codex openspec -type f -print | LC_ALL=C sort)
test "$(printf '%s\0' LICENSE README.md client.go collection.go errors.go example/example.go go.mod go.sum options.go types.go | xargs -0 shasum -a 256 | shasum -a 256 | awk '{print $1}')" = "7d4e55231e08c5d8286bffb921e53adbda9a5d82ced9c4688e4cadb129b1855e"
test -z "$(git diff --name-only 8173f44911360150a5a5a7c6418021d1014fe85b -- LICENSE README.md client.go collection.go errors.go example/example.go go.mod go.sum options.go types.go)"
test "$(shasum -a 256 CLAUDE.md | awk '{print $1}')" = "c69937414d643411e03f3d7982e4bf58e8ceb02802cf1338dada4ea6a630407d"
test "$(shasum -a 256 note | awk '{print $1}')" = "7b2376f605f5d9d9f677eaf1165a2ff911cec9092b6b2763b966a3d6f93e2e74"
test "$(shasum -a 256 .claude/settings.local.json | awk '{print $1}')" = "c67845ef6924b080344bc0010931336106d9df3a700f67f10a22763d07517e3e"
test -z "$(git ls-files --stage -- CLAUDE.md note .claude/settings.local.json)"
test "$(git status --porcelain=v1 --untracked-files=all --ignored=matching -- CLAUDE.md note .claude/settings.local.json | shasum -a 256 | awk '{print $1}')" = "221322afd4137d24875b5fbd369c9028e99f69f1cf155e3c49975996841a3938"
openspec status --change bootstrap-bangumi-collection-go-openspec --json
openspec validate bootstrap-bangumi-collection-go-openspec --strict
openspec validate --all --strict
openspec doctor --json
git diff --check
test -z "$(find .codex/skills/openspec-* openspec -type l -print)"
```

Before the second main-agent acceptance, run the exact complete staged-delta and synchronized-spec gate:

```zsh
bootstrap_expected_staged="$(printf 'A\t%s\n' \
  .codex/skills/openspec-apply-change/SKILL.md \
  .codex/skills/openspec-archive-change/SKILL.md \
  .codex/skills/openspec-explore/SKILL.md \
  .codex/skills/openspec-propose/SKILL.md \
  .codex/skills/openspec-sync-specs/SKILL.md \
  .codex/skills/openspec-update-change/SKILL.md \
  openspec/changes/archive/2026-07-24-bootstrap-bangumi-collection-go-openspec/.openspec.yaml \
  openspec/changes/archive/2026-07-24-bootstrap-bangumi-collection-go-openspec/design.md \
  openspec/changes/archive/2026-07-24-bootstrap-bangumi-collection-go-openspec/proposal.md \
  openspec/changes/archive/2026-07-24-bootstrap-bangumi-collection-go-openspec/specs/external-openspec-bootstrap/spec.md \
  openspec/changes/archive/2026-07-24-bootstrap-bangumi-collection-go-openspec/tasks.md \
  openspec/config.yaml \
  openspec/specs/external-openspec-bootstrap/spec.md | LC_ALL=C sort)"
test "$(git diff --cached --name-status | LC_ALL=C sort)" = "$bootstrap_expected_staged"
test "$(git branch --show-current)" = "codex/v0.1.0-hardening"
test "$(git rev-parse HEAD)" = "8173f44911360150a5a5a7c6418021d1014fe85b"
test "$(git rev-parse refs/heads/codex/v0.1.0-hardening)" = "8173f44911360150a5a5a7c6418021d1014fe85b"
test "$(git rev-parse refs/heads/main)" = "8173f44911360150a5a5a7c6418021d1014fe85b"
test "$(git rev-parse refs/remotes/origin/main)" = "8173f44911360150a5a5a7c6418021d1014fe85b"
git diff --quiet
bootstrap_expected_staged_status="$({
  printf '%s\n' "$bootstrap_expected_staged" | awk -F '\t' '{print "A  " $2}'
  printf '%s\n' '?? CLAUDE.md' '?? note' '!! .claude/settings.local.json'
} | LC_ALL=C sort)"
test "$(git status --porcelain=v1 --untracked-files=all --ignored=matching | LC_ALL=C sort)" = "$bootstrap_expected_staged_status"
test "$(shasum -a 256 CLAUDE.md | awk '{print $1}')" = "c69937414d643411e03f3d7982e4bf58e8ceb02802cf1338dada4ea6a630407d"
test "$(shasum -a 256 note | awk '{print $1}')" = "7b2376f605f5d9d9f677eaf1165a2ff911cec9092b6b2763b966a3d6f93e2e74"
test "$(shasum -a 256 .claude/settings.local.json | awk '{print $1}')" = "c67845ef6924b080344bc0010931336106d9df3a700f67f10a22763d07517e3e"
test -z "$(git ls-files --stage -- CLAUDE.md note .claude/settings.local.json)"
test ! -e openspec/changes/bootstrap-bangumi-collection-go-openspec
test "$(grep -Fxc 'Define the independent OpenSpec governance, baseline ancestry, protected local-file boundary, and validation gates for the external bangumi-collection-go hardening lane.' openspec/specs/external-openspec-bootstrap/spec.md)" = "1"
diff -u \
  <(sed -n '/^### Requirement:/,$p' openspec/changes/archive/2026-07-24-bootstrap-bangumi-collection-go-openspec/specs/external-openspec-bootstrap/spec.md) \
  <(sed -n '/^### Requirement:/,$p' openspec/specs/external-openspec-bootstrap/spec.md)
while IFS= read -r bootstrap_path; do
  test -f "$bootstrap_path" && test ! -L "$bootstrap_path"
  iconv -f UTF-8 -t UTF-8 "$bootstrap_path" >/dev/null
  test "$(head -c 3 "$bootstrap_path" | od -An -tx1 | tr -d ' \n')" != "efbbbf"
  if LC_ALL=C grep -n $'\r' "$bootstrap_path"; then exit 1; fi
  if LC_ALL=C grep -n '[[:blank:]]$' "$bootstrap_path"; then exit 1; fi
  test "$(tail -c 1 "$bootstrap_path" | od -An -tx1 | tr -d ' \n')" = "0a"
done < <(git diff --cached --name-only | LC_ALL=C sort)
openspec validate --all --strict
openspec doctor --json
git diff --cached --check
test "$(printf '%s\0' LICENSE README.md client.go collection.go errors.go example/example.go go.mod go.sum options.go types.go | xargs -0 shasum -a 256 | shasum -a 256 | awk '{print $1}')" = "7d4e55231e08c5d8286bffb921e53adbda9a5d82ced9c4688e4cadb129b1855e"
test "$(git status --porcelain=v1 --untracked-files=all --ignored=matching -- CLAUDE.md note .claude/settings.local.json | shasum -a 256 | awk '{print $1}')" = "221322afd4137d24875b5fbd369c9028e99f69f1cf155e3c49975996841a3938"
```

Only after that acceptance may the finalization owner commit and run this exact post-commit proof:

```zsh
bootstrap_expected_commit="$(printf 'A\t%s\n' \
  .codex/skills/openspec-apply-change/SKILL.md \
  .codex/skills/openspec-archive-change/SKILL.md \
  .codex/skills/openspec-explore/SKILL.md \
  .codex/skills/openspec-propose/SKILL.md \
  .codex/skills/openspec-sync-specs/SKILL.md \
  .codex/skills/openspec-update-change/SKILL.md \
  openspec/changes/archive/2026-07-24-bootstrap-bangumi-collection-go-openspec/.openspec.yaml \
  openspec/changes/archive/2026-07-24-bootstrap-bangumi-collection-go-openspec/design.md \
  openspec/changes/archive/2026-07-24-bootstrap-bangumi-collection-go-openspec/proposal.md \
  openspec/changes/archive/2026-07-24-bootstrap-bangumi-collection-go-openspec/specs/external-openspec-bootstrap/spec.md \
  openspec/changes/archive/2026-07-24-bootstrap-bangumi-collection-go-openspec/tasks.md \
  openspec/config.yaml \
  openspec/specs/external-openspec-bootstrap/spec.md | LC_ALL=C sort)"
test "$(git branch --show-current)" = "codex/v0.1.0-hardening"
test "$(git rev-list --parents -n 1 HEAD | awk '{print NF}')" = "2"
test "$(git rev-parse HEAD^)" = "8173f44911360150a5a5a7c6418021d1014fe85b"
test "$(git rev-parse refs/heads/codex/v0.1.0-hardening)" = "$(git rev-parse HEAD)"
test "$(git rev-parse refs/heads/main)" = "8173f44911360150a5a5a7c6418021d1014fe85b"
test "$(git rev-parse refs/remotes/origin/main)" = "8173f44911360150a5a5a7c6418021d1014fe85b"
test "$(git log -1 --format=%s HEAD)" = "chore: establish external openspec governance"
test "$(git diff-tree --no-commit-id --name-status -r HEAD | LC_ALL=C sort)" = "$bootstrap_expected_commit"
git diff --cached --quiet
git diff --quiet
test ! -e openspec/changes/bootstrap-bangumi-collection-go-openspec
openspec validate --all --strict
openspec doctor --json
test "$(printf '%s\0' LICENSE README.md client.go collection.go errors.go example/example.go go.mod go.sum options.go types.go | xargs -0 shasum -a 256 | shasum -a 256 | awk '{print $1}')" = "7d4e55231e08c5d8286bffb921e53adbda9a5d82ced9c4688e4cadb129b1855e"
test "$(shasum -a 256 CLAUDE.md | awk '{print $1}')" = "c69937414d643411e03f3d7982e4bf58e8ceb02802cf1338dada4ea6a630407d"
test "$(shasum -a 256 note | awk '{print $1}')" = "7b2376f605f5d9d9f677eaf1165a2ff911cec9092b6b2763b966a3d6f93e2e74"
test "$(shasum -a 256 .claude/settings.local.json | awk '{print $1}')" = "c67845ef6924b080344bc0010931336106d9df3a700f67f10a22763d07517e3e"
test -z "$(git ls-files --stage -- CLAUDE.md note .claude/settings.local.json)"
bootstrap_expected_post_status="$(printf '%s\n' '?? CLAUDE.md' '?? note' '!! .claude/settings.local.json' | LC_ALL=C sort)"
test "$(git status --porcelain=v1 --untracked-files=all --ignored=matching | LC_ALL=C sort)" = "$bootstrap_expected_post_status"
test "$(git status --porcelain=v1 --untracked-files=all --ignored=matching -- CLAUDE.md note .claude/settings.local.json | shasum -a 256 | awk '{print $1}')" = "221322afd4137d24875b5fbd369c9028e99f69f1cf155e3c49975996841a3938"
test -z "$(git branch -r --contains HEAD)"
test -z "$(git tag --points-at HEAD)"
```

Any mismatch stops without restoration or ref mutation. The final report records committed, pushed, released, and deployed separately.

## Goals / Non-Goals

**Goals:**

- Make this repository its own repo-local OpenSpec planning home.
- Preserve exact baseline ancestry and a dedicated local hardening branch.
- Retain only the generated root config, six generated Codex OpenSpec skills, and this bootstrap change as the planning candidate.
- Make user-owned local files and baseline product bytes explicit protected inputs.
- Define a narrow apply that verifies and accepts the already generated framework without implementing product behavior.
- Leave an auditable dependency gate for the later `harden-bangumi-collection-go-v0-1-0` change.

**Non-Goals:**

- Editing or testing Go product behavior, DTOs, pagination, rate limiting, retries, errors, ordering, CI, or Go 1.26.
- Modifying the `BangumiStaffStats` repository or importing its OpenSpec root.
- Re-running framework initialization during apply.
- Staging, committing, archiving, pushing, opening a pull request, tagging, releasing, or publishing during this planning candidate.
- Reading user files as hardening requirements, changing their contents or Git classification, or adding them to an allowlist.

## Decisions

### 1. The exact commit, not a moving branch name, is the ancestry authority

The only accepted bootstrap parent is `8173f44911360150a5a5a7c6418021d1014fe85b`. Preflight proves `main` and `origin/main` resolved to it, and the working branch must have it as an ordinary ancestor. During this planning candidate and its later validation-only apply, `HEAD` remains that baseline and no Git ref is mutable.

Using a moving `main`/`origin/main` name or fetching opportunistically was rejected because it would make the external-lane baseline depend on unrelated remote timing. Recreating or force-moving an existing target branch was rejected because it could destroy local work.

### 2. Bootstrap uses one repo-local OpenSpec root

`openspec/config.yaml` declares `schema: spec-driven` and repository-specific `context` plus artifact rules. The context fixes this repository as the external public-collection client lane, names the baseline/branch/protected paths, requires spec-before-apply and explicit owned paths, and separates local development commits from push/PR/tag/release authorization. Proposal/design/tasks rules require the complete boundary fields, exact paths/dependencies, executable acceptance, explicit status dimensions, and remote/publication non-goals. Default example comments without this governance are not an accepted config. This repository does not use a shared store, a nested planning root, or the main repository's OpenSpec. The persistent generated skill inventory is exactly:

- `.codex/skills/openspec-apply-change/SKILL.md`
- `.codex/skills/openspec-archive-change/SKILL.md`
- `.codex/skills/openspec-explore/SKILL.md`
- `.codex/skills/openspec-propose/SKILL.md`
- `.codex/skills/openspec-sync-specs/SKILL.md`
- `.codex/skills/openspec-update-change/SKILL.md`

Each skill must identify OpenSpec `generatedBy: "1.6.0"`. Before archive, the active bootstrap change is the only other OpenSpec content authorized by the candidate. After the approved archive action, the active change must be absent and the only additional OpenSpec content must be its exact date-stamped archive plus `openspec/specs/external-openspec-bootstrap/spec.md`. Codex global slash prompts generated by the initialization command are not repository deliverables and are not refreshed during apply.

A shared external store was rejected because this lane must be independently cloneable and reviewable. Hand-copying selected skill text was rejected because it loses generator provenance and creates drift.

### 3. Protection is content-, state-, and index-exact

The three user-owned paths are checked by SHA-256, by the exact porcelain states `??`, `??`, and `!!`, and by absence from `git ls-files` and the index. Their protected status projection must continue to hash to `221322afd4137d24875b5fbd369c9028e99f69f1cf155e3c49975996841a3938`.

Their presence is tolerated known state, not candidate output. Broad cleanup, reset, checkout restoration, stash, recursive deletion, `git add -A`, and any command that targets them are forbidden. Treating a clean worktree as the only valid precondition was rejected because it would either block on intentional user state or encourage destructive cleanup.

### 4. Product paths are a read-only complement

Every literal baseline-product path is protected. Acceptance requires no worktree/index difference for those ten paths relative to `8173f44911360150a5a5a7c6418021d1014fe85b`, an empty global index during apply, the unchanged literal-list product seal, and a physical/status path audit showing that candidate files occur only under `.codex/skills/openspec-*/` and `openspec/`. The permitted tasks-checkbox diff is not misclassified as a product diff.

The bootstrap apply may change only checkbox markers in this change's `tasks.md`; it accepts the already generated framework bytes rather than regenerating them. Any extra `.codex` content, product change, generated Go file, module-file change, or new path outside the exact candidate envelope fails closed.

### 5. Delegation and lifecycle states remain explicit

All filesystem and Git-branch mutations are delegated subagent work. The main agent may review and amend OpenSpec planning artifacts and performs read-only acceptance. The current handoff is deliberately unstaged and uncommitted.

Bootstrap acceptance records investigated, specified, framework-generated, and locally verified separately from committed, pushed, released, or deployed. After the main agent accepts the unstaged validation-only apply candidate, one finalization subagent reads and follows `openspec-archive-change`, archives/synchronizes only this change into `openspec/changes/archive/2026-07-24-bootstrap-bangumi-collection-go-openspec/**`, and stages only the accepted framework plus exact archive/root-spec output. The synchronized root capability SHALL replace any generated `TBD` Purpose with exactly `Define the independent OpenSpec governance, baseline ancestry, protected local-file boundary, and validation gates for the external bangumi-collection-go hardening lane.`; its requirements remain byte-identical to the archived delta requirements. The subagent then stops for a second main-agent read-only acceptance. Only that acceptance authorizes one local commit with exact subject `chore: establish external openspec governance`; its sole parent is the baseline, its delta is exactly 13 added files, it contains no protected/product path, and pushed/released/deployed remain false. Hardening implementation requires its own approved OpenSpec change and cannot start merely because bootstrap commits.

## Risks / Trade-offs

- **[Risk] Intentional user files look like bootstrap dirt** → Audit an exact protected-path projection and keep it separate from the approved OpenSpec candidate envelope.
- **[Risk] A future `openspec init` refresh changes framework bytes or global prompts** → Do not rerun init during apply; validate the sealed generated inventory and generator version already present.
- **[Risk] A moving remote branch changes the hardening base** → Pin the commit and verify ordinary ancestry without fetching or rewriting refs.
- **[Risk] Empty framework directories are not represented by Git** → Treat `openspec/config.yaml`, the active change, and the generated skills as the persistent root; let OpenSpec create needed directories when lifecycle commands later require them.
- **[Risk] Bootstrap approval is mistaken for product readiness** → Keep all hardening behavior out of this capability and require a separate dependent change.
- **[Risk] Framework generation necessarily preceded its own change** → Treat it as a one-time uncommitted control-plane bootstrap only; approve the complete change before validation-only apply or any local commit, and allow no equivalent exception for product work.

## Migration Plan

1. Verify the exact repository, baseline refs, target branch, tracked/index state, protected hashes, and protected status projection.
2. Verify the already generated config, six skills, active change inventory, and baseline product seal without rerunning init.
3. Complete only this change's validation tasks and keep all output unstaged.
4. Run strict OpenSpec validation, doctor, diff/format/path/index checks, and obtain independent main-agent acceptance.
5. A delegated finalization subagent archives/synchronizes only this change, stages the exact accepted framework/archive/root-spec union, reruns all cached/path/seal checks, and stops for second main-agent acceptance.
6. After that acceptance, the same subagent creates the exact-subject local commit and proves its baseline parent, exact no-product/no-protected delta, clean index/worktree apart from the three protected user paths, and local-only status. Push/PR/tag/release/publication remain separately unauthorized.

Rollback is preservation-first: on any mismatch, stop and report the exact path or seal. Do not reset, clean, stash, checkout, delete, or rewrite user/product state. A later approved rollback may remove only an exact reviewed set of OpenSpec-owned candidate paths.

## Open Questions

None. Product hardening paths and behavior will be decided by the dependent hardening change.
