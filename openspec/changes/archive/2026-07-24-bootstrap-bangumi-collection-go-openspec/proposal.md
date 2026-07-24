## Why

`bangumi-collection-go` needs its own reviewable OpenSpec root before the separately approved v0.1.0 hardening work can begin. The repository currently has only product code and user-owned local files, so the external lane needs an exact baseline, branch, writable-path, and protection contract that prevents bootstrap work from being confused with product implementation.

## What Changes

- Establish this repository as an independent, repo-local OpenSpec planning home using the `spec-driven` schema and Codex skills.
- Govern the hardening branch as `codex/v0.1.0-hardening`, created from exact baseline `8173f44911360150a5a5a7c6418021d1014fe85b` with ordinary ancestry to that commit.
- Seal and protect the existing untracked `CLAUDE.md`, untracked `note`, and ignored `.claude/settings.local.json` by exact content hash and Git status.
- Define a path-exact bootstrap apply that may only retain the already generated OpenSpec root and Codex OpenSpec skills, complete this change's task checkboxes, and run strict validation.
- Explicitly exclude Go product changes, v0.1.0 hardening implementation, staging/committing protected files, remote mutation, and push/PR/tag/release activity.

## Capabilities

### New Capabilities

- `external-openspec-bootstrap`: Governs the independent repository's OpenSpec root, exact baseline and branch ancestry, generated framework inventory, protected local-file seals, path boundaries, validation gates, and handoff to the later hardening change.

### Modified Capabilities

None.

## Impact

- Repository scope: only `/Users/luca/dev/bangumi-collection-go`; the `BangumiStaffStats` repository remains read-only external planning evidence and is not modified by this change.
- Generated planning/framework scope: `openspec/**` and the six `.codex/skills/openspec-*/SKILL.md` files created by OpenSpec 1.6.0 with `--tools codex --profile core`.
- Protected local state: `CLAUDE.md`, `note`, and `.claude/settings.local.json` must keep their approved bytes, untracked/ignored classification, and absence from the Git index.
- Product/API/dependency impact: none. Existing tracked Go source, examples, module files, README, and license remain byte-identical to the baseline.
- Remote impact: none. This change does not fetch, push, open a pull request, tag, release, or publish a Go module version.

| Field | Declaration |
|---|---|
| Status | investigated: complete; specified: pending strict validation and main-agent approval; framework-generated: planning candidate only; implemented: no; verified: planning-only; committed: no; pushed: no; released: no; deployed: no |
| Owner | External OpenSpec bootstrap subagent owns initialization/apply/finalization; the main agent may amend OpenSpec artifacts and performs read-only approval/acceptance |
| Writable / owned paths | Planning/framework candidate: `openspec/config.yaml`; `.codex/skills/openspec-apply-change/SKILL.md`; `.codex/skills/openspec-archive-change/SKILL.md`; `.codex/skills/openspec-explore/SKILL.md`; `.codex/skills/openspec-propose/SKILL.md`; `.codex/skills/openspec-sync-specs/SKILL.md`; `.codex/skills/openspec-update-change/SKILL.md`; and `openspec/changes/bootstrap-bangumi-collection-go-openspec/**`. Apply: checkbox markers only in that change's `tasks.md`. Finalization: the accepted framework files plus the exact archive move and synchronized `openspec/specs/external-openspec-bootstrap/spec.md` |
| Read-only protected inputs | Baseline commit `8173f44911360150a5a5a7c6418021d1014fe85b`; exact tracked product paths `LICENSE`, `README.md`, `client.go`, `collection.go`, `errors.go`, `example/example.go`, `go.mod`, `go.sum`, `options.go`, and `types.go`; `refs/heads/main`; `refs/remotes/origin/main`; `CLAUDE.md`; `note`; `.claude/settings.local.json`; `/Users/luca/dev/BangumiStaffStats/tmp-formal-development/formal-development-master-plan.md` external-lane section |
| Consumes | Fixed external baseline, approved protected-file/status seals, and the formal master plan's external-lane governance |
| Produces | One independent repo-local OpenSpec root, six generated Codex OpenSpec skills, synchronized root capability `external-openspec-bootstrap`, archived bootstrap evidence, and one local bootstrap-governance commit |
| Dependencies | None; this is the external lane's bootstrap prerequisite |
| Deliverables | Exact framework/config inventory, protected and baseline seals, strict validation evidence, accepted archived capability, and clean local commit `chore: establish external openspec governance` |
| Acceptance | Exact branch/ancestry/path/hash/index gates; governance-rich config; strict change/all validation; doctor; byte/symlink/diff checks; main-agent acceptance of the unstaged apply candidate and then the exact staged archive candidate |
| Non-goals | Any Go product change/test, v0.1.0 hardening, dependency update, fetch, push, PR, tag, release, publication, deployment, or main-repository mutation |
| Operations deferred | All remote mutation and public `v0.1.0` publication; any production/deployment activity |
| Mutable refs | Branch creation from the exact baseline is the already completed delegated bootstrap preflight. Planning/apply keep every ref immutable. Only after the second main-agent acceptance may one finalization subagent advance `refs/heads/codex/v0.1.0-hardening` once with the exact local commit; no other ref may move |
| Stop / rollback conditions | Stop on root/branch/HEAD/ancestry, protected seal/state, baseline-product seal, candidate inventory, approval, validation, index, staged-delta, or archive-output mismatch. Preserve all state; never reset, clean, stash, checkout-restore, broadly delete, amend/rebase, or touch a protected/product path |

