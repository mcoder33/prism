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
	defer f.Close()
	_, err = f.WriteString(entry)
	return err
}
