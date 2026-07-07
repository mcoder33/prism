package cli

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strings"

	"github.com/spf13/cobra"

	"github.com/mcoder33/prism/internal/installer"
	"github.com/mcoder33/prism/internal/workflows"
)

func newDoctorCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "doctor [path]",
		Short: "diagnose a prism installation: version drift, stale pointer, prerequisites",
		Args:  cobra.MaximumNArgs(1),
		RunE: func(_ *cobra.Command, args []string) error {
			root, err := projectRootArg(args)
			if err != nil {
				return err
			}
			return runDoctor(root)
		},
	}
}

// runDoctor is read-only: it reports drift and prerequisites, changes nothing.
func runDoctor(projectRoot string) error {
	fmt.Println(bold(fmt.Sprintf("prism doctor — CLI v%s", workflows.Version)))

	// Installed command files, per tool: version drift against the CLI.
	configured := installer.ConfiguredTools(projectRoot)
	if len(configured) == 0 {
		fmt.Println(yellow("  ⚠ no prism commands installed here — run `prism init`"))
	} else {
		for _, t := range configured {
			v := installer.InstalledVersion(projectRoot, t)
			if v == workflows.Version {
				fmt.Printf("  %s %s: v%s\n", green("✔"), t.Name, v)
				continue
			}
			if v == "" {
				v = "unknown"
			}
			fmt.Printf("  %s %s: v%s (CLI is v%s) — run `prism update`\n",
				yellow("⚠"), t.Name, v, workflows.Version)
		}
	}

	// .prism/ git-exclude — only meaningful inside a git repo.
	if st, err := os.Stat(filepath.Join(projectRoot, ".git")); err == nil && st.IsDir() {
		if installer.IsGitExcluded(projectRoot) {
			fmt.Printf("  %s .prism/ is git-excluded\n", green("✔"))
		} else {
			fmt.Printf("  %s .prism/ not in .git/info/exclude — run `prism init` or `prism update`\n", yellow("⚠"))
		}
	}

	prismDir := filepath.Join(projectRoot, installer.PrismDir)

	// CURRENT pointer vs reality.
	if b, err := os.ReadFile(filepath.Join(prismDir, "CURRENT")); err == nil {
		if cur := strings.TrimSpace(string(b)); cur != "" {
			if st, err := os.Stat(filepath.Join(prismDir, cur)); err == nil && st.IsDir() {
				fmt.Printf("  %s active change: %s\n", green("✔"), cur)
			} else {
				fmt.Printf("  %s CURRENT points at %q with no .prism/%s/ — repair with the use command\n",
					yellow("⚠"), cur, cur)
			}
		}
	}

	// Active changes inventory.
	if entries, err := os.ReadDir(prismDir); err == nil {
		var changes []string
		for _, e := range entries {
			if e.IsDir() && e.Name() != "archive" {
				changes = append(changes, e.Name())
			}
		}
		sort.Strings(changes)
		if len(changes) > 0 {
			fmt.Printf("  %s %d active change(s): %s\n", dim("•"), len(changes), strings.Join(changes, ", "))
		}
	}

	// xmllint: drawio validation dependency used by propose/drill/integrate.
	if _, err := exec.LookPath("xmllint"); err == nil {
		fmt.Printf("  %s xmllint present (drawio validation)\n", green("✔"))
	} else {
		fmt.Printf("  %s xmllint not found — drawio validation falls back to a manual parse check\n", dim("•"))
	}

	return nil
}
