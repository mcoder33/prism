# Changelog

## 0.8.0 — 2026-11-17

Distribution and determinism release. Agent Skills became the cross-tool standard (OpenSpec
2.0 went skills-only; Codex deprecated custom prompts), so prism now ships as skills too. Checks
the agent used to do by reading Markdown move into the CLI. Methodology changes are informed by
a survey of mcoder33/codex-skills (`effort`/`checks`, `nasa`, unresolved-forks), OpenSpec 2.x
and mattpocock/skills (`implement-spec`, `grilling`).

### Fixed

- **Codex CLI support.** prism ≤ 0.7 wrote `.codex/prompts/prism-*.md` into the project, but
  Codex only ever read prompts from `~/.codex/prompts/` and has since deprecated them — the
  commands never loaded. Codex is now skills-only (`.agents/skills/prism-*/`, invoked as
  `$prism-<name>`, implicit invocation disabled via `agents/openai.yaml`); `prism update`
  migrates and removes the old files.

### Added

- **Skills delivery** (`--delivery skills` on `init`/`update`, `d` in the TUI): one Agent Skills
  folder per workflow with a bundled `references/conventions.md`, so committed skills work on a
  fresh clone without the CLI. Claude Code, Cursor, Codex, Gemini CLI, Copilot, OpenCode,
  Roo Code and Cline. Commands stay the default where a tool has them.
- **`prism update --delivery`** switches delivery; the other delivery's and legacy generated
  files are removed (stamp-checked — hand-written files are kept). `doctor` reports delivery and
  leftover layouts.
- **`prism status [--json]`** — the change's state table parsed and cross-checked against the
  node files: status mismatches, rows without dirs and vice versa, missing mandatory artifacts,
  untagged Opens, unchecked tasks on applied nodes, branch mismatch. The `status` workflow
  starts from it.
- **`prism check [--json]`** — artifact validation, exit 1 on errors: well-formed drawio XML
  (xmllint is no longer needed), WHEN/THEN in every spec Scenario, atomicity bounds (spec ≤ ~3
  Requirements, tasks ≤ ~12 boxes), plus the status findings. `drill` and `apply` run it.
  The shipped worked example is checked in CI.
- **`verify --checks low|high`** — evidence depth independent of implementation size; default
  `low` for small-tier changes, `high` otherwise.
- **Independent design-conformance review** in `verify`: run by a fresh-context subagent when
  the agent supports one, otherwise labelled `self-review`.
- **`apply --parallel`** (opt-in, standard tier): the ready frontier of parts implemented by
  fresh-context subagents in separate worktrees, merged one at a time, flipped 🔵 only once on
  the feature branch.
- **Forecast** line written by `apply` before the first part; `archive`'s mini-retro compares
  it with the actual diff to calibrate the atomicity bounds.

### Changed

- **Closing a fork** (conventions, Open tags): a `[blocking]` Open closes only on a chosen
  outcome for the whole affected scope — "decide during testing", an owner, a deadline or
  equal-cost options don't close it.
- **Simplest complete design** in `drill`: every mechanism beyond the direct path names the
  requirement or concrete failure that needs it.
- `propose` grill: questions are worded so "yes" accepts the recommended answer.
- `conventions.md` trimmed to stay inside its byte budget: archive sharing and active-change
  handoff moved into `archive` and `use`.
- Templates reference conventions through `{{conventions}}`; lint rejects hardcoded paths and
  unknown placeholders.

## 0.7.0 — 2026-07-17

Methodology sharpening release: split heuristics, escape hatches for stuck decisions, and a
documented style bar for the templates themselves. Informed by a survey of mattpocock/skills
and an internal agent kit.

### Changed

- **`decompose` prefers vertical slices over layers**: 01 is now a *tracer bullet* — the
  thinnest end-to-end slice proving the pipeline — instead of "01 = foundation"; an
  infrastructure-only part must be shared by 2+ parts to earn its slot; wide mechanical
  refactors split as *expand–contract* (add new path → migrate callers → remove old path).
- **`verify` design conformance reports on two separate axes** — *spec* (does the behaviour do
  what was asked) and *design-shape* (does the implementation match the artifacts) — so one
  can't mask the other.
- **`archive` mini-retro classifies lessons** before routing them: *project* → CLAUDE.md/
  AGENTS.md one-liners, *methodology* → candidate feedback for prism itself, *one-off* → drop.
