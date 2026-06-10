Assembles the parts into the overall picture. Run when the main parts are drilled (🟢).

**Read first**: `.prism/conventions.md` (if missing, run `prism update`).

## Preconditions (check, don't assume)

- `README.md` records `**Tier:** small` → integrate is skipped entirely (no root artifacts);
  say so and point to {{cmd:apply}}. Stop.
- Every node in the `README.md` table is 🟢 (or ⏸ deferred — see conventions, Open tags) — else
  {{cmd:drill}} the gaps first. Stop.

## Procedure

1. Read `proposal.md` and `node.md`/`signatures.md` of all parts.
2. Create in the root `.prism/<change>/`:
   - **`integration.drawio`** — overall live-path diagram; nodes annotated with `[NN]` of which parts they belong to;
     gray — out of scope / unchanged. Validate (see conventions, drawio).
   - **`signatures.md`** (combined) — assembled call flow: one pseudocode block "who calls whom",
     a link table, and a type map of what flows between parts.
   - **`tasks.md`** (root) — **only order** (by dependencies) and **cross-cutting**.
     Cross-cutting = work owned by no single part: dependency ordering, migrations spanning
     parts, shared cleanup, integration/e2e tests touching ≥ 2 parts, project-wide checks.
     A task naming exactly one part belongs in that part's `tasks.md` — move it. NO repetition
     of part details. A cross-cutting task that depends on a ⏸ node gets a `(blocked by NN ⏸)`
     suffix — apply skips it, verify reports it as `affected but untested`.
3. Update `README.md` (Phase → **integrate**, links to the combined artifacts).
4. Present inline in chat: the part order with one line each, the root task list, and **every
   remaining Open item collected from all `node.md`** — `[blocking]` ones must be resolved here,
   or explicitly deferred by the user: record `[deferred: <reason> — user, <date>]` in the
   owning `node.md` (a whole node → status ⏸ in `README.md`; tags — in conventions).

   **GATE** — only after the user's reply declare the design package ready. If the user
   deferred a node here, update the already-written root artifacts before declaring ready:
   `(blocked by NN ⏸)` suffixes in root `tasks.md`, gray the node out in `integration.drawio`.

## Next

Design package ready → {{cmd:apply}} for implementation.
End your turn here — do not start applying.
