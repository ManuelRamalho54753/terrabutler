package tests

import (
	"os"
	"path/filepath"
	"testing"
)

// Wrapper da função CheckRequirements do ficheiro main.go, adaptado para testes
func checkRequirements() {
	if os.Getenv("TERRABUTLER_ENABLE") != "true" {
		panic("Terrabutler is not enabled")
	}

	root := os.Getenv("TERRABUTLER_ROOT")
	if root == "" {
		panic("TERRABUTLER_ROOT is not set")
	}

	settingsPath := filepath.Join(root, "configs", "settings.yaml")
	if _, err := os.Stat(settingsPath); os.IsNotExist(err) {
		panic("settings.yaml is missing")
	}
}

// ------------------------------
// TESTES
// ------------------------------

// Testa quando todas as variáveis e ficheiros estão corretamente definidos
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
	err = os.WriteFile(filepath.Join(settingsDir, "settings.yaml"), settingsContent, 0644)
	if err != nil {
		t.Fatalf("Failed to write settings.yaml: %v", err)
	}

	checkRequirements() // Deve passar sem pânico
}

// Testa quando TERRABUTLER_ENABLE está em branco
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

// Testa quando TERRABUTLER_ROOT está em branco
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

// Testa quando o ficheiro settings.yaml está em falta
func TestCheckRequirements_MissingSettingsFile(t *testing.T) {
	t.Setenv("TERRABUTLER_ENABLE", "true")
	t.Setenv("TERRABUTLER_ROOT", t.TempDir())

	// Não cria o ficheiro settings.yaml

	defer func() {
		if r := recover(); r == nil {
			t.Error("Expected panic due to missing settings file")
		}
	}()

	checkRequirements()
}
