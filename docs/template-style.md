# Writing prism templates

The templates under [`templates/`](../templates/) ARE the product — the Go code just installs
them. This guide sets the bar for editing them. It condenses two sources worth reading in
full: Matt Pocock's [writing-great-skills](https://github.com/mattpocock/skills) and John
Ousterhout's *A Philosophy of Software Design* (deep modules).

## Two budgets, enforced by tests

Every template line is paid for twice:

- **Context load** — `conventions.md` is read by *every* command invocation; command bodies are
  read on theirs. `lint_test.go` enforces hard byte budgets (14KB / 20KB). Trim before adding;
  raising a budget is a deliberate decision, not a default.
- **Cognitive load** — the user sits at every GATE. Each extra artifact, question, or status
  they must react to costs attention. A step that doesn't sharpen a decision is sediment.

## Progressive disclosure

Command bodies stay short and procedural; shared definitions live once in `conventions.md` and
are referenced as `(see conventions, <Section>)` — the lint test checks those references
resolve. Duplicating a convention in a command body creates drift: one copy gets edited, the
other keeps being obeyed.

## Leading words

Prefer one pretrained concept over a paragraph of instructions: *tracer bullet*, *spike*,
*expand–contract*, *fog*, *decision-first*. A leading word recruits everything the model
already knows about the concept; a paraphrase recruits nothing. Use the canonical term, bold
it once at introduction, and don't invent private synonyms for it later.

## Positive phrasing

State what to do, not what to avoid — a negation activates the very behaviour it forbids
("don't think of an elephant"). When a prohibition is genuinely needed, pair it with the
replacement behaviour in the same sentence ("never amend the old commit — implement the fix
as a new `fix: NN` commit").

## Checkable completion criteria

Every step and artifact wants a criterion the agent can test itself against, not an adjective:
"decision-complete" is a mood; "a fresh implementer would ask zero questions" is a check.
Same for questions to the user — the fog test in conventions ("could someone answer it as
posed?") is the model.

## The no-op test

A line earns its place only if it changes behaviour versus what the model does by default.
That is model-relative and empirical: settle it by running the command with and without the
line, not by argument. Lines that restate defaults are sediment — they cost budget and dilute
the lines that do work.

## Rejected framings

Alternatives the methodology deliberately rejected — don't reintroduce them in edits:

- **Big upfront design docs** — overwhelm and get skimmed; approval becomes theater. The unit
  of review is a 5–7 line node.
- **Options menus** ("here are 3 approaches, you pick") — offloads the decision to the user
  wholesale. Decision-first: one proposal, rejected alternatives one line each.
- **Layer-first decomposition** ("01 = data layer, 02 = service, 03 = API") — nothing is
  verifiable until the last part lands. Vertical slices; 01 is the tracer bullet.
- **Fixed decomposition depth** — over-designs the obvious, under-designs the murky. Drill is
  per-node, until atomic.
- **Agent-judged approval** — "this looks good, proceeding" defeats the gate. A GATE ends the
  turn; only an explicit user reply resumes.

## Tool neutrality

Templates install into 9 different agents. Name capabilities, not products: "symbol-overview /
find-symbol tools if available, otherwise grep", "your interactive question tool (e.g.
`AskUserQuestion` in Claude Code)". A hard dependency on one agent's feature breaks the other
eight.

## Checklist before a PR

1. `make ci` green — the methodology lint (`internal/workflows/lint_test.go`) checks section
   references, status-glyph consistency, transition legality, `{{cmd:*}}` ids, and byte budgets.
2. Every new rule passed the no-op test at least informally: what would the agent have done
   without it?
3. New shared vocabulary landed in `conventions.md` once, referenced elsewhere.
4. Bump `Version` in `internal/workflows/workflows.go` and add a `CHANGELOG.md` entry —
   installed files are stamped with it and `prism doctor`/`update` key off the drift.
