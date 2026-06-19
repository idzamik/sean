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
	Action string
	Flags  []string
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

	userFlagSet := buildUserFlagSet(userFlags)

	if cmd.Action != "" {
		argv = append(argv, cmd.Action)
	}

	for _, flagTpl := range cmd.Flags {
		resolved := resolvePlaceholders(flagTpl, target, rulesPath, outputPath)
		flagName := extractFlagName(resolved)
		if flagName != "" && userFlagSet[flagName] {
			continue
		}
		argv = append(argv, splitFlagString(flagTpl, resolved)...)
	}
	argv = append(argv, userFlags...)

	return bin, argv
}


func resolvePlaceholders(tpl, target, rules, output string) string {
	r := strings.NewReplacer(
		"{target}", target,
		"{rules}", rules,
		"{output}", output,
	)
	return r.Replace(tpl)
}


func splitFlagString(tpl, resolved string) []string {
	trimmed := strings.TrimSpace(tpl)
	if trimmed == "{target}" || trimmed == "{rules}" || trimmed == "{output}" {
		return []string{resolved}
	}

	parts := strings.SplitN(resolved, " ", 2)
	if len(parts) == 2 {
		return []string{parts[0], parts[1]}
	}
	return []string{resolved}
}


func extractFlagName(s string) string {
	s = strings.TrimSpace(s)
	if !strings.HasPrefix(s, "-") {
		return ""
	}
	s = strings.TrimLeft(s, "-")
	if i := strings.IndexAny(s, " ="); i != -1 {
		s = s[:i]
	}
	return s
}


func buildUserFlagSet(userFlags []string) map[string]bool {
	set := make(map[string]bool)
	for _, f := range userFlags {
		if name := extractFlagName(f); name != "" {
			set[name] = true
		}
	}
	return set
}


func GenerateOutputPath(toolName, dir, format string) (string, error) {
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return "", fmt.Errorf("cannot create output dir %s: %w", dir, err)
	}
	timestamp := time.Now().Format("20060102_150405")
	filename := fmt.Sprintf("%s_%s.%s", toolName, timestamp, format)
	return filepath.Join(dir, filename), nil
}

