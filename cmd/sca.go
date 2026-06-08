package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/idzamik/sean/meta"
)

var (
	scaTarget string
	scaTool   string
	scaOutput string
)

var scaCmd = &cobra.Command{
	Use:   "sca",
	Short: "Run Software Composition Analysis (SCA)",
	Long: `Запускает анализ состава ПО (SCA) над указанной целью.
Оркестратор выбирает подходящий инструмент и сохраняет результат в SBOM-файл.

Примеры:
  ` + meta.AppName + ` sca --target ./myproject
  ` + meta.AppName + ` sca --target ./myproject --tool trivy --output sbom.json`,
	RunE: func(cmd *cobra.Command, args []string) error {
		// TODO: вызов orchestrator.RunSCA(scaTarget, scaTool, scaOutput)
		fmt.Printf("[%s] SCA stub: target=%q tool=%q output=%q\n",
			meta.AppName, scaTarget, scaTool, scaOutput)
		return nil
	},
}

func init() {
	scaCmd.Flags().StringVarP(&scaTarget, "target", "t", ".", "Путь к анализируемому проекту или образу")
	scaCmd.Flags().StringVar(&scaTool, "tool", "", "Инструмент (по умолчанию — из конфига)")
	scaCmd.Flags().StringVarP(&scaOutput, "output", "o", "sbom.json", "Путь для сохранения SBOM-файла")

	rootCmd.AddCommand(scaCmd)
}
