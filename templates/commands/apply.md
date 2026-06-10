Implements the change in code based on the designed artifacts.

**Read first**: `.prism/conventions.md` (if missing, run `prism update`).

## Preconditions (check, don't assume)

- Root `tasks.md` + `integration.drawio` exist — **or** `README.md` records `**Tier:** small`
  with its single node 🟢 (then there are no root artifacts; the part's own `tasks.md` is the
  order) — else run {{cmd:integrate}} first. Stop.
- Every node in `README.md` is 🟢 or 🔵 (⏸ deferred nodes are skipped — see conventions,
  Open tags) — else {{cmd:drill}} the gaps. Stop.
- No `[blocking]` Open left in any `node.md` (`[deferred: …]` and `[minor]` pass) — else
  resolve them with the user now.

## Procedure

1. Identify `<change>`: explicit arg → else `.prism/CURRENT` (the active change) → else ask. Never guess.
2. Read the **root `tasks.md` only** (order by dependencies; for `Tier: small` there is none —
   go straight to the single part). Read each part's artifacts (`spec.md`, `detail.md`,
   `signatures.md`, `tasks.md`) **just-in-time** at the start of that part — not all up front
   (early reading goes stale and bloats context).
3. Branch check, **in this order**: if `README.md` already records `**Branch:**` and the
   current git branch differs — warn and stop (resume on the recorded branch; never switch or
   create branches yourself). Else, if on the default branch — create a feature branch and
   record it in `README.md` under the Tier line as `**Branch:** <name>`. If the branch is old
   and the default branch has moved — merge/rebase it in first, resolving conflicts per the
   design artifacts, not ad hoc. Then update `README.md` Phase → **apply**.
4. Implement **in dependency order**, part by part:
   - write code per `signatures.md`/`detail.md`, using `spec.md` as acceptance criteria;
   - write tests per scenarios from `spec.md` — name each test after its scenario, and record
     each scenario in the part's `spec.md` `## Coverage` section as you check off the task
     (format — conventions, spec.md), so {{cmd:verify}} can trace coverage;
   - mark completed items in the corresponding `tasks.md`: `- [ ]` → `- [x]`.
5. After each logical part, run the project checks (see the project docs — CLAUDE.md/AGENTS.md,
   Makefile, CI config), then make a **separate commit for that part** — one commit per part
   (01, 02, …), with a subject naming the part, e.g. `feat: 02 batch repos + split emit`.
   Right after the commit, flip that node 🟢 → 🔵 in `README.md` (so a killed session resumes
   correctly).
6. **Pause on blocker**: unclear task, design gap, error. On a design gap — return to
   {{cmd:drill}}/{{cmd:decompose}}, don't improvise in code.
7. At the end — cross-cutting items from the root `tasks.md` (cleanup, integration/smoke, green
   project checks) — as their own final commit. Skip items marked `(blocked by NN ⏸)` — leave
   them unchecked; verify reports them as `affected but untested`. (`Tier: small` has no root
   `tasks.md` — skip the step.)
8. **Verify the result**: run {{cmd:verify}} `--skip-static` — thorough QA on the running dev
   environment (functional smoke, concurrency/load checks); static checks and the full test run
   already happened in steps 5/7. Fix blocking findings and re-verify.
9. **Archive the change automatically**: once every task passes, the project checks are green and
   verification is PASS, run {{cmd:archive}} `<change>` to move `.prism/<change>/` →
   `.prism/archive/<change>/`. This keeps `.prism/` showing only active changes.

## Integration with code review

PRISM ends at archive; pushing and review follow your project's practice. Two common shapes:

- **Verify first, then PR**: steps 1–9 above, then push the branch and open the PR/MR.
- **PR first, verify in CI**: push and open the PR after the part commits (step 5); verify runs
  as a CI check; archive after the merge.

Either way, archiving is a local filesystem move — it is independent of the merge.

## Rules

- **One commit per part** (01, 02, …); cross-cutting work is its own final commit. Don't squash
  several parts into one commit. `.prism/` artifacts are not committed.
- **Re-apply after a verify escalation** (the node is 🟢 again but its old commit is in
  history): implement the fix as a **new commit** for the same part (`fix: NN <part>` subject) —
  never amend or rebase the old commit (it may be pushed). Flip 🟢 → 🔵 after the new commit
  as usual.
- Don't mark a task done until its scenarios from `spec.md` pass.
- **Design-as-built**: if the implementation deviates from `signatures.md`/`detail.md`/`spec.md`,
  update that artifact **before the part's commit**. Archived artifacts are the reference examples
  for future changes — they must describe what was actually built.
