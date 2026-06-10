Splits the problem into first-level parts — or drills deeper into one node.

**Read first**: `.prism/conventions.md` (if missing, run `prism update`).

## Preconditions (check, don't assume)

- `proposal.md` exists (or, when drilling deeper, the target `NN-name/node.md`) — else run
  {{cmd:propose}} first. Stop.

## Input

- No argument → split `proposal.md` into top-level parts.
- With argument `<NN-name>` → split that node into sub-nodes (recursive).

## Procedure

1. Read `proposal.md` (or `node.md` of the target node) and, if needed, the real code.
2. Propose a split into **2–4 parts** — a list of "title + 1-line summary" right in chat.
   Check the split before presenting it:
   - 2–4 parts (5+ usually means the proposal itself is too big — say so instead);
   - dependencies are acyclic; number in dependency order (01 = foundation);
   - each part independently testable and reviewable on its own;
   - comparable sizes; no "misc / the rest" bucket;
   - you can name the types/calls flowing between parts now — if you can't, the boundary is wrong.
3. **GATE** — show the split and wait for confirmation (the user can merge/split differently/
   reorder). Propose your split, don't dump a menu.
4. After "ok", create `NN-name/` directories with a small `node.md` in each (template — in
   conventions). When drilling a node — `NN-name/NNa-…/node.md`.
5. Update Phase, the tree and the status table in `README.md` (new parts → ⚪).

### Re-splitting an existing node

Confirm with the user first → delete the stale sub-dirs being replaced → keep sibling numbering
stable (don't renumber untouched parts) → update the `README.md` table.

## Next

Suggest which part to drill first (usually — the core/foundation): {{cmd:drill}} `<NN-name>`.
End your turn here — do not start drilling.
