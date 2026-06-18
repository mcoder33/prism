# Changelog

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
