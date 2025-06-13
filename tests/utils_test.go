package tests

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/montblu/terrabutler/internal/utils"
)

// ------------------------------
// TESTE: PathExists
// ------------------------------
func TestPathExists_ExistingPath(t *testing.T) {
	tmp := t.TempDir()
	if !utils.PathExists(tmp) {
		t.Errorf("Expected PathExists to return true for existing path")
	}
}

func TestPathExists_NonexistentPath(t *testing.T) {
	nonexistent := filepath.Join(os.TempDir(), "does_not_exist")
	if utils.PathExists(nonexistent) {
		t.Errorf("Expected PathExists to return false for nonexistent path")
	}
}

// ------------------------------
// TESTE: GetCurrentEnv
// ------------------------------
func TestGetCurrentEnv_Success(t *testing.T) {
	tmp := t.TempDir()
	tfDir := filepath.Join(tmp, ".terraform")
	os.Mkdir(tfDir, 0755)

	expectedEnv := "dev"
	err := os.WriteFile(filepath.Join(tfDir, "environment"), []byte(expectedEnv), 0644)
	if err != nil {
		t.Fatalf("Could not write environment file: %v", err)
	}

	env, err := utils.GetCurrentEnv(tmp)
	if err != nil {
		t.Errorf("Expected no error, got: %v", err)
	}
	if env != expectedEnv {
		t.Errorf("Expected env %q, got %q", expectedEnv, env)
	}
}

func TestGetCurrentEnv_FileMissing(t *testing.T) {
	tmp := t.TempDir()
	_, err := utils.GetCurrentEnv(tmp)
	if err == nil {
		t.Error("Expected error when environment file is missing")
	}
}

// ------------------------------
// TESTE: GetPaths
// ------------------------------
func TestGetPaths_ReturnsCorrectPaths(t *testing.T) {
	root := "/my/project"
	paths := utils.GetPaths(root)

	expected := map[string]string{
		"root":        root,
		"settings":    filepath.Join(root, "configs", "settings.yml"),
		"inception":   filepath.Join(root, "site_inception"),
		"backends":    filepath.Join(root, "configs", "backends"),
		"variables":   filepath.Join(root, "configs", "variables"),
		"environment": filepath.Join(root, ".terraform", "environment"),
	}

	for key, expectedPath := range expected {
		if paths[key] != expectedPath {
			t.Errorf("Path %q: expected %q, got %q", key, expectedPath, paths[key])
		}
	}
}
