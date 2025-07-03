package tests

import (
	"os"
	"testing"

	"github.com/montblu/terrabutler/internal/cmd"

	"github.com/spf13/cobra"
)

func TestIsValidEnvName(t *testing.T) {
	cases := []struct {
		name     string
		input    string
		expected bool
	}{
		{"valid hyphen", "dev-env", true},
		{"valid underscore", "env_test", true},
		{"valid alphanumeric", "env123", true},
		{"invalid uppercase", "Env", false},
		{"invalid space", "my env", false},
		{"invalid symbol", "env@dev", false},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			result := cmd.IsValidEnvName(tc.input)
			if result != tc.expected {
				t.Errorf("Expected %v, got %v", tc.expected, result)
			}
		})
	}
}

func TestRunEnvNew(t *testing.T) {
	tempDir := t.TempDir()
	os.Chdir(tempDir)
	os.Mkdir("environments", 0755)

	err := cmd.RunEnvNew("valid_env")
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	err = cmd.RunEnvNew("Invalid Env!")
	if err == nil {
		t.Error("Expected error for invalid env name, got nil")
	}
}

func TestEnvListCmd(t *testing.T) {
	tempDir := t.TempDir()
	os.Chdir(tempDir)
	os.MkdirAll("environments/env1", 0755)
	os.MkdirAll("environments/env2", 0755)

	rootCmd := &cobra.Command{Use: "tb"}
	rootCmd.AddCommand(cmd.EnvCmd())
	rootCmd.SetArgs([]string{"env", "list"})

	if err := rootCmd.Execute(); err != nil {
		t.Errorf("env list failed: %v", err)
	}
}

func TestEnvDeleteCmd(t *testing.T) {
	tempDir := t.TempDir()
	os.Chdir(tempDir)
	os.MkdirAll("environments/delete_me", 0755)

	rootCmd := &cobra.Command{Use: "tb"}
	rootCmd.AddCommand(cmd.EnvCmd())
	rootCmd.SetArgs([]string{"env", "delete", "delete_me"})

	if err := rootCmd.Execute(); err != nil {
		t.Errorf("env delete failed: %d", err)
	}

	if _, err := os.Stat("environments/delete_me"); !os.IsNotExist(err) {
		t.Error("Environment was not deleted")
	}
}

func TestEnvSelectCmd(t *testing.T) {
	tempDir := t.TempDir()
	os.Chdir(tempDir)
	os.MkdirAll("environments/myenv", 0755)

	rootCmd := &cobra.Command{Use: "tb"}
	rootCmd.AddCommand(cmd.EnvCmd())
	rootCmd.SetArgs([]string{"env", "select", "myenv"})

	if err := rootCmd.Execute(); err != nil {
		t.Errorf("env select failed: %v", err)
	}

	data, err := os.ReadFile(".terrabutler_env")
	if err != nil || string(data) != "myenv" {
		t.Errorf("Expected 'myenv' selected, got: %s", data)
	}
}

func TestEnvShowCmd(t *testing.T) {
	tempDir := t.TempDir()
	os.Chdir(tempDir)
	os.WriteFile(".terrabutler_env", []byte("showenv"), 0644)

	rootCmd := &cobra.Command{Use: "tb"}
	rootCmd.AddCommand(cmd.EnvCmd())
	rootCmd.SetArgs([]string{"env", "show"})

	if err := rootCmd.Execute(); err != nil {
		t.Errorf("env show failed: %v", err)
	}

}
func TestEnvRenameCmd(t *testing.T) {
	tempDir := t.TempDir()
	os.Chdir(tempDir)
	os.MkdirAll("environments/old_env", 0755)
	os.WriteFile(".terrabutler_env", []byte("old_env"), 0644)

	rootCmd := &cobra.Command{Use: "tb"}
	rootCmd.AddCommand(cmd.EnvCmd())
	rootCmd.SetArgs([]string{"env", "rename", "old_env", "new_env"})

	if err := rootCmd.Execute(); err != nil {
		t.Errorf("env rename failed: %v", err)
	}

	if _, err := os.Stat("environments/new_env"); os.IsNotExist(err) {
		t.Error("New environment not found after rename")
	}

	selected, _ := cmd.GetSelectedEnv()
	if selected != "new_env" {
		t.Errorf("Expected selected env to be 'new_env', got %s", selected)
	}
}
