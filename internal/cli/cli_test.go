package cli

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/mcoder33/prism/internal/adapters"
	"github.com/mcoder33/prism/internal/installer"
	"github.com/mcoder33/prism/internal/workflows"
)

func TestParseToolsFlag(t *testing.T) {
	if got, err := parseToolsFlag("all"); err != nil || len(got) != len(adapters.All) {
		t.Errorf("all = %d tools, err %v", len(got), err)
	}
	if got, err := parseToolsFlag("none"); err != nil || got != nil {
		t.Errorf("none = %v, err %v", got, err)
	}
	got, err := parseToolsFlag("claude, CURSOR")
	if err != nil || len(got) != 2 || got[0].ID != "claude" || got[1].ID != "cursor" {
		t.Errorf("list = %v, err %v", got, err)
	}
	if _, err := parseToolsFlag("claude,unknown"); err == nil {
		t.Error("unknown tool must error")
	} else if !strings.Contains(err.Error(), "unknown tool") {
		t.Errorf("error should name the problem: %v", err)
	}
}

func TestProjectRootArg(t *testing.T) {
	dir := t.TempDir()
	if got, err := projectRootArg([]string{dir}); err != nil || got != dir {
		t.Errorf("explicit dir = %q, err %v", got, err)
	}
	if _, err := projectRootArg([]string{filepath.Join(dir, "missing")}); err == nil {
		t.Error("missing dir must error")
	}
	cwd, _ := os.Getwd()
	if got, err := projectRootArg(nil); err != nil || got != cwd {
		t.Errorf("default = %q, err %v (cwd %q)", got, err, cwd)
	}
}

func TestRunInitWithToolsFlag(t *testing.T) {
	root := t.TempDir()
	if err := runInit(root, "claude"); err != nil {
		t.Fatal(err)
	}
	if _, err := os.Stat(filepath.Join(root, installer.ConventionsPath)); err != nil {
		t.Fatalf("conventions not written: %v", err)
	}
	for _, w := range workflows.All {
		if _, err := os.Stat(filepath.Join(root, adapters.Claude.CommandFile(w.ID))); err != nil {
			t.Errorf("command %q not written: %v", w.ID, err)
		}
	}
}

func TestRunInitToolsNoneIsNoop(t *testing.T) {
	root := t.TempDir()
	if err := runInit(root, "none"); err != nil {
		t.Fatal(err)
	}
	if _, err := os.Stat(filepath.Join(root, installer.PrismDir)); !os.IsNotExist(err) {
		t.Fatal(".prism must not be created for --tools none")
	}
}

// go test runs without a TTY, so the flagless path falls through to detection.
func TestRunInitNonTTYDetectsTools(t *testing.T) {
	if isTTY() {
		t.Skip("requires a non-TTY environment")
	}
	root := t.TempDir()
	if err := os.MkdirAll(filepath.Join(root, ".cursor"), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := runInit(root, ""); err != nil {
		t.Fatal(err)
	}
	if _, err := os.Stat(filepath.Join(root, adapters.Cursor.CommandFile("propose"))); err != nil {
		t.Fatalf("cursor commands not installed: %v", err)
	}
}

func TestRunInitNonTTYNoToolsErrors(t *testing.T) {
	if isTTY() {
		t.Skip("requires a non-TTY environment")
	}
	if err := runInit(t.TempDir(), ""); err == nil {
		t.Fatal("expected error when nothing is detected and not interactive")
	}
}

func TestRunUpdateNothingInstalled(t *testing.T) {
	if err := runUpdate(t.TempDir(), false); err != nil {
		t.Fatal(err)
	}
}

func TestRunUpdateRegeneratesStaleTool(t *testing.T) {
	root := t.TempDir()
	if err := runInit(root, "claude"); err != nil {
		t.Fatal(err)
	}

	// age one generated file so the tool reads as stale
	stale := filepath.Join(root, adapters.Claude.CommandFile(workflows.All[0].ID))
	b, err := os.ReadFile(stale)
	if err != nil {
		t.Fatal(err)
	}
	aged := strings.Replace(string(b), "prism:generated v"+workflows.Version, "prism:generated v0.0.1", 1)
	if err := os.WriteFile(stale, []byte(aged), 0o644); err != nil {
		t.Fatal(err)
	}
	if v := installer.InstalledVersion(root, adapters.Claude); v != "0.0.1" {
		t.Fatalf("setup failed: InstalledVersion = %q", v)
	}

	if err := runUpdate(root, false); err != nil {
		t.Fatal(err)
	}
	if v := installer.InstalledVersion(root, adapters.Claude); v != workflows.Version {
		t.Fatalf("after update InstalledVersion = %q, want %q", v, workflows.Version)
	}
}

func TestRunUpdateUpToDateWithoutForce(t *testing.T) {
	root := t.TempDir()
	if err := runInit(root, "claude"); err != nil {
		t.Fatal(err)
	}
	if err := runUpdate(root, false); err != nil {
		t.Fatal(err)
	}
	if err := runUpdate(root, true); err != nil {
		t.Fatal(err)
	}
}

func TestRunListVariants(t *testing.T) {
	// no .prism at all
	if err := runList(t.TempDir()); err != nil {
		t.Fatal(err)
	}

	// changes + CURRENT + archive
	root := t.TempDir()
	for _, dir := range []string{"rate-limiter", "auth-rework", "archive/old-change"} {
		if err := os.MkdirAll(filepath.Join(root, installer.PrismDir, dir), 0o755); err != nil {
			t.Fatal(err)
		}
	}
	current := filepath.Join(root, installer.PrismDir, "CURRENT")
	if err := os.WriteFile(current, []byte("rate-limiter\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := runList(root); err != nil {
		t.Fatal(err)
	}
}

func TestToolNames(t *testing.T) {
	got := toolNames([]adapters.Tool{adapters.Claude, adapters.Cursor})
	if got != "Claude Code, Cursor" {
		t.Errorf("toolNames = %q", got)
	}
}
