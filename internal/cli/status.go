package cli

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/spf13/cobra"

	"github.com/mcoder33/prism/internal/change"
)

func newStatusCmd() *cobra.Command {
	var changeFlag string
	var asJSON bool
	cmd := &cobra.Command{
		Use:   "status [path]",
		Short: "report a change's state table cross-checked against the files (read-only)",
		Long: "Report the active change (or --change): phase, tier, branch, every node's status,\n" +
			"task and Open counts, and deterministic discrepancies between the README table and\n" +
			"the node files. Read-only. --json emits one JSON document for agents.",
		Args: cobra.MaximumNArgs(1),
		RunE: func(_ *cobra.Command, args []string) error {
			root, err := projectRootArg(args)
			if err != nil {
				return err
			}
			return runStatus(root, changeFlag, asJSON)
		},
	}
	cmd.Flags().StringVar(&changeFlag, "change", "", "change slug under .prism/ (default: .prism/CURRENT)")
	cmd.Flags().BoolVar(&asJSON, "json", false, "emit machine-readable JSON")
	return cmd
}

func runStatus(projectRoot, changeFlag string, asJSON bool) error {
	name, err := change.Resolve(projectRoot, changeFlag)
	if err != nil {
		return emitFailure(asJSON, err)
	}
	c, err := change.Load(projectRoot, name)
	if err != nil {
		return emitFailure(asJSON, err)
	}
	if asJSON {
		return writeJSON(c)
	}
	printChange(c)
	return nil
}

func printChange(c *change.Change) {
	marker := ""
	if c.Current {
		marker = green(" (current)")
	}
	fmt.Println(bold(c.Name) + marker)
	meta := []string{"phase: " + orDash(c.Phase), "tier: " + c.Tier}
	if c.Branch != "" {
		meta = append(meta, "branch: "+c.Branch)
	}
	fmt.Println(dim("  " + strings.Join(meta, " · ")))

	if len(c.Nodes) == 0 {
		fmt.Println(dim("  no nodes yet"))
	} else {
		w := len("NODE")
		for _, n := range c.Nodes {
			w = max(w, len(n.ID))
		}
		fmt.Printf("  %-*s  %-6s  %-6s  %s\n", w, "NODE", "STATUS", "TASKS", "OPEN")
		for _, n := range c.Nodes {
			tasks := "-"
			if n.Tasks.Total > 0 {
				tasks = fmt.Sprintf("%d/%d", n.Tasks.Done, n.Tasks.Total)
			}
			fmt.Printf("  %-*s  %s  %-6s  %s\n", w, n.ID, padCells(orDash(n.Status), 6), tasks, openSummary(n.Open))
		}
	}
	if c.RootTasks != nil {
		fmt.Println(dim(fmt.Sprintf("  root tasks: %d/%d", c.RootTasks.Done, c.RootTasks.Total)))
	}
	printFindings(c.Findings)
}

func printFindings(fs []change.Finding) {
	if len(fs) == 0 {
		fmt.Printf("  %s no discrepancies\n", green("✔"))
		return
	}
	for _, f := range fs {
		icon := yellow("⚠")
		if f.Severity == change.Error {
			icon = red("✖")
		}
		loc := ""
		if f.Path != "" {
			loc = dim(" (" + f.Path + ")")
		}
		fmt.Printf("  %s %s%s %s\n", icon, f.Message, loc, dim("["+f.Code+"]"))
	}
}

func openSummary(o change.OpenCounts) string {
	var parts []string
	add := func(n int, label string) {
		if n > 0 {
			parts = append(parts, fmt.Sprintf("%d %s", n, label))
		}
	}
	add(o.Blocking, "blocking")
	add(o.Minor, "minor")
	add(o.Deferred, "deferred")
	add(o.Untagged, "untagged")
	if len(parts) == 0 {
		return "none"
	}
	return strings.Join(parts, ", ")
}

// padCells pads by display width: status glyphs are emoji, two terminal cells wide.
func padCells(s string, width int) string {
	return s + strings.Repeat(" ", max(width-lipgloss.Width(s), 0))
}

func orDash(s string) string {
	if s == "" {
		return "-"
	}
	return s
}

func writeJSON(v any) error {
	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", "  ")
	enc.SetEscapeHTML(false)
	return enc.Encode(v)
}

// emitFailure keeps the --json contract (one JSON document on stdout) on errors too.
func emitFailure(asJSON bool, err error) error {
	if !asJSON {
		return err
	}
	_ = writeJSON(map[string]any{"error": err.Error()})
	return errSilent
}
