## Task boundary

| Field | Declaration |
|---|---|
| Status | investigated: complete; specified: pending approval; framework-generated: planning candidate only; implemented: no; verified: planning-only; committed: no; pushed: no; released: no; deployed: no |
| Owner | Delegated external bootstrap apply/finalization subagents; main agent amends OpenSpec only and performs read-only approval/acceptance |
| Writable / owned paths | Apply: checkbox markers in this file only. Finalization after acceptance: exact accepted `.codex/skills/openspec-*/SKILL.md`, `openspec/config.yaml`, active-change archive move into `openspec/changes/archive/2026-07-24-bootstrap-bangumi-collection-go-openspec/**`, and `openspec/specs/external-openspec-bootstrap/spec.md` |
| Read-only protected inputs | Exact baseline commit and ten literal product paths; `main`/`origin/main`; `CLAUDE.md`; `note`; `.claude/settings.local.json`; main-repository external-lane master-plan section; all candidate bytes outside checkbox markers during apply |
| Consumes | Main-agent-approved framework/change seals, fixed baseline, protected status/content seals |
| Produces | Validated bootstrap candidate; synchronized/archived capability; exact local governance commit |
| Dependencies | None |
| Deliverables | Complete tasks, two acceptance handoffs, strict-valid root/archive, local commit `chore: establish external openspec governance` |
| Acceptance | Exact commands below; empty index during apply; main acceptance before archive/staging and again before commit; exact parent/subject/delta/clean-state proof |
| Non-goals | Product/hardening code or tests, Go/module changes, remote/publication activity, main-repository mutation |
| Operations deferred | Fetch/push/PR/tag/release/module publication/deploy and all host/service work |
| Mutable refs | Apply: none. After second main-agent acceptance, finalization may advance only `refs/heads/codex/v0.1.0-hardening` once from baseline via the exact local commit |
| Stop / rollback conditions | Stop on root/branch/HEAD/ancestry/approval/seal/path/index/validation/archive/staged-delta mismatch; no reset/clean/stash/checkout restore/broad delete/history rewrite |

## Executable acceptance commands

