package workflows

import (
	"strings"
	"testing"
)

func TestWorkflowIDsUniqueAndKebab(t *testing.T) {
	seen := map[string]bool{}
	for _, w := range All {
		if seen[w.ID] {
			t.Errorf("duplicate workflow id %q", w.ID)
		}
		seen[w.ID] = true
		if w.ID != strings.ToLower(w.ID) || strings.ContainsAny(w.ID, " _") {
			t.Errorf("workflow id %q is not kebab-case", w.ID)
		}
		if w.Title == "" || w.Description == "" {
			t.Errorf("workflow %q: empty title or description", w.ID)
		}
	}
}

func TestBodyReturnsEveryWorkflow(t *testing.T) {
	for _, w := range All {
		body := Body(w.ID)
		if len(body) < 100 {
			t.Errorf("workflow %q: body too short (%d bytes)", w.ID, len(body))
		}
		if strings.HasSuffix(body, "\n") {
			t.Errorf("workflow %q: body should be trimmed of trailing newlines", w.ID)
		}
	}
}

func TestBodyPanicsOnUnknownWorkflow(t *testing.T) {
	defer func() {
		if recover() == nil {
			t.Fatal("Body(\"no-such-workflow\") should panic")
		}
	}()
	Body("no-such-workflow")
}

func TestConventionsContainCoreSections(t *testing.T) {
	c := Conventions()
	for _, want := range []string{"node.md", "drawio", ".prism/", "CURRENT"} {
		if !strings.Contains(c, want) {
			t.Errorf("conventions missing %q", want)
		}
	}
}

func TestVersionIsBareSemver(t *testing.T) {
	if Version == "" {
		t.Fatal("Version is empty")
	}
	if strings.HasPrefix(Version, "v") {
		t.Fatalf("Version %q should not carry a v prefix (it is stamped as v%%s)", Version)
	}
	if got := strings.Count(Version, "."); got != 2 {
		t.Fatalf("Version %q is not x.y.z", Version)
	}
}
