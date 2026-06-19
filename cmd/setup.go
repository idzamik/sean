package cmd

import (
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
	"github.com/idzamik/sean/meta"
	"gopkg.in/yaml.v3"
)

const (
	defaultBinDir    = "/usr/local/bin"
	defaultConfigDir = "/etc/sean"
	defaultDataDir   = "/var/lib/sean"
)

var (
	setupInstall bool
	setupConfig  string
	setupUninstall bool
)

var setupCmd = &cobra.Command{
	Use:   "setup",
	Short: "Install, configure and uninstall " + meta.AppName,
	Long: `Управление установкой и конфигурацией ` + meta.AppName + `.

Примеры:
  sudo ` + meta.AppName + ` setup -i
  ` + meta.AppName + ` setup -c ./manifests/trivy.yaml
  sudo ` + meta.AppName + ` setup -u`,
	RunE: func(cmd *cobra.Command, args []string) error {
		switch {
		case setupInstall:
			return runSetupInstall()
		case setupConfig != "":
			return runSetupConfig(setupConfig)
		case setupUninstall:
			return runSetupUninstall()
		default:
			return cmd.Help()
		}
	},
}


func runSetupInstall() error {
	if os.Getuid() != 0 {
		return fmt.Errorf("setup -i requires root privileges. Run: sudo sean setup -i")
	}

	fmt.Printf("[%s] Starting installation...\n\n", meta.AppName)

	binDest := filepath.Join(defaultBinDir, meta.AppName)
	if err := copySelf(binDest); err != nil {
		return fmt.Errorf("cannot install binary: %w", err)
	}
	fmt.Printf("  ✓ binary      → %s\n", binDest)

	manifestsDir := filepath.Join(defaultConfigDir, "manifests")
	if err := os.MkdirAll(manifestsDir, 0o755); err != nil {
		return fmt.Errorf("cannot create config dir: %w", err)
	}
	fmt.Printf("  ✓ config dir  → %s\n", defaultConfigDir)

	installedPath := filepath.Join(defaultConfigDir, "installed.yaml")
	if _, err := os.Stat(installedPath); os.IsNotExist(err) {
		empty := State{Tools: []StateEntry{}}
		data, _ := yaml.Marshal(empty)
		if err := os.WriteFile(installedPath, data, 0o644); err != nil {
			return fmt.Errorf("cannot create installed.yaml: %w", err)
		}
		fmt.Printf("  ✓ state file  → %s\n", installedPath)
	} else {
		fmt.Printf("  ~ state file  → %s (already exists, skipped)\n", installedPath)
	}

	resultsDir := filepath.Join(defaultDataDir, "results")
	if err := os.MkdirAll(resultsDir, 0o755); err != nil {
		return fmt.Errorf("cannot create data dir: %w", err)
	}
	fmt.Printf("  ✓ data dir    → %s\n", defaultDataDir)

	profilePath := "/etc/profile.d/sean.sh"
	profileContent := fmt.Sprintf("export SEAN_CONFIG_DIR=%s\n", defaultConfigDir)
	if err := os.WriteFile(profilePath, []byte(profileContent), 0o644); err != nil {
		fmt.Printf("  ! env file    → cannot write %s: %v (set manually)\n", profilePath, err)
	} else {
		fmt.Printf("  ✓ env file    → %s\n", profilePath)
	}

	fmt.Printf("\n[%s] Installation complete.\n", meta.AppName)
	fmt.Printf("       Run: %s info\n", meta.AppName)
	fmt.Printf("       Add tools: %s setup -c ./path/to/manifest.yaml\n\n", meta.AppName)
	fmt.Println("  Note: restart shell or run 'source /etc/profile.d/sean.sh' to apply env.")
	return nil
}


