package tests

import (
	"os"
	"os/exec"
	"path/filepath"
	"testing"

	"terrabutler/inception"
	"terrabutler/logger"

	"go.uber.org/zap"
)

func init() {
	// Logger temporário para testes
	l, _ := zap.NewDevelopment()
	logger.Log = l.Sugar()
}

// ------------------------------
// HELPERS
// ------------------------------

// Cria um settings.yml válido e backend.tfvars simulado
func createTestSettingsAndBackend(tempDir string) {
	os.MkdirAll(filepath.Join(tempDir, "configs"), 0755)
	os.WriteFile(filepath.Join(tempDir, "configs", "settings.yml"), []byte(`
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
`), 0644)

	os.MkdirAll(filepath.Join(tempDir, "backends"), 0755)
	os.WriteFile(filepath.Join(tempDir, "backends", "testorg-dev-inception.tfvars"), []byte(""), 0644)

	os.MkdirAll(filepath.Join(tempDir, "site_inception"), 0755)
}

// ------------------------------
// TESTS FOR InitInception
// ------------------------------

func TestInitInception_Success(t *testing.T) {
	tempDir := t.TempDir()
	os.Setenv("TERRABUTLER_ROOT", tempDir)

	createTestSettingsAndBackend(tempDir)

	inception.InitInception()

	envFile := filepath.Join(tempDir, "site_inception", ".terraform", "environment")
	if _, err := os.Stat(envFile); os.IsNotExist(err) {
		t.Errorf("Expected .terraform/environment to be created at %s", envFile)
	}
}

func TestInitInception_AlreadyExists(t *testing.T) {
	tempDir := t.TempDir()
	os.Setenv("TERRABUTLER_ROOT", tempDir)

	createTestSettingsAndBackend(tempDir)

	// Cria manualmente a estrutura de init
	terraformDir := filepath.Join(tempDir, "site_inception", ".terraform")
	os.MkdirAll(terraformDir, 0755)
	os.WriteFile(filepath.Join(terraformDir, "environment"), []byte("dev"), 0644)

	inception.InitInception() // Não deve fazer nada (já existe)
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
	os.MkdirAll(terraformDir, 0755)
	os.WriteFile(filepath.Join(terraformDir, "environment"), []byte("dev"), 0644)

	inception.InceptionInitNeeded() // Não deve falhar
}
