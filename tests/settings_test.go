package tests

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/montblu/terrabutler/internal/settings"
)

// ------------------------------
// TESTEs for LoadSettings
// ------------------------------

func TestLoadSettings_ValidFile(t *testing.T) {
	tmp := t.TempDir()
	yamlContent := `
general:
  organization: test-org
  secrets_key_id: key123

sites:
  ordered: [site1, site2]

environments:
  default:
    domain: test.com
    name: dev
    profile_name: default
    region: eu-west-1
  permanent: [prod]
  temporary:
    secrets:
      firebase_credentials: firebase-key
      mail_password: mail-secret
`
	settingsPath := filepath.Join(tmp, "settings.yml")
	if err := os.WriteFile(settingsPath, []byte(yamlContent), 0644); err != nil {
		t.Fatalf("failed to write settings.yml: %v", err)
	}

	s, err := settings.LoadSettings(settingsPath)
	if err != nil {
		t.Fatalf("LoadSettings failed: %v", err)
	}
	if s.General.Organization != "test-org" {
		t.Errorf("Expected organization 'test-org', got '%s'", s.General.Organization)
	}
}

func TestLoadSettings_InvalidFile(t *testing.T) {
	tmp := t.TempDir()
	settingsPath := filepath.Join(tmp, "invalid.yml")
	if err := os.WriteFile(settingsPath, []byte("invalid: [:"), 0644); err != nil {
		t.Fatalf("failed to write invalid file: %v", err)
	}

	_, err := settings.LoadSettings(settingsPath)
	if err == nil {
		t.Error("Expected LoadSettings to fail on invalid file, but got no error")
	}
}

// ------------------------------
// TESTS FOR ValidateSettings
// ------------------------------

func TestValidateSettings_Success(t *testing.T) {
	s := &settings.Settings{}
	s.General.Organization = "org"
	s.Sites.Ordered = []string{"site1"}
	s.Environments.Default.Name = "dev"

	err := settings.ValidateSettings(s)
	if err != nil {
		t.Errorf("Expected validation success, got error: %v", err)
	}
}

func TestValidateSettings_MissingOrganization(t *testing.T) {
	s := &settings.Settings{}
	s.Sites.Ordered = []string{"site1"}
	s.Environments.Default.Name = "dev"

	err := settings.ValidateSettings(s)
	if err == nil || err.Error() != "organization field in general config is required" {
		t.Errorf("Expected organization error, got: %v", err)
	}
}

func TestValidateSettings_EmptySites(t *testing.T) {
	s := &settings.Settings{}
	s.General.Organization = "org"
	s.Environments.Default.Name = "dev"

	err := settings.ValidateSettings(s)
	if err == nil || err.Error() != "at least one site must be listed in sites.ordered" {
		t.Errorf("Expected sites error, got: %v", err)
	}
}

func TestValidateSettings_MissingDefaultEnvName(t *testing.T) {
	s := &settings.Settings{}
	s.General.Organization = "org"
	s.Sites.Ordered = []string{"site1"}

	err := settings.ValidateSettings(s)
	if err == nil || err.Error() != "default environment name is required" {
		t.Errorf("Expected default env name error, got: %v", err)
	}
}

// ------------------------------
// TESTS FOR GetSettings (panics expected)
// ------------------------------

func TestGetSettings_PanicOnInvalidPath(t *testing.T) {
	os.Setenv("TERRABUTLER_ROOT", "/invalid/path")
	defer func() {
		if r := recover(); r == nil {
			t.Error("Expected panic on invalid path, got none")
		}
	}()
	_ = settings.GetSettings()
}

func TestGetSettings_PanicOnInvalidData(t *testing.T) {
	tmp := t.TempDir()
	os.Setenv("TERRABUTLER_ROOT", tmp)
	configs := filepath.Join(tmp, "configs")
	os.MkdirAll(configs, 0755)
	os.WriteFile(filepath.Join(configs, "settings.yml"), []byte("invalid:["), 0644)

	defer func() {
		if r := recover(); r == nil {
			t.Error("Expected panic on invalid yml data, got none")
		}
	}()

	_ = settings.GetSettings()
}
