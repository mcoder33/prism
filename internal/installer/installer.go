// Package installer detects AI tools in a project and writes prism command files.
package installer

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/mcoder33/prism/internal/adapters"
	"github.com/mcoder33/prism/internal/workflows"
)

const (
	// PrismDir holds all prism artifacts at the project root (git-excluded).
	PrismDir = ".prism"
	// ConventionsPath is the shared methodology file referenced by every command.
	ConventionsPath = PrismDir + "/conventions.md"
)

// DetectTools returns tools whose dot-dirs/config files are present in the project.
func DetectTools(projectRoot string) []adapters.Tool {
	var found []adapters.Tool
	for _, t := range adapters.All {
		for _, p := range t.DetectPaths {
			if _, err := os.Stat(filepath.Join(projectRoot, p)); err == nil {
				found = append(found, t)
				break
			}
		}
	}
	return found
}

// ConfiguredTools returns tools that already have prism-generated files installed
// (either delivery, or a legacy layout awaiting migration).
func ConfiguredTools(projectRoot string) []adapters.Tool {
	var found []adapters.Tool
	for _, t := range adapters.All {
		if len(generatedPaths(projectRoot, t)) > 0 {
			found = append(found, t)
		}
	}
	return found
}

// generatedPaths lists the prism-stamped files present for a tool, across every layout.
func generatedPaths(projectRoot string, t adapters.Tool) []string {
	var found []string
	for _, w := range workflows.All {
		for _, rel := range t.Paths(w) {
			if isGenerated(filepath.Join(projectRoot, rel)) {
				found = append(found, rel)
			}
		}
	}
	return found
}

// InstalledDelivery reports which delivery a tool is installed with. A tool with nothing
// installed (or only legacy files) reports its default delivery.
func InstalledDelivery(projectRoot string, t adapters.Tool) adapters.Delivery {
	for _, d := range []adapters.Delivery{adapters.Skills, adapters.Commands} {
		if !t.Supports(d) {
			continue
		}
		for _, w := range workflows.All {
			if isGenerated(filepath.Join(projectRoot, t.PrimaryFile(d, w.ID))) {
				return d
			}
		}
	}
	return t.DefaultDelivery()
}

// InstalledVersion reads the prism version stamp from the tool's installed primary files
// ("" if none). Legacy-only installs report "" so update treats them as stale.
func InstalledVersion(projectRoot string, t adapters.Tool) string {
	d := InstalledDelivery(projectRoot, t)
	for _, w := range workflows.All {
		b, err := os.ReadFile(filepath.Join(projectRoot, t.PrimaryFile(d, w.ID)))
		if err == nil {
			return adapters.ParseGeneratedVersion(string(b))
		}
	}
	return ""
}

// NeedsMigration reports prism-generated files outside the tool's current delivery
// (the other delivery or a legacy layout) — update removes them.
func NeedsMigration(projectRoot string, t adapters.Tool, d adapters.Delivery) bool {
	keep := map[string]bool{}
	for _, w := range workflows.All {
		for _, f := range t.Render(d, w, workflows.Version) {
			keep[f.Path] = true
		}
	}
	for _, rel := range generatedPaths(projectRoot, t) {
		if !keep[rel] {
			return true
		}
	}
	return false
}

func writeFileEnsured(path, content string) error {
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return err
	}
	return os.WriteFile(path, []byte(content), 0o644)
}

// InstallShared writes .prism/conventions.md (read by command deliveries) and excludes
// .prism/ from git. Skill deliveries bundle their own copy, but active changes and
// `prism status`/`check` still live under .prism/.
func InstallShared(projectRoot string) ([]string, error) {
	content := adapters.ConventionsFile(workflows.Version)
	if err := writeFileEnsured(filepath.Join(projectRoot, ConventionsPath), content); err != nil {
		return nil, fmt.Errorf("write conventions: %w", err)
	}
	if err := AddToGitExclude(projectRoot); err != nil {
		return nil, fmt.Errorf("update .git/info/exclude: %w", err)
	}
	return []string{ConventionsPath}, nil
}

// InstallWorkflow writes one workflow for a tool under a delivery, then removes the
// prism-generated files of that workflow in any other layout (the other delivery, legacy
// paths) — switching delivery or upgrading from an old layout is just a reinstall.
// Files are tool-owned: always overwritten. Returns the written project-relative paths.
func InstallWorkflow(projectRoot string, t adapters.Tool, d adapters.Delivery, w workflows.Workflow) ([]string, error) {
	d = t.Resolve(d)
	files := t.Render(d, w, workflows.Version)
	written := map[string]bool{}
	var out []string
	for _, f := range files {
		if err := writeFileEnsured(filepath.Join(projectRoot, f.Path), f.Content); err != nil {
			return nil, fmt.Errorf("write %s: %w", f.Path, err)
		}
		written[f.Path] = true
		out = append(out, f.Path)
	}
	var stale []string
	for _, rel := range t.Paths(w) {
		if !written[rel] {
			stale = append(stale, rel)
		}
	}
	if _, err := removeGenerated(projectRoot, stale); err != nil {
		return nil, err
	}
	return out, nil
}

