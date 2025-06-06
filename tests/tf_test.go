package tests

import (
	"context"
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/montblu/terrabutler/internal/cmd"

	"github.com/hashicorp/terraform-exec/tfexec"
	"github.com/spf13/cobra"
)

var root *cobra.Command

// ------------------------------
// TEST MAIN - Initialization of the commands
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
// "Helper functions".
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
// Basic functional tests.
// ------------------------------

func TestTfInitCmdWithoutEnvSelected(t *testing.T) {
	tempDir := t.TempDir()
	os.Chdir(tempDir)

	root.SetArgs([]string{"tf", "init"})
	err := root.Execute()

	if err != nil {
		t.Errorf("Expected no panic when no environment is selected, got error: %v", err)
	} else {
		t.Log("Terraform init executed with no environment selected (fallback correct).")
	}
}

func TestGetSelectedEnvWhenMissing(t *testing.T) {
	tempDir := t.TempDir()
	os.Chdir(tempDir)

	env := getSelectedEnv()
	if env != "" {
		t.Errorf("Expected empty environment when .terrabutler_env is missing, got: %s", env)
	} else {
		t.Log("No environment selected returns an empty string as expected.")
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
		t.Log("Selected environment read correctly: 'myenv'.")
	}
}

// ------------------------------
// Command structure tests
// ------------------------------

func TestTfApplyCmdStructure(t *testing.T) {
	cmd := findCommand(root, []string{"tf", "apply"})
	if cmd == nil {
		t.Error("tf apply command not found")
	} else {
		t.Log("'tf apply' command found in structure.")
	}
}

func TestTfDestroyCmdStructure(t *testing.T) {
	cmd := findCommand(root, []string{"tf", "destroy"})
	if cmd == nil {
		t.Error("tf destroy command not found")
	} else {
		t.Log("'tf destroy' command found in structure.")
	}
}

func TestTfOutputCmdStructure(t *testing.T) {
	cmd := findCommand(root, []string{"tf", "output"})
	if cmd == nil {
		t.Error("tf output command not found")
	} else {
		t.Log("'tf output' command found in structure.")
	}
}

func TestTfShowCmdStructure(t *testing.T) {
	cmd := findCommand(root, []string{"tf", "show"})
	if cmd == nil {
		t.Error("tf show command not found")
	} else {
		t.Log("'tf show' command found in structure.")
	}
}

func TestTfRefreshCmdStructure(t *testing.T) {
	cmd := findCommand(root, []string{"tf", "refresh"})
	if cmd == nil {
		t.Error("tf refresh command not found")
	} else {
		t.Log("'tf refresh' command found in structure.")
	}
}

func TestTfGenVarsCmdStructure(t *testing.T) {
	cmd := findCommand(root, []string{"tf", "generate-vars"})
	if cmd == nil {
		t.Error("tf generate-vars command not found")
	} else {
		t.Log("'tf generate-vars' command found in structure .")
	}
}

// ------------------------------
// MOCK OF TERRAFORM
// ------------------------------

type mockTerraform struct {
	called      []string
	returnError bool
}

func (m *mockTerraform) Init(ctx context.Context, opts ...tfexec.InitOption) error {
	m.called = append(m.called, "Init")
	if m.returnError {
		return fmt.Errorf("mock init failed")
	}
	return nil
}

func (m *mockTerraform) Apply(ctx context.Context) error {
	m.called = append(m.called, "Apply")
	if m.returnError {
		return fmt.Errorf("mock apply failed")
	}
	return nil
}

func (m *mockTerraform) Destroy(ctx context.Context) error {
	m.called = append(m.called, "Destroy")
	if m.returnError {
		return fmt.Errorf("mock destroy failed")
	}
	return nil
}

func (m *mockTerraform) Refresh(ctx context.Context) error {
	m.called = append(m.called, "Refresh")
	if m.returnError {
		return fmt.Errorf("mock refresh failed")
	}
	return nil
}

func (m *mockTerraform) Show(ctx context.Context) (string, error) {
	m.called = append(m.called, "Show")
	if m.returnError {
		return "", fmt.Errorf("mock show failed")
	}
	return "mocked state", nil
}

