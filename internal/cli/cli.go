// Package cli wires the prism commands: init, update, list.
package cli

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"

	"github.com/mcoder33/prism/internal/adapters"
	"github.com/mcoder33/prism/internal/workflows"
)

func Execute() {
	root := &cobra.Command{
		Use:   "prism",
		Short: "PRISM — recursive decomposition workflow for AI coding agents",
		Long: "PRISM — recursive decomposition workflow for AI coding agents.\n" +
			"Installs /prism slash commands into a project for the agents you use.",
		Version:       workflows.Version,
		SilenceUsage:  true,
		SilenceErrors: true,
	}

	root.AddCommand(newInitCmd(), newUpdateCmd(), newListCmd())

	if err := root.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, red("error: "+err.Error()))
		os.Exit(1)
	}
}

func projectRootArg(args []string) (string, error) {
	target := "."
	if len(args) > 0 {
		target = args[0]
	}
	abs, err := os.Getwd()
	if err != nil {
		return "", err
	}
	if target != "." {
		if st, err := os.Stat(target); err != nil || !st.IsDir() {
			return "", fmt.Errorf("not a directory: %s", target)
		}
		return target, nil
	}
	return abs, nil
}

func isTTY() bool {
	st, err := os.Stdin.Stat()
	if err != nil {
		return false
	}
	return st.Mode()&os.ModeCharDevice != 0
}

func toolNames(tools []adapters.Tool) string {
	names := make([]string, len(tools))
	for i, t := range tools {
		names[i] = t.Name
	}
	return strings.Join(names, ", ")
}

// Minimal ANSI styling, disabled when stdout is not a terminal.
func styled(code, s string) string {
	if st, err := os.Stdout.Stat(); err != nil || st.Mode()&os.ModeCharDevice == 0 {
		return s
	}
	return "\x1b[" + code + "m" + s + "\x1b[0m"
}

func bold(s string) string   { return styled("1", s) }
func dim(s string) string    { return styled("2", s) }
func red(s string) string    { return styled("31", s) }
func green(s string) string  { return styled("32", s) }
func yellow(s string) string { return styled("33", s) }
func cyan(s string) string   { return styled("36", s) }