// InstallTool writes every workflow for one tool under a delivery ("" = tool default).
func InstallTool(projectRoot string, t adapters.Tool, d adapters.Delivery) ([]string, error) {
	var files []string
	for _, w := range workflows.All {
		written, err := InstallWorkflow(projectRoot, t, d, w)
		if err != nil {
			return nil, err
		}
		files = append(files, written...)
	}
	return files, nil
}

// AddToGitExclude adds .prism/ to .git/info/exclude so artifacts are never committed.
// No-op outside a git repo or when an entry already exists.
func AddToGitExclude(projectRoot string) error {
	gitDir := filepath.Join(projectRoot, ".git")
	if st, err := os.Stat(gitDir); err != nil || !st.IsDir() {
		return nil
	}
	excludeFile := filepath.Join(gitDir, "info", "exclude")
	current, err := os.ReadFile(excludeFile)
	if err == nil {
		for _, line := range strings.Split(string(current), "\n") {
			trimmed := strings.TrimSpace(line)
			if trimmed == ".prism/" || trimmed == ".prism" {
				return nil
			}
		}
	} else if err := os.MkdirAll(filepath.Dir(excludeFile), 0o755); err != nil {
		return err
	}
	entry := ".prism/\n"
	if len(current) > 0 && !strings.HasSuffix(string(current), "\n") {
		entry = "\n" + entry
	}
	f, err := os.OpenFile(excludeFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0o644)
	if err != nil {
		return err
	}
	_, err = f.WriteString(entry)
	if closeErr := f.Close(); err == nil {
		err = closeErr
	}
	return err
}

// isGenerated reports whether the file at path carries the prism:generated stamp,
// i.e. it is a prism-owned file safe to remove.
func isGenerated(path string) bool {
	b, err := os.ReadFile(path)
	if err != nil {
		return false
	}
	return adapters.ParseGeneratedVersion(string(b)) != ""
}

// UninstallTool removes every prism-generated file of one tool — both deliveries and
// legacy layouts. Only files carrying the prism:generated stamp are removed; a same-named
// file the user hand-wrote is left untouched. Emptied prism directories are cleaned up.
// Returns the removed project-relative paths.
func UninstallTool(projectRoot string, t adapters.Tool) ([]string, error) {
	var all []string
	for _, w := range workflows.All {
		all = append(all, t.Paths(w)...)
	}
	return removeGenerated(projectRoot, all)
}

// removeGenerated deletes the stamped files among rels and prunes directories they leave
// empty, walking up to (never including) the top-level tool directory such as .claude/.
func removeGenerated(projectRoot string, rels []string) ([]string, error) {
	var removed []string
	dirs := map[string]bool{}
	for _, rel := range rels {
		abs := filepath.Join(projectRoot, rel)
		if !isGenerated(abs) {
			continue
		}
		if err := os.Remove(abs); err != nil {
			return removed, fmt.Errorf("remove %s: %w", rel, err)
		}
		removed = append(removed, rel)
		for d := filepath.Dir(filepath.FromSlash(rel)); strings.Contains(d, string(filepath.Separator)); d = filepath.Dir(d) {
			dirs[d] = true
		}
	}
	// Deepest first; os.Remove refuses non-empty dirs, so shared dirs survive.
	ordered := make([]string, 0, len(dirs))
	for d := range dirs {
		ordered = append(ordered, d)
	}
	sort.Slice(ordered, func(i, j int) bool { return len(ordered[i]) > len(ordered[j]) })
	for _, d := range ordered {
		_ = os.Remove(filepath.Join(projectRoot, d))
	}
	return removed, nil
}

// UninstallShared removes .prism/conventions.md when it is prism-generated. It never
// touches change artifacts under .prism/ — those are the user's design work.
func UninstallShared(projectRoot string) (bool, error) {
	abs := filepath.Join(projectRoot, ConventionsPath)
	if !isGenerated(abs) {
		return false, nil
	}
	if err := os.Remove(abs); err != nil {
		return false, fmt.Errorf("remove %s: %w", ConventionsPath, err)
	}
	return true, nil
}

// IsGitExcluded reports whether .prism/ is listed in .git/info/exclude.
// Returns false outside a git repo or when the entry is absent.
func IsGitExcluded(projectRoot string) bool {
	b, err := os.ReadFile(filepath.Join(projectRoot, ".git", "info", "exclude"))
	if err != nil {
		return false
	}
	for _, line := range strings.Split(string(b), "\n") {
		trimmed := strings.TrimSpace(line)
		if trimmed == ".prism/" || trimmed == ".prism" {
			return true
		}
	}
	return false
}