### Executable acceptance

Run every block from `/Users/luca/dev/bangumi-collection-go` with `zsh`; every comparison is exact and every command must exit zero.

Before archive/staging, prove the fixed baseline, empty index, exact 12-file active candidate, unchanged product/protected state, and healthy OpenSpec:

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
  test -f "$bootstrap_path"
  test ! -L "$bootstrap_path"
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

After archive/sync and before the second main-agent acceptance, prove the complete staged delta is exactly 13 additions—not merely 13 added paths alongside other staged changes—and prove the synchronized requirement bodies are identical:

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
while IFS= read -r bootstrap_path; do
  test -f "$bootstrap_path"
  test ! -L "$bootstrap_path"
  iconv -f UTF-8 -t UTF-8 "$bootstrap_path" >/dev/null
  test "$(head -c 3 "$bootstrap_path" | od -An -tx1 | tr -d ' \n')" != "efbbbf"
  if LC_ALL=C grep -n $'\r' "$bootstrap_path"; then exit 1; fi
  if LC_ALL=C grep -n '[[:blank:]]$' "$bootstrap_path"; then exit 1; fi
  test "$(tail -c 1 "$bootstrap_path" | od -An -tx1 | tr -d ' \n')" = "0a"
done < <(git diff --cached --name-only | LC_ALL=C sort)
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
test -f openspec/changes/archive/2026-07-24-bootstrap-bangumi-collection-go-openspec/specs/external-openspec-bootstrap/spec.md
test -f openspec/specs/external-openspec-bootstrap/spec.md
test "$(grep -Fxc 'Define the independent OpenSpec governance, baseline ancestry, protected local-file boundary, and validation gates for the external bangumi-collection-go hardening lane.' openspec/specs/external-openspec-bootstrap/spec.md)" = "1"
diff -u \
  <(sed -n '/^### Requirement:/,$p' openspec/changes/archive/2026-07-24-bootstrap-bangumi-collection-go-openspec/specs/external-openspec-bootstrap/spec.md) \
  <(sed -n '/^### Requirement:/,$p' openspec/specs/external-openspec-bootstrap/spec.md)
openspec validate --all --strict
openspec doctor --json
git diff --cached --check
test "$(printf '%s\0' LICENSE README.md client.go collection.go errors.go example/example.go go.mod go.sum options.go types.go | xargs -0 shasum -a 256 | shasum -a 256 | awk '{print $1}')" = "7d4e55231e08c5d8286bffb921e53adbda9a5d82ced9c4688e4cadb129b1855e"
test "$(git status --porcelain=v1 --untracked-files=all --ignored=matching -- CLAUDE.md note .claude/settings.local.json | shasum -a 256 | awk '{print $1}')" = "221322afd4137d24875b5fbd369c9028e99f69f1cf155e3c49975996841a3938"
```

Only after that second acceptance may the finalization subagent commit. Its post-commit proof is fully mechanical:

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

The report records `committed: yes`, `pushed: no`, `released: no`, and `deployed: no` separately. Any mismatch in any block is the stop gate.
