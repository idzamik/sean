package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/yourusername/sean/meta"
)

var (
	resFile   string
	resTool   string
	resType   string
	resSearch string
	resFormat string
)

var resCmd = &cobra.Command{
	Use:   "res",
	Short: "Browse and search analysis results",
	Long: `Просматривает и фильтрует результаты анализа по ключам:
тип анализа (SAST/SCA), инструмент, текстовый поиск.
В будущем — расширенный поиск по severity, правилам, файлам.

Примеры:
  ` + meta.AppName + ` res --file results.sarif
  ` + meta.AppName + ` res --type sast --tool semgrep
  ` + meta.AppName + ` res --search "sql injection"
  ` + meta.AppName + ` res --format table`,
	RunE: func(cmd *cobra.Command, args []string) error {
		// TODO: вызов orchestrator.ShowResults(resFile, resTool, resType, resSearch, resFormat)
		fmt.Printf("[%s] RES stub: file=%q tool=%q type=%q search=%q format=%q\n",
			meta.AppName, resFile, resTool, resType, resSearch, resFormat)
		return nil
	},
}

func init() {
	resCmd.Flags().StringVarP(&resFile, "file", "f", "", "SARIF/SBOM файл для отображения")
	resCmd.Flags().StringVar(&resTool, "tool", "", "Фильтр по инструменту")
	resCmd.Flags().StringVar(&resType, "type", "", "Тип анализа: sast, sca (по умолчанию — все)")
	resCmd.Flags().StringVarP(&resSearch, "search", "s", "", "Текстовый поиск по результатам")
	resCmd.Flags().StringVar(&resFormat, "format", "table", "Формат вывода: table, json, sarif")

	rootCmd.AddCommand(resCmd)
}
