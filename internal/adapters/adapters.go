// Package adapters renders prism workflows into per-AI-tool command files.
package adapters

import (
	"encoding/json"
	"fmt"
	"path"
	"regexp"
	"strings"

	"github.com/mcoder33/prism/internal/workflows"
)

// Delivery is how a workflow reaches the agent.
type Delivery string

const (
	// Commands delivers one slash-command / prompt / workflow file per workflow.
	Commands Delivery = "commands"
	// Skills delivers one Agent Skills folder per workflow: SKILL.md plus a bundled
	// references/conventions.md, so a committed skill needs nothing from .prism/.
	Skills Delivery = "skills"
)

// ParseDelivery validates a --delivery flag value ("" means "tool default").
func ParseDelivery(s string) (Delivery, error) {
	switch d := Delivery(strings.ToLower(strings.TrimSpace(s))); d {
	case "", Commands, Skills:
		return d, nil
	default:
		return "", fmt.Errorf("unknown delivery %q; use %q or %q", s, Commands, Skills)
	}
}

// File is one generated file: project-relative path and full content.
type File struct {
	Path    string
	Content string
}

// Tool describes one supported AI agent: where its generated files live,
// how a workflow is invoked, and how a file is formatted.
type Tool struct {
	// ID is the stable id used in --tools, e.g. "claude".
	ID string
	// Name is the display name, e.g. "Claude Code".
	Name string
	// DetectPaths are project-relative paths whose presence means the tool is used.
	DetectPaths []string
	// CommandFile returns the project-relative path of the generated command file for a
	// workflow. nil when the tool has no command surface (skills-only, e.g. Codex).
	CommandFile func(workflowID string) string
	// CommandRef is how the user invokes a command in this tool, e.g. "/prism:drill".
	CommandRef func(workflowID string) string
	// Format renders the full command file (frontmatter + body) for this tool.
	Format func(w workflows.Workflow, body, version string) string

	// SkillsDir is the project-relative directory the tool scans for Agent Skills
	// ("" = the tool has no project skills support).
	SkillsDir string
	// SkillRef is how the user invokes a prism skill; nil means "/prism-<id>".
	SkillRef func(workflowID string) string
	// SkillArgsNative: the tool substitutes $ARGUMENTS in skill bodies itself.
	SkillArgsNative bool
	// SkillUserOnly marks gated skills user-invoked only (Claude Code's
	// disable-model-invocation); workflows.ModelInvocable ones stay callable by the model.
	SkillUserOnly bool
	// SkillExtras returns extra per-skill files, paths relative to the skill dir.
	SkillExtras func(w workflows.Workflow, version string) []File
	// LegacyFiles lists paths earlier prism versions generated for a workflow; they are
	// removed (when still prism-stamped) on install, update and uninstall.
	LegacyFiles func(workflowID string) []string
}

// Supports reports whether the tool can receive the given delivery.
func (t Tool) Supports(d Delivery) bool {
	switch d {
	case Commands:
		return t.CommandFile != nil
	case Skills:
		return t.SkillsDir != ""
	}
	return false
}

// DefaultDelivery is commands where the tool has them (backward compatible), else skills.
func (t Tool) DefaultDelivery() Delivery {
	if t.Supports(Commands) {
		return Commands
	}
	return Skills
}

// Resolve returns the requested delivery when the tool supports it, else its default.
func (t Tool) Resolve(requested Delivery) Delivery {
	if requested != "" && t.Supports(requested) {
		return requested
	}
	return t.DefaultDelivery()
}

// SkillName is the Agent Skills name (and folder name) of a workflow.
func SkillName(workflowID string) string { return "prism-" + workflowID }

// SkillDir is the project-relative skill folder of a workflow.
func (t Tool) SkillDir(workflowID string) string {
	return t.SkillsDir + "/" + SkillName(workflowID)
}

// SkillConventionsPath is the bundled conventions copy inside a workflow's skill folder.
func (t Tool) SkillConventionsPath(workflowID string) string {
	return t.SkillDir(workflowID) + "/references/conventions.md"
}

