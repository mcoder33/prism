Implements the change in code based on the designed artifacts.

**Read first**: `.prism/conventions.md` (if missing, run `prism update`).

## Procedure

1. Identify `<change>`: explicit arg → else `.prism/CURRENT` (the active change) → else ask. Never guess.
2. Read the root `tasks.md` (**order** by dependencies) and for each part — `spec.md`,
   `detail.md`, `signatures.md`, `tasks.md`.
3. If on the default branch — create a feature branch.
4. Implement **in dependency order**, part by part:
   - write code per `signatures.md`/`detail.md`, using `spec.md` as acceptance criteria;
   - write tests per scenarios from `spec.md`;
   - mark completed items in the corresponding `tasks.md`: `- [ ]` → `- [x]`.
5. After each logical part, run the project checks (see the project docs — CLAUDE.md/AGENTS.md,
   Makefile, CI config), then make a **separate commit for that part** — one commit per part
   (01, 02, …), with a subject naming the part, e.g. `feat: 02 batch repos + split emit`.
6. **Pause on blocker**: unclear task, design gap, error. On a design gap — return to
   {{cmd:drill}}/{{cmd:decompose}}, don't improvise in code.
7. At the end — cross-cutting items from the root `tasks.md` (cleanup, integration/smoke, green
   project checks) — as their own final commit.
8. **Verify the result**: run {{cmd:verify}} `--skip-static` — thorough QA on the running dev
   environment (functional smoke, concurrency/load checks); static checks and the full test run
   already happened in steps 5/7. Fix blocking findings and re-verify.
9. **Archive the change automatically**: once every task passes, the project checks are green and
   verification is PASS, run {{cmd:archive}} `<change>` to move `.prism/<change>/` →
   `.prism/archive/<change>/`. This keeps `.prism/` showing only active changes.

## Rules

- **One commit per part** (01, 02, …); cross-cutting work is its own final commit. Don't squash
  several parts into one commit. `.prism/` artifacts are not committed.
- Don't mark a task done until its scenarios from `spec.md` pass.
