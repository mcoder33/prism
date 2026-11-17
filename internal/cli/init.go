package cli

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"

	"github.com/mcoder33/prism/internal/adapters"
	"github.com/mcoder33/prism/internal/installer"
	"github.com/mcoder33/prism/internal/workflows"
)

func newInitCmd() *cobra.Command {
	var toolsFlag, deliveryFlag string
	cmd := &cobra.Command{
		Use:   "init [path]",
		Short: "install prism slash commands into a project",
		Args:  cobra.MaximumNArgs(1),
		RunE: func(_ *cobra.Command, args []string) error {
			root, err := projectRootArg(args)
			if err != nil {
				return err
			}
			delivery, err := adapters.ParseDelivery(deliveryFlag)
			if err != nil {
				return err
			}
			return runInit(root, toolsFlag, delivery)
		},
	}
	cmd.Flags().StringVar(&toolsFlag, "tools", "",
		fmt.Sprintf("comma-separated tool ids (%s), or \"all\"/\"none\"; omit for interactive selection",
			strings.Join(adapters.IDs(), ", ")))
	cmd.Flags().StringVar(&deliveryFlag, "delivery", "", deliveryHelp)
	return cmd
}

const deliveryHelp = "how workflows reach the agent: \"skills\" (Agent Skills folders with bundled " +
	"conventions — commit them and teammates need nothing else) or \"commands\" (slash-command " +
	"files); omit for each tool's default (commands where supported; Codex is skills-only)"

func runInit(projectRoot, toolsFlag string, delivery adapters.Delivery) error {
	var tools []adapters.Tool
	var err error
	switch {
	case toolsFlag != "":
		tools, err = parseToolsFlag(toolsFlag)
		if err != nil {
			return err
		}
	case isTTY():
		return runInitTUI(projectRoot, delivery)
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
		d := t.Resolve(delivery)
		if _, err := installer.InstallTool(projectRoot, t, d); err != nil {
			return err
		}
		fmt.Printf("  %s %s: %d %s → %s\n", green("✔"), t.Name, len(workflows.All), d, t.Location(d))
		if delivery != "" && d != delivery {
			fmt.Println(dim(fmt.Sprintf("    %s has no %s support — installed as %s", t.Name, delivery, d)))
		}
	}
	first := tools[0]
	fmt.Println()
	fmt.Printf("Try %s in %s to start a change.\n",
		cyan(first.Ref(first.Resolve(delivery), "propose")), first.Name)
	fmt.Println(dim("Restart your IDE/agent if the new commands/skills do not show up."))
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
