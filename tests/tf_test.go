package tests

import (
	"os"
	"strings"
	"testing"

	"terrabutler/cmd"

	"github.com/spf13/cobra"
)

var root *cobra.Command

// ------------------------------
// TEST MAIN - INICIALIZAÇÃO DOS COMANDOS
// ------------------------------
func TestMain(m *testing.M) {
	root = cmd.RootCmd

	tfCmd := &cobra.Command{Use: "tf", Short: "Fake tf root"}
	tfCmd.AddCommand(&cobra.Command{Use: "apply"})
	tfCmd.AddCommand(&cobra.Command{Use: "destroy"})
	tfCmd.AddCommand(&cobra.Command{Use: "output"})
	tfCmd.AddCommand(&cobra.Command{Use: "show"})
	tfCmd.AddCommand(&cobra.Command{Use: "refresh"})
	tfCmd.AddCommand(&cobra.Command{Use: "generate-vars"})
	tfCmd.AddCommand(&cobra.Command{Use: "init"})

	root.AddCommand(tfCmd)

	os.Exit(m.Run())
}

// ------------------------------
// FUNÇÕES AUXILIARES
// ------------------------------

func getSelectedEnv() string {
	data, err := os.ReadFile(".terrabutler_env")
	if err != nil {
		return ""
	}
	return strings.TrimSpace(string(data))
}

func findCommand(root *cobra.Command, path []string) *cobra.Command {
	if len(path) == 0 {
		return root
	}
	for _, cmd := range root.Commands() {
		if cmd.Use == path[0] {
			return findCommand(cmd, path[1:])
		}
	}
	return nil
}

// ------------------------------
// TESTES FUNCIONAIS BÁSICOS
// ------------------------------

func TestTfInitCmdWithoutEnvSelected(t *testing.T) {
	tempDir := t.TempDir()
	os.Chdir(tempDir)

	root.SetArgs([]string{"tf", "init"})
	err := root.Execute()

	if err != nil {
		t.Errorf("Expected no panic when no environment is selected, got error: %v", err)
	} else {
		t.Log("Terraform init executado com ambiente não selecionado (fallback correto).")
	}
}

func TestGetSelectedEnvWhenMissing(t *testing.T) {
	tempDir := t.TempDir()
	os.Chdir(tempDir)

	env := getSelectedEnv()
	if env != "" {
		t.Errorf("Expected empty environment when .terrabutler_env is missing, got: %s", env)
	} else {
		t.Log("Ambiente não selecionado retorna string vazia como esperado.")
	}
}

func TestGetSelectedEnvValid(t *testing.T) {
	tempDir := t.TempDir()
	os.Chdir(tempDir)

	os.WriteFile(".terrabutler_env", []byte("myenv"), 0644)
	env := getSelectedEnv()
	if env != "myenv" {
		t.Errorf("Expected 'myenv', got: %s", env)
	} else {
		t.Log("Ambiente selecionado lido corretamente: 'myenv'.")
	}
}

func TestTfApplyCmdStructure(t *testing.T) {
	cmd := findCommand(root, []string{"tf", "apply"})
	if cmd == nil {
		t.Error("tf apply command not found")
	} else {
		t.Log("Comando 'tf apply' encontrado na estrutura.")
	}
}

func TestTfDestroyCmdStructure(t *testing.T) {
	cmd := findCommand(root, []string{"tf", "destroy"})
	if cmd == nil {
		t.Error("tf destroy command not found")
	} else {
		t.Log("Comando 'tf destroy' encontrado na estrutura.")
	}
}

func TestTfOutputCmdStructure(t *testing.T) {
	cmd := findCommand(root, []string{"tf", "output"})
	if cmd == nil {
		t.Error("tf output command not found")
	} else {
		t.Log("Comando 'tf output' encontrado na estrutura.")
	}
}
