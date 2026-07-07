package installer

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/mcoder33/prism/internal/adapters"
	"github.com/mcoder33/prism/internal/workflows"
)

func TestUninstallToolRemovesGeneratedAndCleansDirs(t *testing.T) {
	root := t.TempDir()
	if _, err := InstallTool(root, adapters.Claude); err != nil {
		t.Fatal(err)
	}

	removed, err := UninstallTool(root, adapters.Claude)
	if err != nil {
		t.Fatal(err)
	}
	if len(removed) != len(workflows.All) {
		t.Errorf("removed %d files, want %d", len(removed), len(workflows.All))
	}
	// The prism command dir should be gone once emptied.
	dir := filepath.Dir(filepath.Join(root, adapters.Claude.CommandFile("propose")))
	if _, err := os.Stat(dir); !os.IsNotExist(err) {
		t.Errorf("empty command dir %s not cleaned up", dir)
	}
}

func TestUninstallToolSkipsUnstampedFiles(t *testing.T) {
	root := t.TempDir()
	if _, err := InstallTool(root, adapters.Claude); err != nil {
		t.Fatal(err)
	}
	// Replace one file with an unstamped, hand-written version.
	hand := filepath.Join(root, adapters.Claude.CommandFile(workflows.All[0].ID))
	if err := os.WriteFile(hand, []byte("no stamp here\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	removed, err := UninstallTool(root, adapters.Claude)
	if err != nil {
		t.Fatal(err)
	}
	if len(removed) != len(workflows.All)-1 {
		t.Errorf("removed %d, want %d (the unstamped file must survive)", len(removed), len(workflows.All)-1)
	}
	if _, err := os.Stat(hand); err != nil {
		t.Errorf("unstamped file was removed: %v", err)
	}
}

func TestUninstallSharedOnlyGenerated(t *testing.T) {
	root := t.TempDir()
	if _, err := InstallShared(root); err != nil {
		t.Fatal(err)
	}
	ok, err := UninstallShared(root)
	if err != nil || !ok {
		t.Fatalf("UninstallShared = %v, %v", ok, err)
	}
	if _, err := os.Stat(filepath.Join(root, ConventionsPath)); !os.IsNotExist(err) {
		t.Error("conventions.md not removed")
	}
	// Absent file: no-op, no error.
	if ok, err := UninstallShared(root); err != nil || ok {
		t.Errorf("second UninstallShared = %v, %v", ok, err)
	}
}

func TestIsGitExcluded(t *testing.T) {
	root := t.TempDir()
	if err := os.MkdirAll(filepath.Join(root, ".git"), 0o755); err != nil {
		t.Fatal(err)
	}
	if IsGitExcluded(root) {
		t.Error("no exclude file yet — should be false")
	}
	if err := AddToGitExclude(root); err != nil {
		t.Fatal(err)
	}
	if !IsGitExcluded(root) {
		t.Error("after AddToGitExclude — should be true")
	}
}
