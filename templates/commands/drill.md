Brings ONE part to "atomic" and generates its artifacts. Argument: `<NN-name>`.

**Read first**: `.prism/conventions.md` (if missing, run `prism update`).

## Preconditions (check, don't assume)

- The node directory `NN-name/` with its `node.md` exists — else run {{cmd:decompose}} first. Stop.

## Procedure

1. Read `NN-name/node.md` and the relevant **real code** (symbol-overview / find-symbol tools
   if available, otherwise grep) — signatures and details must be grounded in facts, not guesses.
   Mark the node 🟡 in `README.md`.
2. **Atomicity check** (criteria — in conventions). Not atomic → propose the sub-split and
   redirect to {{cmd:decompose}} `<NN-name>`. Stop.
3. Write/refresh **`node.md` only**. Present inline in chat:
   - the full `node.md` text;
   - key decisions: chosen + rejected, one line each;
   - which artifacts you will produce and which you'll skip, with a one-line reason each
     (tiers — in conventions);
   - open questions, each tagged `[blocking]` or `[minor]`.

   **GATE** — the user reacts to the digest. `[blocking]` opens must be resolved here.
4. On approval: generate the remaining artifacts of the agreed tier (templates — in conventions):
   - `spec.md` — Requirement/Scenario (these also drive the tests later).
   - `detail.md` — how to implement, decision-complete (algorithm, subtleties, edge-cases, worked
     example). **Optional**: if the node transforms data, sketch its **local** data-mutation chain
     here (text/pseudocode in the worked example) — the same idea as the change-level
     `data-flow.drawio`, at node granularity. No separate file needed.
   - `concept.drawio` — diagram; **mandatory** validation (see conventions, drawio).
   - `signatures.md` — code sketch: signatures + what/why comments (no implementation).
   - `tasks.md` — checklist `- [ ]`.
5. Mark the node 🟢 in `README.md`. Print a one-line summary per artifact + the full `tasks.md`
   inline, report the next candidates, and **end your turn** — do NOT drill another part
   (low cognitive load: one node at a time).

## Rules

- Decision-first: one decision, rationale, rejected alternatives — one line each.
- Follow formatting (bullets/blank lines) from conventions.
- If drilling reveals a flawed model — return to {{cmd:decompose}}, don't patch on the fly.
