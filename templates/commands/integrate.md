Assembles the parts into the overall picture. Run when the main parts are drilled (🟢).

**Read first**: `.prism/conventions.md` (if missing, run `prism update`).

## Procedure

1. Read `proposal.md` and `node.md`/`signatures.md` of all parts.
2. Create in the root `.prism/<change>/`:
   - **`integration.drawio`** — overall live-path diagram; nodes annotated with `[NN]` of which parts they belong to;
     gray — out of scope / unchanged. `xmllint --noout`.
   - **`signatures.md`** (combined) — assembled call flow: one pseudocode block "who calls whom",
     a link table, and a type map of what flows between parts.
   - **`tasks.md`** (root) — **only order** (by dependencies) and **cross-cutting** (cleanup,
     integration test, project checks). NO repetition of part details — those live in `NN/tasks.md`.
3. Update `README.md` (links to the combined artifacts).

## Next

Design package ready → {{cmd:apply}} for implementation.