Run from `/Users/luca/dev/bangumi-collection-go` with `zsh`. Before archive, this exact block implements tasks 1–4 and explicitly enumerates the six owned skill paths plus five active-change artifacts and config:

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
test -z "$(git diff --cached --name-only 8173f44911360150a5a5a7c6418021d1014fe85b -- LICENSE README.md client.go collection.go errors.go example/example.go go.mod go.sum options.go types.go)"
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
test -z "$(find .codex openspec -type l -print)"
```

After archive/sync, tasks 5.2–5.4 use this complete staged-name-status gate; no additional staged modification, deletion, or rename can pass:

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

Only after the second acceptance, the post-commit protocol runs these exact proof commands:

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

## 1. Re-prove the external baseline and protection seals

- [x] 1.1 From canonical root `/Users/luca/dev/bangumi-collection-go`, verify branch `codex/v0.1.0-hardening`, exact `HEAD` `8173f44911360150a5a5a7c6418021d1014fe85b`, matching local/remote-tracking main refs, ordinary ancestry, and an empty index without fetching or moving any ref.
- [x] 1.2 Verify the exact SHA-256 values for `CLAUDE.md`, `note`, and `.claude/settings.local.json`; require their three-path porcelain projection to remain `??`, `??`, `!!` with hash `221322afd4137d24875b5fbd369c9028e99f69f1cf155e3c49975996841a3938`, preserve the global-ignore evidence for the settings file, and prove none has index membership.
- [x] 1.3 Compare only the literal product paths `LICENSE`, `README.md`, `client.go`, `collection.go`, `errors.go`, `example/example.go`, `go.mod`, `go.sum`, `options.go`, and `types.go` against baseline `8173f44911360150a5a5a7c6418021d1014fe85b` in both worktree and index. Recompute their ordered aggregate with `printf '%s\0' LICENSE README.md client.go collection.go errors.go example/example.go go.mod go.sum options.go types.go | xargs -0 shasum -a 256 | shasum -a 256` and require `7d4e55231e08c5d8286bffb921e53adbda9a5d82ced9c4688e4cadb129b1855e`; do not use current `git ls-files` after framework paths become tracked. Stop without reset, clean, stash, checkout restoration, deletion, or overwrite on any mismatch.

## 2. Accept the already generated OpenSpec framework

- [x] 2.1 Verify `openspec/config.yaml` declares the repo-local `spec-driven` schema, repository-specific context, and proposal/design/tasks rules required by the spec; resolves this repository as the only planning root; has no store pointer; and exposes only active change `bootstrap-bangumi-collection-go-openspec`. Reject the default commented example config as incomplete.
- [x] 2.2 Enumerate exactly the six approved `.codex/skills/openspec-*/SKILL.md` files, reject symlinks or extra `.codex` content, and verify every skill declares `generatedBy: "1.6.0"`.
- [x] 2.3 Verify proposal, design, `external-openspec-bootstrap` spec, tasks, and `.openspec.yaml` are the exact main-agent-approved bytes. Do not rerun `openspec init`, refresh global prompts, rewrite config/skills, or create any product/hardening file.

## 3. Run strict, path-exact bootstrap validation

- [x] 3.1 Run `openspec status --change bootstrap-bangumi-collection-go-openspec --json`, `openspec validate bootstrap-bangumi-collection-go-openspec --strict`, `openspec validate --all --strict`, and `openspec doctor --json`; require artifact completeness, healthy root status, and zero validation failure.
- [x] 3.2 Run trailing-whitespace, UTF-8 BOM, CR, final-LF, symlink, and `git diff --check` audits across `.codex/skills/openspec-*/**` and `openspec/**`; require all checks to pass without formatter or generator writes.
- [x] 3.3 Enumerate physical and Git-status paths and prove the candidate envelope contains only `openspec/config.yaml`, the six generated skill files, and `openspec/changes/bootstrap-bangumi-collection-go-openspec/**`, alongside the three separately protected user paths; require no tracked product diff, ignored candidate residue, staged path, or ref mutation.

## 4. Stop at the delegated acceptance gate

- [x] 4.1 During apply, change only checkbox markers in this tasks file, then re-run every baseline, framework, strict, path, product, protected-file, and empty-index gate against the approved seals.
- [x] 4.2 Produce exact SHA-256 inventories for the generated framework and active change, record investigated/specified/framework-generated/locally-verified status, and explicitly record committed/pushed/released/deployed as false.
- [x] 4.3 Stop with all candidate output unstaged, uncommitted, and unarchived for the first main-agent read-only acceptance; do not run `git add`, commit, archive, fetch, push, PR, tag, release, publish, or modify the `BangumiStaffStats` repository.
- [x] 4.4 Treat the later `harden-bangumi-collection-go-v0-1-0` proposal as a separate post-acceptance change with its own writable paths and tests; do not implement or apply hardening in this bootstrap.

## 5. Archive and local commit after two acceptance gates

- [x] 5.1 Only after the main agent accepts the exact unstaged apply candidate, hand off its path/hash inventory and baseline/protected seals to one quiescent finalization subagent. That subagent SHALL first read and follow the repository's `openspec-archive-change` skill, re-prove the branch/baseline/index/product/protection gates, and archive/synchronize only `bootstrap-bangumi-collection-go-openspec` into `openspec/changes/archive/2026-07-24-bootstrap-bangumi-collection-go-openspec/**`.
- [x] 5.2 Replace any generated root-spec `TBD` Purpose with exactly `Define the independent OpenSpec governance, baseline ancestry, protected local-file boundary, and validation gates for the external bangumi-collection-go hardening lane.`, prove its requirements byte-identical to the archived delta requirements, then stage only the six accepted generated skill files, `openspec/config.yaml`, five archived change files, and `openspec/specs/external-openspec-bootstrap/spec.md`: exactly 13 added paths. Never use `git add -A`; reject every product/protected/extra path, rename outside the archive move, or byte drift.
- [x] 5.3 Run strict all validation, doctor, cached diff check, exact staged-path/hash inventory, baseline-product comparison, protected content/status/index seals, and prove `HEAD` remains `8173f44911360150a5a5a7c6418021d1014fe85b`. Stop with the exact archive candidate staged for a second main-agent read-only acceptance; do not commit before it.
- [x] 5.4 Before the second acceptance, prove commit readiness without moving a ref: the only staged bytes are the exact accepted 13-path candidate, `HEAD` is still the baseline, the intended sole parent/subject/delta are mechanically fixed, all strict/doctor/product/protection gates pass, and the exact post-commit proof commands in the executable-acceptance section above are recorded for the finalization handoff.

### Post-commit protocol (not a task checkbox)

Only after the second main-agent acceptance may the quiescent finalization subagent create one local commit with exact subject `chore: establish external openspec governance`. It then proves the sole parent is `8173f44911360150a5a5a7c6418021d1014fe85b`; the delta is exactly the accepted 13 added framework/archive/root-spec paths; the active change is absent; strict validation and doctor pass; the index is empty; only `CLAUDE.md`, `note`, and ignored `.claude/settings.local.json` remain as tolerated local status; and the product/protection seals are unchanged. It records committed: yes; pushed: no; released: no; deployed: no, then stops. It must not amend, rebase, fetch, push, open a PR, tag, release, publish, or start hardening apply.
