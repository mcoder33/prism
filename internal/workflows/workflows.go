// Package workflows is the registry of prism commands and their template bodies.
package workflows

import (
	"fmt"
	"strings"

	"github.com/mcoder33/prism/templates"
)

// Version of the prism CLI; stamped into every generated file.
// Overridable at build time: -ldflags "-X .../internal/workflows.Version=x.y.z".
var Version = "0.2.0"

// Workflow is one prism command installed into agents as a slash command.
type Workflow struct {
	// ID is the stable command id; also the file/slash-command suffix.
	ID string
	// Title is the human title used in frontmatter, e.g. "PRISM: Propose".
	Title string
	// Description is the one-liner shown by agents in command pickers.
	Description string
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
		Description: "Grill on requirements, survey best practices, pick a strategy + data-flow, " +
			"then write the seed (proposal + concept) for a new decomposition change.",
	},
	{
		ID:          "decompose",
		Title:       "PRISM: Decompose",
		Description: "Split the proposal (or a node) into a few small digestible node.md parts. Recursive.",
	},
	{
		ID:    "drill",
		Title: "PRISM: Drill",
		Description: "Drill ONE part to atomic and generate its artifact set " +
			"(spec, detail, concept.drawio, signatures, tasks).",
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
		Description: "Implement the change in code per the tasks, in dependency order, " +
			"marking tasks done and running checks.",
	},
	{
		ID:    "verify",
		Title: "PRISM: Verify",
		Description: "Thorough post-implementation verification on a running dev environment — " +
			"full test suite, blocking static checks, diff-driven functional and browser smoke, " +
			"targeted concurrency/parallelism checks, load corner-cases, ability to fix findings " +
			"and re-verify, final report with recommendations. Project-agnostic: commands and " +
			"entry points are detected from repository configs.",
	},
	{
		ID:          "archive",
		Title:       "PRISM: Archive",
		Description: "Archive a completed change — move .prism/<change>/ to .prism/archive/<change>/.",
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
