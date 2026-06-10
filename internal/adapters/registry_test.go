package adapters

import (
	"strings"
	"testing"

	"github.com/mcoder33/prism/internal/workflows"
)

func TestByIDKnownAndUnknown(t *testing.T) {
	for _, want := range All {
		got, ok := ByID(want.ID)
		if !ok || got.Name != want.Name {
			t.Errorf("ByID(%q) = %q, %v", want.ID, got.Name, ok)
		}
	}
	if _, ok := ByID("no-such-tool"); ok {
		t.Error("ByID should reject unknown ids")
	}
}

func TestIDsMatchAll(t *testing.T) {
	ids := IDs()
	if len(ids) != len(All) {
		t.Fatalf("IDs() len = %d, want %d", len(ids), len(All))
	}
	for i, tl := range All {
		if ids[i] != tl.ID {
			t.Errorf("IDs()[%d] = %q, want %q", i, ids[i], tl.ID)
		}
	}
}

func TestParseGeneratedVersionAbsent(t *testing.T) {
	if got := ParseGeneratedVersion("# just a file\nno stamp here\n"); got != "" {
		t.Errorf("ParseGeneratedVersion = %q, want empty", got)
	}
}

func TestResolveCommandRefs(t *testing.T) {
	body := "run {{cmd:drill}} then {{cmd:decompose}}"
	if got := ResolveCommandRefs(body, Claude); got != "run `/prism:drill` then `/prism:decompose`" {
		t.Errorf("claude refs = %q", got)
	}
	if got := ResolveCommandRefs(body, Cursor); got != "run `/prism-drill` then `/prism-decompose`" {
		t.Errorf("cursor refs = %q", got)
	}
}

func TestYamlQuoteEscaping(t *testing.T) {
	if got := yamlQuote(`say "hi" \ bye`); got != `"say \"hi\" \\ bye"` {
		t.Errorf("yamlQuote = %s", got)
	}
}

func TestGeminiGuardsTripleQuoteDelimiter(t *testing.T) {
	w := workflows.All[0]
	out := Gemini.Format(w, "text with ''' inside", "0.0.0")
	// the body's ''' must be neutralized so only the two TOML delimiters remain
	if got := strings.Count(out, "'''"); got != 2 {
		t.Errorf("gemini output has %d ''' delimiters, want exactly 2:\n%s", got, out)
	}
}

func TestDetectPathsDeclaredForEveryTool(t *testing.T) {
	for _, tl := range All {
		if len(tl.DetectPaths) == 0 {
			t.Errorf("%s: no DetectPaths", tl.ID)
		}
		if tl.CommandFile("x") == "" || tl.CommandRef("x") == "" {
			t.Errorf("%s: empty CommandFile/CommandRef", tl.ID)
		}
	}
}
