
package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/idzamik/sean/orchestrator"
)


func registerToolCommands(root *cobra.Command) {
	state, err := LoadState()
	if err != nil {
		fmt.Fprintf(os.Stderr, "warn: cannot load installed.yaml: %v\n", err)
		return
	}

	requestedTool := ""
	if len(os.Args) > 1 {
		requestedTool = os.Args[1]
	}

	for _, entry := range state.Tools {
		entry := entry // захват переменной для замыкания

		if requestedTool != "" && entry.Name != requestedTool {
			registerStubCommand(root, entry.Name)
			continue
		}

		manifest, err := LoadManifest(entry.Manifest)
		if err != nil {
			fmt.Fprintf(os.Stderr, "warn: cannot load manifest for %q: %v\n", entry.Name, err)
			registerStubCommand(root, entry.Name)
			continue
		}

		registerToolCommand(root, manifest)
	}
}


func registerStubCommand(root *cobra.Command, toolName string) {
	cmd := &cobra.Command{
		Use:   toolName,
		Short: fmt.Sprintf("Run %s analysis", toolName),
		RunE: func(cmd *cobra.Command, args []string) error {
			return fmt.Errorf("tool %q manifest could not be loaded", toolName)
		},
	}
	root.AddCommand(cmd)
}


func registerToolCommand(root *cobra.Command, manifest *Manifest) {
	toolCmd := &cobra.Command{
		Use:   manifest.Name + " [target] [flags...]",
		Short: fmt.Sprintf("Run %s analysis", manifest.Name),
		Long: fmt.Sprintf(
			"Run %s using manifest configs/manifests/%s.yaml.\n\n"+
				"Default command : %s\n"+
				"Results saved to: %s/<timestamp>.%s\n\n"+
				"Extra flags are passed directly to the tool and override manifest defaults.",
			manifest.Name,
			manifest.Name,
			manifest.DefaultCommand,
			manifest.Output.Dir,
			manifest.Output.Format,
		),
		DisableFlagParsing: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			return runToolCommand(manifest, args)
		},
	}

	root.AddCommand(toolCmd)
}


func runToolCommand(manifest *Manifest, args []string) error {
	if len(args) == 0 {
		return fmt.Errorf(
			"usage: sean %s <target> [flags...]\n\nExample: sean %s ./src",
			manifest.Name, manifest.Name,
		)
	}

	target, userFlags := parseTargetAndFlags(args)

	if target == "" {
		return fmt.Errorf("target path is required: sean %s <target> [flags...]", manifest.Name)
	}

	adapter := NewManifestAdapter(manifest)
	return orchestrator.RunAnalyser(adapter, target, userFlags)
}


func parseTargetAndFlags(args []string) (target string, flags []string) {
	for _, arg := range args {
		if len(arg) > 0 && arg[0] == '-' {
			flags = append(flags, arg)
		} else if target == "" {
			target = arg
		} else {
			flags = append(flags, arg)
		}
	}
	return
}


