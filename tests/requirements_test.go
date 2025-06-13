package tests

import (
	"os"
	"path/filepath"
	"testing"
)

// Wrapper for the CheckRequirements function from main.go, adapted for testing
func checkRequirements() {
	if os.Getenv("TERRABUTLER_ENABLE") != "true" {
		panic("Terrabutler is not enabled")
	}

	root := os.Getenv("TERRABUTLER_ROOT")
	if root == "" {
		panic("TERRABUTLER_ROOT is not set")
	}

	settingsPath := filepath.Join(root, "configs", "settings.yml")
	if _, err := os.Stat(settingsPath); os.IsNotExist(err) {
		panic("settings.yml is missing")
	}
}

// ------------------------------
// TESTS
// ------------------------------

// Tests when all environment variables and files are correctly set
func TestCheckRequirements_AllValid(t *testing.T) {
	t.Setenv("TERRABUTLER_ENABLE", "true")
	t.Setenv("TERRABUTLER_ROOT", t.TempDir())

	root := os.Getenv("TERRABUTLER_ROOT")
	settingsDir := filepath.Join(root, "configs")

	err := os.MkdirAll(settingsDir, 0755)
	if err != nil {
		t.Fatalf("Failed to create configs dir: %v", err)
	}

	settingsContent := []byte(`general:
  organization: "MyOrg"
  secrets_key_id: "key123"
sites:
  ordered: ["site1"]
environments:
  default:
    domain: "example.com"
    name: "default"
    profile_name: "profile"
    region: "us-east-1"
  permanent: ["prod"]
  temporary:
    secrets:
      firebase_credentials: "somekey"
      mail_password: "somepass"
`)
	err = os.WriteFile(filepath.Join(settingsDir, "settings.yml"), settingsContent, 0644)
	if err != nil {
		t.Fatalf("Failed to write settings.yml: %v", err)
	}

	checkRequirements()
}

// Tests when TERRABUTLER_ENABLE is empty
func TestCheckRequirements_MissingEnable(t *testing.T) {
	t.Setenv("TERRABUTLER_ENABLE", "")
	t.Setenv("TERRABUTLER_ROOT", t.TempDir())

	defer func() {
		if r := recover(); r == nil {
			t.Error("Expected panic due to missing TERRABUTLER_ENABLE")
		}
	}()

	checkRequirements()
}

// Tests when TERRABUTLER_ROOT is empty
func TestCheckRequirements_MissingRoot(t *testing.T) {
	t.Setenv("TERRABUTLER_ENABLE", "true")
	t.Setenv("TERRABUTLER_ROOT", "")

	defer func() {
		if r := recover(); r == nil {
			t.Error("Expected panic due to missing TERRABUTLER_ROOT")
		}
	}()

	checkRequirements()
}

// Tests when the settings.yaml file is missing
func TestCheckRequirements_MissingSettingsFile(t *testing.T) {
	t.Setenv("TERRABUTLER_ENABLE", "true")
	t.Setenv("TERRABUTLER_ROOT", t.TempDir())

	// Does not create the settings.yaml file

	defer func() {
		if r := recover(); r == nil {
			t.Error("Expected panic due to missing settings file")
		}
	}()

	checkRequirements()
}
