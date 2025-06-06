package tests

import (
	"os"
	"os/exec"
	"path/filepath"
	"testing"

	"github.com/montblu/terrabutler/internal/inception"
	"github.com/montblu/terrabutler/internal/logger"

	"go.uber.org/zap"
)

func init() {
	// Temporary logger for tests
	l, _ := zap.NewDevelopment()
	logger.Log = l.Sugar()
}

// ------------------------------
// HELPERS
// ------------------------------

func createTestSettingsAndBackend(tempDir string, t *testing.T) {
	t.Helper()

	settings := `Erro ao carregar ficheiro settings.yml:")
general:
  organization: testorg
  secrets_key_id: dummy-key
sites:
  ordered: ["site1"]
environments:
  default:
    domain: "test.com"
    name: "dev"
    profile_name: "default"
    region: "us-west-1"
  permanent: []
  temporary:
    secrets:
      firebase_credentials: "abc"
      mail_password: "123"
`

	settingsPath := filepath.Join(tempDir, "configs")
	if err := os.MkdirAll(settingsPath, 0755); err != nil {
		t.Fatalf("failed to create configs dir: %v", err)
	}
	if err := os.WriteFile(filepath.Join(settingsPath, "settings.yml"), []byte(settings), 0644); err != nil {
		t.Fatalf("failed to write settings.yml: %v", err)
	}

	backendPath := filepath.Join(tempDir, "backends")
	if err := os.MkdirAll(backendPath, 0755); err != nil {
		t.Fatalf("failed to create backends dir: %v", err)
	}
	if err := os.WriteFile(filepath.Join(backendPath, "testorg-dev-inception.tfvars"), []byte(""), 0644); err != nil {
		t.Fatalf("failed to write backend tfvars: %v", err)
	}

	sitePath := filepath.Join(tempDir, "site_inception")
	if err := os.MkdirAll(sitePath, 0755); err != nil {
		t.Fatalf("failed to create site_inception dir: %v", err)
	}
}

// ------------------------------
// TESTS FOR InitInception
// ------------------------------

func TestInitInception_Success(t *testing.T) {
	tempDir := t.TempDir()
	os.Setenv("TERRABUTLER_ROOT", tempDir)

	createTestSettingsAndBackend(tempDir, t)

	inception.InitInception()

	envFile := filepath.Join(tempDir, "site_inception", ".terraform", "environment")
	if _, err := os.Stat(envFile); os.IsNotExist(err) {
		t.Errorf("Expected .terraform/environment to be created at %s", envFile)
	}
}

func TestInitInception_AlreadyExists(t *testing.T) {
	tempDir := t.TempDir()
	os.Setenv("TERRABUTLER_ROOT", tempDir)

	createTestSettingsAndBackend(tempDir, t)

	terraformDir := filepath.Join(tempDir, "site_inception", ".terraform")
	if err := os.MkdirAll(terraformDir, 0755); err != nil {
		t.Fatalf("failed to create .terraform dir: %v", err)
	}
	if err := os.WriteFile(filepath.Join(terraformDir, "environment"), []byte("dev"), 0644); err != nil {
		t.Fatalf("failed to create environment file: %v", err)
	}

	inception.InitInception() // Should not do anything, already exists
}

// ------------------------------
// TESTS FOR InceptionInitNeeded
// ------------------------------

func TestInceptionInitNeeded_FailsWhenMissing(t *testing.T) {
	if os.Getenv("TEST_CHILD") == "1" {
		tempDir := t.TempDir()
		os.Setenv("TERRABUTLER_ROOT", tempDir)
		inception.InceptionInitNeeded()
		return
	}

	cmd := exec.Command(os.Args[0], "-test.run=TestInceptionInitNeeded_FailsWhenMissing")
	cmd.Env = append(os.Environ(), "TEST_CHILD=1")
	err := cmd.Run()

	if err == nil {
		t.Fatal("Expected failure due to missing .terraform/environment, but got nil")
	}
	if exitErr, ok := err.(*exec.ExitError); !ok || exitErr.ExitCode() == 0 {
		t.Fatal("Expected non-zero exit code due to missing .terraform/environment")
	}
}

func TestInceptionInitNeeded_SuccessWhenExists(t *testing.T) {
	tempDir := t.TempDir()
	os.Setenv("TERRABUTLER_ROOT", tempDir)

	terraformDir := filepath.Join(tempDir, "site_inception", ".terraform")
	if err := os.MkdirAll(terraformDir, 0755); err != nil {
		t.Fatalf("failed to create .terraform dir: %v", err)
	}
	if err := os.WriteFile(filepath.Join(terraformDir, "environment"), []byte("dev"), 0644); err != nil {
		t.Fatalf("failed to create environment file: %v", err)
	}

	inception.InceptionInitNeeded()
}
