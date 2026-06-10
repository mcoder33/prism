Read-only: reports where the active change stands and the single next action. Changes nothing —
the re-entry point after a session break.

**Read first**: `.prism/conventions.md` (if missing, run `prism update`).

## Procedure

1. **Resolve `<change>`**: explicit arg → else `.prism/CURRENT` → else list active changes
   (sub-dirs of `.prism/` excluding `archive/`) and pick via an interactive question. If `.prism/`
   has no active changes at all, say so and suggest {{cmd:propose}}. Never guess.
2. **Read `README.md` of the change, then cross-check it against reality** — don't silently
   trust the table:
   - `NN-name/` directories present vs the table rows;
   - expected artifacts per node for its tier (node.md/tasks.md always; spec/detail/drawio/signatures
     per the agreed set);
   - unchecked boxes (`- [ ]`) in the root and per-part `tasks.md`;
   - `git log --oneline` on the current branch — a part marked 🔵 should have its commit.
   Report every discrepancy (e.g. table says 🔵 but no commit; 🟢 but `detail.md` missing).
3. **Print inline**: the Phase line, the node table (with statuses and Open counts), unchecked
   task counts, and any `[blocking]` Open items or discrepancies as blockers.
4. **Recommend exactly ONE next command with its argument** (e.g. {{cmd:drill}} `02-renderer`,
   or {{cmd:integrate}}). Do not run it — end your turn.

## Rules

- Strictly read-only: no file writes, no status flips, no fixes — only report.
- If `README.md` and reality disagree, reality wins; recommend the command that reconciles them.
