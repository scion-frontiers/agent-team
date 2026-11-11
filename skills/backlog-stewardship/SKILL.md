---
name: backlog-stewardship
description: >-
  Enrich issue intake, assess reported bugs, propose related implementation bundles,
  and maintain completion traceability when assigned backlog stewardship. Complements
  a triage role; does not grant priority, dispatch, closure or recurring-maintenance authority.
---

# Backlog Stewardship

Use the assignment's repositories, landing target, board policy and decision owners.
A standing backlog lead has continuing intake and bookkeeping duties; a one-off
triage assessment does not acquire those duties by loading this skill.

## Intake and evidence

Turn short reports into actionable issues without silently expanding their scope.
Capture the reported behavior or desired outcome, expected behavior, affected
workflow, evidence, acceptance criteria, exclusions and dependencies. Ask only for
missing facts that change the assessment; keep unknowns explicit. Agent handoffs
should identify existing issue URLs, owner (or unowned), current decision, work
location, next step and blockers so the steward need not reconstruct identity
from a feature name. Keep private operational context out of public issue text.

Before filing, search relevant issues in all states and current board membership.
For engineering claims, inspect the relevant code at a named revision and check
merged and in-flight work. Follow likely prior-work links into PRs, branches or
available research; report unavailable sources and search limits. Board absence
alone is not evidence that the work is new.

Separate observations, code-based inference and reproduction results. Record the
source revision, supporting links and what was actually exercised in the triage
narrative. A plausible cause is not a reproduced bug. Use Needs-info when evidence
cannot support a verdict or size; offer the smallest useful next validation step.
Do not turn inaccessible evidence into rejection or a confident bug claim.
Recommend duplicate/superseded treatment with a canonical issue link; closing it
still follows the project's authority rules.

The assessment and proposed fields express the same judgment. Distinguish proposed
priority from an authorized commitment. Repo policy owns field names, taxonomy,
review markers and escalation destinations; mechanics skills own API writes.
Product scope and priority go to the product decision owner, ownership disputes
to the coordination owner, and restricted disclosures to the designated reviewer.

## Shared implementation bundles

A bundle is one implementation stream, often one worker and one PR, covering
several independently traceable issues. An epic groups a broader outcome and may
span many streams; shared parentage alone is not a reason to bundle its children.

Propose a bundle when small related changes share files, setup or validation and
one owner can reduce conflicts and PR overhead. Preserve each issue's scope and
acceptance criteria. Exclude independent work, incompatible release/review needs,
uncertain large investigations and dependencies that cannot land together.
The implementation owner decides feasibility and ordering; the steward records it.

Record included issue URLs, exclusions, owner, dependency order, shared validation,
and a per-issue acceptance-to-PR mapping. Keep an external blocker visible rather
than hiding it under the bundle parent. Split a bundle if scope grows beyond one
coherent stream. Do not relabel a deliverable bundle as an epic solely because it
has children, or reparent existing issues without resolving their current hierarchy.

Track completion per issue, including partial delivery and residual scope. A merged
bundle PR, a closing keyword or a closed parent does not prove that every child is
done. Verify the intended landing target and obtain the project's required owner
confirmation for each issue. Keep a parent open while its required scope remains;
link follow-up work explicitly rather than losing it in a completion summary.

## Writes and recovery

Know which operations the assignment permits and which capabilities are available.
Issue write and board write are separate. Do not request credential values or copy
credentials into guidance or state. Preserve settled human decisions according to
board policy; missing metadata does not authorize reprioritization.

External writes are not one transaction. Record issue/item identifiers, intended
values and successful steps. After a timeout, re-read actual state before retrying:
issue creation and comments can duplicate even when field assignment is repeatable.
Before repeating a field write, check for intervening human edits. Reconcile only
the missing authorized steps, then read back and report any partial result or
capability failure. Stop on conflicting state rather than retrying indefinitely.

## Standing responsibilities and restart

For an assigned standing service, keep a durable checkpoint of source revision,
last processed landing revision, issue owners, pending decisions, partial writes
and explicit suspensions. Verify the recipient can read it across restarts; use
the project's continuity/durability mechanism rather than assuming a local path
is shared. Never put credential material in the checkpoint.

On restart, recover that checkpoint, current authority and suspension state;
verify the intended account/repositories, read pending messages and reconcile
tracked issues against actual remote state. A missing checkpoint calls for bounded
reconstruction with stated uncertainty, not replaying every write.

Describe maintenance duties separately: intake triage, missing metadata, completion
reconciliation, archival and optional backlog refill have different effects.
For recurring work, declare the logical service, responsible owner, target, duty
source, cadence/timezone, permitted writes and enabled/paused intent in project
configuration. No cadence or automatic sweep is implied by this skill. Startup
should verify existing arrangements; an empty visibility-limited schedule list is
not proof of absence. Missing, conflicting or paused schedules require resolution
by their authorized owner, not automatic duplication, retargeting or resumption.
Report actionable changes or failures; do not create periodic status-only traffic.
