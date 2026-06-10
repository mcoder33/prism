Archives a completed decomposition change. {{cmd:apply}} runs this automatically as its last
step; invoke it manually to archive a change that was applied without auto-archive, or to tidy up.

**Read first**: `.prism/conventions.md` (if missing, run `prism update`).

**Input**: optionally a `<change>` name. If omitted, infer from context; if ambiguous you MUST
prompt for selection (never auto-guess).

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
   If there's no `tasks.md`, proceed without a task warning.

3. **Perform the archive.**
   - Create the archive dir if missing: `mkdir -p .prism/archive`.
   - **Collision check**: if `.prism/archive/<change>/` already exists, stop and report it
     (suggest renaming the existing archive); do not overwrite.
   - Move the whole folder: `mv .prism/<change> .prism/archive/<change>`.
   - If `.prism/CURRENT` points to this change, clear it (no active change after archiving).

4. **Display summary** — change name, archive location, and any warnings (incomplete tasks).

## Rules

- Always prompt for selection if the change is not provided/inferable; never auto-guess.
- Don't block on warnings — inform and confirm, then proceed.
- `.prism/` artifacts are not committed; archiving is a filesystem move only (no commit).
- Don't overwrite an existing archive entry.