- **Conventions**: an Open item must pass the **fog test** (a precisely formulated question,
  not a vague unease); `detail.md` gains a checkable completion criterion ("a fresh
  implementer would ask zero questions").

### Added

- **Spike escape hatch** in `propose` and `drill`: a question neither the code nor the user
  can settle — only an experiment can — becomes a small throwaway prototype answering exactly
  that one question; the verdict is recorded, the code discarded.
- **Frontier rounds** in the `propose` grill: on the user's request, ask every question whose
  prerequisite decisions are already made per round, instead of strictly one at a time.
- **`docs/template-style.md`** — the style bar for template contributions: context/cognitive
  budgets, progressive disclosure, leading words, positive phrasing, checkable completion
  criteria, the no-op test, rejected framings, tool neutrality.

## 0.6.0 — 2026-07-07

Distribution and installation-lifecycle release: prism is now installable without a Go
toolchain, ships two more agent adapters, and grows the diagnose/uninstall/abandon paths that
were missing from the install and change lifecycles.

### Added

- **`prism doctor`** — read-only diagnosis of an installation: per-tool version drift vs the CLI,
  `.prism/` git-exclude status, a stale `CURRENT` pointer, active-change inventory, and whether
  `xmllint` (drawio validation) is present. Changes nothing.
- **`prism uninstall`** — removes prism-generated command files (identified by their
  `prism:generated` stamp, so hand-written files at the same path are left alone) and cleans up
  emptied command dirs. Keeps `.prism/` design work by default; `--shared` also drops
  `conventions.md`; `--tools` scopes to specific agents.
- **Two agent adapters**: **Roo Code** (`.roo/commands/`) and **Cline**
  (`.clinerules/workflows/`, invoked as `/prism-<cmd>.md`) — 9 agents supported.
- **Prebuilt binaries**: GoReleaser config + a tag-triggered release workflow, plus `install.sh`
  for a no-Go-toolchain install on Linux/macOS. Homebrew tap wiring is included but commented
  until the tap repo exists.
- **`--abandon` mode** for the archive command — the missing exit for a dead-end change: archives
  it with an `abandoned` status and reason, skipping the task-completion check and mini-retro.
- CI now runs the test suite on a **Linux/macOS/Windows** matrix (the path and git-exclude logic
  is cross-platform-sensitive).

### Changed

- `conventions.md` documents that `.prism/CURRENT` is **one pointer per checkout** and that
  parallel changes want a git worktree each.
- README gains a **Team setup** section (which files to commit vs the git-excluded `.prism/`) and
  documents the new install options and commands.

## 0.5.1 — 2026-07-07

The best-practices survey in `propose` is now unconditional and chat-only: general reasoning
about the problem class always happens first and shapes everything downstream, instead of being
an optional persisted artifact.

### Changed

- **Best-practices step is always-on** — the interactive "survey or skip" question is gone;
  `propose` always leads with a short breadth survey (3–6 bullets) before any strategy is on
  the table.
- **Console-only** — the survey lives in chat and is never persisted (`concept.md` no longer
  has a `## Best practices` section); its value is in shaping what comes next, not in an
  artifact.
- The survey explicitly **feeds downstream**: the step-2 strategy choice cites the practice
  that tips it, and the grill's recommended answers are grounded in the step-1 survey and the
  code.

## 0.5.0 — 2026-06-18

`propose` is now decision-first end to end: it leads with a proposal and refines it, instead of
interrogating the user before anything concrete is on the table.

### Changed

- **Reordered `propose`**: best practices now come **first** (skippable breadth) and feed an
  **initial proposal** (strategy + data-flow sketch + draft seed) — one concrete thing to react
  to. The requirements **grill** moves *after* that, as a one-question-at-a-time loop that
  **evolves** the proposal, ending in a single confirm gate. Previously the grill ran first and
  the proposal was assembled only afterwards.
- **One propose gate** for both tiers (was: separate strategy + data-flow gates). The grill is
  the iterative reaction loop; the gate confirms the whole evolved proposal at once.
- `data-flow.drawio` is hand-crafted at the confirm gate (once the direction is settled) rather
  than mid-flow, so the grill can't strand a reworked diagram. `small` tier still keeps the data
  flow as text in `concept.md`.

## 0.4.0 — 2026-06-10

A methodology overhaul: the flow now scales to the change, state is a formal machine, and
mid-flow reality (deferrals, scope changes, failed verification) has first-class paths.

### Added

- **Change tiers** (`small | standard`): a small change is one atomic node — merged propose
  gate, single-part decompose, `integrate` skipped entirely; two gates instead of seven.
- **Formal state model**: legal status transitions with owners and exact flip moments;
  new ⏸ (deferred) status; write-ordering rules that make crash-resume deterministic.
- **Open tags**: `[blocking]` / `[minor]` / `[deferred: reason — user, date]` with a defined
  un-defer path (drill confirms, unparks, re-drills).
- **Revision rules**: amendment mode in `propose` for mid-flow scope changes; applied (🔵)
  nodes are immutable history; small→standard promotion paths.
- **Status repair**: `status` cross-checks the table against reality and offers a gated
  repair plan (with exceptions protecting verify escalations and reverted parts).
- **Verify hardening**: defined fix cycles, `FAIL (design)` escalation (🔵 → 🟡 + re-drill
  route), scenario→test mapping via a `## Coverage` section in spec.md, persistent
  `.prism/<change>/verify.md` report with run history, headless-run fallbacks.
- **Branch tracking**: `apply` records `**Branch:**` in the change README; `use`/`status`/
  `apply` warn on mismatch (never switch themselves).
- **Archive lifecycle**: `-rN` re-archive on collision, hotfix-via-new-change rule, optional
  committed `.prism/archive/` for team reference, active-change handoff notes.
- **Worked example**: `.prism/archive/example-json-list/` — a complete small-tier change,
  every artifact in its final state; linked from the README "First time?" section.
- **Methodology lint** (`internal/workflows/lint_test.go`): cross-references, status glyphs,
  transition legality, `{{cmd:*}}` ids and size budgets are now CI-checked.

### Changed

- `conventions.md` restructured (Change tiers, Open tags, Statuses and transitions,
  Revision rules as dedicated sections); `node.md` template lists the full glyph set.
- `apply` documents integration with code review (verify-then-PR / PR-then-verify-in-CI).

## 0.3.0

Hardened the methodology: gates, atomicity criteria, state model, status command.

## 0.2.0

Rewritten in Go; interactive TUI for `prism init`; 7 agent adapters.

## 0.1.0

Initial release.
