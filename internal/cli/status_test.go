package cli

import (
	"encoding/json"
	"errors"
	"io"
	"os"
	"path/filepath"
	"testing"

	"github.com/mcoder33/prism/internal/adapters"
	"github.com/mcoder33/prism/internal/installer"
	"github.com/mcoder33/prism/internal/workflows"
)

// captureStdout runs fn and returns what it printed to os.Stdout.
func captureStdout(t *testing.T, fn func() error) (string, error) {
	t.Helper()
	r, w, err := os.Pipe()
	if err != nil {
		t.Fatal(err)
	}
	orig := os.Stdout
	os.Stdout = w
	runErr := fn()
	os.Stdout = orig
	_ = w.Close()
	b, _ := io.ReadAll(r)
	return string(b), runErr
}

func writeChange(t *testing.T, root string, files map[string]string) {
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

const demoReadme = "# demo\n\n**Phase:** propose → **decompose**\n**Tier:** small\n\n" +
	"| Node | Title | Status | Open |\n|---|---|---|---|\n| 01-a | a | ⚪ | none |\n"

func TestStatusJSONIsOneDocument(t *testing.T) {
	root := t.TempDir()
	writeChange(t, root, map[string]string{
		".prism/CURRENT":           "demo",
		".prism/demo/README.md":    demoReadme,
		".prism/demo/01-a/node.md": "**What:** a\n\n**Status:** ⚪\n\n**Open:** none\n",
	})
	out, err := captureStdout(t, func() error { return runStatus(root, "", true) })
	if err != nil {
		t.Fatal(err)
	}
	var got struct {
		Change string `json:"change"`
		Phase  string `json:"phase"`
		Tier   string `json:"tier"`
		Nodes  []struct {
			ID     string `json:"id"`
			Status string `json:"status"`
		} `json:"nodes"`
		Findings []any `json:"findings"`
	}
	if err := json.Unmarshal([]byte(out), &got); err != nil {
		t.Fatalf("stdout is not one JSON document: %v\n%s", err, out)
	}
	if got.Change != "demo" || got.Phase != "decompose" || got.Tier != "small" ||
		len(got.Nodes) != 1 || got.Nodes[0].Status != "⚪" || len(got.Findings) != 0 {
		t.Errorf("status = %+v", got)
	}
}

func TestStatusJSONFailureKeepsContract(t *testing.T) {
	out, err := captureStdout(t, func() error { return runStatus(t.TempDir(), "", true) })
	if !errors.Is(err, errSilent) {
		t.Fatalf("err = %v, want errSilent", err)
	}
	var got map[string]string
	if json.Unmarshal([]byte(out), &got) != nil || got["error"] == "" {
		t.Fatalf("failure is not a JSON error document: %s", out)
	}
}

func TestCheckFailsOnBrokenDiagram(t *testing.T) {
	root := t.TempDir()
	writeChange(t, root, map[string]string{
		".prism/demo/README.md":        demoReadme,
		".prism/demo/01-a/node.md":     "**What:** a\n\n**Status:** ⚪\n\n**Open:** none\n",
		".prism/demo/data-flow.drawio": "<mxfile><diagram>",
	})
	if _, err := captureStdout(t, func() error { return runCheck(root, "", false) }); err == nil {
		t.Fatal("check should fail on a malformed drawio")
	}
	// fixed diagram → clean
	writeChange(t, root, map[string]string{".prism/demo/data-flow.drawio": "<mxfile><diagram/></mxfile>"})
	if _, err := captureStdout(t, func() error { return runCheck(root, "", false) }); err != nil {
		t.Fatalf("check on a clean change: %v", err)
	}
}

func TestCheckWithNoChangesIsNoop(t *testing.T) {
	if _, err := captureStdout(t, func() error { return runCheck(t.TempDir(), "", false) }); err != nil {
		t.Fatal(err)
	}
}

func TestInitDeliverySkillsAndUpdateSwitchBack(t *testing.T) {
	root := t.TempDir()
	if _, err := captureStdout(t, func() error { return runInit(root, "claude,codex", adapters.Skills) }); err != nil {
		t.Fatal(err)
	}
	first := workflows.All[0].ID
	if _, err := os.Stat(filepath.Join(root, adapters.Claude.SkillDir(first), "SKILL.md")); err != nil {
		t.Fatalf("claude skill missing: %v", err)
	}
	if _, err := os.Stat(filepath.Join(root, adapters.Codex.SkillDir(first), "SKILL.md")); err != nil {
		t.Fatalf("codex skill missing: %v", err)
	}

	// switch back to commands: claude moves, codex stays skills (it has no commands)
	if _, err := captureStdout(t, func() error { return runUpdate(root, false, adapters.Commands) }); err != nil {
		t.Fatal(err)
	}
	if d := installer.InstalledDelivery(root, adapters.Claude); d != adapters.Commands {
		t.Errorf("claude delivery = %s, want commands", d)
	}
	if _, err := os.Stat(filepath.Join(root, adapters.Claude.SkillDir(first))); !os.IsNotExist(err) {
		t.Error("claude skill dir not removed after switching to commands")
	}
	if d := installer.InstalledDelivery(root, adapters.Codex); d != adapters.Skills {
		t.Errorf("codex delivery = %s, want skills", d)
	}
}

func TestUpdateKeepsInstalledDelivery(t *testing.T) {
	root := t.TempDir()
	if _, err := captureStdout(t, func() error { return runInit(root, "cursor", adapters.Skills) }); err != nil {
		t.Fatal(err)
	}
	if _, err := captureStdout(t, func() error { return runUpdate(root, true, "") }); err != nil {
		t.Fatal(err)
	}
	if d := installer.InstalledDelivery(root, adapters.Cursor); d != adapters.Skills {
		t.Errorf("update without --delivery changed delivery to %s", d)
	}
}
