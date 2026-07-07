# Changelog

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
