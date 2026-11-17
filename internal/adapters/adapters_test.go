package adapters

import (
	"strings"
	"testing"

	"github.com/mcoder33/prism/internal/workflows"
)

func TestEveryWorkflowHasTemplateBody(t *testing.T) {
	for _, w := range workflows.All {
		if len(workflows.Body(w.ID)) < 100 {
			t.Errorf("workflow %q: template body too short", w.ID)
		}
	}
}

func TestRenderLeavesNoPlaceholders(t *testing.T) {
	for _, a := range All {
		for _, d := range []Delivery{Commands, Skills} {
			if !a.Supports(d) {
				continue
			}
			for _, w := range workflows.All {
				for _, f := range a.Render(d, w, "0.0.0") {
					if strings.Contains(f.Content, "{{cmd:") || strings.Contains(f.Content, ConventionsPlaceholder) {
						t.Errorf("%s/%s/%s %s: leftover placeholder", a.ID, d, w.ID, f.Path)
					}
				}
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
		for _, d := range []Delivery{Commands, Skills} {
			if !a.Supports(d) {
				continue
			}
			for _, f := range a.Render(d, w, "1.2.3") {
				if got := ParseGeneratedVersion(f.Content); got != "1.2.3" {
					t.Errorf("%s %s: ParseGeneratedVersion = %q, want 1.2.3", a.ID, f.Path, got)
				}
			}
		}
	}
}

func TestSkillFollowsAgentSkillsFormat(t *testing.T) {
	w := workflows.All[0]
	for _, a := range All {
		if !a.Supports(Skills) {
			continue
		}
		files := a.Render(Skills, w, "1.2.3")
		skill := files[0]
		if want := a.SkillsDir + "/prism-" + w.ID + "/SKILL.md"; skill.Path != want {
			t.Fatalf("%s: skill path = %s, want %s", a.ID, skill.Path, want)
		}
		// frontmatter must open the file; name must equal the folder name.
		if !strings.HasPrefix(skill.Content, "---\nname: prism-"+w.ID+"\ndescription: ") {
			t.Errorf("%s: SKILL.md frontmatter malformed:\n%.120s", a.ID, skill.Content)
		}
		conv := a.SkillConventionsPath(w.ID)
		if !strings.Contains(skill.Content, "`"+conv+"`") {
			t.Errorf("%s: SKILL.md does not point at its bundled conventions %s", a.ID, conv)
		}
		found := false
		for _, f := range files {
			if f.Path == conv && strings.Contains(f.Content, "## Gates") {
				found = true
			}
		}
		if !found {
			t.Errorf("%s: bundled conventions file missing", a.ID)
		}
	}
}

func TestCodexIsSkillsOnlyWithDollarRefs(t *testing.T) {
	if Codex.Supports(Commands) || Codex.DefaultDelivery() != Skills {
		t.Fatal("codex must be skills-only")
	}
	if got := Codex.Resolve(Commands); got != Skills {
		t.Errorf("codex Resolve(commands) = %s, want skills", got)
	}
	if got := Codex.Ref(Skills, "drill"); got != "$prism-drill" {
		t.Errorf("codex ref = %s", got)
	}
	if got := Codex.SkillDir("drill"); got != ".agents/skills/prism-drill" {
		t.Errorf("codex skill dir = %s", got)
	}
	var verify workflows.Workflow
	for _, w := range workflows.All {
		if w.ID == "verify" {
			verify = w
		}
	}
	files := Codex.Render(Skills, verify, "0.0.0")
	if strings.Contains(files[0].Content, "$ARGUMENTS") {
		t.Error("codex skill still contains $ARGUMENTS, which Codex does not expand")
	}
	hasPolicy := false
	for _, f := range files {
		if strings.HasSuffix(f.Path, "agents/openai.yaml") && strings.Contains(f.Content, "allow_implicit_invocation: true") {
			hasPolicy = true
		}
	}
	if !hasPolicy {
		t.Error("codex verify skill must stay implicitly invocable (apply chains into it)")
	}
	if got := Codex.LegacyFiles("drill"); len(got) != 1 || got[0] != ".codex/prompts/prism-drill.md" {
		t.Errorf("codex legacy = %v", got)
	}
}

func TestClaudeSkillKeepsArgumentsAndIsUserInvoked(t *testing.T) {
	var verify workflows.Workflow
	for _, w := range workflows.All {
		if w.ID == "verify" {
			verify = w
		}
	}
	skill := Claude.Render(Skills, verify, "0.0.0")[0].Content
	if !strings.Contains(skill, "$ARGUMENTS") {
		t.Error("claude skill should keep $ARGUMENTS (expanded natively)")
	}
	if strings.Contains(skill, "disable-model-invocation") {
		t.Error("claude verify skill must stay model-invocable: apply runs it as a sub-step")
	}
	for _, w := range workflows.All {
		got := strings.Contains(Claude.Render(Skills, w, "0.0.0")[0].Content, "disable-model-invocation: true")
		if got == w.ModelInvocable {
			t.Errorf("%s: user-only = %v, want %v", w.ID, got, !w.ModelInvocable)
		}
	}
	if !strings.Contains(skill, "`/prism-") {
		t.Error("claude skill refs should use /prism-<id>")
	}
}

func TestCommandsUseSharedConventions(t *testing.T) {
	out := Claude.Render(Commands, workflows.All[0], "0.0.0")[0].Content
	if !strings.Contains(out, "`"+SharedConventionsPath+"`") {
		t.Errorf("command body should point at %s", SharedConventionsPath)
	}
}

func TestParseDelivery(t *testing.T) {
	for in, want := range map[string]Delivery{"": "", "skills": Skills, "Commands": Commands} {
		if got, err := ParseDelivery(in); err != nil || got != want {
			t.Errorf("ParseDelivery(%q) = %q, %v", in, got, err)
		}
	}
	if _, err := ParseDelivery("both"); err == nil {
		t.Error("ParseDelivery should reject unknown values")
	}
}
