// Package installer detects AI tools in a project and writes prism command files.
package installer

import (
	"fmt"
	"os"
	"path/filepath"
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

// ConfiguredTools returns tools that already have prism command files installed.
func ConfiguredTools(projectRoot string) []adapters.Tool {
	var found []adapters.Tool
	for _, t := range adapters.All {
		for _, w := range workflows.All {
			if _, err := os.Stat(filepath.Join(projectRoot, t.CommandFile(w.ID))); err == nil {
				found = append(found, t)
				break
			}
		}
	}
	return found
}

// InstalledVersion reads the prism version stamp from a tool's generated files ("" if none).
func InstalledVersion(projectRoot string, t adapters.Tool) string {
	for _, w := range workflows.All {
		b, err := os.ReadFile(filepath.Join(projectRoot, t.CommandFile(w.ID)))
		if err == nil {
			return adapters.ParseGeneratedVersion(string(b))
		}
	}
	return ""
}

func writeFileEnsured(path, content string) error {
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return err
	}
	return os.WriteFile(path, []byte(content), 0o644)
}

// InstallShared writes .prism/conventions.md (shared by all tools) and excludes .prism/ from git.
func InstallShared(projectRoot string) ([]string, error) {
	content := adapters.GeneratedStamp(workflows.Version) + "\n\n" + workflows.Conventions() + "\n"
	if err := writeFileEnsured(filepath.Join(projectRoot, ConventionsPath), content); err != nil {
		return nil, fmt.Errorf("write conventions: %w", err)
	}
	if err := AddToGitExclude(projectRoot); err != nil {
		return nil, fmt.Errorf("update .git/info/exclude: %w", err)
	}
	return []string{ConventionsPath}, nil
}

// InstallCommand writes one command file for a tool, returning its project-relative path.
// Files are tool-owned: always overwritten.
func InstallCommand(projectRoot string, t adapters.Tool, w workflows.Workflow) (string, error) {
	body := adapters.ResolveCommandRefs(workflows.Body(w.ID), t)
	rendered := t.Format(w, body, workflows.Version)
	rel := t.CommandFile(w.ID)
	if err := writeFileEnsured(filepath.Join(projectRoot, rel), rendered); err != nil {
		return "", fmt.Errorf("write %s: %w", rel, err)
	}
	return rel, nil
}

// InstallTool writes all command files for one tool.
func InstallTool(projectRoot string, t adapters.Tool) ([]string, error) {
	files := make([]string, 0, len(workflows.All))
	for _, w := range workflows.All {
		rel, err := InstallCommand(projectRoot, t, w)
		if err != nil {
			return nil, err
		}
		files = append(files, rel)
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

// UninstallTool removes all prism-generated command files for one tool. Only files
// carrying the prism:generated stamp are removed — a same-named file the user
// hand-wrote is left untouched. Now-empty prism command directories are cleaned up.
// Returns the removed project-relative paths.
func UninstallTool(projectRoot string, t adapters.Tool) ([]string, error) {
	var removed []string
	dirs := map[string]bool{}
	for _, w := range workflows.All {
		rel := t.CommandFile(w.ID)
		abs := filepath.Join(projectRoot, rel)
		if !isGenerated(abs) {
			continue
		}
		if err := os.Remove(abs); err != nil {
			return removed, fmt.Errorf("remove %s: %w", rel, err)
		}
		removed = append(removed, rel)
		dirs[filepath.Dir(abs)] = true
	}
	// Drop now-empty command dirs (e.g. .claude/commands/prism/). os.Remove refuses
	// non-empty dirs, so a directory shared with other files is left intact.
	for d := range dirs {
		_ = os.Remove(d)
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
