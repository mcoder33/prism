package cli

import (
	"fmt"

	"github.com/spf13/cobra"

	"gitlab.gidfinance.tech/zadolbator/prism/internal/adapters"
	"gitlab.gidfinance.tech/zadolbator/prism/internal/installer"
	"gitlab.gidfinance.tech/zadolbator/prism/internal/workflows"
)

func newUpdateCmd() *cobra.Command {
	var force bool
	cmd := &cobra.Command{
		Use:   "update [path]",
		Short: "refresh previously installed prism command files",
		Args:  cobra.MaximumNArgs(1),
		RunE: func(_ *cobra.Command, args []string) error {
			root, err := projectRootArg(args)
			if err != nil {
				return err
			}
			return runUpdate(root, force)
		},
	}
	cmd.Flags().BoolVar(&force, "force", false, "regenerate even if versions match")
	return cmd
}

func runUpdate(projectRoot string, force bool) error {
	configured := installer.ConfiguredTools(projectRoot)
	if len(configured) == 0 {
		fmt.Println(yellow("No prism commands found in this project. Run `prism init` first."))
		return nil
	}

	stale := configured
	if !force {
		stale = nil
		for _, t := range configured {
			if installer.InstalledVersion(projectRoot, t) != workflows.Version {
				stale = append(stale, t)
			}
		}
	}

	if len(stale) == 0 {
		fmt.Println(green(fmt.Sprintf("All tools are up to date (v%s).", workflows.Version)) +
			dim(" Use --force to regenerate anyway."))
	} else {
		if _, err := installer.InstallShared(projectRoot); err != nil {
			return err
		}
		for _, t := range stale {
			from := installer.InstalledVersion(projectRoot, t)
			if from == "" {
				from = "unknown"
			}
			if _, err := installer.InstallTool(projectRoot, t); err != nil {
				return err
			}
			fmt.Printf("  %s %s: %s → v%s\n", green("✔"), t.Name, from, workflows.Version)
		}
	}

	configuredIDs := map[string]bool{}
	for _, t := range configured {
		configuredIDs[t.ID] = true
	}
	var newTools []adapters.Tool
	for _, t := range installer.DetectTools(projectRoot) {
		if !configuredIDs[t.ID] {
			newTools = append(newTools, t)
		}
	}
	if len(newTools) > 0 {
		fmt.Println(dim(fmt.Sprintf("Detected but not configured: %s — run `prism init` to add them.",
			toolNames(newTools))))
	}
	return nil
}
