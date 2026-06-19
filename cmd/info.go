package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/idzamik/sean/meta"
	"github.com/idzamik/sean/orchestrator"
)

var (
	infoTool    string
	infoVerbose bool
)

var infoCmd = &cobra.Command{
	Use:   "info",
	Short: "Show installed tools and system state",
	Long: `Выводит информацию об установленных инструментах анализа:
бинари, доступные команды, директории результатов.

Примеры:
  ` + meta.AppName + ` info
  ` + meta.AppName + ` info --tool trivy
  ` + meta.AppName + ` info --verbose`,
	RunE: func(cmd *cobra.Command, args []string) error {
    state, err := LoadState()
    if err != nil {
        return fmt.Errorf("cannot load state: %w", err)
    }

    var manifests []*Manifest
    for _, entry := range state.Tools {
        if infoTool != "" && entry.Name != infoTool {
            continue
        }
        m, err := LoadManifest(entry.Manifest)
        if err != nil {
            fmt.Fprintf(os.Stderr, "warn: cannot load manifest for %q: %v\n", entry.Name, err)
            continue
        }
        manifests = append(manifests, m)
    }

    infoManifests := make([]orchestrator.InfoManifest, len(manifests))
		for i, m := range manifests {
		    infoManifests[i] = NewManifestAdapter(m)
		}
		return orchestrator.Info(infoManifests, infoVerbose)
	},
}

func init() {
	infoCmd.Flags().StringVar(&infoTool, "tool", "", "Детальная информация о конкретном инструменте")
	infoCmd.Flags().BoolVarP(&infoVerbose, "verbose", "v", false, "Подробный вывод")

	rootCmd.AddCommand(infoCmd)
}
