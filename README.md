# Scion Agent Team

Welcome to the **Scion Agent Team** repository. This is the central repository for the standard collection of agent templates and shared skills used by the Scion agentic orchestration framework. 

These templates and skills enable autonomous, multi-agent teams to coordinate, design, build, test, and maintain complex software projects.

---

## 🏗️ Repository Architecture

This repository is organized into two primary pillars:

```
agent-team/
├── templates/         # Agent role definitions (instructions, system prompts, default skills)
│   ├── coordinator/
│   ├── developer/
│   └── ...
└── skills/            # Reusable skill banks (command-line tools, specific action scripts)
    ├── agent-recovery/
    ├── scheduler/
    └── ...
```

1. **Templates (`/templates`)**: Standard definitions for various agent roles (e.g., Coordinator, Developer, Architect, QA Tester). Each template contains:
   - `scion-agent.yaml`: Configures the agent type, description, model harness, and imported skills.
   - `agents.md`: High-level role instructions and operational protocols for the agent.
   - `system-prompt.md`: Base system prompt establishing behavioral constraints, environment paths, and execution context.
2. **Skills (`/skills`)**: Independent, reusable capability modules (scripts, tools, and manuals) that can be imported by any agent to perform specialized actions (e.g., `agent-recovery`, `scheduler`).

---

## 🤖 Agent Templates

Each folder under `/templates` defines a distinct agent persona with specialized capabilities.

| Role / Template | Description | Key Focus Area |
| :--- | :--- | :--- |
| **`coordinator`** | Project coordinator; manages agent teams, delegates implementation, and drives progress autonomously. | Decomposing requirements, orchestrating workers, status tracking |
| **`architect`** | System designer; defines software architecture, APIs, data schemas, and implementation blueprints. | Technical design, API specs, component breakdown |
| **`developer`** | Software engineer; implements features, refactors code, and writes unit/integration tests. | Feature development, code quality, unit testing |
| **`web-dev`** | Web frontend engineer; specializes in styling, responsive UI/UX, and component-driven web code. | UI layouts, frontend component architecture, responsive styling |
| **`web-builder`** | Full-stack web application builder; crafts end-to-end interactive apps and user experiences from specs. | Complete web app assembly, integration, deployment readiness |
| **`qa-tester`** | Quality Assurance specialist; designs test cases, validates user acceptance criteria, and runs manual/E2E tests. | Acceptance testing, functional verification, test reporting |
| **`test-engineer`** | Test infrastructure specialist; maintains test harnesses, writes integration/E2E suites, and builds CI pipelines. | E2E automation, CI/CD integration, test suite performance |
| **`code-reviewer`** | Senior reviewer; evaluates pull requests for correctness, readability, security, and performance. | PR quality control, architectural alignment, compliance |
| **`security-auditor`** | Security specialist; reviews code and infrastructure configs for security vulnerabilities and compliance issues. | Static analysis, threat modeling, vulnerability detection |
| **`investigator`** | Root cause debugger; traces complex regressions, analyzes core dumps, and resolves production incidents. | Regressions, dump analysis, telemetry investigation |
| **`doc-writer`** | Technical writer; maintains project documentation, API references, guides, and user manuals. | Code docstring coverage, technical guides, system diagrams |
| **`release-notes`** | Release notes compiler; synthesizes git commits, pull requests, and changelogs into customer-facing release notes. | Release automation, user-facing summaries, milestone tracking |
| **`eng-manager`** | Engineering supervisor; tracks team throughput, manages roadmap risk, and optimizes delivery processes. | Team throughput, risk mitigation, alignment |
| **`researcher`** | Domain researcher; crawls, synthesizes, and compiles structured market, domain, or technology surveys. | Market analysis, tech stack research, trend reporting |

---

## 🛠️ Shared Skills

Shared skills are modular capabilities that can be dynamically linked into any agent's toolbox.

| Skill | Purpose | Key Tools Provided |
| :--- | :--- | :--- |
| **`agent-recovery`** | Automatically detects and recovers Scion agents from stuck or transient error states. | Health-check diagnostics, auto-restarts, token refreshes |
| **`scheduler`** | Schedules future events and cron-style recurring messages/actions across the team. | Cron expression parsing, task queuing, scheduled message triggers |
| **`pr-code-review`** | Connects to GitHub API to pull files, post inline PR comments, and automate review workflows. | GitHub PR integration, automated lints, review comments |
| **`release-notes-daily`** | Compiles daily digests of work finished, pull requests merged, and open-source activities. | Daily progress synthesis, commit parser, changelog compiler |
| **`changelog-parallel-backfill`** | Efficiently populates historical project changelogs in parallel across many historical branches. | Parallel execution orchestrator, backfill scripts |
| **`docs-update`** | Evaluates documentation health and automates documentation updates across the repository. | Doc health scanner, markdown links formatter, content sync |

---

## 🚀 How to Use

### Starting an Agent from this Repository
You can run any agent in this repository directly using the Scion CLI by referencing the GitHub URI:

```shell
scion start coordinator-agent --type gh://scion-frontiers/agent-team/templates/coordinator
```

### Developing Locally
If you are modifying templates or skills locally:

1. Clone this repository:
   ```shell
   git clone https://github.com/scion-frontiers/agent-team.git
   cd agent-team
   ```
2. Start an agent using your local template path to test your changes:
   ```shell
   scion start test-agent --type ./templates/coordinator
   ```

---

## 🤝 Contributing

Contributions of new agent roles and reusable skills are highly welcome! Please see [docs/contributing.md](docs/contributing.md) and [docs/code-of-conduct.md](docs/code-of-conduct.md) for details on our contributor license agreement and community standards.

### Contribution Checklist
- All code modifications or script edits must include appropriate copyright headers (see `docs/contributing.md`).
- Ensure any reusable skill scripts in `/skills` are fully documented with a corresponding `SKILL.md` file.
- Verify that your agent templates under `/templates` contain a descriptive `scion-agent.yaml`, along with robust `agents.md` and `system-prompt.md` documents.
