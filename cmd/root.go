package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/yourusername/sean/meta"
)

var rootCmd = &cobra.Command{
	Use:   meta.AppName,
	Short: meta.AppDescription,
	Long:  meta.AppLongDescription,
	RunE: func(cmd *cobra.Command, args []string) error {
		return cmd.Help()
	},
}

// versionCmd — secman version
var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print " + meta.AppName + " version",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("%s (%s) version %s\n", meta.AppName, meta.AppFullName, meta.AppVersion)
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "%s error: %v\n", meta.AppName, err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
