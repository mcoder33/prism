Selects the **active change** — the prism analog of the current git branch. The choice is
persisted in `.prism/CURRENT` (like `.git/HEAD`), so it survives across sessions and **every other
prism command and all design writing target it by default**, without re-specifying the name.

**Read first**: `.prism/conventions.md` (if missing, run `prism update`).

## Input / actions

- {{cmd:use}} (no arg) — **interactive picker** (the default, primary mode).
- {{cmd:use}} `<change>` — **start / switch** directly to a known slug (no prompt).
- {{cmd:use}} `stop` — **stop**: clear the active change (like a detached HEAD).

## Procedure

1. **No arg → interactive picker.** Build the option list with your interactive question tool
   (e.g. `AskUserQuestion` in Claude Code):
   - one option per active change (sub-dirs of `.prism/` except `archive/`), marking the current
     one (from `.prism/CURRENT`) as `(current)`;
   - **`➕ New change…`** — routes straight to {{cmd:propose}} to create one (see step 4);
   - **`⏹ Stop (no active change)`** — clears the pointer (step 5).

   (If active changes outnumber the picker's option slots, show the most recently modified and note
   that a name can be passed directly.) Whatever the user picks → run the matching step below.

2. **start / switch** (`<change>` given, or picked from the list):
   - Verify `.prism/<change>/` exists; if not, fall back to the interactive picker (step 1).
   - Write the slug to `.prism/CURRENT` (create `.prism/` if missing), replacing any previous value.
   - Confirm: `active change → <change>`.

3. **show**: when the user just wants the current value, print `.prism/CURRENT` (or "none").

4. **➕ New change** (picked): invoke {{cmd:propose}} — it creates `.prism/<new>/` and sets it
   active (writes `.prism/CURRENT`) automatically, so no extra switch is needed.

5. **stop** (`stop`, or picked): remove `.prism/CURRENT`. Confirm: `no active change`.

## Rules

- `.prism/CURRENT` is the single source of truth for "which change am I working on". When a
  prism command gets no explicit `<change>`, it reads `.prism/CURRENT`; if that is also
  empty, it asks (never guesses).
- {{cmd:propose}} sets the new change active automatically; {{cmd:archive}} clears it if it
  archived the active change.
- `.prism/` (including `CURRENT`) is not committed — it's local working state.
