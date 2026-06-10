package cli

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/spf13/cobra"

	"gitlab.gidfinance.tech/zadolbator/prism/internal/installer"
)

func newListCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "list [path]",
		Short: "list active prism changes in a project",
		Args:  cobra.MaximumNArgs(1),
		RunE: func(_ *cobra.Command, args []string) error {
			root, err := projectRootArg(args)
			if err != nil {
				return err
			}
			return runList(root)
		},
	}
}

func runList(projectRoot string) error {
	prismDir := filepath.Join(projectRoot, installer.PrismDir)
	entries, err := os.ReadDir(prismDir)
	if err != nil {
		fmt.Println(yellow("No .prism/ directory here. Start a change with the propose command in your agent."))
		return nil
	}

	current := ""
	if b, err := os.ReadFile(filepath.Join(prismDir, "CURRENT")); err == nil {
		current = strings.TrimSpace(string(b))
	}

	var changes []string
	for _, e := range entries {
		if e.IsDir() && e.Name() != "archive" {
			changes = append(changes, e.Name())
		}
	}
	sort.Strings(changes)

	if len(changes) == 0 {
		fmt.Println(dim("No active changes."))
	} else {
		fmt.Println(bold("Active changes:"))
		for _, name := range changes {
			marker := ""
			if name == current {
				marker = green(" (current)")
			}
			fmt.Printf("  %s%s\n", name, marker)
		}
	}

	if archived, err := os.ReadDir(filepath.Join(prismDir, "archive")); err == nil {
		n := 0
		for _, e := range archived {
			if e.IsDir() {
				n++
			}
		}
		if n > 0 {
			fmt.Println(dim(fmt.Sprintf("Archived: %d (.prism/archive/)", n)))
		}
	}
	return nil
}
