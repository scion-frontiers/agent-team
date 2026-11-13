---
name: software-engineering-process
description: >-
  How multi-agent engineering work is owned, sized, delegated, and verified from request
  to delivery. Covers the ownership tree, the independent code-review invariant, sizing by
  delegation depth, the code-review cycle, and briefing. Use when you are responsible for
  delivering an issue you will not implement yourself — classifying incoming work,
  deciding which roles to instantiate, creating and supervising agents, or governing a
  code-review cycle.
---

# Software Engineering Process

## Overview

How an engineering issue moves from request to delivered change when more than one agent is
involved: **who owns delivery, who creates whom, what each role owes, and when the work is
actually done.**

How to design, implement, and review is not here — read the role definitions and their own
skills for that.

## When to use

Use when you are **responsible for delivering an issue you will not implement yourself**:

- classifying an incoming issue and deciding how much process it warrants
- deciding which roles to instantiate, and creating those agents
- supervising agents you created and deciding whether to accept their output
- governing a code-review cycle
- deciding whether you are finished

**Do not use** as a developer, reviewer, architect, or investigator executing an assigned
task. Read your brief and your own role material instead.

---

## 1. Ownership

**Exactly one agent owns delivery of an issue at any moment.** That agent is the
**issue owner**. By default the issue owner is the coordinator. The issue owner:

1. creates the agents it needs,
2. briefs them,
3. monitors them,
4. accepts or rejects what comes back,
5. remains the issue owner until the issue is delivered.

