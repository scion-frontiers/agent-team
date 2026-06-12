---
name: pr-code-review
description: >-
  Perform a rigorous code review of a single pull request, focused exclusively on the
  delta introduced by the PR. Covers reviewer mindset, review categories, output format,
  and container-friendly git tips for getting a clean PR diff. Use when reviewing a PR,
  running a presubmit review, or auditing a branch before merging.
---

# PR Code Review

## Reviewer Mindset

**Review only the delta introduced by the PR.** Do not comment on technical debt, style issues, or architectural flaws in lines that were not modified by this branch.

## Priorities

1. **Logic & correctness.** Edge cases, race conditions, off-by-one errors, incorrect handling of error states. Does the code achieve its stated intent?
2. **Architecture & patterns.** Does the change align with existing project structure, idioms, and conventions?
3. **Security.** New attack surface from the change — unsanitized inputs, insecure concurrency handling, credential exposure, missing authz checks.
4. **Efficiency.** O(n) regressions or unnecessary allocations introduced by the new code paths. Do not litigate efficiency of code that wasn't touched.

## Constraints

- **Ignore technical debt** outside the diff.
- **Verify dependencies but don't review them.** If the new code calls into an existing function, verify the signature to confirm the call is valid — but do not review the existing function itself.
- **Tone.** Objective, technically precise, with "Suggested Fix" code blocks where they help.

## Output Format

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
<Breaking changes? Removed fields? New required fields?>

## Final Verdict
APPROVE | REWORK
```

Findings should be severity-tagged so the dispatching coordinator can decide which require a fix-cycle.

## Getting the PR Diff in a Container

These tips address friction in a fresh container where local refs are sparse.

### Prefer `gh` over raw `git` for the PR delta

`gh` computes the merge-base server-side:

```bash
gh pr diff <num>
gh pr view <num> --json files,additions,deletions,baseRefName,headRefName,body,author,state
```

### Never assume `main` exists locally

A freshly checked-out PR via `gh pr checkout` populates only the PR branch:

```bash
git fetch origin main:main
git diff origin/main...HEAD
```

### Compute the merge-base explicitly

Don't trust `A..B` to mean "what changed in this PR":

```bash
BASE=$(git merge-base origin/main HEAD)
git log --oneline "$BASE..HEAD"
git diff "$BASE...HEAD"
```

### Cap potentially-huge git output

```bash
git log --oneline --max-count=30 "$BASE..HEAD"
```

### Verify the PR's actual base ref

Some PRs target branches other than `main`:

```bash
BASE_REF=$(gh pr view <num> --json baseRefName --jq .baseRefName)
git fetch origin "$BASE_REF"
BASE=$(git merge-base "origin/$BASE_REF" HEAD)
```

### One-liner to start any PR review

```bash
gh pr checkout <num> && gh pr diff <num> > /tmp/pr.diff && wc -l /tmp/pr.diff
```
