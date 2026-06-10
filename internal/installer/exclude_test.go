package installer

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/mcoder33/prism/internal/adapters"
	"github.com/mcoder33/prism/internal/workflows"
)

func firstWorkflow(t *testing.T) workflows.Workflow {
	t.Helper()
	if len(workflows.All) == 0 {
		t.Fatal("no workflows registered")
	}
	return workflows.All[0]
}

func TestDetectToolsEmptyProject(t *testing.T) {
	if got := DetectTools(t.TempDir()); len(got) != 0 {
		t.Fatalf("DetectTools on empty dir = %v", got)
	}
}

func TestInstalledVersionEmptyWhenNotInstalled(t *testing.T) {
	if v := InstalledVersion(t.TempDir(), adapters.Claude); v != "" {
		t.Fatalf("InstalledVersion = %q, want empty", v)
	}
}

func TestInstallSharedWritesStampedConventions(t *testing.T) {
	root := t.TempDir()
	files, err := InstallShared(root)
	if err != nil {
		t.Fatal(err)
	}
	if len(files) != 1 || files[0] != ConventionsPath {
		t.Fatalf("InstallShared files = %v", files)
	}
	b, err := os.ReadFile(filepath.Join(root, ConventionsPath))
	if err != nil {
		t.Fatal(err)
	}
	content := string(b)
	if !strings.Contains(content, "prism:generated v") {
		t.Error("conventions missing generated stamp")
	}
	if !strings.Contains(content, "node.md") {
		t.Error("conventions missing methodology content")
	}
}

func TestAddToGitExcludeOutsideGitRepo(t *testing.T) {
	root := t.TempDir()
	if err := AddToGitExclude(root); err != nil {
		t.Fatalf("should no-op outside a git repo: %v", err)
	}
	if _, err := os.Stat(filepath.Join(root, ".git")); !os.IsNotExist(err) {
		t.Fatal(".git must not be created")
	}
}

func TestAddToGitExcludeCreatesInfoDir(t *testing.T) {
	root := t.TempDir()
	mustMkdir(t, filepath.Join(root, ".git"))

	if err := AddToGitExclude(root); err != nil {
		t.Fatal(err)
	}
	b, err := os.ReadFile(filepath.Join(root, ".git", "info", "exclude"))
	if err != nil {
		t.Fatal(err)
	}
	if string(b) != ".prism/\n" {
		t.Fatalf("exclude = %q", b)
	}
}

func TestAddToGitExcludeAppendsMissingNewline(t *testing.T) {
	root := t.TempDir()
	mustMkdir(t, filepath.Join(root, ".git", "info"))
	excludeFile := filepath.Join(root, ".git", "info", "exclude")
	if err := os.WriteFile(excludeFile, []byte("vendor"), 0o644); err != nil {
		t.Fatal(err)
	}

	if err := AddToGitExclude(root); err != nil {
		t.Fatal(err)
	}
	b, err := os.ReadFile(excludeFile)
	if err != nil {
		t.Fatal(err)
	}
	if string(b) != "vendor\n.prism/\n" {
		t.Fatalf("exclude = %q", b)
	}
}

func TestInstallCommandSingleFile(t *testing.T) {
	root := t.TempDir()
	w := firstWorkflow(t)
	rel, err := InstallCommand(root, adapters.Claude, w)
	if err != nil {
		t.Fatal(err)
	}
	if rel != adapters.Claude.CommandFile(w.ID) {
		t.Fatalf("rel = %q", rel)
	}
	if _, err := os.Stat(filepath.Join(root, rel)); err != nil {
		t.Fatal(err)
	}
}
