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
     per the agreed set) — ⏸ nodes are exempt: they only need `node.md` with a `[deferred: …]` line;
   - unchecked boxes (`- [ ]`) in the root and per-part `tasks.md` (none for `Tier: small` root);
   - `git log --oneline` on the current branch — a part marked 🔵 should have its commit;
   - if `README.md` records `**Branch:**`, compare with the current git branch — a mismatch is
     a warning (never switch branches).
   Report every discrepancy (e.g. table says 🔵 but no commit; 🟢 but `detail.md` missing).
3. **Print inline**: the Phase line, the node table (with statuses and Open counts), unchecked
   task counts, and any `[blocking]` Open items or discrepancies as blockers (`[deferred: …]`
   and `[minor]` items are not blockers — list them as notes).
4. **Recommend exactly ONE next command with its argument** (e.g. {{cmd:drill}} `02-renderer`,
   or {{cmd:integrate}}). Do not run it — end your turn.

## Rules

- Read-only by default: no fixes, no status flips during steps 1–4 — only report.
- If `README.md` and reality disagree, reality wins. When step 2 finds discrepancies, print a
  **Repair plan**: per discrepancy, the exact `README.md`/`.prism/CURRENT` edit that makes the
  table match reality (commit exists → 🔵; tier artifacts missing → demote to 🟡; node dir
  missing → drop the row; stale `CURRENT` → clear it). Exception: a 🟡 node whose `node.md`
  carries `[blocking] verify escalation: …` or `[minor] reverted <short-sha>` stays 🟡 even
  though a commit exists — never promote it back to 🔵. **GATE** — only on the user's explicit "repair" reply, perform those
  `README.md`/`CURRENT` edits and nothing else. Never touch code, design artifacts, or git.