// Ref is how the user invokes a workflow under a delivery.
func (t Tool) Ref(d Delivery, workflowID string) string {
	if d == Skills || t.CommandRef == nil {
		if t.SkillRef != nil {
			return t.SkillRef(workflowID)
		}
		return "/" + SkillName(workflowID)
	}
	return t.CommandRef(workflowID)
}

// Location is a short human description of where a delivery's files land.
func (t Tool) Location(d Delivery) string {
	if d == Skills {
		return t.SkillsDir + "/prism-*/"
	}
	return path.Dir(t.CommandFile(workflows.All[0].ID)) + "/"
}

// SharedConventionsPath is the git-excluded conventions copy read by command deliveries.
const SharedConventionsPath = ".prism/conventions.md"

// ConventionsFile renders the stamped conventions document.
func ConventionsFile(version string) string {
	return GeneratedStamp(version) + "\n\n" + workflows.Conventions() + "\n"
}

// Render produces every file of one workflow for this tool under a delivery.
func (t Tool) Render(d Delivery, w workflows.Workflow, version string) []File {
	body := ResolveRefs(workflows.Body(w.ID), t, d, w.ID)
	if d != Skills {
		return []File{{Path: t.CommandFile(w.ID), Content: t.Format(w, body, version)}}
	}
	if !t.SkillArgsNative {
		body = strings.ReplaceAll(body, "`$ARGUMENTS`", "the text after the skill name")
		body = strings.ReplaceAll(body, "$ARGUMENTS", "the text after the skill name")
	}
	dir := t.SkillDir(w.ID)
	files := []File{
		{Path: dir + "/SKILL.md", Content: t.skillMarkdown(w, body, version)},
		{Path: t.SkillConventionsPath(w.ID), Content: ConventionsFile(version)},
	}
	if t.SkillExtras != nil {
		for _, f := range t.SkillExtras(w, version) {
			files = append(files, File{Path: dir + "/" + f.Path, Content: f.Content})
		}
	}
	return files
}

// Paths lists every project-relative path prism may have generated for a workflow in this
// tool — both deliveries plus legacy layouts. Used for detection, migration and uninstall.
func (t Tool) Paths(w workflows.Workflow) []string {
	var ps []string
	for _, d := range []Delivery{Commands, Skills} {
		if t.Supports(d) {
			for _, f := range t.Render(d, w, "0.0.0") {
				ps = append(ps, f.Path)
			}
		}
	}
	if t.LegacyFiles != nil {
		ps = append(ps, t.LegacyFiles(w.ID)...)
	}
	return ps
}

// PrimaryFile is the file whose presence marks a workflow as installed under a delivery.
func (t Tool) PrimaryFile(d Delivery, workflowID string) string {
	if d == Skills {
		return t.SkillDir(workflowID) + "/SKILL.md"
	}
	return t.CommandFile(workflowID)
}

// skillMarkdown renders SKILL.md per the Agent Skills format: name must equal the folder.
func (t Tool) skillMarkdown(w workflows.Workflow, body, version string) string {
	lines := []string{
		"---",
		"name: " + SkillName(w.ID),
		"description: " + yamlQuote(w.Description),
	}
	if t.SkillUserOnly && !w.ModelInvocable {
		lines = append(lines, "disable-model-invocation: true")
	}
	lines = append(lines,
		"---",
		GeneratedStamp(version),
		"",
		"# "+w.Title,
		"",
		body,
		"",
	)
	return strings.Join(lines, "\n")
}

var cmdRefPattern = regexp.MustCompile(`\{\{cmd:([a-z-]+)\}\}`)

// ConventionsPlaceholder in a template body names the conventions file to read.
const ConventionsPlaceholder = "{{conventions}}"

// ResolveCommandRefs resolves a body for the tool's command delivery.
func ResolveCommandRefs(body string, t Tool) string {
	return ResolveRefs(body, t, Commands, "")
}

// ResolveRefs replaces {{cmd:<id>}} with the tool's invocation under the delivery, and
// {{conventions}} with the conventions path that delivery ships: the shared
// .prism/conventions.md for commands, the skill's own bundled copy for skills.
func ResolveRefs(body string, t Tool, d Delivery, workflowID string) string {
	conv := SharedConventionsPath
	if d == Skills {
		conv = t.SkillConventionsPath(workflowID)
	}
	body = strings.ReplaceAll(body, ConventionsPlaceholder, "`"+conv+"`")
	return cmdRefPattern.ReplaceAllStringFunc(body, func(m string) string {
		id := cmdRefPattern.FindStringSubmatch(m)[1]
		return "`" + t.Ref(d, id) + "`"
	})
}

