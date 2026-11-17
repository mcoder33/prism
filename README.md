<div align="center">

# 🔷 PRISM

**Stop letting your AI agent one-shot 2000-line changes it doesn't understand.**

*Recursive decomposition workflow for AI coding agents — turn a vague idea into an implemented,
verified change through small, confirmable steps.*

[![CI](https://github.com/mcoder33/prism/actions/workflows/ci.yml/badge.svg)](https://github.com/mcoder33/prism/actions/workflows/ci.yml)
[![Go](https://img.shields.io/badge/Go-1.26+-00ADD8?logo=go&logoColor=white)](https://go.dev)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](LICENSE)
[![PRs Welcome](https://img.shields.io/badge/PRs-welcome-brightgreen.svg)](https://github.com/mcoder33/prism/pulls)
[![Works with: Claude Code](https://img.shields.io/badge/Claude_Code-✓-d97757)](#supported-tools)
[![Works with: Cursor](https://img.shields.io/badge/Cursor-✓-000000)](#supported-tools)
[![Works with: 9 agents](https://img.shields.io/badge/+7_more_agents-✓-7c4ff4)](#supported-tools)

</div>

<div align="center">

![prism init demo](assets/demo.gif)

</div>

## The problem

You ask an AI agent for a non-trivial feature. It either:

- **dives straight into code** and you discover its (wrong) design decisions in a 1500-line diff, or
- **writes a 10-page design doc** that you skim, approve without really reading, and that drifts
  from reality by step three.

Both fail the same way: **the unit of review is too big**. Humans are great at judging one small,
concrete decision — and terrible at auditing a wall of text or a mega-diff.

## The PRISM way

PRISM splits design into **small recursive nodes** — one digestible `node.md` at a time, each with
an explicit decision gate. The agent proposes, you react, it drills deeper. Nothing gets implemented
until every node is **atomic** — small enough to be obviously right.

```mermaid
flowchart LR
    P[propose] --> D[decompose]
    D --> R[drill]
    R -->|until atomic| R
    R --> I[integrate]
    I --> A[apply]
    A --> V[verify]
    V --> X[archive]
```

The core principles, baked into every command:

- **Decision-first, not analysis-first.** Not *"here are 3 options, you weigh them"* —
  but *"I propose X because Y; rejected B/C in one line."* One thing to react to.
- **One small node at a time.** You never face the whole analysis at once.
- **Recursive decomposition** until a node becomes obvious.
- **Grounded in real code** — the agent reads your codebase before designing, never in a vacuum.
- **No big upfront docs.** They overwhelm and get skimmed. Forbidden by convention.
- **Sized to the change.** A `small`-tier change is one atomic node: two gates total, no
  decomposition theater, `integrate` skipped entirely. The full flow is for changes that earn it.

Every node carries an explicit status (⚪ proposed → 🟡 in progress → 🟢 atomic → 🔵 applied →
✅ verified, plus ⏸ for user-deferred parts) with formally defined transitions — a killed
session resumes exactly where it stopped, and `status` offers a gated repair when the state
table and reality disagree.

## Quick start

```bash
# Go toolchain:
go install github.com/mcoder33/prism@latest

# or, no Go required — prebuilt binary (Linux/macOS):
curl -fsSL https://raw.githubusercontent.com/mcoder33/prism/main/install.sh | sh
# (Windows: grab the .zip from the Releases page.)

cd your-project
prism init                     # interactive TUI — pick your agents (press d to pick delivery)
prism init --delivery skills   # or: Agent Skills folders, conventions bundled — commit and share
```

Then, inside your agent (shown as Claude Code commands; skills are `/prism-propose`, and
`$prism-propose` in Codex):

```
/prism:propose      # describe the problem, get a 1-screen seed proposal (+ change tier)
/prism:decompose    # split into 2-4 small parts (a single part for small-tier changes)
/prism:drill        # drill ONE part until atomic (spec, detail, diagram, tasks)
/prism:integrate    # cross-part wiring: diagram + combined signatures (skipped for small)
/prism:apply        # implement part-by-part, one commit per part, checks after each
/prism:verify       # pedantic QA on a running dev env: tests, smoke, concurrency, load
/prism:status       # where am I — phase, node table vs reality, the one next action
```

The CLI also reads what the agent writes: `prism status` cross-checks the state table against
the node files, and `prism check` validates the artifacts (drawio XML, spec scenarios, Open
tags, atomicity bounds) — deterministic checks the agent runs instead of eyeballing Markdown.

All design artifacts live in `.prism/` at the repo root — **git-excluded automatically**,
it's local working state. (Optionally a project can commit `.prism/archive/` as shared
reference — see the conventions installed by `prism init`.)

## First time with PRISM?

1. Skim [the principles](#the-prism-way) above.
2. Open the worked example — a complete small-tier change, every artifact in its final state:
   [`.prism/archive/example-json-list/`](.prism/archive/example-json-list/). Start with its
   [`README.md`](.prism/archive/example-json-list/README.md) (the status table), then
   [`proposal.md`](.prism/archive/example-json-list/proposal.md) and the
   [`01-json-flag/`](.prism/archive/example-json-list/01-json-flag/) node.
3. Run `/prism:propose` in your agent and let the gates walk you through your first change.

## Inside a drilled node

`/prism:drill` takes **one part** and brings it to atomic. It writes the `node.md` digest first
and **stops at a gate** — you react to the digest and the proposed artifact set before anything
else is generated (trivial nodes skip the diagram and spec). Here's what that looks like for a
`01-token-bucket` node of a rate-limiter change — three artifacts carry the weight:

**`node.md` — the unit of review.** 5-7 lines. You react to *this*, not to a design doc:

<p align="center"><img src="assets/drill-node.png" width="680" alt="node.md: What / Logic / Guarantees / Input-output digest, status, open questions"></p>

**`concept.drawio` — a real diagram, not ASCII art.** The agent hand-crafts mxGraph XML and
validates it with `prism check`. Open it in draw.io or the VS Code extension:

<p align="center"><img src="assets/drill-concept.png" width="760" alt="concept.drawio: request flows through middleware to Limiter.Allow, lazy refill, allow/deny branches"></p>

**`signatures.md` — the API before any code.** Signatures + what/why comments, no implementation.
Catch a bad interface here, where changing it costs nothing:

<p align="center"><img src="assets/drill-signatures.png" width="680" alt="signatures.md: Limiter interface, Decision struct, bucket internals with rationale comments"></p>

Plus `spec.md` (Requirement/Scenario — they drive the tests later) and `tasks.md` (the checklist
`apply` executes). If a node turns out too big — drill redirects to `decompose` instead of
patching on the fly.

## What you get per change

```
.prism/<change>/
├── proposal.md          Why / What / Constraints / Decisions / Non-goals  (< 1 screen)
├── concept.md           candidate + chosen/rejected strategies
├── data-flow.drawio     conceptual data-mutation chain
├── 01-parser/           ← a node
│   ├── node.md          5-7 line digest — the unit of review
│   ├── spec.md          requirements as Requirement/Scenario
│   ├── detail.md        decision-complete implementation plan
│   ├── signatures.md    code sketch: signatures + what/why
│   └── tasks.md         checklist the agent executes
├── 02-renderer/         ← another node (drills into 02a, 02b... if needed)
├── integration.drawio   how the parts connect
└── tasks.md             root order + cross-cutting concerns only
```

Switch between concurrent changes like git branches: `/prism:use` persists the active change
in `.prism/CURRENT`, and every command targets it automatically.

## Supported tools

One `prism init` installs the workflow for every agent you use — same methodology, shared
`.prism/` state, so you can **propose in Claude Code and apply in Cursor**. Two deliveries:

- **commands** (default where the tool has them) — one slash-command file per workflow; they
  read the shared `.prism/conventions.md`.
- **skills** (`--delivery skills`) — one [Agent Skills](https://agentskills.io) folder per
  workflow, `prism-<name>/SKILL.md` plus a bundled `references/conventions.md`. Self-contained:
  commit them and teammates need nothing else. In Claude Code and Codex the skills are
  user-invoked only — the model never starts a gated workflow on its own.

| Tool | Commands → invocation | Skills → invocation |
|---|---|---|
| **Claude Code** | `.claude/commands/prism/` → `/prism:propose` | `.claude/skills/` → `/prism-propose` |
| **Cursor** | `.cursor/commands/` → `/prism-propose` | `.cursor/skills/` → `/prism-propose` |
| **Codex CLI** | — (skills only) | `.agents/skills/` → `$prism-propose` |
| **Gemini CLI** | `.gemini/commands/prism/` → `/prism:propose` | `.gemini/skills/` → `/prism-propose` |
| **GitHub Copilot** | `.github/prompts/` → `/prism-propose` | `.github/skills/` → `/prism-propose` |
| **Windsurf** | `.windsurf/workflows/` → `/prism-propose` | — |
| **OpenCode** | `.opencode/command/` → `/prism-propose` | `.opencode/skills/` → `/prism-propose` |
| **Roo Code** | `.roo/commands/` → `/prism-propose` | `.roo/skills/` → `/prism-propose` |
| **Cline** | `.clinerules/workflows/` → `/prism-propose.md` | `.cline/skills/` → `/prism-propose` |

**Codex is skills-only**: Codex deprecated custom prompts and only ever read them from
`~/.codex/prompts/`, so the project-level `.codex/prompts/` files prism ≤ 0.7 wrote were never
loaded. `prism update` migrates them to `.agents/skills/` and removes the old files.

Switch an installed project with `prism update --delivery skills|commands` — the other
delivery's generated files are removed; hand-written files are never touched.

Adding a tool is one `adapters.Tool` value ([internal/adapters/adapters.go](internal/adapters/adapters.go)) —
paths, naming, frontmatter format, skills folder. PRs welcome.

## Team setup

**Simplest: skills.** Run `prism init --delivery skills` and commit the generated skill folders.
Each bundles its own conventions, so a teammate's clone works with no CLI installed. Only your
in-flight change artifacts stay local in `.prism/`.

**With commands**, two kinds of files have two different fates — this trips people up:

- **Generated command files** (`.claude/commands/prism/`, `.cursor/commands/`, …) are normal repo
  files. **Commit them** so every teammate gets the slash commands on clone, without installing anything.
- **`.prism/`** (the `conventions.md` the commands read, plus your in-flight change artifacts) is
  **git-excluded** — local working state. On a fresh clone it isn't there.

So a teammate who clones a repo with committed commands still needs the shared conventions locally.
The one-time bootstrap: install the CLI and run **`prism update`** (or `prism init`) once — it
regenerates `.prism/conventions.md` in place. Run **`prism doctor`** to confirm everything lines up
(versions match, `.prism/` is git-excluded, no stale `CURRENT`). If you'd rather not have teammates
install the CLI at all, commit `.prism/conventions.md` too (see the archive-sharing note in the
installed conventions for how to un-exclude selectively).

## How is this different from OpenSpec / spec-kit?

[OpenSpec](https://github.com/Fission-AI/OpenSpec) and [spec-kit](https://github.com/github/spec-kit)
pioneered spec-driven AI development; PRISM follows them in shipping as Agent Skills.
The difference is the shape of the artifact:

|  | Spec-driven (OpenSpec, spec-kit) | **PRISM** |
|---|---|---|
| Unit of review | a spec document | a **5-7 line node** |
| Structure | flat: spec → tasks | **recursive tree**, drill until atomic |
| Gates | approve the generated document | **hard stop per node**; a fork is closed only by a chosen outcome |
| Depth control | fixed | per-node — drill only where it's murky |
| Post-implementation | done at merge | **`verify`** at `--checks low\|high`, independent design review |

If your changes are small, spec-driven is enough. PRISM earns its keep on changes where the
design itself is the risk: refactors, new subsystems, concurrency, integrations.

## CLI reference

```bash
prism init [path] [--tools claude,cursor|all] [--delivery skills|commands]
                                                # install (TUI without --tools)
prism update [path] [--force] [--delivery …]    # regenerate after a CLI upgrade / switch delivery
prism list [path]                               # list active changes in .prism/
prism status [path] [--change X] [--json]       # state table cross-checked against the files
prism check [path] [--change X] [--json]        # validate artifacts; exits 1 on errors
prism doctor [path]                             # diagnose: version drift, old layouts, stale pointer
prism uninstall [path] [--tools …] [--shared]   # remove generated files (keeps .prism/ work)
```

Generated files are **tool-owned**: each carries a `prism:generated v<version>` stamp and is
overwritten wholesale by `prism update` — edit [templates/](templates/), not the output.
How to write template text (budgets, style, rejected framings):
[docs/template-style.md](docs/template-style.md).

Single static binary, zero runtime dependencies, templates embedded via `go:embed`.

## Development

```bash
make ci          # everything CI runs: test -race, vet, lint, vuln, mod-check
make test        # just the tests
make cover       # tests + coverage summary
make install     # build + put `prism` on PATH
make help        # list all targets

prism init --tools claude /tmp/sandbox        # smoke test
```

CI (GitHub Actions) runs the same five checks on every push and PR — the test job across a
Linux/macOS/Windows matrix; `lint` and `vuln` auto-install their tools, so a fresh clone only
needs Go. Pushing a `v*` tag runs [GoReleaser](https://goreleaser.com) to publish cross-platform
binaries to the release (see [`.goreleaser.yaml`](.goreleaser.yaml)).

Layout:

```
main.go                      entrypoint → internal/cli
templates/                   tool-neutral command bodies + shared conventions (go:embed)
internal/
  workflows/                 command registry (id, title, description) + Version
  adapters/                  per-tool: command + skill paths, invocation naming, frontmatter
  installer/                 detection, delivery migration, writing, git-exclude
  change/                    reads .prism/<change>/: state table, findings (status, check)
  cli/                       cobra commands + bubbletea TUI
```

Stack: Go stdlib + [cobra](https://github.com/spf13/cobra) +
[bubbletea](https://github.com/charmbracelet/bubbletea)/[lipgloss](https://github.com/charmbracelet/lipgloss).

## License

[MIT](LICENSE)

---

<div align="center">

**If PRISM saves you from one bad mega-diff, consider a ⭐ — it helps others find it.**

</div>
