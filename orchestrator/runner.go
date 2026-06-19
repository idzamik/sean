package orchestrator

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"
)

type ManifestCommand struct {
	Bin    string
	Action string   // подкоманда инструмента (например "scan")
	Flags  []string // шаблоны флагов: "--config {rules}", "{target}", "--metrics=off"
}

type ManifestOutput struct {
	Format string
	Dir    string
}


func BuildCommand(
	cmd ManifestCommand,
	target, rulesPath, outputPath string,
	userFlags []string,
) (bin string, argv []string) {
	bin = cmd.Bin

	// Множество имён флагов, переданных пользователем (без "--" и значения)
	userFlagSet := buildUserFlagSet(userFlags)

	// Подкоманда инструмента (например "scan")
	if cmd.Action != "" {
		argv = append(argv, cmd.Action)
	}

	// Параметры из манифеста
	for _, flagTpl := range cmd.Flags {
		// Подставляем плейсхолдеры
		resolved := resolvePlaceholders(flagTpl, target, rulesPath, outputPath)

		// Проверяем, не переопределил ли пользователь этот флаг
		flagName := extractFlagName(resolved)
		if flagName != "" && userFlagSet[flagName] {
			continue // пользовательский флаг имеет приоритет
		}

		// Разбиваем строку на аргументы
		argv = append(argv, splitFlagString(flagTpl, resolved)...)
	}

	// Пользовательские флаги добавляются в конец
	argv = append(argv, userFlags...)

	return bin, argv
}

// resolvePlaceholders подставляет плейсхолдеры в строку флага.
func resolvePlaceholders(tpl, target, rules, output string) string {
	r := strings.NewReplacer(
		"{target}", target,
		"{rules}", rules,
		"{output}", output,
	)
	return r.Replace(tpl)
}

// splitFlagString разбивает строку флага на аргументы.
//
// Правила:
//   - Если строка — чистый плейсхолдер ("{target}") — возвращаем resolved как один аргумент
//     (защита от пробелов в путях)
//   - "--flag=value" — один аргумент (знак = не является разделителем)
//   - "--flag value" — два аргумента (разбиваем по первому пробелу)
//   - "--flag" — один аргумент (булевый флаг)
func splitFlagString(tpl, resolved string) []string {
	// Чистый плейсхолдер → один аргумент, не разбивать
	trimmed := strings.TrimSpace(tpl)
	if trimmed == "{target}" || trimmed == "{rules}" || trimmed == "{output}" {
		return []string{resolved}
	}

	// Разбиваем по первому пробелу
	parts := strings.SplitN(resolved, " ", 2)
	if len(parts) == 2 {
		return []string{parts[0], parts[1]}
	}
	return []string{resolved}
}

// extractFlagName возвращает имя флага без "--" и значения.
// "--config tools/semgrep/rules" → "config"
// "--sarif-output /path/to.sarif" → "sarif-output"
// "--metrics=off" → "metrics"
// "tools/semgrep/rules" (позиционный) → ""
func extractFlagName(s string) string {
	s = strings.TrimSpace(s)
	if !strings.HasPrefix(s, "-") {
		return "" // позиционный аргумент
	}
	s = strings.TrimLeft(s, "-")
	// Убираем значение после пробела или =
	if i := strings.IndexAny(s, " ="); i != -1 {
		s = s[:i]
	}
	return s
}

// buildUserFlagSet строит множество имён флагов из пользовательских аргументов.
// ["--config", "./rules", "--metrics=off"] → {"config": true, "metrics": true}
func buildUserFlagSet(userFlags []string) map[string]bool {
	set := make(map[string]bool)
	for _, f := range userFlags {
		if name := extractFlagName(f); name != "" {
			set[name] = true
		}
	}
	return set
}

// GenerateOutputPath создаёт директорию результатов и возвращает путь к файлу.
// Формат имени: <name>_<YYYYMMDD_HHMMSS>.<format>
func GenerateOutputPath(toolName, dir, format string) (string, error) {
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return "", fmt.Errorf("cannot create output dir %s: %w", dir, err)
	}
	timestamp := time.Now().Format("20060102_150405")
	filename := fmt.Sprintf("%s_%s.%s", toolName, timestamp, format)
	return filepath.Join(dir, filename), nil
}