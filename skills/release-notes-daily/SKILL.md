---
name: release-notes-daily
description: >-
  Generate high-quality, semantic release notes from git history for a single day or short
  date range. Synthesizes raw commits into Features / Fixes / Breaking Changes, ranked by
  significance, with noise filtered out. Use when writing release notes for one day or a
  small range. For large multi-day backfills, use changelog-parallel-backfill instead.
---

# Release Notes — Daily / Short-Range

You are an expert Release Engineer and analytical git historian. Your job is to read a
day's commits and produce a single coherent release-notes file that synthesizes the
underlying changes — not a raw commit dump.

## Inputs

- A specific date (`YYYY-MM-DD`) or a small range (a few days).
- Pacific Time (PT) is the default timezone unless told otherwise.
- **Omit "today"** — only fully completed calendar days are eligible.

## Workflow

### 1. Extract commits for the day

```bash
git log \
  --since="YYYY-MM-DD 00:00:00 -0800" \
  --until="YYYY-MM-DD 23:59:59 -0800" \
  --pretty=format:"%h - %an - %s"
```

Adjust the timezone offset if not Pacific.

### 2. Identify significance

Scan the commit messages. Look for:

- Major features (new APIs, new subsystems)
- Security-relevant changes
- Breaking changes: keywords like `BREAKING`, "Drop support", "Migrate", "Remove", major version bumps in package files
- Vague messages ("Fixed the thing", "Updated DB") that hide larger changes

For anything that looks significant or vague, run `git show --stat <hash>` (or `git show <hash>` for the full diff) to deduce what actually changed. **Don't guess.**

### 3. Synthesize

Group related commits — e.g. multiple "wip", "fixup", "more X" commits for the same
feature collapse into one bullet. **Filter out the noise:** chores, minor typos,
dependency bumps (unless major), non-functional refactors, CI-only changes.

Rank items within each section by significance, most impactful first.

### 4. Write the file

**Path:** `changelog/YYYY-MM-DD-changelog.md`

**Format:**

```markdown
# Release Notes (Month Day, Year)

<1–2 sentence summary describing the day's theme, e.g.
 "This day focused on hardening the agent dispatch path and resolving fork-PR review feedback.">

## ⚠️ BREAKING CHANGES
* **[Component]:** Description of the breaking change and what users/developers need to do to migrate.

## 🚀 Features
* **[High Impact Feature]:** Synthesized description of the feature (consolidated from commits X, Y, Z).
* **[Medium Impact Feature]:** Synthesized description...

## 🐛 Fixes
* **[High Impact Fix]:** Description of the critical bug resolved.
* **[Lower Impact Fix]:** ...
```

Omit empty sections. If there were no breaking changes, drop the `⚠️ BREAKING CHANGES` heading entirely.

## Constraints

- **Read-only by default.** Do not modify code, commit history, or tags unless explicitly
  asked. The only write you should make is the changelog file itself.
- **Precision with git.** Always translate vague user dates ("last week", "since v2.0")
  into exact git-compatible date strings or tags before querying.
- **Security.** Never quote secrets, API keys, or sensitive payloads that might be
  embedded in commit messages or diffs.
- **Brief explanations.** Before running git commands, say in one sentence what you're
  about to do (e.g. "Pulling commits for 2026-06-05 in PT.").

## Handling Overwhelming Volumes

If `git log` returns 500+ commits, don't process them raw. Adapt:

- `git shortlog --since=... --until=...` for a grouped overview
- `git log --merges` to focus on merge commits only
- Ask the user whether to focus on a specific directory or to summarize at a higher level

## Common Pitfalls

- Forgetting the timezone offset — the same commit can fall on different "days" depending on offset.
- Letting raw commit messages through as bullets — synthesize, don't dump.
- Treating a "version bump" as a chore — it can be a breaking-change signal.
