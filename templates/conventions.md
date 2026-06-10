# prism — shared conventions

> Read this first when running any prism command. It defines the artifact formats,
> layout, and rules so the individual commands stay short.
>
> Command naming depends on your agent: namespaced (`/prism:use`) or flat (`/prism-use`).
> This document uses the `prism:<name>` form; map it to your agent's naming.

Methodology: **recursive decomposition of a problem into small nodes**, laid out across files/directories.
We move node-by-node with gates. Large upfront docs are forbidden — they overwhelm and get skimmed.
Full reference example: any archived change under `.prism/archive/` — mirror it.

## Principles

- **Decision-first, not analysis-first.** Not "here are 3 options, weigh them" — but "I propose X,
  because Y; rejected B/C in one line". One thing to react to.
- **One small node at a time.** At any moment there's one `node.md` in front of you, not the whole analysis.
- **Progressive disclosure.** Core idea at the top; deep breakdown — in a separate file/section, on request.
- **Decompose recursively** until a node becomes obvious ("atomic").
- **Ground in real code** (symbol-overview / find-symbol tools if available, otherwise grep) —
  don't design in a vacuum.

## Gates

A GATE is a hard stop. When a command step is marked **GATE**:

1. Print the thing to review **inline in chat** — the user must not have to open files.
2. **End your turn.** Do not create or modify any file, do not begin the next step.
3. Resume only on an explicit user reply. Your own judgement never counts as approval.

## Change tiers

Set at the end of `prism:propose` (its gate includes the tier); recorded in `README.md` under
the Phase line as `**Tier:** small | standard`. No Tier line (older change) = standard.

- **small** — the change is expected to be ONE atomic node. propose merges the strategy and
  data-flow gates into one (data flow may be text/pseudocode in `concept.md`, no
  `data-flow.drawio`). decompose may conclude "single part" (a lone `01-…/`). integrate is
  SKIPPED entirely — no root artifacts; the part's own `tasks.md` is the order. Flow:
  drill → apply directly.
- **standard** — full flow: 2–4 parts, full integrate.

Escape hatch: if drilling a small change trips the atomicity bounds, promote it to standard
(announce it, run `prism:decompose` properly). Never demote mid-flow.

## Atomicity — when to stop drilling

A node is atomic when ALL of these hold:

- **One responsibility** — its What line needs no "and" joining unrelated behaviours.
- **Decision-complete** — no open choice left that would change a signature (`[minor]` opens only).
- **Independently testable** — its spec scenarios can pass with sibling nodes stubbed.
- **Bounded** — est. ≤ ~400 changed LOC · ≤ ~5 files · tasks.md ≤ ~12 boxes · detail.md ≤ 1 screen.

The numbers are trip-wires, not validation rules: if two or more size bounds trip, or any of the
first three fail — decompose further instead of writing artifacts.

## Node artifact tiers

`node.md` + `tasks.md` are mandatory for every node. The rest scale with node complexity:

- **Trivial** (rename/move, mechanical change, no real decisions) — skip `spec.md` and `concept.drawio`.
- **Standard** — add `detail.md` + `signatures.md`.
- **Complex** (branching logic, invariants, concurrency, new API) — full set incl. `concept.drawio`.

The drill GATE proposes the set with a one-line reason per skipped artifact; the user confirms.

## Layout

All artifacts live under `.prism/` at the repo root. The directory is **created automatically**
if missing, and `.prism/` is added to `.git/info/exclude` (artifacts are never committed).
`<change>` is a short kebab-case slug of the problem.

- **Active** change: `.prism/<change>/`.
- **Archived** change (after `prism:apply` completes): `.prism/archive/<change>/` — `apply` moves
  the whole change folder there automatically once all tasks pass.

### Sharing archives (optional, per-project decision)

