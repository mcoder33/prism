Splits the problem into first-level parts — or drills deeper into one node.

**Read first**: `.prism/conventions.md` (if missing, run `prism update`).

## Input

- No argument → split `proposal.md` into top-level parts.
- With argument `<NN-name>` → split that node into sub-nodes (recursive).

## Procedure

1. Read `proposal.md` (or `node.md` of the target node) and, if needed, the real code.
2. Propose a split into **several parts** — a list of "title + 1-line summary" right in chat.
   Each part should be graspable in one sitting.
3. **Show the split and wait for confirmation** (can merge/split differently/reorder).
   This is a decision-point — propose your split, don't dump a menu.
4. After "ok", create `NN-name/` directories with a small `node.md` in each (template — in conventions).
   When drilling a node — `NN-name/NNa-…/node.md`.
5. Update the tree and status table in `README.md` (new parts → ⚪).

## Next

Suggest which part to drill first (usually — the core/foundation): {{cmd:drill}} `<NN-name>`.
