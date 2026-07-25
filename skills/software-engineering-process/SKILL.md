---
name: software-engineering-process
description: >-
  How multi-agent engineering work is owned, sized, delegated, and verified from request
  to delivery. Covers the ownership tree, the independent-review invariant, sizing by
  delegation depth, the review cycle, and briefing. Use when you are responsible for
  delivering an issue you will not implement yourself — classifying incoming work,
  deciding which roles to instantiate, creating and supervising agents, or governing a
  review cycle.
---

# Software Engineering Process

## Overview

This skill governs how an engineering task moves from request to delivered change when
more than one agent is involved. It defines **who owns delivery, who creates whom, what
each role owes, and when the work is actually done.**

It deliberately does **not** teach how to design, implement, or review. Those live in the
role definitions and their own skills. This is the spine; the roles carry the depth.

## When to use

Use when you are **responsible for delivering work you will not do yourself**:

- classifying an incoming issue and deciding how much process it warrants
- deciding which roles to instantiate, and creating those agents
- supervising agents you created and deciding whether to accept their output
- governing a review cycle
- deciding whether you are finished

**Do not use** as a developer, reviewer, architect, or investigator executing an assigned
task. Read your brief and your own role material instead.

---

## 1. Ownership

**Exactly one agent owns delivery of a work item at any moment.** By default that is the
coordinator. The owner:

1. creates the agents it needs,
2. briefs them,
3. monitors them,
4. accepts or rejects what comes back,
5. remains the owner until the work is delivered.

An owner may delegate a **bounded subtree** — an engineering manager owns implementation
and review; a sub-coordinator owns an entire large issue.

> **Delegation transfers the authority to create agents. It does not transfer
> accountability.**

That sentence is what an engineering manager mechanically *is*: the thing that takes
agent-creation for the implement-and-review subtree off the coordinator's plate. Ownership
does not transfer. The root owns the outcome all the way through.

### Why the tree exists

The tree is not an organisational metaphor. **It is the notification topology.**

| Mechanism | Consequence |
|---|---|
| Creating an agent **auto-subscribes** you to its notifications | You are the one told when it completes or stalls |
| Ancestry chains (`root → parent → child`) grant **transitive access** | Any ancestor can manage any descendant, at any depth |
| Subscriptions can also be project-scoped | Non-creators *can* observe, but only if wired deliberately |

So "own the delivery" is mechanical, not moral: **you are the one who will be told when it
breaks. If you look away, nobody is looking.**

The practical corollary — **supervision follows creation by default.** Delegating oversight
to an agent that did not create the workers is the trap: it looks fine, and silently
nobody is subscribed.

### Definition of done for an owner

> **Your task is not complete while your subtree is live. Delivery is the deliverable, not
> dispatch.**

Dispatching agents and reporting that you dispatched them is not completion. Signal
completion when the work is delivered, not when it is handed out.

*(Status-signal mechanics — `blocked`, `task_completed`, `ask_user` — are already injected
into every agent. This section is only about what state you are actually in.)*

---

## 2. The invariant

Everything about this process scales down with the size of the work **except independent
review.**

> **For any item under work, the author and the reviewer are never the same agent.**

Stated as a separation constraint, not a headcount, because a coordinator may supervise
several items at once and counting agents in the tree gives the wrong answer. The
coordinator does not count toward it; the constraint is scoped to the item.

**Why it cannot be relaxed:** the reviewer's job is *not letting bad code through*. An
agent that just argued itself into a design cannot then hold that goal against its own
work. This is role conflict, not context size — a bigger context window does not fix it.

---

## 3. Sizing

Sizing is **depth of the ownership tree** — not stage count, not agent count.

| Tier | Tree | Depth |
|---|---|---|
| Small | coordinator → developer, reviewer | 1 |
| Chunky | coordinator → investigator, architect, engineering manager<br>engineering manager → developers, reviewers | 2 |
| Large | coordinator → sub-coordinator → engineering manager → developers | 3 |

The coordinator is the invariant root and is present at every size. Sizing asks only how
many layers hang beneath it.

Changing the look of a single button does not need an investigator, an architect, an
engineering manager, developers *and* reviewers. It needs a developer and a reviewer.

### Two rules that must ship together

1. **Size up when uncertain.** Estimation is unreliable; over-sizing is the cheaper error.
2. **Every orchestrating role must collapse itself when handed less than its tier
   assumes.** An engineering manager given a trivial design relays it to a single
   developer instead of manufacturing phases.

Rule 2 is what makes rule 1 safe. Without it, over-sizing produces **ceremony inflation** —
an agent told it is an engineering manager performs engineering management, because that
is what the role description says the role does.

### Cost of supervision

A coordinator sustains roughly **2–3 issues concurrently**, or tens across a session when
supervision is light. It scales because it holds very little detail per issue: it
supervises rather than does.

> Coordinator concurrency is a direct function of how much per-issue context the process
> forces it to hold.

Anything that makes the coordinator carry design detail or review state per issue reduces
how many issues it can carry at all. Keep coordinator obligations thin on purpose.

---

## 4. The work

Roles under supervision, each owing an artifact to the agent that created it. Instantiate
only the roles the size calls for.

