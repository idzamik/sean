// Package orchestrator содержит всю бизнес-логику CLI.
// Команды из cmd/ только разбирают флаги и вызывают функции этого пакета.
// Оркестратор решает, какой модуль из analyzers/ запустить,
// исходя из параметров, манифестов и state/installed.yaml.
package orchestrator

import (
	"fmt"

	"github.com/idzamik/sean/meta"
)

func RunAnalyser(target, tool, output, rules string) error {
	fmt.Printf("[%s/orchestrator] RunAnalyser target=%q tool=%q output=%q rules=%q\n", meta.AppName, target, tool, output, rules)
	// TODO: выбрать SAST-инструмент → analyzers/sast-<tool> → нормализовать в SARIF → сохранить
	return nil
}

// func RunSCA(target, tool, output string) error {
// 	fmt.Printf("[%s/orchestrator] RunSCA target=%q tool=%q output=%q\n", meta.AppName, target, tool, output)
// 	// TODO: выбрать SCA-инструмент → analyzers/sca-<tool> → получить SBOM → сохранить
// 	return nil
// }

// +
func Info(tool string, verbose bool) error {
	fmt.Printf("[%s/orchestrator] Info tool=%q verbose=%v\n", meta.AppName, tool, verbose)
	// TODO: прочитать state/installed.yaml → вывести инструменты и их статусы
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

// // -+
// func ListTools() error {
// 	fmt.Printf("[%s/orchestrator] ListTools\n", meta.AppName)
// 	// TODO: прочитать manifests/*.yaml → пересечь с state/installed.yaml → вывести таблицу
// 	return nil
// }
