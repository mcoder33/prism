package cli

import (
	"errors"
	"fmt"

	"github.com/spf13/cobra"

	"github.com/mcoder33/prism/internal/change"
)

// errSilent signals a failing exit whose output was already printed.
var errSilent = errors.New("")

func newCheckCmd() *cobra.Command {
	var changeFlag string
	var asJSON bool
	cmd := &cobra.Command{
		Use:   "check [path]",
		Short: "validate change artifacts: drawio XML, spec scenarios, state table vs files",
		Long: "Validate a change (--change, else .prism/CURRENT, else every active change):\n" +
			"well-formed .drawio XML (no xmllint needed), WHEN/THEN steps in every spec Scenario,\n" +
			"tagged Open items, atomicity bounds, and the README table against the node files.\n" +
			"Exits 1 when any error-level finding exists. Read-only.",
		Args: cobra.MaximumNArgs(1),
		RunE: func(_ *cobra.Command, args []string) error {
			root, err := projectRootArg(args)
			if err != nil {
				return err
			}
			return runCheck(root, changeFlag, asJSON)
		},
	}
	cmd.Flags().StringVar(&changeFlag, "change", "",
		"change under .prism/, e.g. my-change or archive/my-change (default: CURRENT, else all active)")
	cmd.Flags().BoolVar(&asJSON, "json", false, "emit machine-readable JSON")
	return cmd
}

func runCheck(projectRoot, changeFlag string, asJSON bool) error {
	var names []string
	switch {
	case changeFlag != "":
		names = []string{changeFlag}
	case change.Current(projectRoot) != "":
		names = []string{change.Current(projectRoot)}
	default:
		names = change.Active(projectRoot)
	}
	if len(names) == 0 {
		if asJSON {
			return writeJSON(map[string]any{"changes": []any{}, "errors": 0})
		}
		fmt.Println(dim("No active changes to check."))
		return nil
	}

	var results []*change.Change
	errorsTotal := 0
	for _, name := range names {
		c, err := change.Load(projectRoot, name)
		if err != nil {
			return emitFailure(asJSON, err)
		}
		c.Check(projectRoot)
		errorsTotal += c.Errors()
		results = append(results, c)
	}

	if asJSON {
		if err := writeJSON(map[string]any{"changes": results, "errors": errorsTotal}); err != nil {
			return err
		}
	} else {
		for _, c := range results {
			fmt.Println(bold(c.Name))
			printFindings(c.Findings)
		}
	}
	if errorsTotal > 0 {
		if asJSON {
			return errSilent
		}
		return fmt.Errorf("%d error(s) found", errorsTotal)
	}
	return nil
}
