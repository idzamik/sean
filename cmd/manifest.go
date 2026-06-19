package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

type State struct {
	Tools []StateEntry `yaml:"tools"`
}

type StateEntry struct {
	Name     string `yaml:"name"`
	Manifest string `yaml:"manifest"`
}

type Manifest struct {
	Version        string             `yaml:"version"`
	Name           string             `yaml:"name"`
	Bin            string             `yaml:"bin"`
	DefaultCommand string             `yaml:"default_command"`
	Rules          string             `yaml:"rules"`
	Commands       map[string]Command `yaml:"commands"`
	Output         OutputConfig       `yaml:"output"`
}

type Command struct {
	Bin    string   `yaml:"bin"`
	Action string   `yaml:"action"`
	Flags  []string `yaml:"flags"`
}

type OutputConfig struct {
	Format string `yaml:"format"`
	Dir    string `yaml:"dir"`
}

func configBase() string {
	if dir := os.Getenv("SEAN_CONFIG_DIR"); dir != "" {
		return dir
	}
	return "configs"
}

func resolveConfigPath(rel string) string {
	return filepath.Join(configBase(), rel)
}

func LoadState() (*State, error) {
	path := resolveConfigPath("installed.yaml")
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("cannot read state file %s: %w", path, err)
	}
	var state State
	if err := yaml.Unmarshal(data, &state); err != nil {
		return nil, fmt.Errorf("cannot parse state file %s: %w", path, err)
	}
	return &state, nil
}

func LoadManifest(relativePath string) (*Manifest, error) {
	path := resolveConfigPath(relativePath)
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("cannot read manifest %s: %w", path, err)
	}
	var manifest Manifest
	if err := yaml.Unmarshal(data, &manifest); err != nil {
		return nil, fmt.Errorf("cannot parse manifest %s: %w", path, err)
	}
	return &manifest, nil
}

func LoadManifestByName(toolName string) (*Manifest, error) {
	state, err := LoadState()
	if err != nil {
		return nil, fmt.Errorf("cannot load state: %w", err)
	}
	for _, entry := range state.Tools {
		if entry.Name == toolName {
			return LoadManifest(entry.Manifest)
		}
	}
	return nil, fmt.Errorf("tool %q not found in installed.yaml", toolName)
}