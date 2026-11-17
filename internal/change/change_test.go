package change

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// write lays out files under root; keys are slash paths.
func write(t *testing.T, root string, files map[string]string) {
	t.Helper()
	for rel, content := range files {
		p := filepath.Join(root, filepath.FromSlash(rel))
		if err := os.MkdirAll(filepath.Dir(p), 0o755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(p, []byte(content), 0o644); err != nil {
			t.Fatal(err)
		}
	}
}

const readme = `# demo

**Phase:** propose → decompose → **drill** → integrate → apply → verify
**Tier:** standard

| Node      | Title  | Status | Open    |
|-----------|--------|--------|---------|
| 01-parser | parser | 🟢     | 1 minor |
| 02-render | render | 🟡     | none    |
`

func node(status, open string) string {
	return "# NN — x\n\n- **What:** thing\n\n**Status:** " + status + "\n\n**Open:** " + open + "\n"
}

func codes(fs []Finding) string {
	var cs []string
	for _, f := range fs {
		cs = append(cs, f.Code)
	}
	return strings.Join(cs, ",")
}

func TestLoadCleanChange(t *testing.T) {
	root := t.TempDir()
	write(t, root, map[string]string{
		".prism/CURRENT":                  "demo\n",
		".prism/demo/README.md":           readme,
		".prism/demo/01-parser/node.md":   node("🟢", "\n- [minor] naming of the flag"),
		".prism/demo/01-parser/tasks.md":  "- [x] 1.1 a\n- [ ] 1.2 b\n",
		".prism/demo/02-render/node.md":   node("🟡", "none"),
		".prism/demo/02-render/sub/x.txt": "not a node dir",
	})
	name, err := Resolve(root, "")
	if err != nil || name != "demo" {
		t.Fatalf("Resolve = %q, %v", name, err)
	}
	c, err := Load(root, name)
	if err != nil {
		t.Fatal(err)
	}
	if c.Phase != "drill" || c.Tier != "standard" || !c.Current {
		t.Errorf("phase/tier/current = %s/%s/%v", c.Phase, c.Tier, c.Current)
	}
	if len(c.Nodes) != 2 {
		t.Fatalf("nodes = %d", len(c.Nodes))
	}
	p := c.Nodes[0]
	if p.Status != "🟢" || p.NodeStatus != "🟢" || p.Tasks != (Count{Done: 1, Total: 2}) || p.Open.Minor != 1 {
		t.Errorf("node 01 = %+v", p)
	}
	if len(c.Findings) != 0 {
		t.Errorf("unexpected findings: %v", c.Findings)
	}
}

func TestLoadReportsDiscrepancies(t *testing.T) {
	root := t.TempDir()
	write(t, root, map[string]string{
		".prism/demo/README.md": readme,
		// 01: table says 🟢 but node.md says 🟡, tasks.md missing, untagged + blocking open
		".prism/demo/01-parser/node.md": node("🟡", "\n- something feels off\n- [blocking] which codec?"),
		// 02 has no directory at all; 03 has a directory but no row
		".prism/demo/03-extra/node.md": node("⚪", "none"),
	})
	c, err := Load(root, "demo")
	if err != nil {
		t.Fatal(err)
	}
	got := codes(c.Findings)
	for _, want := range []string{"row_without_dir", "dir_without_row", "status_mismatch", "open_untagged",
		"blocking_open_past_drill", "missing_artifact"} {
		if !strings.Contains(got, want) {
			t.Errorf("missing finding %s in %s", want, got)
		}
	}
	if c.Errors() == 0 {
		t.Error("Errors() = 0")
	}
}

func TestResolveNeverGuessesBetweenChanges(t *testing.T) {
	root := t.TempDir()
	write(t, root, map[string]string{
		".prism/a/README.md": readme,
		".prism/b/README.md": readme,
	})
	if _, err := Resolve(root, ""); err == nil || !strings.Contains(err.Error(), "a, b") {
		t.Fatalf("Resolve should refuse to pick: %v", err)
	}
	write(t, root, map[string]string{".prism/b/README.md": readme})
	if n, _ := Resolve(root, "b"); n != "b" {
		t.Errorf("explicit name ignored: %q", n)
	}
}

func TestLoadStaleCurrent(t *testing.T) {
	root := t.TempDir()
	write(t, root, map[string]string{".prism/CURRENT": "gone\n"})
	_, err := Load(root, "gone")
	if err == nil || !strings.Contains(err.Error(), "stale") {
		t.Fatalf("expected stale CURRENT error, got %v", err)
	}
}

func TestCheckSpec(t *testing.T) {
	spec := `# Spec

## Requirement: A
### Scenario: ok
- WHEN x
- THEN y
### Scenario: no then
- WHEN x
## Requirement: B
## Requirement: C
## Requirement: D

## Coverage
- Scenario: ok → test x
`
	fs := CheckSpec(spec)
	got := codes(fs)
	if got != "scenario_incomplete,atomicity_spec" {
		t.Fatalf("findings = %s (%v)", got, fs)
	}
	if !strings.Contains(fs[0].Message, `"no then"`) || !strings.Contains(fs[0].Message, "THEN") {
		t.Errorf("message = %s", fs[0].Message)
	}
}

func TestValidateDrawio(t *testing.T) {
	dir := t.TempDir()
	cases := map[string]string{
		"ok.drawio":    `<mxfile><diagram><mxGraphModel><root><mxCell id="0"/></root></mxGraphModel></diagram></mxfile>`,
		"amp.drawio":   `<mxfile><diagram name="a & b"/></mxfile>`,
		"open.drawio":  `<mxfile><diagram>`,
		"wrong.drawio": `<svg></svg>`,
	}
	for name, content := range cases {
		write(t, dir, map[string]string{name: content})
	}
	if msg := ValidateDrawio(filepath.Join(dir, "ok.drawio")); msg != "" {
		t.Errorf("ok.drawio: %s", msg)
	}
	for _, bad := range []string{"amp.drawio", "open.drawio", "wrong.drawio"} {
		if ValidateDrawio(filepath.Join(dir, bad)) == "" {
			t.Errorf("%s should fail validation", bad)
		}
	}
}

func TestCheckFlagsBrokenArtifacts(t *testing.T) {
	root := t.TempDir()
	write(t, root, map[string]string{
		".prism/demo/README.md":          readme,
		".prism/demo/data-flow.drawio":   `<mxfile><diagram>`,
		".prism/demo/01-parser/node.md":  node("🟢", "none"),
		".prism/demo/01-parser/tasks.md": "- [ ] a\n",
		".prism/demo/01-parser/spec.md":  "## Requirement: A\n### Scenario: s\n- THEN y\n",
		".prism/demo/02-render/node.md":  node("🟡", "none"),
	})
	c, err := Load(root, "demo")
	if err != nil {
		t.Fatal(err)
	}
	c.Check(root)
	got := codes(c.Findings)
	if !strings.Contains(got, "drawio_invalid") || !strings.Contains(got, "scenario_incomplete") {
		t.Fatalf("findings = %s", got)
	}
	for _, f := range c.Findings {
		if f.Code == "scenario_incomplete" && f.Path != "01-parser/spec.md" {
			t.Errorf("finding path = %q", f.Path)
		}
	}
}

// The shipped worked example is the reference every user mirrors: it must stay clean.
func TestShippedExampleIsClean(t *testing.T) {
	root, err := filepath.Abs("../..")
	if err != nil {
		t.Fatal(err)
	}
	c, err := Load(root, "archive/example-json-list")
	if err != nil {
		t.Fatal(err)
	}
	c.Check(root)
	if len(c.Findings) != 0 {
		t.Fatalf("example has findings: %v", c.Findings)
	}
	if c.Tier != "small" || len(c.Nodes) != 1 || c.Nodes[0].Status != "✅" {
		t.Errorf("example parsed wrong: %+v", c)
	}
}
