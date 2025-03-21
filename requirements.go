package main

import (
	"fmt"
	"os"
	"path/filepath"
)

func checkRequirements() {
	// Check if TERRABUTLER_ENABLE is set to "true"
	if os.Getenv("TERRABUTLER_ENABLE") != "true" {
		fmt.Println("Terrabutler is not currently enabled on this folder.")
		os.Exit(1)
	}

	// Check if TERRABUTLER_ROOT is set and exists
	root := os.Getenv("TERRABUTLER_ROOT")
	if root == "" || !pathExists(root) {
		fmt.Println("Terrabutler can't determine the root folder of your project or it doesn't exist.")
		os.Exit(1)
	}

	// Check if settings.yml exists in configs folder
	settingsPath := filepath.Join(root, "configs", "settings.yml")
	if !pathExists(settingsPath) {
		fmt.Println("Terrabutler can't find your settings file.")
		os.Exit(1)
	}
}

// pathExists checks if a given path exists
func pathExists(path string) bool {
	_, err := os.Stat(path)
	return !os.IsNotExist(err)
}
