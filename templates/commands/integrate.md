Assembles the parts into the overall picture. Run when the main parts are drilled (🟢).

**Read first**: `.prism/conventions.md` (if missing, run `prism update`).

## Preconditions (check, don't assume)

- Every node in the `README.md` table is 🟢 (or the user explicitly deferred it) — else
  {{cmd:drill}} the gaps first. Stop.

## Procedure

1. Read `proposal.md` and `node.md`/`signatures.md` of all parts.
2. Create in the root `.prism/<change>/`:
   - **`integration.drawio`** — overall live-path diagram; nodes annotated with `[NN]` of which parts they belong to;
     gray — out of scope / unchanged. Validate (see conventions, drawio).
   - **`signatures.md`** (combined) — assembled call flow: one pseudocode block "who calls whom",
     a link table, and a type map of what flows between parts.
   - **`tasks.md`** (root) — **only order** (by dependencies) and **cross-cutting** (cleanup,
     integration test, project checks). NO repetition of part details — those live in `NN/tasks.md`.
3. Update `README.md` (Phase → **integrate**, links to the combined artifacts).
4. Present inline in chat: the part order with one line each, the root task list, and **every
   remaining Open item collected from all `node.md`** — `[blocking]` ones must be resolved or
   explicitly deferred by the user here.

   **GATE** — only after the user's reply declare the design package ready.

## Next

Design package ready → {{cmd:apply}} for implementation.
End your turn here — do not start applying.
