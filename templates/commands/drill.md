Brings ONE part to "atomic" and generates its artifacts. Argument: `<NN-name>`.

**Read first**: `.prism/conventions.md` (if missing, run `prism update`).

## Procedure

1. Read `NN-name/node.md` and the relevant **real code** (symbol-overview / find-symbol tools
   if available, otherwise grep) — signatures and details must be grounded in facts, not guesses.
2. **Decide: atomic or drill further?**
   - Too large / too many branches → propose sub-decomposition and redirect to {{cmd:decompose}} `<NN-name>`. Stop.
   - Manageable scope → continue.
3. Generate the node artifacts (templates — in conventions):
   - `spec.md` — Requirement/Scenario (these also drive the tests later).
   - `detail.md` — how to implement, decision-complete (algorithm, subtleties, edge-cases, worked example).
     **Optional**: if the node transforms data, sketch its **local** data-mutation chain here
     (text/pseudocode in the worked example) — the same idea as the change-level `data-flow.drawio`,
     at node granularity. No separate file needed.
   - `concept.drawio` — diagram; **mandatory** `xmllint --noout`.
   - `signatures.md` — code sketch: signatures + what/why comments (no implementation).
   - `tasks.md` — checklist `- [ ]`.
4. Mark the node 🟢 in `README.md`.
5. **Stop after one part** (low cognitive load). Report the next candidates.

## Rules

- Decision-first: one decision, rationale, rejected alternatives — one line each.
- Follow formatting (bullets/blank lines) from conventions.
- If drilling reveals a flawed model — return to {{cmd:decompose}}, don't patch on the fly.
