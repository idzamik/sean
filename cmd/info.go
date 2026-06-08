package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/idzamik/sean/meta"
)

var (
	infoTool    string
	infoVerbose bool
)

var infoCmd = &cobra.Command{
	Use:   "info",
	Short: "Show system and installed tools information",
	Long: `Выводит информацию о системе и установленных инструментах анализа:
версии, статусы, текущие активные анализы.

Примеры:
  ` + meta.AppName + ` info
  ` + meta.AppName + ` info --tool semgrep
  ` + meta.AppName + ` info --verbose`,
	RunE: func(cmd *cobra.Command, args []string) error {
		// TODO: вызов orchestrator.Info(infoTool, infoVerbose)
		fmt.Printf("[%s] INFO stub: tool=%q verbose=%v\n",
			meta.AppName, infoTool, infoVerbose)
		return nil
	},
}

func init() {
	infoCmd.Flags().StringVar(&infoTool, "tool", "", "Детальная информация о конкретном инструменте")
	infoCmd.Flags().BoolVarP(&infoVerbose, "verbose", "v", false, "Подробный вывод")

	rootCmd.AddCommand(infoCmd)
}
