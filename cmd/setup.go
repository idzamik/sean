package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/idzamik/sean/meta"
)

var setupCmd = &cobra.Command{
	Use:   "setup",
	Short: "Install and uninstall analysis tools",
	Long: `Управляет установкой и деинсталляцией инструментов анализа.
Все инструменты устанавливаются из локальных дистрибутивов (без интернета)
на основе YAML-манифестов из папки manifests/.

Примеры:
  ` + meta.AppName + ` setup list
  ` + meta.AppName + ` setup install semgrep
  ` + meta.AppName + ` setup uninstall trivy`,
	RunE: func(cmd *cobra.Command, args []string) error {
		return cmd.Help()
	},
}

var setupInstallCmd = &cobra.Command{
	Use:   "install <tool>",
	Short: "Install a tool from its local distributive",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		toolName := args[0]
		// TODO: вызов orchestrator.Install(toolName)
		fmt.Printf("[%s] SETUP INSTALL stub: tool=%q\n", meta.AppName, toolName)
		return nil
	},
}

var setupUninstallCmd = &cobra.Command{
	Use:   "uninstall <tool>",
	Short: "Uninstall a previously installed tool",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		toolName := args[0]
		// TODO: вызов orchestrator.Uninstall(toolName)
		fmt.Printf("[%s] SETUP UNINSTALL stub: tool=%q\n", meta.AppName, toolName)
		return nil
	},
}

var setupListCmd = &cobra.Command{
	Use:   "list",
	Short: "List available tools from manifests",
	RunE: func(cmd *cobra.Command, args []string) error {
		// TODO: вызов orchestrator.ListTools()
		fmt.Printf("[%s] SETUP LIST stub: listing tools from manifests/\n", meta.AppName)
		return nil
	},
}

func init() {
	setupCmd.AddCommand(setupInstallCmd)
	setupCmd.AddCommand(setupUninstallCmd)
	setupCmd.AddCommand(setupListCmd)
	rootCmd.AddCommand(setupCmd)
}
