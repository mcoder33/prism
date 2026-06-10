Starts a new "change" in the recursive decomposition methodology.

**Read first**: `.prism/conventions.md` (formats, layout, rules; if missing, run `prism update`).

## Amendment mode (existing change)

If `proposal.md` already exists for the active change, this is an **amendment**, not a new
change: present the DELTA to Why/What/Decisions (**GATE**), update `proposal.md` + `concept.md`,
then route to {{cmd:decompose}} for the affected parts only (revision rules — in conventions).
Skip steps 2–4 unless the strategy itself is in question.

## Step 1 — grill the user on requirements

Before writing anything, **interview the user on the plan/task until you reach shared understanding**.
Principles (in the spirit of grill-me):

- Walk the decision tree, resolving dependencies between decisions **one at a time**.
- **Ask questions one at a time.** For each question give **your recommended answer**.
- If a question can be answered by **reading the code — read the code** (symbol tools/grep), don't ask.
- For discrete choices use your interactive question tool (e.g. `AskUserQuestion` in Claude Code),
  with the recommended option first, marked as such; for open-ended formulations — plain chat.

Cover at minimum: **problem and why** · **scope and non-goals** · **hard constraints/invariants**
(find them in code, tests, docs, history) · **success criteria** · **affected code/files**.

Stop when the *problem* is clear and no ambiguities remain. (The *solution* — strategy and data
flow — is fixed in steps 2–4 below, not here.)

> **Tier**: by the end of step 1 you know the change's size — propose its tier
> (`small | standard`; criteria — conventions, Change tiers) and have the user confirm it at
> the next gate; the confirmed value becomes the `**Tier:**` line in `README.md` (step 5).
> For **small**: merge steps 3 and 4 into ONE gate (strategy + data flow together; data flow
> may be text/pseudocode in `concept.md` — no `data-flow.drawio` needed).

## Step 2 — best practices (skippable)

> Setup: pick the kebab-case `<change>` slug and create `.prism/<change>/` now (create `.prism/`
> if missing; add it to `.git/info/exclude`; write the slug to `.prism/CURRENT`) so steps 2–4 can
> persist `concept.md` / `data-flow.drawio` in place. `proposal.md` / `README.md` come in step 5.

Give the user **breadth**: how is this *class* of problem solved in general?

- First ask (interactive question): do the survey, or skip (recommended option = the one that fits
  — a routine/well-understood task can skip; a wide solution space benefits from it).
- If doing it: a short, decision-first survey — **3–6 bullets**, each `practice — when it applies
  [source]`. You may use web search / docs lookup for current practices, plus your own knowledge.
  Don't dump essays; give the lay of the land.
- Persist into `concept.md` under `## Best practices` (mark `> Skipped — user opted out.` if skipped).

## Step 3 — candidate strategies

Pick the **high-level approach** before any decomposition (the solution space is often wide).

- Present **a few** strategies — **decision-first**, not a menu to weigh: recommend one, the
  others in **one line each**. Ground them in real code (symbol tools/grep), not in a vacuum.
- Invite the user to describe **their own** strategy too.
- **GATE** — this is a decision-point: the user picks (or supplies) the strategy. Don't proceed
  until fixed.
- Record in `concept.md` under `## Candidate strategies` + `## Chosen strategy` (chosen + rejected
  one line each).

## Step 4 — data-mutation schema

For the **chosen** strategy, sketch the chain of **how data changes** end-to-end.

- `data-flow.drawio` — nodes labelled with project types / pseudocode, edges = transformations.
  This is the conceptual data flow, distinct from `integration.drawio` (call graph, made at integrate).
- Hand-craft mxGraph; **after writing always** validate (see conventions, drawio — xmllint or fallback).
- **GATE** — conceptual gate: confirm the user likes the idea. If not → back to step 3 (don't
  patch on the fly).

## Step 5 — write the seed

When the strategy and data flow are confirmed (the change dir was created in the step 2 setup;
`concept.md` + `data-flow.drawio` were filled in steps 2–4):

1. Write `proposal.md` — **short, < 1 screen**, openspec-style:
   `## Why` · `## What` · `## Constraints & Invariants` · `## Decisions` · `## Non-goals`.
   In bullets. `## Why` = 2–4 bullets: the concrete pain or trigger, who hits it, and the cost
   of not doing it — no mission statements. `## Decisions` reflects the **chosen strategy** +
   invariants; the depth (best-practices, candidates, rejected) lives in `concept.md`, the data
   chain in `data-flow.drawio` — not here.
2. Write `README.md` — per the template in conventions: Phase line (bold **propose**), Tier
   line (`small | standard`, agreed at the gate), empty status table (no parts yet), links to
   `concept.md` + `data-flow.drawio` (if made).
3. Confirm `.prism/CURRENT` points at `<change>` (set in step 2) so subsequent prism commands target it.

## Next

Tell the user: proposal + concept ready → {{cmd:decompose}} to split into parts.
End your turn here — do not start decomposing.
