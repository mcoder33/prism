package cli

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/mcoder33/prism/internal/adapters"
	"github.com/mcoder33/prism/internal/installer"
)

func newUninstallCmd() *cobra.Command {
	var toolsFlag string
	var shared bool
	cmd := &cobra.Command{
		Use:   "uninstall [path]",
		Short: "remove prism command files (design artifacts under .prism/ are kept)",
		Args:  cobra.MaximumNArgs(1),
		RunE: func(_ *cobra.Command, args []string) error {
			root, err := projectRootArg(args)
			if err != nil {
				return err
			}
			return runUninstall(root, toolsFlag, shared)
		},
	}
	cmd.Flags().StringVar(&toolsFlag, "tools", "",
		"comma-separated tool ids to remove, or \"all\"; omit to remove every installed tool")
	cmd.Flags().BoolVar(&shared, "shared", false,
		"also remove .prism/conventions.md (kept by default so active changes stay usable)")
	return cmd
}

func runUninstall(projectRoot, toolsFlag string, shared bool) error {
	var tools []adapters.Tool
	if toolsFlag != "" {
		parsed, err := parseToolsFlag(toolsFlag)
		if err != nil {
			return err
		}
		tools = parsed
	} else {
		tools = installer.ConfiguredTools(projectRoot)
	}
	if len(tools) == 0 {
		fmt.Println(yellow("No prism commands found to remove."))
		return nil
	}

	total := 0
	for _, t := range tools {
		removed, err := installer.UninstallTool(projectRoot, t)
		if err != nil {
			return err
		}
		if len(removed) > 0 {
			fmt.Printf("  %s %s: removed %d command file(s)\n", green("✔"), t.Name, len(removed))
			total += len(removed)
		}
	}
	if total == 0 {
		fmt.Println(yellow("Nothing removed — no prism-generated files matched."))
	}

	if shared {
		ok, err := installer.UninstallShared(projectRoot)
		if err != nil {
			return err
		}
		if ok {
			fmt.Printf("  %s removed %s\n", green("✔"), installer.ConventionsPath)
		}
	} else {
		fmt.Println(dim("Kept .prism/ (conventions + change artifacts). Pass --shared to also remove conventions.md."))
	}
	return nil
}
