---
name: agent-state-continuity
description: >-
  How to maintain a scratch state file that survives container restarts and agent
  deletion. Covers the agent-name namespace rule and scratchpad storage. Use when
  your role requires persistent state across sessions (coordinators, engineering
  managers, or any orchestrator that tracks workstreams).
---

# Agent State Continuity

## The state file

Keep your state file on the project scratchpad on the shared volume (`<scratchpad>/projects/<slug>/<agent-name>-state.md`, typically under `/scion-volumes/scratchpad/`). This is the only copy — do not keep a workspace-root copy, because the state file should never exist inside a git working tree.

Substitute your own agent name into the path, so an agent named `acme-web-coord` writes to `<scratchpad>/projects/<slug>/acme-web-coord-state.md`.

**The namespace is your agent name, not your role.** One fixed name per role is one name for every instance of that role, so the second agent of any role on one project writes the same file — last writer wins, no error and no edit event, and read-before-write does not help because by then the corruption looks like the reader's own file. The agent name is the only defence against this collision.

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

## When no shared volume exists

If the scratchpad volume is unavailable, follow `artifact-durability` → **When only container-local storage is available**. Do not write the state file to a path inside the git working tree as a fallback — state files should not be committed to any repository.
