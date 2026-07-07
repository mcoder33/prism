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

func TestRunDoctorEmptyProject(t *testing.T) {
	// No .prism, nothing installed — must not error.
	if err := runDoctor(t.TempDir()); err != nil {
		t.Fatal(err)
	}
}

func TestRunDoctorInstalledProject(t *testing.T) {
	root := t.TempDir()
	if err := runInit(root, "claude"); err != nil {
		t.Fatal(err)
	}
	if err := runDoctor(root); err != nil {
		t.Fatal(err)
	}
}

func TestRunDoctorStaleCurrent(t *testing.T) {
	root := t.TempDir()
	prismDir := filepath.Join(root, installer.PrismDir)
	if err := os.MkdirAll(prismDir, 0o755); err != nil {
		t.Fatal(err)
	}
	// CURRENT names a change that does not exist on disk.
	if err := os.WriteFile(filepath.Join(prismDir, "CURRENT"), []byte("ghost-change\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := runDoctor(root); err != nil {
		t.Fatal(err)
	}
}

func TestRunUninstallRemovesGeneratedFiles(t *testing.T) {
	root := t.TempDir()
	if err := runInit(root, "claude"); err != nil {
		t.Fatal(err)
	}
	// Sanity: files exist before.
	first := filepath.Join(root, adapters.Claude.CommandFile(workflows.All[0].ID))
	if _, err := os.Stat(first); err != nil {
		t.Fatalf("setup: command not installed: %v", err)
	}

	if err := runUninstall(root, "", false); err != nil {
		t.Fatal(err)
	}
	for _, w := range workflows.All {
		p := filepath.Join(root, adapters.Claude.CommandFile(w.ID))
		if _, err := os.Stat(p); !os.IsNotExist(err) {
			t.Errorf("command %q still present after uninstall", w.ID)
		}
	}
	// conventions.md is kept by default.
	if _, err := os.Stat(filepath.Join(root, installer.ConventionsPath)); err != nil {
		t.Error("conventions.md should be kept without --shared")
	}
}

func TestRunUninstallSharedRemovesConventions(t *testing.T) {
	root := t.TempDir()
	if err := runInit(root, "claude"); err != nil {
		t.Fatal(err)
	}
	if err := runUninstall(root, "all", true); err != nil {
		t.Fatal(err)
	}
	if _, err := os.Stat(filepath.Join(root, installer.ConventionsPath)); !os.IsNotExist(err) {
		t.Error("conventions.md should be removed with --shared")
	}
}

func TestRunUninstallKeepsHandWrittenFile(t *testing.T) {
	root := t.TempDir()
	if err := runInit(root, "claude"); err != nil {
		t.Fatal(err)
	}
	// Overwrite one generated file with an unstamped, hand-written one.
	hand := filepath.Join(root, adapters.Claude.CommandFile(workflows.All[0].ID))
	if err := os.WriteFile(hand, []byte("# my own command, no stamp\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := runUninstall(root, "", false); err != nil {
		t.Fatal(err)
	}
	b, err := os.ReadFile(hand)
	if err != nil {
		t.Fatalf("hand-written file was removed: %v", err)
	}
	if !strings.Contains(string(b), "my own command") {
		t.Error("hand-written file content changed")
	}
}

func TestRunUninstallNothingInstalled(t *testing.T) {
	if err := runUninstall(t.TempDir(), "", false); err != nil {
		t.Fatal(err)
	}
}
