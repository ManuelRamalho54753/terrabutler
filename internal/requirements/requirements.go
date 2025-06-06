package requirements

import (
	"os"
	"path/filepath"

	"github.com/montblu/terrabutler/internal/logger"
	"github.com/montblu/terrabutler/internal/settings"
	"github.com/montblu/terrabutler/internal/utils"
)

func CheckRequirements() {
	// Check if TERRABUTLER_ENABLE is set to "true"
	if os.Getenv("TERRABUTLER_ENABLE") != "true" {
		logger.Log.Fatal("Terrabutler is not currently enabled on this folder.")
	}

	// Check if TERRABUTLER_ROOT is set and exists
	root := os.Getenv("TERRABUTLER_ROOT")
	if root == "" || !utils.PathExists(root) {
		logger.Log.Fatal("Terrabutler can't determine the root folder of your project or it doesn't exist.")
	}

	// Check if settings.yaml exists
	settingsPath := filepath.Join(root, "configs", "settings.yml")
	if !utils.PathExists(settingsPath) {
		logger.Log.Fatal("Terrabutler can't find your settings file.")
	}

	// Load and validate settings using Koanf
	s, err := settings.LoadSettings(settingsPath)
	if err != nil {
		logger.Log.Fatal("Error loading settings", zapError(err))
	}

	if err := settings.ValidateSettings(s); err != nil {
		logger.Log.Fatal("Invalid settings", zapError(err))
	}
}

func zapError(err error) any {
	return map[string]any{"error": err.Error()}
}
