package installer

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/mcoder33/prism/internal/adapters"
	"github.com/mcoder33/prism/internal/workflows"
)

func exists(t *testing.T, root, rel string) bool {
	t.Helper()
	_, err := os.Stat(filepath.Join(root, rel))
	return err == nil
}

func TestInstallSkillsWritesBundledConventions(t *testing.T) {
	root := t.TempDir()
	if _, err := InstallTool(root, adapters.Claude, adapters.Skills); err != nil {
		t.Fatal(err)
	}
	for _, w := range workflows.All {
		if !exists(t, root, adapters.Claude.SkillDir(w.ID)+"/SKILL.md") {
			t.Errorf("%s: SKILL.md missing", w.ID)
		}
		if !exists(t, root, adapters.Claude.SkillConventionsPath(w.ID)) {
			t.Errorf("%s: bundled conventions missing", w.ID)
		}
		if exists(t, root, adapters.Claude.CommandFile(w.ID)) {
			t.Errorf("%s: command file written under skills delivery", w.ID)
		}
	}
	if d := InstalledDelivery(root, adapters.Claude); d != adapters.Skills {
		t.Errorf("InstalledDelivery = %s", d)
	}
	if v := InstalledVersion(root, adapters.Claude); v != workflows.Version {
		t.Errorf("InstalledVersion = %q", v)
	}
}

func TestSwitchDeliveryRemovesOldFilesKeepsHandWritten(t *testing.T) {
	root := t.TempDir()
	if _, err := InstallTool(root, adapters.Claude, adapters.Commands); err != nil {
		t.Fatal(err)
	}
	// a hand-written command next to prism's must survive the switch
	hand := filepath.Join(root, ".claude/commands/prism/mine.md")
	if err := os.WriteFile(hand, []byte("my own command\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	if _, err := InstallTool(root, adapters.Claude, adapters.Skills); err != nil {
		t.Fatal(err)
	}
	for _, w := range workflows.All {
		if exists(t, root, adapters.Claude.CommandFile(w.ID)) {
			t.Errorf("%s: old command file not removed on switch", w.ID)
		}
	}
	if _, err := os.Stat(hand); err != nil {
		t.Errorf("hand-written command removed: %v", err)
	}
	if NeedsMigration(root, adapters.Claude, adapters.Skills) {
		t.Error("NeedsMigration after a clean switch")
	}
}

func TestCodexLegacyPromptsMigrateToSkills(t *testing.T) {
	root := t.TempDir()
	w := workflows.All[0]
	legacy := adapters.Codex.LegacyFiles(w.ID)[0]
	stamped := adapters.GeneratedStamp("0.7.0") + "\nold prompt\n"
	if err := os.MkdirAll(filepath.Join(root, filepath.Dir(legacy)), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(root, legacy), []byte(stamped), 0o644); err != nil {
		t.Fatal(err)
	}

	cfg := ConfiguredTools(root)
	if len(cfg) != 1 || cfg[0].ID != "codex" {
		t.Fatalf("ConfiguredTools = %v, want codex via its legacy prompt", cfg)
	}
	if InstalledVersion(root, adapters.Codex) != "" {
		t.Error("legacy-only install must read as stale (empty version)")
	}
	if !NeedsMigration(root, adapters.Codex, adapters.Skills) {
		t.Error("legacy prompt should need migration")
	}

	if _, err := InstallTool(root, adapters.Codex, ""); err != nil {
		t.Fatal(err)
	}
	if exists(t, root, legacy) {
		t.Error("legacy .codex/prompts file not removed")
	}
	if exists(t, root, ".codex/prompts") {
		t.Error("emptied .codex/prompts dir not pruned")
	}
	if !exists(t, root, ".agents/skills/prism-"+w.ID+"/agents/openai.yaml") {
		t.Error("codex skill policy file missing")
	}
}

func TestUninstallSkillsPrunesDirs(t *testing.T) {
	root := t.TempDir()
	if _, err := InstallTool(root, adapters.Codex, adapters.Skills); err != nil {
		t.Fatal(err)
	}
	removed, err := UninstallTool(root, adapters.Codex)
	if err != nil {
		t.Fatal(err)
	}
	if want := 3 * len(workflows.All); len(removed) != want {
		t.Errorf("removed %d files, want %d", len(removed), want)
	}
	if exists(t, root, ".agents/skills") {
		t.Error(".agents/skills not pruned after uninstall")
	}
	if !exists(t, root, ".agents") {
		t.Error("top-level tool dir must be kept")
	}
	if len(ConfiguredTools(root)) != 0 {
		t.Error("tool still configured after uninstall")
	}
}

func TestSkillBodiesPointAtTheirOwnConventions(t *testing.T) {
	root := t.TempDir()
	if _, err := InstallTool(root, adapters.Cursor, adapters.Skills); err != nil {
		t.Fatal(err)
	}
	b, err := os.ReadFile(filepath.Join(root, adapters.Cursor.SkillDir("drill"), "SKILL.md"))
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(b), adapters.Cursor.SkillConventionsPath("drill")) {
		t.Error("cursor drill skill does not reference its bundled conventions")
	}
	if strings.Contains(string(b), "`.prism/conventions.md`") {
		t.Error("skill still depends on the git-excluded .prism/conventions.md")
	}
}