func runSetupConfig(srcPath string) error {
	srcPath = filepath.Clean(srcPath)
	if _, err := os.Stat(srcPath); os.IsNotExist(err) {
		return fmt.Errorf("manifest file not found: %s", srcPath)
	}

	m, err := loadManifestFromPath(srcPath)
	if err != nil {
		return fmt.Errorf("invalid manifest: %w", err)
	}
	if m.Name == "" {
		return fmt.Errorf("manifest must have a 'name' field")
	}
	if m.Bin == "" {
		return fmt.Errorf("manifest must have a 'bin' field")
	}

	configDir := configBase()
	manifestsDir := filepath.Join(configDir, "manifests")
	if err := os.MkdirAll(manifestsDir, 0o755); err != nil {
		return fmt.Errorf("cannot access manifests dir %s: %w", manifestsDir, err)
	}

	destName := m.Name + ".yaml"
	destPath := filepath.Join(manifestsDir, destName)
	if err := copyFile(srcPath, destPath); err != nil {
		return fmt.Errorf("cannot copy manifest: %w", err)
	}
	fmt.Printf("[%s] Manifest installed → %s\n", meta.AppName, destPath)

	installedPath := filepath.Join(configDir, "installed.yaml")
	state, err := loadOrCreateState(installedPath)
	if err != nil {
		return fmt.Errorf("cannot load state: %w", err)
	}

	relManifestPath := filepath.Join("manifests", destName)

	for i, entry := range state.Tools {
		if entry.Name == m.Name {
			state.Tools[i].Manifest = relManifestPath
			fmt.Printf("[%s] Updated existing entry for %q in installed.yaml\n", meta.AppName, m.Name)
			return saveState(installedPath, state)
		}
	}

	state.Tools = append(state.Tools, StateEntry{
		Name:     m.Name,
		Manifest: relManifestPath,
	})
	fmt.Printf("[%s] Registered %q in installed.yaml\n", meta.AppName, m.Name)
	return saveState(installedPath, state)
}


func runSetupUninstall() error {
	if os.Getuid() != 0 {
		return fmt.Errorf("setup -u requires root privileges. Run: sudo sean setup -u")
	}

	fmt.Printf("[%s] Uninstalling...\n\n", meta.AppName)

	targets := []string{
		filepath.Join(defaultBinDir, meta.AppName),
		defaultConfigDir,
		defaultDataDir,
		"/etc/profile.d/sean.sh",
	}

	for _, path := range targets {
		if err := os.RemoveAll(path); err != nil {
			fmt.Printf("  ! failed to remove %s: %v\n", path, err)
		} else {
			fmt.Printf("  ✓ removed %s\n", path)
		}
	}

	fmt.Printf("\n[%s] Uninstall complete.\n", meta.AppName)
	return nil
}


func copySelf(dest string) error {
	self, err := os.Executable()
	if err != nil {
		return fmt.Errorf("cannot resolve current binary path: %w", err)
	}
	self, err = filepath.EvalSymlinks(self)
	if err != nil {
		return fmt.Errorf("cannot resolve symlinks: %w", err)
	}
	return copyFile(self, dest)
}

func copyFile(src, dest string) error {
	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()

	info, err := in.Stat()
	if err != nil {
		return err
	}

	out, err := os.OpenFile(dest, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, info.Mode())
	if err != nil {
		return err
	}
	defer out.Close()

	if _, err := io.Copy(out, in); err != nil {
		return err
	}
	return out.Sync()
}

func loadManifestFromPath(path string) (*Manifest, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var m Manifest
	if err := yaml.Unmarshal(data, &m); err != nil {
		return nil, err
	}
	return &m, nil
}

func loadOrCreateState(path string) (*State, error) {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return &State{Tools: []StateEntry{}}, nil
	}
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var state State
	if err := yaml.Unmarshal(data, &state); err != nil {
		return nil, err
	}
	return &state, nil
}

func saveState(path string, state *State) error {
	data, err := yaml.Marshal(state)
	if err != nil {
		return fmt.Errorf("cannot serialize state: %w", err)
	}
	if err := os.WriteFile(path, data, 0o644); err != nil {
		return fmt.Errorf("cannot write state file: %w", err)
	}
	return nil
}

func init() {
	setupCmd.Flags().BoolVarP(&setupInstall, "install", "i", false, "Initialize directories and install binary")
	setupCmd.Flags().StringVarP(&setupConfig, "config", "c", "", "Register a tool manifest (path to .yaml file)")
	setupCmd.Flags().BoolVarP(&setupUninstall, "uninstall", "u", false, "Remove all installed files and directories")
	rootCmd.AddCommand(setupCmd)
}