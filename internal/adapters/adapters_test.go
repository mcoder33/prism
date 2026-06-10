package adapters

import (
	"strings"
	"testing"

	"gitlab.gidfinance.tech/zadolbator/prism/internal/workflows"
)

func TestEveryWorkflowHasTemplateBody(t *testing.T) {
	for _, w := range workflows.All {
		if len(workflows.Body(w.ID)) < 100 {
			t.Errorf("workflow %q: template body too short", w.ID)
		}
	}
}

func TestCommandRefsResolveForEveryAdapter(t *testing.T) {
	for _, a := range All {
		for _, w := range workflows.All {
			body := ResolveCommandRefs(workflows.Body(w.ID), a)
			if strings.Contains(body, "{{cmd:") {
				t.Errorf("%s/%s: leftover {{cmd:}} placeholder", a.ID, w.ID)
			}
		}
	}
}

func TestNamespacedVsFlatNaming(t *testing.T) {
	if got := Claude.CommandFile("drill"); got != ".claude/commands/prism/drill.md" {
		t.Errorf("claude path = %s", got)
	}
	if got := Claude.CommandRef("drill"); got != "/prism:drill" {
		t.Errorf("claude ref = %s", got)
	}
	if got := Cursor.CommandFile("drill"); got != ".cursor/commands/prism-drill.md" {
		t.Errorf("cursor path = %s", got)
	}
	if got := Cursor.CommandRef("drill"); got != "/prism-drill" {
		t.Errorf("cursor ref = %s", got)
	}
}

func TestGeminiUsesArgsPlaceholder(t *testing.T) {
	var verify workflows.Workflow
	for _, w := range workflows.All {
		if w.ID == "verify" {
			verify = w
		}
	}
	out := Gemini.Format(verify, workflows.Body("verify"), "0.0.0")
	if !strings.Contains(out, "{{args}}") {
		t.Error("gemini output missing {{args}}")
	}
	if strings.Contains(out, "$ARGUMENTS") {
		t.Error("gemini output still contains $ARGUMENTS")
	}
	if !strings.Contains(out, "\ndescription = ") {
		t.Error("gemini output missing TOML description")
	}
}

func TestGeneratedFilesCarryParseableVersionStamp(t *testing.T) {
	w := workflows.All[0]
	for _, a := range All {
		out := a.Format(w, "body", "1.2.3")
		if got := ParseGeneratedVersion(out); got != "1.2.3" {
			t.Errorf("%s: ParseGeneratedVersion = %q, want 1.2.3", a.ID, got)
		}
	}
}
