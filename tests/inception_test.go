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
	// Create a temporary logger just for testing
	l, _ := zap.NewDevelopment()
	logger.Log = l.Sugar()
}

// ------------------------------
// TESTS FOR InitInception
// ------------------------------

// Test that InitInception creates the site_inception directory when it does not exist
func TestInitInception_Success(t *testing.T) {
	tempDir := t.TempDir()
	os.Setenv("TERRABUTLER_ROOT", tempDir)

	sitePath := filepath.Join(tempDir, "site_inception")
	os.RemoveAll(sitePath)

	inception.InitInception()

	if _, err := os.Stat(sitePath); os.IsNotExist(err) {
		t.Errorf("Expected site_inception to be created at %s", sitePath)
	}
}

// Test that InitInception does nothing if the site_inception directory already exists
func TestInitInception_AlreadyExists(t *testing.T) {
	tempDir := t.TempDir()
	os.Setenv("TERRABUTLER_ROOT", tempDir)

	sitePath := filepath.Join(tempDir, "site_inception")
	os.MkdirAll(sitePath, 0755)

	inception.InitInception() // Should log info but not fail
}

// ------------------------------
// TESTS FOR InceptionInitNeeded
// ------------------------------

// Test that InceptionInitNeeded fails (exits) if site_inception directory is missing
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
		t.Fatal("Expected failure due to missing site_inception, but got nil")
	}
	if exitErr, ok := err.(*exec.ExitError); !ok || exitErr.ExitCode() == 0 {
		t.Fatal("Expected non-zero exit code due to missing site_inception")
	}
}

// Test that InceptionInitNeeded succeeds if site_inception directory exists
func TestInceptionInitNeeded_SuccessWhenExists(t *testing.T) {
	tempDir := t.TempDir()
	os.Setenv("TERRABUTLER_ROOT", tempDir)

	sitePath := filepath.Join(tempDir, "site_inception")
	os.MkdirAll(sitePath, 0755)

	inception.InceptionInitNeeded() // Should pass without error
}
