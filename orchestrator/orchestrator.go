
package orchestrator

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"sort"
	"strings"

	"github.com/idzamik/sean/meta"
)


type ToolManifest interface {
	GetName() string
	GetDefaultCommand() string
	GetBin() string
	GetRules() string
	GetCommand(name string) (ManifestCommand, bool)
	GetOutput() ManifestOutput
}


type InfoManifest interface {
	GetName() string
	GetBin() string
	GetDefaultCommand() string
	GetRules() string
	GetCommandNames() []string
	GetOutput() ManifestOutput
}


type InfoEntry struct {
	Name           string
	Bin            string
	DefaultCommand string
	Rules          string
	Commands       []string
	OutputDir      string
	OutputFormat   string
}


func RunAnalyser(manifest ToolManifest, target string, userFlags []string) error {
	logger := log.New(os.Stderr, fmt.Sprintf("[%s] ", manifest.GetName()), 0)

	cmdName := manifest.GetDefaultCommand()
	manifestCmd, ok := manifest.GetCommand(cmdName)
	if !ok {
		return fmt.Errorf("command %q not found in manifest for tool %q", cmdName, manifest.GetName())
	}

	out := manifest.GetOutput()
	outputPath, err := GenerateOutputPath(manifest.GetName(), out.Dir, out.Format)
	if err != nil {
		return fmt.Errorf("cannot prepare output path: %w", err)
	}

	bin, argv := BuildCommand(manifestCmd, target, manifest.GetRules(), outputPath, userFlags)

	logger.Printf("running: %s %v", bin, argv)
	logger.Printf("output : %s", outputPath)

	cmd := exec.Command(bin, argv...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("tool %q exited with error: %w", manifest.GetName(), err)
	}

	if _, err := os.Stat(outputPath); os.IsNotExist(err) {
		logger.Printf("warn: output file not found at %s (tool may use a different output path)", outputPath)
	} else {
		logger.Printf("done: results saved to %s", outputPath)
	}

	return nil
}


func Info(manifests []InfoManifest, verbose bool) error {
	fmt.Printf("%s (%s) v%s\n", meta.AppName, meta.AppFullName, meta.AppVersion)
	fmt.Printf("Installed tools: %d\n\n", len(manifests))

	entries := make([]InfoEntry, 0, len(manifests))
	for _, m := range manifests {
		names := m.GetCommandNames()
		sort.Strings(names)
		entries = append(entries, InfoEntry{
			Name:           m.GetName(),
			Bin:            m.GetBin(),
			DefaultCommand: m.GetDefaultCommand(),
			Rules:          m.GetRules(),
			Commands:       names,
			OutputDir:      m.GetOutput().Dir,
			OutputFormat:   m.GetOutput().Format,
		})
	}

	for _, e := range entries {
		fmt.Printf("  %-14s bin: %-40s commands: %s\n",
			e.Name,
			e.Bin,
			strings.Join(e.Commands, ", "),
		)

		if verbose {
			pad := strings.Repeat(" ", 16)
			fmt.Printf("%s default command : %s\n", pad, e.DefaultCommand)
			if e.Rules != "" {
				fmt.Printf("%s rules           : %s\n", pad, e.Rules)
			}
			fmt.Printf("%s output dir       : %s\n", pad, e.OutputDir)
			fmt.Printf("%s output format    : %s\n", pad, e.OutputFormat)
			fmt.Println()
		}
	}

	if !verbose {
		fmt.Println("\nResults directories:")
		for _, e := range entries {
			fmt.Printf("  %s\n", e.OutputDir)
		}
	}

	return nil
}


func ShowResults(file, tool, analysisType, search, format string) error {
	fmt.Printf("[%s/orchestrator] ShowResults file=%q tool=%q type=%q search=%q format=%q\n", meta.AppName, file, tool, analysisType, search, format)
	// TODO: открыть SARIF/SBOM → применить фильтры → отформатировать вывод
	return nil
}
