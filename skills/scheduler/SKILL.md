---
name: scheduler
description: >-
  Create and manage scheduled events and recurring schedules via the Scion CLI.
  Covers one-shot timers, cron-based recurring schedules, the message-to-orchestrator
  pattern for task dispatch, and lifecycle management (pause/resume/delete).
---

# Scheduler

Schedule one-shot or recurring events that send messages to agents at specified times. Use schedules for reminders, periodic checks, maintenance sweeps, and any work that should happen on a timer.

All schedules are **project-scoped** — an agent can only create and manage schedules within its own project.

## CLI Commands

### Create a one-shot event

```bash
# Fire in 30 minutes
scion schedule create --non-interactive --type message \
  --agent my-agent --message "Time to check the build" --in 30m

# Fire at a specific time (ISO 8601, UTC)
scion schedule create --non-interactive --type message \
  --agent my-agent --message "Deploy window open" --at 2026-06-10T14:00:00Z
```

Flags:
- `--type message` — required (only supported type)
- `--in <duration>` — relative delay (e.g. `30m`, `2h`). Mutually exclusive with `--at`
- `--at <ISO8601>` — absolute UTC time. Mutually exclusive with `--in`
- `--agent <name>` — target agent name (required)
- `--message <text>` — message body (required)
- `--interrupt` — deliver as an interrupt (agent sees it immediately)

### Create a recurring schedule

```bash
scion schedule create-recurring --non-interactive \
  --name daily-health-check \
  --cron "0 6 * * *" \
  --type message \
  --agent monitor-agent \
  --message "Run daily health check"
```

Additional flags:
- `--name <name>` — unique name within the project (required)
- `--cron <expr>` — 5-field cron expression in **UTC** (required)

### Manage schedules

```bash
scion schedule list --non-interactive --format json
scion schedule list --non-interactive --show recurring
scion schedule list --non-interactive --show events
scion schedule get <id> --non-interactive --format json
scion schedule pause <id> --non-interactive
scion schedule resume <id> --non-interactive
scion schedule delete <id> --non-interactive
scion schedule cancel <id> --non-interactive          # one-shot events only
scion schedule history <id> --non-interactive --format json
```

## Message-to-Orchestrator Pattern

To schedule creation of new agents on a timer, message a long-lived orchestrator rather than directly dispatching agents:

```bash
scion schedule create-recurring --non-interactive \
  --name weekly-audit \
  --cron "0 14 * * 1" \
  --type message \
  --agent orchestrator \
  --message "Start weekly audit"
```

The orchestrator receives the message and creates/manages task agents as needed. This gives full control over agent lifecycle.

## Common Patterns

### Self-reminder

```bash
scion schedule create --non-interactive --type message \
  --agent "$(scion whoami --non-interactive --format json | jq -r .name)" \
  --message "Check if deployment completed" \
  --in 15m
```

### Periodic monitoring

```bash
scion schedule create-recurring --non-interactive \
  --name ci-monitor \
  --cron "*/10 * * * *" \
  --type message \
  --agent ci-watcher \
  --message "Poll CI pipeline status"
```

## Gotchas

- **Cron is UTC.** There is no timezone configuration. Convert from local time.
- **`cancel` vs `pause` vs `delete`:** `cancel` is for one-shot events. `pause` and `delete` are for recurring schedules. Using the wrong one errors.
- **Schedule names must be unique** within the project.
- **Always use `--non-interactive` and `--format json`** when calling from agent code.
