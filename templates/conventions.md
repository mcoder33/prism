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

### Current change pointer (`.prism/CURRENT`)

`.prism/CURRENT` holds the slug of the **active change** — the prism analog of the current git
branch (`.git/HEAD`). It is set/switched/cleared by `prism:use`, set by `prism:propose`, and
cleared by `prism:archive` when the active change is archived. Persisted on disk, so it survives
across sessions.

**Every prism command resolves `<change>` from `.prism/CURRENT` when no explicit name is
given** (if it's also empty, the command asks — it never guesses). So once you run `prism:use <change>`,
all subsequent decomposition work and design writing target that change automatically.

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

**Status:** ⚪ … | 🟡 … | 🟢 …

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

## Statuses (in README)

⚪ not started · 🟡 in progress · 🟢 atomic, artifacts ready · 🔵 applied (committed) · ✅ verified

### README.md (change root)

```
# <change>

**Phase:** propose → **decompose** → drill → integrate → apply → verify

| Node  | Title | Status | Open    |
|-------|-------|--------|---------|
| 01-…  | …     | 🟢     | 1 minor |

Links: [concept.md](concept.md) · [data-flow.drawio](data-flow.drawio)
```

Bold the **current** phase. Every command updates Phase and the table as its **last** file
write — the table is the resume point for future sessions (`prism:status` reads it); if it
lies, resume lies.

## Formatting (mandatory)

- Each item on its own line (bulleted lists `-`, blank lines between sections).
- Do NOT chain multiple `**Label:**` entries without separators — they collapse into one paragraph and are hard to read.

## drawio

- Hand-craft mxGraph (`<mxfile><diagram><mxGraphModel><root>…`). Nodes `vertex="1"`, edges `edge="1"`.
- Avoid raw `&`, `<`, `>` in labels.
- **After writing always** validate: `xmllint --noout <file>.drawio`. If xmllint is missing:
  `python3 -c "import xml.dom.minidom,sys; xml.dom.minidom.parse(sys.argv[1])" <file>.drawio`;
  if neither is available — re-read the XML carefully for unclosed tags and raw entities.

Two root diagrams, different purpose (don't conflate):

- **`data-flow.drawio`** — how **data** mutates (conceptual chain; nodes = project types/pseudocode,
  edges = transformations). Made in `prism:propose`, before decomposition.
- **`integration.drawio`** — how **parts** connect / the call graph (who calls whom, `[NN]` annotations).
  Made in `prism:integrate`, once parts are drilled.
