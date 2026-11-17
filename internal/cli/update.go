package cli

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/mcoder33/prism/internal/adapters"
	"github.com/mcoder33/prism/internal/installer"
	"github.com/mcoder33/prism/internal/workflows"
)

func newUpdateCmd() *cobra.Command {
	var force bool
	var deliveryFlag string
	cmd := &cobra.Command{
		Use:   "update [path]",
		Short: "refresh previously installed prism command files",
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
			return runUpdate(root, force, delivery)
		},
	}
	cmd.Flags().BoolVar(&force, "force", false, "regenerate even if versions match")
	cmd.Flags().StringVar(&deliveryFlag, "delivery", "",
		"switch installed tools to \"skills\" or \"commands\" (old files are removed); omit to keep each tool's current delivery")
	return cmd
}

func runUpdate(projectRoot string, force bool, delivery adapters.Delivery) error {
	configured := installer.ConfiguredTools(projectRoot)
	if len(configured) == 0 {
		fmt.Println(yellow("No prism commands or skills found in this project. Run `prism init` first."))
		return nil
	}

	// target delivery per tool: an explicit --delivery wins, else keep what is installed.
	target := func(t adapters.Tool) adapters.Delivery {
		if delivery != "" {
			return t.Resolve(delivery)
		}
		return installer.InstalledDelivery(projectRoot, t)
	}

	var stale []adapters.Tool
	for _, t := range configured {
		d := target(t)
		if force || installer.InstalledVersion(projectRoot, t) != workflows.Version ||
			installer.InstalledDelivery(projectRoot, t) != d || installer.NeedsMigration(projectRoot, t, d) {
			stale = append(stale, t)
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
			} else {
				from = "v" + from
			}
			fromD := installer.InstalledDelivery(projectRoot, t)
			d := target(t)
			if _, err := installer.InstallTool(projectRoot, t, d); err != nil {
				return err
			}
			change := ""
			if fromD != d {
				change = fmt.Sprintf(", %s → %s", fromD, d)
			}
			fmt.Printf("  %s %s: %s → v%s (%s%s)\n", green("✔"), t.Name, from, workflows.Version, t.Location(d), change)
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
