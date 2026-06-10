package installer

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"gitlab.gidfinance.tech/zadolbator/prism/internal/adapters"
	"gitlab.gidfinance.tech/zadolbator/prism/internal/workflows"
)

func TestInstallToolAndDetectBack(t *testing.T) {
	root := t.TempDir()
	mustMkdir(t, filepath.Join(root, ".claude"))
	mustMkdir(t, filepath.Join(root, ".cursor"))

	detected := DetectTools(root)
	if got := toolIDs(detected); got != "claude,cursor" {
		t.Fatalf("DetectTools = %s, want claude,cursor", got)
	}
	if len(ConfiguredTools(root)) != 0 {
		t.Fatal("ConfiguredTools should be empty before install")
	}

	if _, err := InstallShared(root); err != nil {
		t.Fatal(err)
	}
	if _, err := InstallTool(root, adapters.Claude); err != nil {
		t.Fatal(err)
	}

	if _, err := os.Stat(filepath.Join(root, ConventionsPath)); err != nil {
		t.Fatalf("conventions not written: %v", err)
	}
	for _, w := range workflows.All {
		if _, err := os.Stat(filepath.Join(root, adapters.Claude.CommandFile(w.ID))); err != nil {
			t.Errorf("command file for %q not written: %v", w.ID, err)
		}
	}
	if got := toolIDs(ConfiguredTools(root)); got != "claude" {
		t.Fatalf("ConfiguredTools = %s, want claude", got)
	}
	if v := InstalledVersion(root, adapters.Claude); v != workflows.Version {
		t.Fatalf("InstalledVersion = %q, want %q", v, workflows.Version)
	}

	drill, err := os.ReadFile(filepath.Join(root, ".claude/commands/prism/drill.md"))
	if err != nil {
		t.Fatal(err)
	}
	for _, want := range []string{"`/prism:decompose`", ".prism/conventions.md"} {
		if !strings.Contains(string(drill), want) {
			t.Errorf("drill.md missing %q", want)
		}
	}
}

func TestAddToGitExcludeIdempotent(t *testing.T) {
	root := t.TempDir()
	mustMkdir(t, filepath.Join(root, ".git", "info"))
	excludeFile := filepath.Join(root, ".git", "info", "exclude")
	if err := os.WriteFile(excludeFile, []byte("# comment\n"), 0o644); err != nil {
		t.Fatal(err)
	}

	for range 2 {
		if err := AddToGitExclude(root); err != nil {
			t.Fatal(err)
		}
	}

	b, err := os.ReadFile(excludeFile)
	if err != nil {
		t.Fatal(err)
	}
	if got := strings.Count(string(b), ".prism/"); got != 1 {
		t.Fatalf(".prism/ appears %d times in exclude, want 1:\n%s", got, b)
	}
}

func mustMkdir(t *testing.T, path string) {
	t.Helper()
	if err := os.MkdirAll(path, 0o755); err != nil {
		t.Fatal(err)
	}
}

func toolIDs(tools []adapters.Tool) string {
	ids := make([]string, len(tools))
	for i, tl := range tools {
		ids[i] = tl.ID
	}
	return strings.Join(ids, ",")
}
