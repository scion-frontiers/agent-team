---
name: artifact-durability
description: >-
  Where a produced artifact has to end up so that it survives the agent that produced it,
  and what to do when no durable destination is available. Owns the property definition
  (reachable after your container is deleted), the list of durable destinations, and the
  fallback ladder. Use when you are about to write a deliverable to a path, when you are
  deciding whether your work is safe, or before you signal completion.
---

# Artifact Durability

## The property

> **An artifact is durable when it is reachable after your container has been deleted.**

Not written. Not saved. Not "in the scratchpad." Reachable by someone else, later, once you
are gone.

You will not be consulted before your container is deleted, so every check you run on your
own work has to test that property and nothing weaker.

## Why this rule is stated as a property and not a place

This defect has been fixed at least six times in this repo, and each fix reintroduced it
one layer down, because each one tested something *easy to observe* instead of the thing
that *determines survival*:

| The check that was written | What it actually tests | What survives |
|---|---|---|
| "don't delete an agent with uncommitted work" | a commit exists locally | the commit is **pushed** |
| "committing prevents losing work" | a commit exists locally | the commit is **pushed** |
| "or written to the scratchpad" | a **path** was used | the path is **outside the container** |

**A location satisfies a check. Only a property satisfies the purpose.** A path is evidence
of durability only when you know why that path outlives you — and one of the two paths this
repo has historically called "the scratchpad" does not.

## Durable destinations

Exactly three things survive your deletion:

1. **A pushed git ref on a remote.** Push, then confirm the remote has it. A commit that
   never left the container dies with it, and it passes an uncommitted-work check on the
   way out.
2. **Storage mounted from outside the container** — e.g. the shared scratchpad volume at
   `/scion-volumes/scratchpad/`, or a hosting bucket you have published to and verified.
   Durable because the deletion cannot reach it, not because of how the path is spelled.
3. **The body of a message you have already sent.** Content inline in the message, not a
   filepath pointing at storage only you can see.

Everything else dies with you: the working tree, an unpushed commit, `/tmp`, a gitignored
file, and `/workspace/.scratch/`.

## When only container-local storage is available

`/workspace/.scratch/` is a real fallback for scratch, and an environment with no shared
volume is real. It is **not** a delivery destination, for two reasons — and the second is
the one that gets missed:

- it does not survive your container, and
- it is **not shared with any other agent**, so it cannot carry a handoff either.

If the shared volume is unavailable, do not write your deliverable there and signal
completion. Do the first of these that is available to you:

1. **Commit it to your own work branch and push it.** This requires a tracked path inside
   the repo — not `.scratch/`, which is gitignored. Push your own work branch, never the
   integration branch.
2. **Send the content in a message** to the agent that dispatched you, inline in the body.
   Say in the message that this is the delivery, because no durable file exists.
3. **Declare a blocker** and withhold the completion signal until it is resolved.

Never settle this silently by picking whichever location is easiest to write to. Choosing a
non-durable path without saying so looks identical, from outside, to having delivered.

## A gitignored deliverable is invisible to "commit and push"

Measured: with `.scratch/` in `.gitignore`, a file written under it leaves `git status`
reporting a **clean tree** and `git add -A` staging **nothing**.

So a completion rule of the form *"commit and push anything still sitting in the working
tree"* cannot see such a deliverable, and does not fail — it reports success. If your
deliverable is a file rather than a commit, name that file and confirm its destination
directly. Do not infer that it is safe from a clean working tree.

## Hardening against one hazard is not durability

When you protect an artifact, **name the hazard you did not cover.** Adding a file to
`.gitignore` protects it from `git clean` — a real fix for a real incident — while leaving
it exposed to container deletion, and while actively preventing it from being committed.
Two hazards; that entry closes one and deepens the other.

Protecting something against the hazard you just met makes it *feel* protected, and that
feeling is what stops the search for the second hazard.

## Related

- Durability-push versus integration-push, and who may merge to shared ground:
  `templates/eng-manager/agents.md`.
- Who may delete an agent, and what must be true first: `scion-agent-manage` →
  **Agent Lifecycle**.
