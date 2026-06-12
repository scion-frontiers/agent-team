---
name: agent-recovery
description: >-
  Recover Scion agents from common stuck states: transient API errors, hub token expiry,
  container crashes, context exhaustion, and broker disconnects. Distinguishes "send a
  continue" from "needs recreation" and explains the underlying mechanics. Use when an
  agent appears stuck, errors repeat in its log, or its phase/activity state looks
  inconsistent.
---

# Agent Recovery

A practical guide for unstucking agents. The recurring trap is treating all "stuck" states the same — many are recoverable with a single message, others need recreation.

## Quick Triage Table

| Symptom | First action |
|---|---|
| Transient API error in agent log | `scion message <agent> "continue"` (do NOT recreate) |
| `LIMITS_EXCEEDED` state | `scion message <agent> "continue"` |
| Container crash (exit 255, `Exited` status) | Recreate the agent |
| Hub auth 401 repeating every 30s | Send any message to trigger token refresh |
| Context at 100% (memory limits) | Send raw `/clear` sequence (see below) |
| Phase stuck at `created` (lastSeen zero) | Wait a few minutes, then delete and recreate |
| Phase/activity desync (`starting` + `completed`) | Check for duplicate process; recreate |
| `gh` CLI returns 401 | See GitHub Auth section — this is not an agent state issue |
| Rebase fails "unrelated histories" | `git fetch --unshallow` inside the agent workspace |
| Interactive prompt blocking agent | `scion message <agent> --raw "ENTER"` or `--raw "0"` to dismiss |

## Hub Token Issues

### Crypto error after hub restart

The hub regenerated its signing key on restart and can't verify existing tokens.

**Symptom:** `failed to verify token: error in cryptographic primitive`

**Recovery:** Send any message to trigger a control-channel token refresh:
```bash
scion message <agent-name> "continue"
```

### Genuine token expiry

The token expires because refresh attempts also fail (refresh requires a valid token).

**Symptom:** `Token refresh failed: token refresh failed with status 401` repeating every 30s.

**Recovery:** Same — send any message to trigger refresh. If that fails, recreate the agent.

## API Error vs Container Crash

- **Transient API error** → `scion message <agent> "continue"`. Do not recreate.
- **Container crash** → look for `exit 255` or `Exited` in `scion list`. These need recreation.

The wrong choice (recreating an agent with a transient error) destroys the agent's in-memory state and any uncommitted work.

## Context Clear for Memory-Limited Agents

When an agent's context approaches 100%:

```bash
scion message <agent> --raw "/"
scion message <agent> --raw "clear"
scion message <agent> --raw "ENTER"
```

Use `scion look <agent>` first to verify the screen state before sending raw input.

## Agents Stuck in 'Created' Phase

**Symptom:** `phase: created` persists for 5+ minutes, `lastSeen` is zero, `scion look` returns 404.

**Cause:** CLI timeout vs broker dispatch delay mismatch. The broker was under load and dispatch exceeded the CLI timeout window.

**Recovery:**
1. Wait a few minutes and retry — broker load may have cleared.
2. If still stuck, delete and recreate.
3. Reduce concurrent agent count if it keeps happening.

## Recovery Decision Flow

```
Agent appears stuck
    ↓
scion look <agent>     ← always start here
    ↓
What does the screen show?
    ├── Transient API error / LIMITS_EXCEEDED
    │     → scion message <agent> "continue"
    ├── Auth 401 in log
    │     → scion message <agent> "continue" (triggers token refresh)
    ├── Container exited / status Exited
    │     → scion delete + recreate
    ├── Context at 100%
    │     → raw "/" + "clear" + "ENTER"
    ├── Phase mismatch (starting + completed)
    │     → check for duplicate process; recreate
    └── Anything ambiguous
          → inspect logs for more detail
```

## Anti-Patterns

- **Recreating on first sign of trouble.** Most stuck states are recoverable with "continue". Recreation destroys uncommitted work.
- **Sending raw input blind.** Always `scion look` first to see the current screen state.
- **Ignoring shallow clone issues.** If rebase fails with "unrelated histories", the fix is `git fetch --unshallow`, not force-push or branch recreation.
