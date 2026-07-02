---
name: pr-code-review
description: >-
  Perform a rigorous code review of a single pull request, focused exclusively on the
  delta introduced by the PR. Covers reviewer mindset, review categories, output format,
  and the containerized-git tips needed to get a clean PR diff in a fresh agent workspace.
  Use when reviewing a PR, running a presubmit review, or auditing a branch before merging.
---

# PR Code Review

## Reviewer Mindset

You are a Senior Staff Software Engineer and Security Researcher performing a rigorous
review of a single pull request. **Review only the delta introduced by the PR.** Do not
comment on technical debt, style issues, or architectural flaws in lines that were not
modified by this branch.

## Priorities

1. **Logic & correctness.** Edge cases, race conditions, off-by-one errors, incorrect
   handling of error states. Does the code achieve its stated intent?
2. **Architecture & patterns.** Does the change align with existing project structure,
   idioms, and dependency-injection patterns?
3. **Security.** New attack surface from the change — unsanitized inputs, insecure
   concurrency handling, credential exposure, missing authz checks.
4. **Efficiency.** O(n) regressions or unnecessary allocations introduced by the new code
   paths. Do not litigate efficiency of code that wasn't touched.

## Constraints

- **Ignore technical debt** outside the diff.
- **Verify dependencies but don't review them.** If the new code calls into an existing
  function, verify the signature in the base branch to confirm the call is valid — but do
  not review the existing function itself.
- **Tone.** Objective, technically precise, with "Suggested Fix" code blocks where they help.

## Output Format

Write the review to the project's designated review output directory with this structure:

```markdown
# PR <N>: <title> — Review

## Executive Summary
<2 sentences. State the risk level: LOW / MEDIUM / HIGH / CRITICAL.>

## Critical Issues
<Blocking bugs or security flaws. If none: "None.">

## High
<Significant issues to fix before merge.>

## Medium
<Improvements for readability, robustness, or performance.>

## Low
<Minor suggestions, nits.>

## Positive Feedback
<Well-implemented logic, clever optimizations.>

## Test Coverage
<Are new code paths covered? Gaps?>

## Backward Compatibility
<Wire-format changes? Removed fields? New required fields without omitempty?>

## Final Verdict
APPROVE | REWORK
```

Findings should be severity-tagged so the dispatching coordinator can
decide which require a fix-cycle before submission.

---

## Getting the PR Diff in a Container

These tips address friction in a fresh container where the repo was checked out via
`gh pr checkout` and local refs are sparse.

### Prefer `gh` over raw `git` for the PR delta

`gh` computes the merge-base server-side and returns only the PR's hunks regardless of
local-ref state:

```bash
gh pr diff <num>
gh pr view <num> --json files,additions,deletions,baseRefName,headRefName,body,author,state
```

This sidesteps every issue below and is the recommended first move.

### Never assume `main` exists locally

A freshly checked-out PR via `gh pr checkout` populates only the PR branch. `main` is
`origin/main`, not `main`. Either materialize it or reference the remote:

```bash
git fetch origin main:main         # create local main tracking origin/main
git diff origin/main...HEAD        # or reference the remote ref directly
```

A naïve `git log main..HEAD` silently returns zero lines when `main` is unknown.

### Compute the merge-base explicitly

Don't trust `A..B` to mean "what changed in this PR" — on a stale branch it will include
upstream commits since divergence. Compute the merge-base:

```bash
BASE=$(git merge-base origin/main HEAD)
git log --oneline "$BASE..HEAD"
git diff "$BASE...HEAD"
```

Dot-syntax reminder:
- `git diff A...B` — changes on B's side since divergence (the PR-shaped diff).
- `git log A..B` — commits reachable from B but not A (the PR's commit list).

### Cap potentially-huge `git` output

Long-lived branches can produce 100KB+ of `git log` output. Always bound it:

```bash
git log --oneline --max-count=30 "$BASE..HEAD"
git log --oneline "$BASE..HEAD" -- web/                # scope by path
```

### Suspect "empty output" before re-running

If `git log <range>` returns zero lines, the first hypothesis should be "the revision
spec doesn't resolve the way I think," not "there's nothing to show":

```bash
git rev-parse main 2>&1
git rev-parse origin/main 2>&1
```

### Verify the PR's actual base ref

Some PRs target release branches or feature integration branches, not `main`:

```bash
BASE_REF=$(gh pr view <num> --json baseRefName --jq .baseRefName)
git fetch origin "$BASE_REF"
BASE=$(git merge-base "origin/$BASE_REF" HEAD)
```

### One-liner to start any PR review

```bash
gh pr checkout <num> && gh pr diff <num> > /tmp/pr.diff && wc -l /tmp/pr.diff
```

This gives an immediately reviewable artifact (`/tmp/pr.diff`) without any of the
local-ref pitfalls.