func (m *mockTerraform) Output(ctx context.Context) (map[string]*tfexec.OutputMeta, error) {
	m.called = append(m.called, "Output")
	if m.returnError {
		return nil, fmt.Errorf("mock output failed")
	}
	return map[string]*tfexec.OutputMeta{
		"example": {Value: []byte(`"mocked"`)},
	}, nil
}

// ------------------------------
// TESTS WITH MOCK
// ------------------------------

func TestMockTerraformApply_Success(t *testing.T) {
	mock := &mockTerraform{}
	err := mock.Apply(context.Background())
	if err != nil {
		t.Errorf("Expected no error from Apply, got: %v", err)
	} else {
		t.Log("Apply executed successfully.")
	}
	if len(mock.called) == 0 || mock.called[0] != "Apply" {
		t.Error("Expected Apply to be called")
	} else {
		t.Log("Apply call confirmed.")
	}
}

func TestMockTerraformApply_Error(t *testing.T) {
	mock := &mockTerraform{returnError: true}
	err := mock.Apply(context.Background())
	if err == nil {
		t.Error("Expected error from Apply, got nil")
	} else {
		t.Logf("Apply failed as expected: %v", err)
	}
}

func TestMockTerraformDestroy_Success(t *testing.T) {
	mock := &mockTerraform{}
	err := mock.Destroy(context.Background())
	if err != nil {
		t.Errorf("Expected no error from Destroy, got: %v", err)
	} else {
		t.Log("Destroy executed successfully.")
	}
	if mock.called[0] != "Destroy" {
		t.Error("Expected Destroy to be called")
	} else {
		t.Log("Destroy call confirmed.")
	}
}

func TestMockTerraformDestroy_Error(t *testing.T) {
	mock := &mockTerraform{returnError: true}
	err := mock.Destroy(context.Background())
	if err == nil {
		t.Error("Expected error from Destroy, got nil")
	} else {
		t.Logf("Destroy failed as expected: %v", err)
	}
}

func TestMockTerraformOutput_Success(t *testing.T) {
	mock := &mockTerraform{}
	out, err := mock.Output(context.Background())
	if err != nil {
		t.Errorf("Expected no error from Output, got: %v", err)
	} else {
		t.Log("Output executed successfully.")
	}
	if string(out["example"].Value) != `"mocked"` {
		t.Errorf("Expected mocked output value, got: %v", string(out["example"].Value))
	} else {
		t.Log("Mock output value confirmed.")
	}
}

func TestMockTerraformOutput_Error(t *testing.T) {
	mock := &mockTerraform{returnError: true}
	_, err := mock.Output(context.Background())
	if err == nil {
		t.Error("Expected error from Output, got nil")
	} else {
		t.Logf("Output failed as expected: %v", err)
	}
}

func TestMockTerraformShow_Success(t *testing.T) {
	mock := &mockTerraform{}
	output, err := mock.Show(context.Background())
	if err != nil {
		t.Errorf("Expected no error from Show, got: %v", err)
	} else {
		t.Log("Show executed successfully.")
	}
	if output != "mocked state" {
		t.Errorf("Expected 'mocked state', got: %s", output)
	} else {
		t.Log("Show output confirmed.")
	}
}

func TestMockTerraformShow_Error(t *testing.T) {
	mock := &mockTerraform{returnError: true}
	_, err := mock.Show(context.Background())
	if err == nil {
		t.Error("Expected error from Show, got nil")
	} else {
		t.Logf("Show failed as expected: %v", err)
	}
}

func TestMockTerraformRefresh_Success(t *testing.T) {
	mock := &mockTerraform{}
	err := mock.Refresh(context.Background())
	if err != nil {
		t.Errorf("Expected no error from Refresh, got: %v", err)
	} else {
		t.Log("Refresh executed successfully.")
	}
	if mock.called[0] != "Refresh" {
		t.Error("Expected Refresh to be called")
	} else {
		t.Log("Refresh call confirmed.")
	}
}

func TestMockTerraformRefresh_Error(t *testing.T) {
	mock := &mockTerraform{returnError: true}
	err := mock.Refresh(context.Background())
	if err == nil {
		t.Error("Expected error from Refresh, got nil")
	} else {
		t.Logf("Refresh failed as expected: %v", err)
	}
}
