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
    ├── code-review/
    └── ...
```

1. **Templates (`/templates`)**: Standard definitions for various agent roles (e.g., Coordinator, Developer, Architect, QA Tester). Each template contains:
   - `scion-agent.yaml`: Configures the agent type, description, model harness, and imported skills.
   - `agents.md`: High-level role instructions and operational protocols for the agent.
   - `system-prompt.md`: Base system prompt establishing behavioral constraints, environment paths, and execution context.
2. **Skills (`/skills`)**: Independent, reusable capability modules (scripts, tools, and manuals) that can be imported by any agent to perform specialized actions (e.g., `code-review`).

---

## 🤖 Agent Templates

Each folder under `/templates` defines a distinct agent persona with specialized capabilities.

| Role / Template | Description | Key Focus Area |
| :--- | :--- | :--- |
| **`coordinator`** | Project coordinator; manages agent teams, delegates implementation, and drives progress autonomously. | Decomposing requirements, orchestrating workers, status tracking |
| **`architect`** | System designer; defines software architecture, APIs, data schemas, and implementation blueprints. | Technical design, API specs, component breakdown |
| **`developer`** | Software engineer; implements features, refactors code, and writes unit/integration tests. | Feature development, code quality, unit testing |
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

**`skills/` is a public registry, not an internal detail of the templates in this repository.** Every directory under `skills/` is addressable by URI from any template — one of ours, or one maintained by someone else entirely.

| Skill | Purpose | Key Tools Provided |
| :--- | :--- | :--- |
| **`code-review`** | Review method for a single change: the five axes, the severity vocabulary, the output format, and how to get a clean PR diff in a container. | Five-axis framework, severity labels, review template, `gh`/`git` diff recipes |
| **`release-notes-daily`** | Compiles daily digests of work finished, pull requests merged, and open-source activities. | Daily progress synthesis, commit parser, changelog compiler |
| **`changelog-parallel-backfill`** | Efficiently populates historical project changelogs in parallel across many historical branches. | Parallel execution orchestrator, backfill scripts |
| **`docs-update`** | Evaluates documentation health and automates documentation updates across the repository. | Doc health scanner, markdown links formatter, content sync |
| **`artifact-durability`** | Where a produced artifact has to end up to survive the agent that produced it, and what to do when no durable destination is available. | Durability property definition, durable-destination list, fallback ladder |
| **`software-engineering-process`** | How multi-agent engineering work is owned, sized, delegated, and verified from request to delivery. | Ownership tree, sizing tiers, the review cycle, briefing rules |
| **`gcs-artifact-publishing`** | Publishes human-reviewable artifacts to the shared GCS artifact bucket and returns a public link. | Bucket/auth setup, `gsutil` upload recipes, public-access verification, HTML renderers |

The table lists every directory under `skills/`.

### Referencing a skill by URI

A template imports a skill by listing its URI under `skills:` in that template's `scion-agent.yaml`. The last path segment is the bundle's **directory name**, and it omits the `skills/` prefix:

```yaml
skills:
  - uri: "gh://scion-frontiers/agent-team/code-review"
```

That is the live declaration in `templates/code-reviewer/scion-agent.yaml`, and it resolves to `skills/code-review/` in this repository. The same form works from a template maintained anywhere — a consuming template does not have to live here.

### A bundle nothing here references is not an unused bundle

Some bundles above are referenced by no template in this repository. **That is an expected state, not a defect.** Because the registry is addressable from outside, a bundle none of our own templates import may still be imported by templates maintained elsewhere.

So, for anyone doing a tidy-up: **do not read "not referenced by any template in this repo" as evidence that a bundle is unused, unmaintained, or safe to delete or rename.** The references that keep a bundle alive are not visible from inside this repository, and this repository deliberately does not track who consumes it.

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
