# PRISM

**Recursive decomposition workflow for AI coding agents.**

PRISM turns a vague task into an implemented, verified change through a sequence of small,
confirmable design steps: propose → decompose → drill → integrate → apply → verify → archive.
Instead of one large design document, the problem is split recursively into small `node.md`
files — one digestible node at a time, with explicit decision gates.

Like [OpenSpec](https://github.com/Fission-AI/OpenSpec), PRISM is installed per-project as a set
of slash commands for the AI tools you use. All design artifacts live in `.prism/` at the repo
root (git-excluded — local working state, never committed).

## Install

Single static Go binary, templates embedded:

```bash
go install gitlab.gidfinance.tech/zadolbator/prism@latest
# or from a local checkout
go install .
```

Then, in any project:

```bash
prism init            # interactive: pick the AI tools to install commands for
prism init --tools claude,cursor   # non-interactive
prism init --tools all
```

`prism init`:

- detects which AI tools are present in the project (by their dot-dirs) and pre-selects them;
- writes the shared conventions to `.prism/conventions.md`;
- writes the slash-command files for each selected tool;
- adds `.prism/` to `.git/info/exclude`.

## Supported tools

| Tool | Commands location | Invocation |
|---|---|---|
| Claude Code | `.claude/commands/prism/*.md` | `/prism:propose` |
| Cursor | `.cursor/commands/prism-*.md` | `/prism-propose` |
| Codex CLI | `.codex/prompts/prism-*.md` | `/prism-propose` |
| Gemini CLI | `.gemini/commands/prism/*.toml` | `/prism:propose` |
| GitHub Copilot | `.github/prompts/prism-*.prompt.md` | `/prism-propose` |
| Windsurf | `.windsurf/workflows/prism-*.md` | `/prism-propose` |
| OpenCode | `.opencode/command/prism-*.md` | `/prism-propose` |

Adding a tool = one `adapters.Tool` value in `internal/adapters/adapters.go`
(file path + slash-command naming + frontmatter format).

## Workflow commands

| Command | Purpose |
|---|---|
| `use` | Select the active change (like `git checkout`); persisted in `.prism/CURRENT` |
| `propose` | Grill on requirements, survey best practices, pick a strategy + data-flow, write the seed |
| `decompose` | Split the proposal (or a node) into a few small `node.md` parts; recursive |
| `drill` | Bring ONE part to atomic; generate spec/detail/diagram/signatures/tasks |
| `integrate` | Cross-part artifacts: integration diagram + combined signatures + root tasks |
| `apply` | Implement part-by-part in dependency order, one commit per part, run checks |
| `verify` | Pedantic post-implementation QA on a running dev env (tests, concurrency, smoke, load) |
| `archive` | Move a finished change to `.prism/archive/` |

The methodology itself (artifact formats, layout, rules) is in `templates/conventions.md`,
installed into each project as `.prism/conventions.md` — every command reads it first.

## CLI

```bash
prism init [path] [--tools <list|all|none>]   # install/extend commands in a project
prism update [path] [--force]                 # regenerate installed command files after a CLI upgrade
prism list [path]                             # list active changes in .prism/
```

Generated files are **tool-owned**: each carries a `prism:generated v<version>` stamp and is
overwritten wholesale by `prism update` — don't hand-edit them. To change the workflow texts,
edit `templates/` here, bump `workflows.Version`, reinstall the binary and run `prism update`
in consuming projects.

## Development

```bash
go test ./...
go vet ./...
go install .                                  # build + put `prism` on PATH
prism init --tools claude /path/to/project    # smoke
```

Stack: Go (stdlib + cobra for the CLI, charmbracelet/huh for the interactive multi-select),
templates shipped via `go:embed` — the binary is self-contained.

Layout:

```
main.go                      entrypoint → internal/cli
templates/
  embed.go                   go:embed of the markdown below
  conventions.md             shared methodology, installed as .prism/conventions.md
  commands/*.md              tool-neutral command bodies ({{cmd:<id>}} placeholders for cross-refs)
internal/
  workflows/workflows.go     workflow registry (id, title, description) + Version
  adapters/adapters.go       per-tool adapters: paths, slash-command naming, frontmatter
  installer/installer.go     detection, rendering, writing, git-exclude
  cli/                       cobra commands: init / update / list
```
