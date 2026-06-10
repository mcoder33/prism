Implements the change in code based on the designed artifacts.

**Read first**: `.prism/conventions.md` (if missing, run `prism update`).

## Preconditions (check, don't assume)

- Root `tasks.md` + `integration.drawio` exist — else run {{cmd:integrate}} first. Stop.
- Every node in `README.md` is 🟢 or 🔵 (or the user explicitly deferred it) — else
  {{cmd:drill}} the gaps. Stop.
- No `[blocking]` Open left in any `node.md` — else resolve them with the user now.

## Procedure

1. Identify `<change>`: explicit arg → else `.prism/CURRENT` (the active change) → else ask. Never guess.
2. Read the **root `tasks.md` only** (order by dependencies). Read each part's artifacts
   (`spec.md`, `detail.md`, `signatures.md`, `tasks.md`) **just-in-time** at the start of that
   part — not all up front (early reading goes stale and bloats context).
3. If on the default branch — create a feature branch. Update `README.md` Phase → **apply**.
4. Implement **in dependency order**, part by part:
   - write code per `signatures.md`/`detail.md`, using `spec.md` as acceptance criteria;
   - write tests per scenarios from `spec.md` — name each test after its scenario (or record the
     scenario→test mapping when checking off the task), so {{cmd:verify}} can trace coverage;
   - mark completed items in the corresponding `tasks.md`: `- [ ]` → `- [x]`.
5. After each logical part, run the project checks (see the project docs — CLAUDE.md/AGENTS.md,
   Makefile, CI config), then make a **separate commit for that part** — one commit per part
   (01, 02, …), with a subject naming the part, e.g. `feat: 02 batch repos + split emit`.
   Right after the commit, flip that node 🟢 → 🔵 in `README.md` (so a killed session resumes
   correctly).
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
- **Design-as-built**: if the implementation deviates from `signatures.md`/`detail.md`/`spec.md`,
  update that artifact **before the part's commit**. Archived artifacts are the reference examples
  for future changes — they must describe what was actually built.