// GeneratedStampText is the stamp sentence; wrapped per file syntax (HTML/TOML/YAML comment).
func GeneratedStampText(version string) string {
	return fmt.Sprintf(
		"prism:generated v%s — managed by the prism CLI, do not edit (run `prism update` to regenerate)",
		version,
	)
}

// GeneratedStamp marks a Markdown file as tool-owned and records the producing version.
func GeneratedStamp(version string) string {
	return "<!-- " + GeneratedStampText(version) + " -->"
}

var generatedVersionPattern = regexp.MustCompile(`prism:generated v([0-9][^\s]*)`)

// ParseGeneratedVersion extracts the version from a previously generated file ("" if absent).
func ParseGeneratedVersion(content string) string {
	m := generatedVersionPattern.FindStringSubmatch(content)
	if m == nil {
		return ""
	}
	return m[1]
}

func yamlQuote(s string) string {
	s = strings.ReplaceAll(s, `\`, `\\`)
	s = strings.ReplaceAll(s, `"`, `\"`)
	return `"` + s + `"`
}

// markdownWithFrontmatter renders a command file with YAML frontmatter (Claude Code style).
func markdownWithFrontmatter(w workflows.Workflow, body, version string, extra ...string) string {
	lines := []string{
		"---",
		"name: " + yamlQuote(w.Title),
		"description: " + yamlQuote(w.Description),
	}
	lines = append(lines, extra...)
	lines = append(lines,
		"---",
		GeneratedStamp(version),
		"",
		body,
		"",
	)
	return strings.Join(lines, "\n")
}

// plainMarkdown renders a command file for tools without frontmatter support.
func plainMarkdown(w workflows.Workflow, body, version string) string {
	return strings.Join([]string{
		"# " + w.Title,
		"",
		"> " + w.Description,
		"",
		GeneratedStamp(version),
		"",
		body,
		"",
	}, "\n")
}

var Claude = Tool{
	ID:          "claude",
	Name:        "Claude Code",
	DetectPaths: []string{".claude"},
	CommandFile: func(id string) string { return ".claude/commands/prism/" + id + ".md" },
	CommandRef:  func(id string) string { return "/prism:" + id },
	Format: func(w workflows.Workflow, body, version string) string {
		return markdownWithFrontmatter(w, body, version,
			"category: Workflow",
			"tags: [workflow, design, prism]",
		)
	},
	SkillsDir:       ".claude/skills",
	SkillArgsNative: true,
	// gated workflows have side effects: the user starts them; sub-steps stay model-callable.
	SkillUserOnly: true,
}

var Cursor = Tool{
	ID:          "cursor",
	Name:        "Cursor",
	DetectPaths: []string{".cursor"},
	CommandFile: func(id string) string { return ".cursor/commands/prism-" + id + ".md" },
	CommandRef:  func(id string) string { return "/prism-" + id },
	Format:      plainMarkdown,
	SkillsDir:   ".cursor/skills",
}

// Codex is skills-only: custom prompts are deprecated upstream, and Codex only ever read
// them from ~/.codex/prompts/ — never from a project's .codex/prompts/. Project skills live
// in .agents/skills/.
var Codex = Tool{
	ID:          "codex",
	Name:        "Codex CLI",
	DetectPaths: []string{".codex", ".agents"},
	SkillsDir:   ".agents/skills",
	// Codex invokes skills with $; a bare /prism-<id> is not recognized.
	SkillRef: func(id string) string { return "$" + SkillName(id) },
	SkillExtras: func(w workflows.Workflow, version string) []File {
		return []File{{
			Path: "agents/openai.yaml",
			Content: strings.Join([]string{
				"# " + GeneratedStampText(version),
				"interface:",
				"  display_name: " + yamlQuote(w.Title),
				"policy:",
				fmt.Sprintf("  allow_implicit_invocation: %t", w.ModelInvocable),
				"",
			}, "\n"),
		}}
	},
	LegacyFiles: func(id string) []string { return []string{".codex/prompts/prism-" + id + ".md"} },
}

var Gemini = Tool{
	ID:          "gemini",
	Name:        "Gemini CLI",
	DetectPaths: []string{".gemini"},
	CommandFile: func(id string) string { return ".gemini/commands/prism/" + id + ".toml" },
	CommandRef:  func(id string) string { return "/prism:" + id },
	Format: func(w workflows.Workflow, body, version string) string {
		// TOML literal multi-line strings process no escapes; guard the delimiter
		// (zero-width space between quotes). Gemini CLI injects invocation args
		// via {{args}}, not $ARGUMENTS.
		safeBody := strings.ReplaceAll(body, "'''", "''\u200b'")
		safeBody = strings.ReplaceAll(safeBody, "$ARGUMENTS", "{{args}}")
		desc, _ := json.Marshal(w.Description)
		return strings.Join([]string{
			"# " + GeneratedStampText(version),
			"description = " + string(desc),
			"prompt = '''",
			"# " + w.Title,
			"",
			safeBody,
			"'''",
			"",
		}, "\n")
	},
	SkillsDir: ".gemini/skills",
}

var Copilot = Tool{
	ID:          "copilot",
	Name:        "GitHub Copilot",
	DetectPaths: []string{".github/copilot-instructions.md", ".github/prompts"},
	CommandFile: func(id string) string { return ".github/prompts/prism-" + id + ".prompt.md" },
	CommandRef:  func(id string) string { return "/prism-" + id },
	Format: func(w workflows.Workflow, body, version string) string {
		return markdownWithFrontmatter(w, body, version)
	},
	SkillsDir: ".github/skills",
}

var Windsurf = Tool{
	ID:          "windsurf",
	Name:        "Windsurf",
	DetectPaths: []string{".windsurf"},
	CommandFile: func(id string) string { return ".windsurf/workflows/prism-" + id + ".md" },
	CommandRef:  func(id string) string { return "/prism-" + id },
	Format: func(w workflows.Workflow, body, version string) string {
		return markdownWithFrontmatter(w, body, version)
	},
}

var OpenCode = Tool{
	ID:          "opencode",
	Name:        "OpenCode",
	DetectPaths: []string{".opencode", "opencode.json"},
	CommandFile: func(id string) string { return ".opencode/command/prism-" + id + ".md" },
	CommandRef:  func(id string) string { return "/prism-" + id },
	Format: func(w workflows.Workflow, body, version string) string {
		return markdownWithFrontmatter(w, body, version)
	},
	SkillsDir: ".opencode/skills",
}

var RooCode = Tool{
	ID:          "roo",
	Name:        "Roo Code",
	DetectPaths: []string{".roo", ".roomodes"},
	CommandFile: func(id string) string { return ".roo/commands/prism-" + id + ".md" },
	CommandRef:  func(id string) string { return "/prism-" + id },
	Format: func(w workflows.Workflow, body, version string) string {
		return markdownWithFrontmatter(w, body, version)
	},
	SkillsDir: ".roo/skills",
}

var Cline = Tool{
	ID:          "cline",
	Name:        "Cline",
	DetectPaths: []string{".clinerules"},
	CommandFile: func(id string) string { return ".clinerules/workflows/prism-" + id + ".md" },
	// Cline invokes a workflow by its filename, extension included.
	CommandRef: func(id string) string { return "/prism-" + id + ".md" },
	Format:     plainMarkdown,
	SkillsDir:  ".cline/skills",
}

var All = []Tool{Claude, Cursor, Codex, Gemini, Copilot, Windsurf, OpenCode, RooCode, Cline}

// ByID returns the adapter with the given id, or false.
func ByID(id string) (Tool, bool) {
	for _, t := range All {
		if t.ID == id {
			return t, true
		}
	}
	return Tool{}, false
}

// IDs lists all known tool ids, for help texts and error messages.
func IDs() []string {
	ids := make([]string, len(All))
	for i, t := range All {
		ids[i] = t.ID
	}
	return ids
}
