package main

import (
	"fmt"
	"os"
	"path/filepath"
	"terrabutler/utils"
)

func checkRequirements() {
	// Check if TERRABUTLER_ENABLE is set to "true"
	if os.Getenv("TERRABUTLER_ENABLE") != "true" {
		fmt.Println("Terrabutler is not currently enabled on this folder.")
		os.Exit(1)
	}

	// Check if TERRABUTLER_ROOT is set and exists
	root := os.Getenv("TERRABUTLER_ROOT")
	if root == "" || !utils.PathExists(root) {
		fmt.Println("Terrabutler can't determine the root folder of your project or it doesn't exist.")
		os.Exit(1)
	}

	// Check if settings.yml exists
	settingsPath := filepath.Join(root, "configs", "settings.yml")
	if !utils.PathExists(settingsPath) {
		fmt.Println("Terrabutler can't find your settings file.")
		os.Exit(1)
	}

	// Load and validate settings using Koanf
	settings, err := LoadSettings(settingsPath)
	if err != nil {
		fmt.Println("Error loading settings:", err)
		os.Exit(1)
	}

	if err := ValidateSettings(settings); err != nil {
		fmt.Println("Invalid settings:", err)
		os.Exit(1)
	}
}
