package utils

import (
	"fmt"
	"os"
	"path/filepath"
)

// PathExists checks if a file or directory exists
func PathExists(path string) bool {
	_, err := os.Stat(path)
	return !os.IsNotExist(err)
}

// GetCurrentEnv reads the current environment from the .terraform/environment file
func GetCurrentEnv(root string) (string, error) {
	envFile := filepath.Join(root, ".terraform", "environment")
	data, err := os.ReadFile(envFile)
	if err != nil {
		return "", fmt.Errorf("could not read current environment: %w", err)
	}
	return string(data), nil
}

// GetPaths returns a map with all key project paths based on TERRABUTLER_ROOT
func GetPaths(root string) map[string]string {
	return map[string]string{
		"root":        root,
		"settings":    filepath.Join(root, "configs", "settings.yaml"),
		"inception":   filepath.Join(root, "site_inception"),
		"backends":    filepath.Join(root, "configs", "backends"),
		"variables":   filepath.Join(root, "configs", "variables"),
		"environment": filepath.Join(root, ".terraform", "environment"),
	}
}