`.git/info/exclude` is per-clone: by default archives die with the machine. To share archived
designs as team reference, remove `.prism` from `.git/info/exclude` and add to `.gitignore`:

    .prism/*
    !.prism/archive/

Active changes stay local either way; durable lessons go to committed agent docs via the
archive mini-retro regardless.

Handing off an ACTIVE change to a teammate: copy `.prism/<change>/` to them out of band
(temporary branch, archive file); they drop it into their `.prism/`, run `prism:use <change>`,
then `prism:status` to resume.

### Current change pointer (`.prism/CURRENT`)

`.prism/CURRENT` holds the slug of the **active change** — the prism analog of the current git
branch (`.git/HEAD`). Set/switched/cleared by `prism:use`, set by `prism:propose`, cleared by
`prism:archive` when the active change is archived. Persisted on disk across sessions.

**Every prism command resolves `<change>` from `.prism/CURRENT` when no explicit name is
given** (if it's also empty, the command asks — it never guesses). So once you run `prism:use <change>`,
all subsequent decomposition work and design writing target that change automatically.

If `.prism/CURRENT` names a slug with no matching `.prism/<change>/` directory, the pointer is
**stale**: say so, don't invent state, and route to `prism:use` — the only command that repairs
the pointer (rewrites or clears it).

```
.prism/<change>/
├── proposal.md          seed: Why / What / Constraints+Invariants / Decisions / Non-goals  (< 1 screen)
├── concept.md           best practices + candidate strategies + chosen/rejected (made in prism:propose)
├── data-flow.drawio     conceptual data-mutation chain for the chosen strategy (xmllint)
├── README.md            tree map + status table + cycle rules
├── NN-name/             node (part) — artifact set per tier (see Node artifact tiers)
│   ├── node.md          5–7 line digest (always)
│   ├── spec.md          requirements (Requirement/Scenario)
│   ├── detail.md        how to implement (decision-complete)
│   ├── concept.drawio   diagram (mxGraph)
│   ├── signatures.md    code sketch (signatures + what/why comments)
│   └── tasks.md         checklist (always)
│       └── NNa-…/       sub-nodes when drilling further (same structure)
├── integration.drawio   overall diagram of how parts connect
├── signatures.md        combined call-graph: who calls whom + types flowing between parts
└── tasks.md             root: order + cross-cutting only (NO repetition of part details)
```

## Artifact templates

### node.md (digest, 5–7 lines)
```
# NN — <name>

- **What:** …
- **Logic:** …
- **Guarantees:** … (invariant, if any)
- **Input → output:** … → …

**Status:** <ONE of ⚪ 🟡 🟢 🔵 ✅ ⏸ — mirrors this node's row in README.md>

**Open:** …
```

### concept.md (made in prism:propose, before the seed)
```
# Concept — <change>

## Best practices
- <practice> — <when it applies> [source]
- …
(or: > Skipped — user opted out.)

## Candidate strategies
- **A. <name>** — RECOMMENDED — <2–3 lines: idea, why it fits>
- **B. <name>** — <one line, why secondary/rejected>
- **C. <name>** — <one line>
(+ user's own strategy if provided)

## Chosen strategy
<A> — because <…>. Rejected: B <one line> · C <one line>.

## Data flow
See `data-flow.drawio` — chain of how data mutates under the chosen strategy.
```

### spec.md (openspec-style, ≤ ~3 Requirements — more means the node isn't atomic)
```
# Spec — NN <name>

## Requirement: <name>
<subject> SHALL <behaviour>.

### Scenario: <name>
- WHEN <condition>
- THEN <expected result>

## Coverage
(appended by prism:apply as tasks are checked off; read and updated by prism:verify)
- Scenario: <name> → test `path/to/file::TestName`
- Scenario: <name> → smoke: <probe description>
- Scenario: <name> → untested: <reason>
```

### detail.md (≤ 1 screen)
How to implement, decision-complete: algorithm/structures, subtleties, edge-cases, worked example.
"Open (minor)" — only if genuinely open.

### signatures.md (code sketch, no implementation)
Signatures in a code block + **what/why** comments above each. Mark reused vs changed.

### tasks.md (checklist)
```
# Tasks — NN <name>

## 1. <group>
- [ ] 1.1 <step>
- [ ] 1.2 <step>
```

## Statuses and transitions (in README)

⚪ not started · 🟡 in progress · 🟢 atomic, artifacts ready · 🔵 applied (committed) ·
✅ verified · ⏸ deferred (see Open tags)

Legal transitions and their owners — no other flips are valid:

| Transition | Owner | When (exact moment) |
|------------|-------|---------------------|
| (new row) → ⚪ | decompose | after the split gate, when `NN-name/node.md` is created |
| ⚪ → 🟡 | drill | step 1, BEFORE reading code (crash-resume marker) |
| 🟡 → 🟢 | drill | after all tier artifacts are written |
| 🟢 → 🔵 | apply | immediately AFTER that part's commit succeeds |
| 🔵 → ✅ | verify | on overall PASS, all applied nodes at once |
| any → ⏸ | drill / integrate | the user explicitly defers the node at a gate |
| ⏸ → 🟡 | drill | un-defer: the user resumes the node (see Open tags) |
| 🟢 → 🟡 | drill | re-drill: design rework before apply |
| 🔵 → 🟡 | verify | escalation: confirmed design-level defect in an applied part |

Exception: `prism:status` repair (gated by the user's explicit approval) may perform
additional reality-matching flips — it reconciles the table with reality, it doesn't do work.

### README.md (change root)

```
# <change>

**Phase:** propose → **decompose** → drill → integrate → apply → verify
**Tier:** standard

| Node  | Title | Status | Open    |
|-------|-------|--------|---------|
| 01-…  | …     | 🟢     | 1 minor |

Links: [concept.md](concept.md) · [data-flow.drawio](data-flow.drawio)
```

The **Phase** line shows the furthest phase that has STARTED; each command bolds its own phase
**on entry** (together with a crash-resume flip when it has one, e.g. drill's ⚪→🟡). A command
invoked as a sub-step (e.g. apply running verify) still advances Phase — it has started.
**Write ordering**: status flips that protect resume (⚪→🟡, 🟢→🔵) and the Phase bold happen
at the moment shown in the transition table / on entry; links and Open counts are the command's
**last** file write. If a session dies, the table must already tell the truth — it is the
resume point (`prism:status` reads it); if it lies, resume lies.

`**Branch:** <name>` — added under the Tier line by `prism:apply` when it creates the feature
branch; other commands compare it against the current git branch and warn on mismatch (they
never switch branches themselves).

## Open tags

Every Open item carries exactly one tag:

- `[blocking]` — must be resolved before apply touches the owning part.
- `[minor]` — won't change a signature; may ride along into apply.
- `[deferred: <one-line reason> — user, <date>]` — was `[blocking]`; the user parked it at a
  gate. Non-blocking for preconditions. Only the USER can defer — the agent records it
  verbatim, never self-defers.

Deferring a WHOLE node = status ⏸ in `README.md` + a `[deferred: …]` line in its `node.md`.
⏸ nodes are skipped by apply, excluded from verify's design conformance, and listed as
warnings by archive.

**Un-defer**: run `prism:drill <NN>` on the ⏸ node — drill confirms the un-defer with the
user, removes the `[deferred: …]` line and flips ⏸ → 🟡 as its entry marker, then re-drills:
the parked question must be answered, not just unparked.

## Revision rules (scope changes mid-flow)

- Requirements changed after decompose/integrate → re-enter `prism:propose` in **amendment
  mode**: present the DELTA to Why/What/Decisions (gate), update `proposal.md` + `concept.md`,
  then route to `prism:decompose` for the affected parts only.
- 🔵 (applied) nodes are immutable history: a change to an applied part becomes a NEW node
  (next free NN, or a sub-node). Exception: an explicit `git revert` of the part's commit
  flips 🔵 → 🟡; record it as `[minor] reverted <short-sha>` in the node's `node.md` Open.
- Affected ⚪/🟡/🟢 nodes: re-drill (🟢 → 🟡) or delete the row if obsolete. Untouched parts
  keep their numbers and statuses.
- An amendment touching more than one part → re-run `prism:integrate` (call graph and root
  tasks must match the new shape).
- An amendment that splits a `Tier: small` change into 2+ parts promotes it to `standard`:
  update the Tier line in `README.md`; integrate becomes required.

## Formatting (mandatory)

- Each item on its own line (bulleted lists `-`, blank lines between sections).
- Do NOT chain multiple `**Label:**` entries without separators — they collapse into one paragraph and are hard to read.

## drawio

- Hand-craft mxGraph (`<mxfile><diagram><mxGraphModel><root>…`). Nodes `vertex="1"`, edges `edge="1"`.
- Avoid raw `&`, `<`, `>` in labels.
- **After writing always** validate: `xmllint --noout <file>.drawio` (fallback: parse with
  python3 `xml.dom.minidom`, or re-read carefully for unclosed tags and raw entities).

Two root diagrams, different purpose (don't conflate):

- **`data-flow.drawio`** — how **data** mutates (conceptual chain; nodes = project types/pseudocode,
  edges = transformations). Made in `prism:propose`, before decomposition.
- **`integration.drawio`** — how **parts** connect / the call graph (who calls whom, `[NN]` annotations).
  Made in `prism:integrate`, once parts are drilled.
