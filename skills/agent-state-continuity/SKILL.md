---
name: agent-state-continuity
description: >-
  How to maintain a scratch state file that survives container restarts and agent
  deletion. Covers the agent-name namespace rule, git-ignore verification, the
  git-clean hazard, and shared-volume mirroring. Use when your role requires
  persistent state across sessions (coordinators, engineering managers, or any
  orchestrator that tracks workstreams).
---

# Agent State Continuity

## The state file

Keep a scratch state file at `.<agent-name>-state.md` in the workspace root as your working copy of project state — substitute your own agent name, so an agent named `acme-web-coord` keeps `.acme-web-coord-state.md`.

Two things about that file are load-bearing:

- **The namespace is your agent name, not your role.** One fixed name per role is one name for every instance of that role, so the second agent of any role to share a workspace overwrites the first one's file — no error, no conflict.
- **Confirm it is actually ignored where you are standing.** `git check-ignore -q .<agent-name>-state.md` answers it: exit 0 ignored, exit 1 not ignored, exit 128 not a git repository. Whether a workspace ignores this name is a property of that workspace, and no template can promise it on that workspace's behalf. If it comes back unignored, get it ignored before you rely on the file — an unignored state file is one `git add .` away from being committed, and in a public repository that is a disclosure, not a mess. Where the workspace root is not a git repository the check does not apply and there is nothing to lose.

## Suggested structure

```markdown
# [Role] State

## Last Updated
[timestamp]

## Active Workstreams
- [workstream]: [status, current agent, blockers]

## Pending Tasks
- [task description]: [priority, dependencies, assigned agent]

## Completed This Session
- [what was done, by whom, commit hashes or branch/PR]

## Decisions Made
- [decision]: [rationale]

## Notes for Next Session
- [anything the next session needs to know]
```

Read this file at the start of every session to restore context. Update it at significant milestones and before signaling completion.

## The file cannot carry continuity on its own

Whether `.gitignore` matches it is a property of the workspace rather than of this template — the `git check-ignore` above is what tells you. Where it is matched, that is deliberate and must stay: the entry is what stops `git clean` deleting the file, after a working-tree reset on 2026-07-27 destroyed unlisted content. It also means the file is never committed and never pushed, so the entry closes the `git clean` hazard and deepens the deletion one. Where it is not matched, `git clean` is live against it and getting it ignored is the fix. Under either answer the file does not survive the container, and a session boundary is precisely the case where the container is gone.

## Mirror to the shared volume

Whenever you update this file, also write the same state to the project scratchpad on the shared volume (`<scratchpad>/projects/<slug>/<agent-name>-state.md`), which outlives you and is where your successor should look first. **Your agent name belongs in that path too, not only in the local one.** One fixed path per project slug on a shared mount is the same collision one layer out: two agents of the same role on one project write the same file, last writer wins, no error and no edit event, and read-before-write does not help because by then the corruption looks like the reader's own file. There is no `.gitignore` on the scratchpad volume to fall back on, so the name is the only defence.

Treat the workspace-root copy as the fast local cache and the shared-volume copy as the record. If no shared volume exists, follow `artifact-durability` → **When only container-local storage is available** — do not un-ignore this file to solve it.