For coordinators, acceptance and oversight follow
[Coordinator Boundary](https://github.com/scion-frontiers/agent-team/blob/main/templates/coordinator/agents.md#coordinator-boundary) in the coordinator role instructions.
This includes dedicated issue owners using that role. Technical assessment stays
with competent workers/reviewers; the coordinator confirms the required acceptance
stamp names its assessor, deliverable/revision, outcome, and supporting evidence
reference. It does not inspect that evidence or repeat the assessment. Delivery,
disposition, sequencing, and escalation accountability remain with the coordinator.

### How an issue gets its owner

An issue is owned in one of two ways, depending on size:

- **Small code issues (depth 1):** the coordinator is the issue owner directly — it creates a
  developer and reviewer(s) and remains the issue owner itself. This corresponds to the
  Small tier in the Sizing table (§3).
- **Chunky and Large issues (depth 3):** the coordinator creates a **dedicated issue
  owner** — an agent using the coordinator template, briefed as a delegated
  sub-coordinator responsible for owning the SEP for that issue. The dedicated issue
  owner creates the roles the issue needs — investigator, architect, and engineering
  manager — and supervises them. This corresponds to the Chunky and Large tiers in the
  Sizing table (§3).

In both cases, everything that follows about the issue owner applies equally: the
responsibilities are the same regardless of whether the coordinator fills the role itself
or delegates it to a dedicated agent.

**Only the top-level coordinator delegates issue ownership.** The issue-owner hierarchy
is exactly two layers deep: coordinator, and optionally a dedicated issue owner beneath
it. A dedicated issue owner does not further delegate issue ownership.

Within an issue, the issue owner may delegate **implementation subtrees** — for example,
an engineering manager owns implementation and review for a set of developers. A
delegated agent creates, supervises, and accepts *inside* its subtree, but the issue
owner still accepts that subtree's deliverable.

> **A dedicated issue owner owns the full process for its issue** — this is precisely
> what delegation of issue ownership means. The coordinator delegates both authority
> and accountability for driving the SEP within that issue. What the coordinator
> retains is **project-level tracking**: it tracks all issues, monitors their dedicated
> owners, and remains responsible for the portfolio, not the per-issue process.

### Supervision follows creation

**The tree is the notification topology.**

| Mechanism | Consequence |
|---|---|
| Creating an agent **auto-subscribes** you to its notifications | You are the one told when it completes or stalls |
| Ancestry chains (`root → parent → child`) grant **transitive access** | Any ancestor can manage any descendant, at any depth |
| Subscriptions can also be project-scoped | Non-creators *can* observe, but only if wired deliberately |

**You are the one who will be told when it breaks. If you look away, nobody is looking.**

**The agent that creates workers is the one subscribed to their notifications.** If you
delegate oversight to an agent that did not create the workers, nobody is subscribed —
it looks wired but is silently disconnected.

### Definition of done for an issue owner

> **Your work is not complete while your subtree is live. Delivery is the deliverable, not
> dispatch.**

Dispatching agents and reporting that you dispatched them is not completion. Signal
completion when the work is delivered, not when it is handed out.

---

## 2. The invariant

Scale implementation process with the work while preserving independent code review.
The review invariant and cycle apply to code changes, not prose-only artifacts.

> **For code changes requiring code review, the author and reviewer are
> different agents.** The constraint is scoped to the issue, not the tree.

Prose-only design documents, skills, instruction templates and documentation generally
go to the user for review and approval. They do not require the agent code-review
cycle, fresh reviewers or successive review rounds. A user may explicitly request
a review when useful. Author validation is not independent review.

Classify the actual change, not its file extension: executable examples, scripts and
behavior-changing code still require code review. In a mixed change, review the code
with its relevant prose context; accompanying prose does not create another mandatory
review chain.

Architecture planning follows direct human design discussion and feedback. Do not
apply the developer/reviewer cycle to a design artifact, appoint a mandatory
architecture reviewer, or require successive reviewer signoffs before the human
can discuss or accept the design. The architect revises the design in response to
that discussion; the issue owner checks scope and delivery without adding a review
round. This does not remove independent review of the later code changes.

---

## 3. Sizing

Sizing is **depth of the ownership tree** — not stage count, not agent count.
The implementation trees below include code reviewers; prose-only work follows §2
without creating reviewer roles to fit a tier.

| Tier | Tree | Depth |
|---|---|---|
| Small | coordinator → developer, reviewer | 1 |
| Chunky | coordinator → dedicated issue owner<br>dedicated issue owner → investigator, architect, engineering manager<br>engineering manager → developers, reviewers | 3 |
| Large | coordinator → dedicated issue owner<br>dedicated issue owner → investigator, architect, engineering manager<br>engineering manager → developers, reviewers | 3 |

The coordinator is the root at every size. Sizing is how many layers hang beneath it.

> **Note:** Chunky and Large are structurally identical in tree shape and depth.
> Both place a dedicated issue owner between the coordinator and the working roles.
> The distinction between Chunky and Large is **scope and complexity of the issue**,
> not tree structure — both warrant a dedicated issue owner, but a Large issue
> typically involves more phases, more agents, or longer timelines than a Chunky one.

**Depth is not change magnitude and not dispatch route.** Those are separate axes, they
carry their own labels elsewhere, and those labels overlap with these tiers. They also
correlate — a depth-3 tree is usually a large change — which is exactly what makes merging
them look safe. It is not: a large change can still be one developer and one reviewer, and
neither axis fixes the depth of the tree. All three may legitimately coexist. **Do not
consolidate them into this one.**

Changing the look of a single button does not need an investigator, an architect, an
engineering manager, developers *and* reviewers. It needs a developer and a reviewer.

### Sizing rules

1. **Size up when uncertain.** Estimation is unreliable; over-sizing is the cheaper error.
2. **Every orchestrating role must collapse itself when handed less than its tier
   assumes.** An engineering manager given a trivial design relays it to a single
   developer instead of manufacturing phases.

**Never apply rule 1 without rule 2.** Over-sizing without collapse produces **ceremony
inflation**: an agent told it is an engineering manager performs engineering management,
because that is what the role description says the role does.

### Cost of supervision

A coordinator sustains roughly **2–3 issues concurrently**, and **tens of agents across a
session**. It scales because it holds very little detail per issue: it supervises rather
than does.

> **Keep execution of delivery details with the workers** (for example, implementation
> detail for a coding task). Design rationale, review findings, and delivery state belong
> in artifacts and scratchpad files. Coordinators retrieve administrative state and
> acceptance stamps under Coordinator Boundary; workers read specialist detail when needed.

Holding that state yourself is what takes you from three concurrent issues down to one.

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

The issue owner monitors each of these and accepts or rejects the artifact. An artifact
that comes back thin, off-brief, or stalled is the issue owner's problem to catch.

---

## 5. The review cycle

This cycle applies to code changes. Prose-only artifacts, including architecture
planning, follow §2 without this cycle or its fresh-reviewer and round machinery.
Run per pull request or per developer. The reviewer reviews the entire proposed change
against its base and intended behavior and inspects surrounding code needed to assess effects.
A whole-branch review may supplement, not replace, that scope.

### Independent assessment and fix verification

Start with a full blind independent review. Give the reviewer the exact base and
head revisions, requirements, relevant project guidance and actual verification
evidence. Do not supply previous findings, verdicts, author rebuttals or the
owner's disposition history before the reviewer records its assessment, and do not
leave earlier review records where the reviewer reads or writes its own. Blind does
not mean deprived of requirements or code context. If prior review content is
unavoidably visible, disclose that limit rather than claiming a blind pass.

After the assessment is recorded, the owner reconciles findings with history.
A targeted independent check of a previous fix may use that history; label it
fix verification, not a full blind review. When fix verification spans a rebase,
also give the previously reviewed base and head. Keep the developer across rounds.
Use a fresh reviewer for each full blind pass; a reviewer who already knows the
findings may verify fixes but cannot claim a blind pass on that change.

### Disposition and convergence

The owner decides each finding once in the existing review record, with a short
reason: fix now, decline as not required, or defer outside this delivery. A fixed
item links its verification. A real deferred defect gets an existing backlog
reference and owner; optional polish may be declined without creating a task.
Group repeated instances of the same cause. Reopen a disposition only for new
evidence or an explicit human instruction. No separate form is required.

Critical/Required findings under `code-review` remain actionable at every pass;
their classification must name a material failure or an unmet requirement.
Disputed blocking classification is resolved on evidence by the owner or escalated
to the existing decision-maker; it is not silently downgraded. Deferring an accepted
blocking defect needs the existing human exception authority. Security holds obey
the project's explicit security threshold; do not bypass it by relabelling a
below-threshold security finding as Required or withholding signoff.

New lower-priority findings from a later blind pass default to decline as not
required for this delivery, or defer if there is a useful separate defect to track.
They do not automatically generate fixes, more reviewers, or a new approval gate.
On the first pass too, non-blocking recommendations do not have to be implemented.
Human instructions and explicit acceptance requirements still govern.

After a material fix, run a fresh full blind review of the resulting complete
change. Also independently verify the original failure is fixed; that check may
occur after the blind assessment is recorded. A lower-priority fix that the owner
chooses to include receives proportionate independent verification of the edit
and its effects. It does not alone require another full pass. If the edit changes
a contract, behavior or risk materially, treat it as a material change and obtain
a full blind pass, subject to the rebase rule below. A human may request a full
pass at any point.

### Rebases

Rebasing onto main, changed commit IDs or a changed base do not by themselves
require another review. Re-review when conflict or semantic resolution appears
to introduce defect risk; a conflict-free rebase can also expose a risky semantic
interaction. Use the existing owner's engineering judgment during the work; a coordinator owner
obtains that technical risk assessment from the responsible worker. This does not
require a separate reviewer, checklist, approval or permission to decide this question.
When risk warrants re-review, obtain a fresh full blind assessment of the complete
resulting change, then reconcile history. Do not turn it into fix-only review.
Run relevant verification and required project checks when earlier evidence no
longer applies. Record the resulting revision and verification limits in the
existing handoff. Review, verification and human merge authority remain distinct.

### Completion

The review cycle ends when a full independent assessment covers the proposed
change, subsequent edits have the verification/review required above, no accepted
blocking finding remains unresolved, all other findings have a finite disposition,
and required checks and human approvals are satisfied. APPROVE may carry declined
or deferred non-blocking findings. Zero comments and perfectly clean successive
passes are not completion requirements. Record the reviewed head and any later
verified edits or rebase. Apply the rebase rule above: a new head/base alone does
not invalidate review; risky resolution does require re-review.

Do not start a new round just because non-blocking findings remain. If the same
material defect survives a fix attempt, or the owner cannot resolve a material
classification dispute, stop automatic dispatch and escalate with the evidence,
what was tried and a recommended next action. A decision to continue must name
the changed approach. This replaces a fixed round count; no round count authorizes
shipping an unresolved blocker or skipping an independent assessment.

## 6. Context discipline

For code implementation, separate agents per phase **are** separate context windows.
**Do not collapse phases into one agent to save time** — what you lose is the context
reset, not just the division of labour. This does not add phases or reviewers to
prose-only work (§2).

- A role boundary is a context reset.
- A fresh reviewer is a fresh window.
- Per-PR review scope keeps attention bounded.

Agent performance degrades well before the context window is full. Treat roughly **40%
utilisation** as the point where quality starts costing more than the saved handoff.

---

## 7. Briefing

**Every agent you create needs a brief, and no agent starts work without one.** An agent
you created that is doing the wrong thing is usually an agent you briefed badly.

For brief structure and file-based handoff, see `scion-agent-manage`. For engineering
briefs:

- **Point at sources instead of restating them.** Cite sources such as the spec section,
  mock scene ID and test names taken from the change. A paraphrase drifts from its source,
  and the recipient has to check it against that source anyway (`code-review` → Review
  Process).
- **Use `code-review`'s severity labels** in code-review briefs rather than a different set.

You can draft the next agent's brief while waiting for the decision that gates it, so
dispatch follows the decision directly.

---

## 8. Common rationalizations

| Rationalization | Reality |
|---|---|
| "I'm the engineering manager, so I should break this into phases." | A trivial design goes straight to one developer. Orchestration you did not need is waste you introduced. |
| "I'm the architect, so I owe a design document." | Handed a one-line fix, you owe a sentence and a decision. |
| "This code change is small enough that I can review my own work." | Independent code review still applies; prose-only artifacts follow §2. |
| "I handed it off, so it's their problem now." | You delegated the work, not the delivery. |
| "My task was to dispatch the developers, and I dispatched them." | Delivery is the deliverable, not dispatch. You are not done while your subtree is live. |
| "I spawned the agents, so I'm done until they report." | An agent that never reports is exactly the case you exist to catch. Silence is not progress. |
| "The last reviewer approved it, so the follow-up fix is fine." | Later edits need proportionate verification; material changes need a full blind pass, with the §5 rebase exception. |
| "These findings are non-blocking, so I can ignore them." | Record a finite disposition; optional polish need not become a task or another round. |
| "I'll review the whole branch at once, it's more efficient." | Per-PR or per-developer. Whole-branch reviews dilute the attention the cycle depends on. |
| "I already know how to review code / plan work / design this." | You have generic priors, not this project's methodology. Read the role material. |
| "The brief is long, so I'll pass it inline." | Write it to the scratchpad and pass the filepath — an inlined brief puts backticks and `$variables` into a shell command. See `scion-cli-operations` → **Shell Safety for Task Prompts**. |

## 9. Red flags

- You created an agent and have not heard from it, and have not gone to check.
- You are about to report completion while an agent you created is still running.
- The same agent authored and code-reviewed a code change.
- Fix verification is presented as a full blind review.
- An architect is re-running investigation that was already delivered to it.
- You are relaying messages between two agents that could talk directly.
- Your tier says "engineering manager" and the design is three lines long.
- Review findings were recorded as non-blocking and never revisited.

## 10. Verification

Before reporting delivery, coordinators confirm current status and required acceptance
stamps under Coordinator Boundary; workers perform the technical verification:

- [ ] Every agent you created has reported, been accepted, or been explicitly stopped.
- [ ] Code changes received independent code review; prose-only artifacts followed §2 without a mandatory agent reviewer loop.
- [ ] For code changes, the review-cycle completion condition in §5 holds, including its rebase rule.
- [ ] Required verification and human approvals are satisfied; remaining gaps and
      non-blocking dispositions are visible.
- [ ] No subtree is still live.

**Run this immediately before you report — not as you go.** A box ticked earlier states
something about a moment that has passed: agents you created have since finished, stalled,
or spawned others. **Re-check; do not recall.**

## 11. Related material

| For | See |
|---|---|
| EM brief template for issues with a dedicated issue owner | `references/em-brief-template.md` in this skill |
| Shell metacharacters and absolute paths in task prompts | `scion-cli-operations` → **Shell Safety for Task Prompts** |
| Briefing structure, model override, agent recovery, stall handling | `scion-agent-manage` |
| How to do architecture, design-artifact shape | architect role definition |
| Work slicing, parallel-vs-sequential technique | engineering-manager role definition |
| Review methodology and severity | `code-review` |

## Engineering-manager lifecycle

Core Scion's `scion-agent-manage` → **Agent Lifecycle** owns generic lead
protection, delegated-supervisor authority, unanswered-human-question protection,
durable deliverables and bottom-up teardown. This engineering workflow owns the
mapping of its templates to those categories.

A lead may start and retire the engineering managers it created for bounded
implementation assignments. Such an EM maps to the core delegated-supervisor
category, not to lead merely because it supervises developers. The creating lead
may delete it only after its subtree is torn down, its work has landed under the
agreed engineering workflow, its deliverables are verified durable and it has no
open question to a human. An EM not started by a lead retains the human-only
deletion fallback. Being an issue owner does not by itself grant authority to
delete another lead or another lead's EM.

If an EM is the first agent on the project/workstream, coordinates that overall
workstream, or is human-designated as its lead, core lead protection takes
precedence: only explicit human instruction naming the workstream authorizes
its deletion. A role label, completion signal or readiness message cannot remove
that protection. The coordinator/issue-owner and EM templates consume this
mapping; they do not redefine core lifecycle policy.
