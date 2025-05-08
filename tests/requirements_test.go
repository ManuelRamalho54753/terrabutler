package tests

import (
	"os"
	"os/exec"
	"path/filepath"
	"testing"
)

// Função auxiliar para criar um ambiente de teste válido
func createValidEnvironment(t *testing.T) string {
	t.Helper()

	tmp := t.TempDir()
	os.Setenv("TERRABUTLER_ENABLE", "true")
	os.Setenv("TERRABUTLER_ROOT", tmp)

	configsPath := filepath.Join(tmp, "configs")
	os.Mkdir(configsPath, 0755)

	settingsContent := []byte("environment: dev\nversion: 1.0.0\n") // Ajusta conforme o schema esperado
	err := os.WriteFile(filepath.Join(configsPath, "settings.yaml"), settingsContent, 0644)
	if err != nil {
		t.Fatalf("Erro ao criar settings.yaml: %v", err)
	}

	return tmp
}

// ------------------------------
// TESTE: TERRABUTLER_ENABLE não está definido
// ------------------------------
func TestCheckRequirements_DisabledEnv(t *testing.T) {
	if os.Getenv("TEST_CHILD") == "1" {
		checkRequirements()
		return
	}

	cmd := exec.Command(os.Args[0], "-test.run=TestCheckRequirements_DisabledEnv")
	cmd.Env = append(os.Environ(), "TERRABUTLER_ENABLE=false", "TEST_CHILD=1")

	err := cmd.Run()
	if exitError, ok := err.(*exec.ExitError); ok && !exitError.Success() {
		t.Log("Falhou corretamente por TERRABUTLER_ENABLE estar desativado")
	} else {
		t.Fatal("Esperava que falhasse por TERRABUTLER_ENABLE estar desativado")
	}
}

// ------------------------------
// TESTE: TERRABUTLER_ROOT não está definido ou inválido
// ------------------------------
func TestCheckRequirements_InvalidRoot(t *testing.T) {
	if os.Getenv("TEST_CHILD") == "1" {
		os.Setenv("TERRABUTLER_ENABLE", "true")
		os.Setenv("TERRABUTLER_ROOT", "/caminho/inexistente")
		checkRequirements()
		return
	}

	cmd := exec.Command(os.Args[0], "-test.run=TestCheckRequirements_InvalidRoot")
	cmd.Env = append(os.Environ(), "TEST_CHILD=1")

	err := cmd.Run()
	if exitError, ok := err.(*exec.ExitError); ok && !exitError.Success() {
		t.Log("Falhou corretamente por TERRABUTLER_ROOT inválido")
	} else {
		t.Fatal("Esperava falhar por TERRABUTLER_ROOT inválido")
	}
}

// ------------------------------
// TESTE: settings.yaml ausente
// ------------------------------
func TestCheckRequirements_MissingSettingsFile(t *testing.T) {
	if os.Getenv("TEST_CHILD") == "1" {
		tmp := t.TempDir()
		os.Setenv("TERRABUTLER_ENABLE", "true")
		os.Setenv("TERRABUTLER_ROOT", tmp)

		os.Mkdir(filepath.Join(tmp, "configs"), 0755)

		checkRequirements()
		return
	}

	cmd := exec.Command(os.Args[0], "-test.run=TestCheckRequirements_MissingSettingsFile")
	cmd.Env = append(os.Environ(), "TEST_CHILD=1")

	err := cmd.Run()
	if exitError, ok := err.(*exec.ExitError); ok && !exitError.Success() {
		t.Log("Falhou corretamente por ausência do settings.yaml")
	} else {
		t.Fatal("Esperava falhar por settings.yaml ausente")
	}
}

// ------------------------------
// TESTE: settings.yaml inválido (simulado com conteúdo inválido)
// ------------------------------
func TestCheckRequirements_InvalidSettings(t *testing.T) {
	if os.Getenv("TEST_CHILD") == "1" {
		tmp := t.TempDir()
		os.Setenv("TERRABUTLER_ENABLE", "true")
		os.Setenv("TERRABUTLER_ROOT", tmp)

		confPath := filepath.Join(tmp, "configs")
		os.Mkdir(confPath, 0755)
		os.WriteFile(filepath.Join(confPath, "settings.yaml"), []byte("invalid_yaml: [::"), 0644)

		checkRequirements()
		return
	}

	cmd := exec.Command(os.Args[0], "-test.run=TestCheckRequirements_InvalidSettings")
	cmd.Env = append(os.Environ(), "TEST_CHILD=1")

	err := cmd.Run()
	if exitError, ok := err.(*exec.ExitError); ok && !exitError.Success() {
		t.Log("Falhou corretamente por settings.yaml inválido")
	} else {
		t.Fatal("Esperava falhar por settings.yaml inválido")
	}
}

// ------------------------------
// TESTE: tudo correto (ambiente válido)
// ------------------------------
func TestCheckRequirements_Success(t *testing.T) {
	tmp := createValidEnvironment(t)

	if os.Getenv("TEST_CHILD") == "1" {
		checkRequirements()
		t.Log("checkRequirements executado com sucesso em ambiente válido")
		return
	}

	cmd := exec.Command(os.Args[0], "-test.run=TestCheckRequirements_Success")
	cmd.Env = append(os.Environ(), "TEST_CHILD=1", "TERRABUTLER_ROOT="+tmp)

	err := cmd.Run()
	if err != nil {
		t.Fatalf("Esperava sucesso com ambiente válido, mas falhou: %v", err)
	}
}
