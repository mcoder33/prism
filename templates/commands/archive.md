Archives a completed decomposition change. {{cmd:apply}} runs this automatically as its last
step; invoke it manually to archive a change that was applied without auto-archive, or to tidy up.

**Read first**: `.prism/conventions.md` (if missing, run `prism update`).

**Input**: optionally a `<change>` name. If omitted, infer from context; if ambiguous you MUST
prompt for selection (never auto-guess).

## Abandon mode (`--abandon`)

`--abandon` (or the user clearly saying "drop / kill / give up on this change") archives a change
that will **not** be finished — the missing exit for a dead-end idea, so it stops cluttering
`prism:list` and the `prism:use` picker without a manual `rm`. It differs from a normal archive:

- **Skip the task-completion check** (step 2) — an abandoned change is expected to be incomplete.
- Record the reason: add `**Status:** abandoned — <one-line reason>, <date>` under the Phase line
  in the change's `README.md` **before** moving it (so the archive states why it died).
- Then run steps 3 and 5 (move + collision `-rN` + clear `.prism/CURRENT` if it points here).
  **Skip the mini-retro** (step 4) — nothing shipped.
- Confirm once before moving (it is not resumable in place afterwards); never abandon without an
  explicit user go-ahead.

## Procedure

1. **Resolve the change.**
   Explicit arg → else `.prism/CURRENT` (the active change) → else list active changes (sub-dirs of
   `.prism/` **excluding `archive/`**) and pick via an interactive question (e.g. `AskUserQuestion`
   in Claude Code). Do not guess.

2. **Check task completion.**
   Scan the change's root `tasks.md` and every part `tasks.md` for unchecked items (`- [ ]` vs
   `- [x]`). If any remain incomplete:
   - show the count and which file(s) they're in;
   - ask the user (interactive question) to confirm they still want to archive;
   - proceed only on confirmation.
   If there's no `tasks.md`, proceed without a task warning. Also list any ⏸ deferred nodes
   as warnings (they were intentionally skipped — see conventions, Open tags).

3. **Perform the archive.**
   - Create the archive dir if missing: `mkdir -p .prism/archive`.
   - **Collision check**: if `.prism/archive/<change>/` already exists, propose archiving as
     `<change>-r2` (next free `-rN` suffix) and confirm via an interactive question. Never
     overwrite, never silently rename the old one.
   - Move the whole folder: `mv .prism/<change> .prism/archive/<change>`.
   - If `.prism/CURRENT` points to this change, clear it (no active change after archiving).

4. **Mini-retro (optional, keep tiny).** If during apply/verify the design had to change
   (deviations recorded under the design-as-built rule), classify each lesson first — one
   lesson per line, with its why — and route it by kind:
   - **project** (about this codebase: a convention, a trap, a constraint) → propose it as one
     of **≤ 3 one-line additions** to the committed agent docs (CLAUDE.md/AGENTS.md) — the
     user decides;
   - **methodology** (a prism gate, artifact, or command misfired or proved unnecessary) →
     report it to the user as candidate feedback for prism itself (an issue/PR upstream), it
     earns no CLAUDE.md line;
   - **one-off** (won't recur) → drop it.
   Don't write a retro file into `.prism/` (it's git-excluded and would die locally).

5. **Display summary** — change name, archive location, and any warnings (incomplete tasks).

## Rules

- Always prompt for selection if the change is not provided/inferable; never auto-guess.
- Don't block on warnings — inform and confirm, then proceed.
- `.prism/` artifacts are not committed by default; archiving is a filesystem move only
  (no commit — unless the project opted into sharing archives, see conventions).
- Don't overwrite an existing archive entry.
- **Hotfix to an archived change**: don't move it back. Start a new change ({{cmd:propose}})
  whose `proposal.md` links the archived folder under `## Why` — the archive stays an
  immutable as-built record.
