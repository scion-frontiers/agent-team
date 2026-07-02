---
name: changelog-parallel-backfill
description: >-
  Fill a large gap of missing daily changelog entries (e.g. several weeks) by orchestrating
  a pool of release-notes sub-agents in parallel. Maintains a steady concurrency, replenishes
  the pool as agents finish, and handles stalled agents. Use when the gap is more than a few
  days — for one-off days use the release-notes-daily skill directly.
---

# Changelog — Parallel Backfill

When the missing-changelog window is large (a week or more), running release-notes
generation serially is too slow. This skill orchestrates a **pool of sub-agents**, one
per day, with steady concurrency.

## Workflow Overview

1. **Sync main** so the local `main` is up-to-date with the remote.
2. **Gap analysis** — identify the missing dates.
3. **Branching** — create a dedicated branch (e.g. `changelog-update-june`).
4. **Parallel orchestration** — launch a pool of `release-notes` sub-agents (default pool size: 10).
5. **Shared workspace** — use the `--workspace` flag pointing to the project root so agents can write directly to the shared filesystem.
6. **Asynchronous management** — monitor for completion and replenish the pool as slots open.
7. **Verification & finalize** — confirm files exist, then commit and push.

## Procedural Steps

### 1. Gap analysis

Identify the latest existing changelog file:

```bash
ls -v changelog/ | tail -n 1
```

Then enumerate missing dates between that file's date and the target end date.
**Always omit "today"** — only fully-completed calendar days are eligible. All date math
is Pacific Time (PT) unless specified otherwise.

### 2. Launch agents

Use unique names (e.g. `rn-YYYY-MM-DD`) and an explicit task prompt naming the **full
absolute path** of the target file:

```bash
scion start "rn-2026-06-01" \
  "Generate release notes for 2026-06-01 and write to <project-root>/changelog/2026-06-01-changelog.md" \
  --type release-notes \
  --workspace <project-root> \
  --non-interactive
```

The `--workspace` flag is **mandatory** — without it the isolated agent can't
write to the shared filesystem and you'll lose its output.

### 3. Maintain the pool

Target concurrency (e.g. 10 agents at a time). On completion:

- **Completion signal:** `state-change` notification — `agent:<name> has reached a state of COMPLETED`.
- **Replenishment:** Delete the finished agent and launch the next date.

```bash
scion delete <completed-agent-name> --non-interactive
```

**Don't wait for a full batch of 10 to finish before launching more.** Replenish as soon
as *any* agent completes — that's the whole point of the pool.

### 4. Handle stalled agents

If `scion list` shows an agent in `waiting_for_input`:

```bash
scion look <agent-name> --non-interactive   # inspect the screen
scion message <agent-name> "Yes" --non-interactive   # respond if needed
```

For other failure modes (LIMITS_EXCEEDED, transient API errors, container crashes), see the agent-recovery skill.

### 5. Finalize the branch

Once all dates are covered:

1. Verify all files are present under the project's changelog directory.
2. Stage the new files.
3. Commit with a descriptive message (e.g. `changelog: backfill 2026-06-01 to 2026-06-20`).
4. Push the branch to the remote.

## Best Practices

- **Task specificity** — the prompt must include the **full absolute path** to the target file.
  Relative paths produce surprising results once the agent runs inside its sandbox.
- **Cleanup** — always delete sub-agents after their work is verified. Otherwise the
  project's agent list quickly becomes useless.
- **Pool size** — 10 is a reasonable default; lower it if the hub or broker is under load.
- **One agent per day** — don't try to batch multiple dates into a single agent; it
  defeats the parallelism and complicates failure recovery.

## Cross-References

- For the per-day synthesis rules, see `release-notes-daily`.
