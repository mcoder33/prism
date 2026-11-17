// Package workflows is the registry of prism commands and their template bodies.
package workflows

import (
	"fmt"
	"strings"

	"github.com/mcoder33/prism/templates"
)

// Version of the prism CLI; stamped into every generated file.
// Overridable at build time: -ldflags "-X .../internal/workflows.Version=x.y.z".
var Version = "0.8.1"

// Workflow is one prism command installed into agents as a slash command.
type Workflow struct {
	// ID is the stable command id; also the file/slash-command suffix.
	ID string
	// Title is the human title used in frontmatter, e.g. "PRISM: Propose".
	Title string
	// Description is the one-liner shown by agents in command pickers.
	Description string
	// ModelInvocable: another workflow invokes this one as a sub-step (apply → verify →
	// archive) or it is read-only (status), so skill deliveries must not make it user-only.
	ModelInvocable bool
}

var All = []Workflow{
	{
		ID:    "use",
		Title: "PRISM: Use",
		Description: "Select the active change (like git checkout) via an interactive picker — " +
			"switch, stop, or \"+ New change\" (→ propose). All prism commands then default to it.",
	},
	{
		ID:    "propose",
		Title: "PRISM: Propose",
		Description: "Survey best practices in chat first, propose a strategy + data-flow grounded " +
			"in them, grill on requirements, then write the seed (proposal + concept) for a new " +
			"decomposition change.",
	},
	{
		ID:             "status",
		ModelInvocable: true,
		Title:          "PRISM: Status",
		Description: "Where the active change stands — phase, node table cross-checked against " +
			"reality, blockers, the single next command; offers a gated repair plan when the table " +
			"and reality disagree. The re-entry point after a break.",
	},
	{
		ID:          "decompose",
		Title:       "PRISM: Decompose",
		Description: "Split the proposal (or a node) into a few small digestible node.md parts. Recursive.",
	},
	{
		ID:    "drill",
		Title: "PRISM: Drill",
		Description: "Drill ONE part to atomic: node.md digest first (gate), then the agreed " +
			"artifact set (spec, detail, concept.drawio, signatures, tasks).",
	},
	{
		ID:    "integrate",
		Title: "PRISM: Integrate",
		Description: "Produce the cross-part artifacts — integration.drawio + combined signatures.md " +
			"+ overall tasks.md.",
	},
	{
		ID:    "apply",
		Title: "PRISM: Apply",
		Description: "Implement the change in code per the tasks, in dependency order (or the " +
			"ready frontier in parallel worktrees with --parallel), one commit per part, marking " +
			"tasks done and running checks.",
	},
	{
		ID:             "verify",
		ModelInvocable: true,
		Title:          "PRISM: Verify",
		Description: "Post-implementation verification on a running dev environment at a chosen " +
			"depth (--checks low|high) — full test suite, blocking static checks, diff-driven " +
			"functional and browser smoke, concurrency checks, load corner-cases at high, an " +
			"independent design-conformance review, fixes with re-verification, and a final " +
			"report. Project-agnostic: commands and entry points are detected from repository configs.",
	},
	{
		ID:             "archive",
		ModelInvocable: true,
		Title:          "PRISM: Archive",
		Description:    "Archive a completed change — move .prism/<change>/ to .prism/archive/<change>/.",
	},
}

// Body returns the tool-neutral template body of a workflow.
func Body(id string) string {
	b, err := templates.FS.ReadFile("commands/" + id + ".md")
	if err != nil {
		panic(fmt.Sprintf("workflow template %q not embedded: %v", id, err))
	}
	return strings.TrimRight(string(b), "\n")
}

// Conventions returns the shared methodology text installed as .prism/conventions.md.
func Conventions() string {
	b, err := templates.FS.ReadFile("conventions.md")
	if err != nil {
		panic(fmt.Sprintf("conventions template not embedded: %v", err))
	}
	return strings.TrimRight(string(b), "\n")
}
