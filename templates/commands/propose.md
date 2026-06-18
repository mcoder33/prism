Starts a new "change" in the recursive decomposition methodology.

**Read first**: `.prism/conventions.md` (formats, layout, rules; if missing, run `prism update`).

## Amendment mode (existing change)

If `proposal.md` already exists for the active change, this is an **amendment**, not a new
change: present the DELTA to Why/What/Decisions (**GATE**), update `proposal.md` + `concept.md`,
then route to {{cmd:decompose}} for the affected parts only (revision rules — in conventions).
Re-grill (step 3) the affected scope only; revisit best practices / strategy (steps 1–2) only if
the strategy itself is in question.

## Step 1 — best practices (skippable)

> Setup: pick the kebab-case `<change>` slug and create `.prism/<change>/` now (create `.prism/`
> if missing; add it to `.git/info/exclude`; write the slug to `.prism/CURRENT`) so steps 1–3 can
> persist `concept.md` in place. `proposal.md` / `README.md` come in step 4.

**Lead with breadth**, before committing to any proposal: how is this *class* of problem solved
in general?

- First ask (interactive question): do the survey, or skip (recommended option = the one that fits
  — a routine/well-understood task can skip; a wide solution space benefits from it).
- If doing it: a short, decision-first survey — **3–6 bullets**, each `practice — when it applies
  [source]`. You may use web search / docs lookup for current practices, plus your own knowledge.
  Don't dump essays; give the lay of the land.
- Persist into `concept.md` under `## Best practices` (mark `> Skipped — user opted out.` if skipped).
- These practices **feed the proposal** in step 2 — cite them when one tips a strategy choice.

## Step 2 — initial proposal (decision-first)

Put **one concrete proposal** on the table for the user to react to — grounded in the best
practices above and in **real code** (symbol-overview / find-symbol tools if available, otherwise
grep), never in a vacuum. This is a first cut, not the final word: the grill (step 3) evolves it.

- **Chosen strategy** + a few rejected alternatives in **one line each** — decision-first, not a
  menu to weigh. Invite the user to supply their own strategy too.
- **Data-mutation chain** — sketch how data changes end-to-end (conceptually, in chat + a one-line
  pointer in `concept.md`). The `data-flow.drawio` itself is hand-crafted later, at the confirm
  gate, once the direction is settled — so the grill can't strand a reworked diagram.
- **Draft seed** — a first pass at Why / What / Constraints & Invariants / Decisions / Non-goals,
  inline in chat (short — it becomes `proposal.md` at step 4).
- Persist `concept.md` now: `## Candidate strategies` + `## Chosen strategy` (chosen + rejected,
  one line each) + a `## Data flow` pointer.

> **Tier**: the initial proposal tells you the change's size — propose its tier
> (`small | standard`; criteria — conventions, Change tiers). The user confirms it at the gate
> (step 4); the confirmed value becomes the `**Tier:**` line in `README.md`. For **small**: keep
> the data flow as text in `concept.md` — no `data-flow.drawio`.

## Step 3 — grill: evolve the proposal

Now **interview the user to refine the proposal until you reach shared understanding** — one
question at a time, each one reshaping the step-2 draft. Principles (in the spirit of grill-me):

- Walk the decision tree, resolving dependencies between decisions **one at a time**.
- **Ask questions one at a time.** For each question give **your recommended answer**.
- If a question can be answered by **reading the code — read the code** (symbol tools/grep), don't ask.
- For discrete choices use your interactive question tool (e.g. `AskUserQuestion` in Claude Code),
  with the recommended option first, marked as such; for open-ended formulations — plain chat.

Cover at minimum: **problem and why** · **scope and non-goals** · **hard constraints/invariants**
(find them in code, tests, docs, history) · **success criteria** · **affected code/files**.

As answers land, **update the draft seed and `concept.md`** — strategy and data flow may shift,
that's the point of the loop. If an answer overturns the chosen strategy, revise step 2's proposal
rather than patching around it. Stop when the proposal is stable and no ambiguities remain.

## Step 4 — GATE: confirm and write the seed

**GATE** — present the evolved proposal as a whole (strategy + data flow + seed + tier) and wait
for confirmation. The user may still redirect the strategy → loop back to step 2/3; don't patch on
the fly. On approval:

1. Write `proposal.md` — **short, < 1 screen**, openspec-style:
   `## Why` · `## What` · `## Constraints & Invariants` · `## Decisions` · `## Non-goals`.
   In bullets. `## Why` = 2–4 bullets: the concrete pain or trigger, who hits it, and the cost
   of not doing it — no mission statements. `## Decisions` reflects the **chosen strategy** +
   invariants; the depth (best-practices, candidates, rejected) lives in `concept.md`, the data
   chain in the data flow — not here.
2. For **standard** tier: hand-craft `data-flow.drawio` — nodes labelled with project types /
   pseudocode, edges = transformations (distinct from `integration.drawio`, the call graph made at
   integrate). **After writing always** validate (see conventions, drawio). For **small**: skip it
   — the data flow stays text in `concept.md`.
3. Write `README.md` — per the template in conventions: Phase line (bold **propose**), Tier
   line (`small | standard`, agreed at this gate), empty status table (no parts yet), links to
   `concept.md` + `data-flow.drawio` (if made).
4. Confirm `.prism/CURRENT` points at `<change>` (set in step 1) so subsequent prism commands target it.

## Next

Tell the user: proposal + concept ready → {{cmd:decompose}} to split into parts.
End your turn here — do not start decomposing.
