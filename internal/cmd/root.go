package cmd

import (
	"fmt"
	"os"

	"github.com/montblu/terrabutler/internal/requirements"
	"github.com/montblu/terrabutler/internal/settings"
	"github.com/montblu/terrabutler/internal/tf"
	"github.com/spf13/cobra"
)

// RootCmd defines the base command
var RootCmd = &cobra.Command{
	Use:   "terrabutler",
	Short: "Terrabutler CLI",
	Long:  "Terrabutler is a CLI tool for managing environments and Terraform commands.",
}

// Register all subcommands
func init() {
	RootCmd.AddCommand(tf.TfCmd)
	RootCmd.AddCommand(envSwitchCmd)

}

// Execute runs the root command after validating settings
func Execute() {
	validateSettings()
	requirements.CheckRequirements()

	if err := RootCmd.Execute(); err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}
}

// validateSettings ensures TERRABUTLER_ROOT and settings file are present
func validateSettings() {
	root := os.Getenv("TERRABUTLER_ROOT")
	if root == "" {
		fmt.Println("Environment variable TERRABUTLER_ROOT is not set.")
		fmt.Println("Use: export TERRABUTLER_ROOT=/path/to/repository")
		os.Exit(1)
	}

	settingsPath := fmt.Sprintf("%s/internal/configs/settings.yml", root)
	_, err := settings.LoadSettings(settingsPath)
	if err != nil {
		fmt.Println("Error loading settings.yml file:")
		fmt.Println("  ", err)
		os.Exit(1)
	}
}
