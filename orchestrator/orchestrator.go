// Package orchestrator содержит всю бизнес-логику CLI.
// Команды из cmd/ только разбирают флаги и вызывают функции этого пакета.
// Оркестратор решает, какой модуль из analyzers/ запустить,
// исходя из параметров, манифестов и state/installed.yaml.
package orchestrator

import (
	"fmt"
	"log"
	"os"
	"os/exec"

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


func RunAnalyser(manifest ToolManifest, target string, userFlags []string) error {
	logger := log.New(os.Stderr, fmt.Sprintf("[%s] ", manifest.GetName()), 0)

	// 1. Определяем команду
	cmdName := manifest.GetDefaultCommand()
	manifestCmd, ok := manifest.GetCommand(cmdName)
	if !ok {
		return fmt.Errorf("command %q not found in manifest for tool %q", cmdName, manifest.GetName())
	}

	// 2. Генерируем путь к файлу результата
	out := manifest.GetOutput()
	outputPath, err := GenerateOutputPath(manifest.GetName(), out.Dir, out.Format)
	if err != nil {
		return fmt.Errorf("cannot prepare output path: %w", err)
	}

	// 3. Собираем аргументы
	bin, argv := BuildCommand(manifestCmd, target, manifest.GetRules(), outputPath, userFlags)

	logger.Printf("running: %s %v", bin, argv)
	logger.Printf("output : %s", outputPath)

	// 4. Запускаем инструмент
	cmd := exec.Command(bin, argv...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("tool %q exited with error: %w", manifest.GetName(), err)
	}

	// 5. Проверяем что файл результата создан
	if _, err := os.Stat(outputPath); os.IsNotExist(err) {
		logger.Printf("warn: output file not found at %s (tool may use a different output path)", outputPath)
	} else {
		logger.Printf("done: results saved to %s", outputPath)
	}

	return nil
}


// +
func Info(manifest ToolManifest, verbose bool) error {
	fmt.Printf("Tool   : %s\n", manifest.GetName())
	fmt.Printf("Bin    : %s\n", manifest.GetBin())
	fmt.Printf("Default: %s\n", manifest.GetDefaultCommand())
	fmt.Printf("Rules  : %s\n", manifest.GetRules())
	out := manifest.GetOutput()
	fmt.Printf("Output : %s/*.%s\n", out.Dir, out.Format)
	return nil
}

// +
func ShowResults(file, tool, analysisType, search, format string) error {
	fmt.Printf("[%s/orchestrator] ShowResults file=%q tool=%q type=%q search=%q format=%q\n", meta.AppName, file, tool, analysisType, search, format)
	// TODO: открыть SARIF/SBOM → применить фильтры → отформатировать вывод
	return nil
}

// // -+
// func Install(toolName string) error {
// 	fmt.Printf("[%s/orchestrator] Install tool=%q\n", meta.AppName, toolName)
// 	// TODO: загрузить манифест → deploy (copy/unzip/script) → записать в state/installed.yaml
// 	return nil
// }

// // -+
// func Uninstall(toolName string) error {
// 	fmt.Printf("[%s/orchestrator] Uninstall tool=%q\n", meta.AppName, toolName)
// 	// TODO: проверить state → удалить файлы → убрать запись из state/installed.yaml
// 	return nil
// }