| Role | Receives | Owes |
|---|---|---|
| Investigator | the issue | findings — facts, no design decisions |
| Architect | the investigation findings | an agreed design, iterated with a human. Reads the findings; does not re-research |
| Engineering manager | the agreed design | implementation phases; the parallel-vs-sequential call; a supervised subtree of developers and reviewers |
| Developer | a brief and a phase | a change, ready for review |
| Reviewer | a change | a verdict, held to *do not let bad code through* |

The owner monitors each of these and accepts or rejects the artifact. An artifact that
comes back thin, off-brief, or stalled is the owner's problem to catch.

---

## 5. The review cycle

Run per **pull request or per developer** — not across the whole branch at once. A single
final whole-branch review is sometimes warranted, but it does not replace the per-unit
cycle.

Each round: review → fixes → **new reviewer**. Repeat until a review comes back clean.

### The asymmetry is deliberate

| | Across rounds | Why |
|---|---|---|
| Reviewer | **replaced** | Sheds context. Avoids anchoring on its own prior review; attention decays across long reviews |
| Developer | **retained** | Keeps context. Holds the thread of what was tried and why |

**Compact the adversary, preserve the author.** Do not generalise this to "fresh context is
always better" and start rotating developers.

One reviewer commonly misses items it raised itself on re-read. That is the failure the
rotation exists to catch — not reviewer incompetence.

**Cap: 6 rounds**, then escalate. This is a safety valve against pathological
non-convergence, not an expected path. Hitting it regularly means something upstream is
broken.

**Non-blocking is not optional.** It means "will not block the merge review." It does not
mean "skip it before shipping."

---

## 6. Context discipline

Separate agents per phase **are** separate context windows. That is a real benefit of the
structure, not just division of labour — which is why phases should not be collapsed into
one agent to save time.

- A role boundary is a context reset.
- A fresh reviewer is a fresh window.
- Per-PR review scope keeps attention bounded.

Agent performance degrades well before the context window is full. Treat roughly **40%
utilisation** as the point where quality starts costing more than the saved handoff.

---

## 7. Briefing

Every agent you create needs a brief. Write the brief to a **shared scratchpad file and
pass the filepath** — do not inline a long brief into the creation command.

A brief states:

| Section | Content |
|---|---|
| Task | what to do, in one or two sentences |
| Context | what has already been decided, and where to read it |
| Boundaries | what is explicitly out of scope |
| Deliverable | what artifact is owed, and in what shape |
| Reporting | who to report to, and when |

For shell-escaping rules when passing prompts, see the `scion-cli-operations` skill —
do not improvise quoting.

---

## 8. Common rationalizations

| Rationalization | Reality |
|---|---|
| "I'm the engineering manager, so I should break this into phases." | A trivial design goes straight to one developer. Orchestration you did not need is waste you introduced. |
| "I'm the architect, so I owe a design document." | Handed a one-line fix, you owe a sentence and a decision. |
| "This change is small enough that I can review my own work." | You cannot hold both "make this work" and "do not let this through." |
| "I handed it off, so it's their problem now." | You delegated the work, not the delivery. |
| "My task was to dispatch the developers, and I dispatched them." | Delivery is the deliverable, not dispatch. You are not done while your subtree is live. |
| "I spawned the agents, so I'm done until they report." | An agent that never reports is exactly the case you exist to catch. Silence is not progress. |
| "The last reviewer approved it, so the follow-up fix is fine." | Fresh reviewer each round. Reviewers miss items on re-read of their own review. |
| "These findings are non-blocking, so I'll ship and fix later." | Non-blocking means it will not block merge review. It does not mean skip. |
| "I'll review the whole branch at once, it's more efficient." | Per-PR or per-developer. Whole-branch reviews dilute the attention the cycle depends on. |
| "I already know how to review code / plan work / design this." | You have generic priors, not this project's methodology. Read the role material. |
| "The brief is long, so I'll pass it inline." | Write it to the scratchpad and pass the filepath. |

## 9. Red flags

- You created an agent and have not heard from it, and have not gone to check.
- You are about to report completion while an agent you created is still running.
- The same agent wrote and reviewed a change.
- A reviewer is reviewing its own previous review's fixes.
- An architect is re-running investigation that was already delivered to it.
- You are relaying messages between two agents that could talk directly.
- Your tier says "engineering manager" and the design is three lines long.
- Review findings were recorded as non-blocking and never revisited.

## 10. Verification

Before reporting delivery:

- [ ] Every agent you created has reported, been accepted, or been explicitly stopped.
- [ ] Every change was reviewed by an agent that did not write it.
- [ ] The final review round came back clean — not "clean except".
- [ ] Non-blocking findings were resolved or consciously deferred with a reason.
- [ ] No subtree is still live.

## 11. Related material

| For | See |
|---|---|
| Shell escaping when passing prompts | `scion-cli-operations` |
| Model override, agent recovery, stall handling | `scion-agent-manage`, `agent-recovery` |
| How to do architecture, design-artifact shape | architect role definition |
| Work slicing, parallel-vs-sequential technique | engineering-manager role definition |
| Review methodology and severity | code-reviewer role definition, `pr-code-review` |

Status-signal mechanics (`blocked`, `task_completed`, `ask_user`) are injected into every
agent automatically and are not restated here.
