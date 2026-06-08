package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/yourusername/sean/meta"
)

var (
	sastTarget string
	sastTool   string
	sastOutput string
	sastRules  string
)

var sastCmd = &cobra.Command{
	Use:   "sast",
	Short: "Run Static Application Security Testing (SAST)",
	Long: `Запускает статический анализ исходного кода (SAST) над указанной целью.
Оркестратор выбирает подходящий инструмент и сохраняет результат в SARIF-файл.

Примеры:
  ` + meta.AppName + ` sast --target ./src
  ` + meta.AppName + ` sast --target ./src --tool semgrep --rules p/security-audit
  ` + meta.AppName + ` sast --target ./src --output results.sarif`,
	RunE: func(cmd *cobra.Command, args []string) error {
		// TODO: вызов orchestrator.RunSAST(sastTarget, sastTool, sastOutput, sastRules)
		fmt.Printf("[%s] SAST stub: target=%q tool=%q output=%q rules=%q\n",
			meta.AppName, sastTarget, sastTool, sastOutput, sastRules)
		return nil
	},
}

func init() {
	sastCmd.Flags().StringVarP(&sastTarget, "target", "t", ".", "Путь к анализируемому исходному коду")
	sastCmd.Flags().StringVar(&sastTool, "tool", "", "Инструмент (по умолчанию — из конфига)")
	sastCmd.Flags().StringVarP(&sastOutput, "output", "o", "results.sarif", "Путь для сохранения SARIF-файла")
	sastCmd.Flags().StringVar(&sastRules, "rules", "", "Набор правил анализа (специфично для инструмента)")

	rootCmd.AddCommand(sastCmd)
}
