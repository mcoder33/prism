package cli

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/charmbracelet/huh"
	"github.com/spf13/cobra"

	"github.com/mcoder33/prism/internal/adapters"
	"github.com/mcoder33/prism/internal/installer"
	"github.com/mcoder33/prism/internal/workflows"
)

func newInitCmd() *cobra.Command {
	var toolsFlag string
	cmd := &cobra.Command{
		Use:   "init [path]",
		Short: "install prism slash commands into a project",
		Args:  cobra.MaximumNArgs(1),
		RunE: func(_ *cobra.Command, args []string) error {
			root, err := projectRootArg(args)
			if err != nil {
				return err
			}
			return runInit(root, toolsFlag)
		},
	}
	cmd.Flags().StringVar(&toolsFlag, "tools", "",
		fmt.Sprintf("comma-separated tool ids (%s), or \"all\"/\"none\"; omit for interactive selection",
			strings.Join(adapters.IDs(), ", ")))
	return cmd
}

func runInit(projectRoot, toolsFlag string) error {
	var tools []adapters.Tool
	var err error
	switch {
	case toolsFlag != "":
		tools, err = parseToolsFlag(toolsFlag)
		if err != nil {
			return err
		}
	case isTTY():
		tools, err = selectToolsInteractive(projectRoot)
		if err != nil {
			return err
		}
	default:
		tools = installer.DetectTools(projectRoot)
		if len(tools) == 0 {
			return fmt.Errorf("no AI tools detected and not running interactively; pass --tools <list|all>")
		}
	}

	if len(tools) == 0 {
		fmt.Println(yellow("No tools selected — nothing to do."))
		return nil
	}

	sharedFiles, err := installer.InstallShared(projectRoot)
	if err != nil {
		return err
	}

	fmt.Println()
	fmt.Println(bold(fmt.Sprintf("prism v%s installed into %s", workflows.Version, projectRoot)))
	fmt.Println(dim(fmt.Sprintf("  shared: %s (.prism/ is git-excluded)", strings.Join(sharedFiles, ", "))))
	for _, t := range tools {
		files, err := installer.InstallTool(projectRoot, t)
		if err != nil {
			return err
		}
		fmt.Printf("  %s %s: %d commands → %s/\n", green("✔"), t.Name, len(files), filepath.Dir(files[0]))
	}
	fmt.Println()
	fmt.Printf("Try %s in %s to start a change.\n", cyan(tools[0].CommandRef("propose")), tools[0].Name)
	fmt.Println(dim("Restart your IDE/agent if slash commands do not show up."))
	return nil
}

func parseToolsFlag(flag string) ([]adapters.Tool, error) {
	switch flag {
	case "all":
		return adapters.All, nil
	case "none":
		return nil, nil
	}
	var tools []adapters.Tool
	for _, raw := range strings.Split(flag, ",") {
		id := strings.ToLower(strings.TrimSpace(raw))
		t, ok := adapters.ByID(id)
		if !ok {
			return nil, fmt.Errorf("unknown tool %q; known tools: %s, or \"all\"/\"none\"",
				id, strings.Join(adapters.IDs(), ", "))
		}
		tools = append(tools, t)
	}
	return tools, nil
}

func selectToolsInteractive(projectRoot string) ([]adapters.Tool, error) {
	detected := map[string]bool{}
	for _, t := range installer.DetectTools(projectRoot) {
		detected[t.ID] = true
	}
	configured := map[string]bool{}
	for _, t := range installer.ConfiguredTools(projectRoot) {
		configured[t.ID] = true
	}

	options := make([]huh.Option[string], 0, len(adapters.All))
	for _, t := range adapters.All {
		label := t.Name
		switch {
		case configured[t.ID]:
			label += " (installed)"
		case detected[t.ID]:
			label += " (detected)"
		}
		options = append(options, huh.NewOption(label, t.ID).Selected(configured[t.ID] || detected[t.ID]))
	}

	var picked []string
	prompt := huh.NewMultiSelect[string]().
		Title("Which AI tools should get the prism commands?").
		Options(options...).
		Value(&picked)
	if err := huh.NewForm(huh.NewGroup(prompt)).Run(); err != nil {
		return nil, err
	}

	tools := make([]adapters.Tool, 0, len(picked))
	for _, id := range picked {
		t, _ := adapters.ByID(id)
		tools = append(tools, t)
	}
	return tools, nil
}
