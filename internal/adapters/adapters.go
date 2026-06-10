// Package adapters renders prism workflows into per-AI-tool command files.
package adapters

import (
	"encoding/json"
	"fmt"
	"regexp"
	"strings"

	"github.com/mcoder33/prism/internal/workflows"
)

// Tool describes one supported AI agent: where its command files live,
// how a command is invoked, and how a command file is formatted.
type Tool struct {
	// ID is the stable id used in --tools, e.g. "claude".
	ID string
	// Name is the display name, e.g. "Claude Code".
	Name string
	// DetectPaths are project-relative paths whose presence means the tool is used.
	DetectPaths []string
	// CommandFile returns the project-relative path of the generated file for a workflow.
	CommandFile func(workflowID string) string
	// CommandRef is how the user invokes a workflow in this tool, e.g. "/prism:drill".
	CommandRef func(workflowID string) string
	// Format renders the full command file (frontmatter + body) for this tool.
	Format func(w workflows.Workflow, body, version string) string
}

var cmdRefPattern = regexp.MustCompile(`\{\{cmd:([a-z-]+)\}\}`)

// ResolveCommandRefs replaces {{cmd:<id>}} placeholders with the tool-specific slash command.
func ResolveCommandRefs(body string, t Tool) string {
	return cmdRefPattern.ReplaceAllStringFunc(body, func(m string) string {
		id := cmdRefPattern.FindStringSubmatch(m)[1]
		return "`" + t.CommandRef(id) + "`"
	})
}

// GeneratedStamp marks a file as tool-owned and records the producing version.
func GeneratedStamp(version string) string {
	return fmt.Sprintf(
		"<!-- prism:generated v%s — managed by the prism CLI, do not edit (run `prism update` to regenerate) -->",
		version,
	)
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
}

var Cursor = Tool{
	ID:          "cursor",
	Name:        "Cursor",
	DetectPaths: []string{".cursor"},
	CommandFile: func(id string) string { return ".cursor/commands/prism-" + id + ".md" },
	CommandRef:  func(id string) string { return "/prism-" + id },
	Format:      plainMarkdown,
}

var Codex = Tool{
	ID:          "codex",
	Name:        "Codex CLI",
	DetectPaths: []string{".codex"},
	CommandFile: func(id string) string { return ".codex/prompts/prism-" + id + ".md" },
	CommandRef:  func(id string) string { return "/prism-" + id },
	Format:      plainMarkdown,
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
			fmt.Sprintf("# prism:generated v%s — managed by the prism CLI, do not edit (run `prism update` to regenerate)", version),
			"description = " + string(desc),
			"prompt = '''",
			"# " + w.Title,
			"",
			safeBody,
			"'''",
			"",
		}, "\n")
	},
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
}

var All = []Tool{Claude, Cursor, Codex, Gemini, Copilot, Windsurf, OpenCode}

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
