package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/montblu/terrabutler/internal/settings"
	"github.com/montblu/terrabutler/internal/utils"
	"github.com/montblu/terrabutler/internal/variables"
	"github.com/spf13/cobra"
)

var envCmd = &cobra.Command{
	Use:   "env",
	Short: "Manage environments",
}

var envNewCmd = &cobra.Command{
	Use:   "new [name]",
	Short: "Create a new environment",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		name := args[0]
		if err := RunEnvNew(name); err != nil {
			fmt.Println(err.Error())
			os.Exit(1)
		}
		fmt.Println("Environment created successfully:", name)
	},
}

func RunEnvNew(name string) error {
	if !IsValidEnvName(name) {
		return fmt.Errorf("Invalid environment name. Use only lowercase letters, numbers, hyphens or underscores.")
	}

	root := os.Getenv("TERRABUTLER_ROOT")
	paths := utils.GetPaths(root)
	settingsPath := paths["settings"]

	cfg, err := settings.LoadSettings(settingsPath)
	if err != nil {
		return fmt.Errorf("Error loading settings: %v", err)
	}

	sites := cfg.Sites.Ordered
	if Contains(sites, "inception") {
		sites = Remove(sites, "inception")
	}

	defaultFiles := filepath.Join(paths["root"], "internal", "configs", "default_tf_files")

	for _, site := range sites {
		sitePath := filepath.Join(paths["environments"], name, site)
		err := os.MkdirAll(sitePath, 0755)
		if err != nil {
			return fmt.Errorf("Failed to create site directory: %v", err)
		}

		err = copyTerraformFiles(defaultFiles, sitePath)
		if err != nil {
			return fmt.Errorf("Failed to copy Terraform files: %v", err)
		}
	}

	variables.GenerateVarFiles(name)
	return nil
}

func copyTerraformFiles(srcDir, dstDir string) error {
	files, err := os.ReadDir(srcDir)
	if err != nil {
		return err
	}
	for _, file := range files {
		srcPath := filepath.Join(srcDir, file.Name())
		dstPath := filepath.Join(dstDir, file.Name())
		srcData, err := os.ReadFile(srcPath)
		if err != nil {
			return err
		}
		err = os.WriteFile(dstPath, srcData, 0644)
		if err != nil {
			return err
		}
	}
	return nil
}

var envListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all environments",
	Run: func(cmd *cobra.Command, args []string) {
		files, err := os.ReadDir("./environments")
		if err != nil {
			fmt.Println("Failed to list environments:", err)
			return
		}
		fmt.Println("Environments:")
		for _, file := range files {
			if file.IsDir() {
				fmt.Println("-", file.Name())
			}
		}
	},
}

var envDeleteCmd = &cobra.Command{
	Use:   "delete [name]",
	Short: "Delete an environment",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		name := args[0]
		path := filepath.Join("./environments", name)
		if _, err := os.Stat(path); os.IsNotExist(err) {
			fmt.Println("Environment does not exist:", name)
			return
		}
		err := os.RemoveAll(path)
		if err != nil {
			fmt.Println("Failed to delete environment:", err)
			return
		}
		fmt.Println("Environment deleted successfully:", name)
	},
}

var envSelectCmd = &cobra.Command{
	Use:   "select [name]",
	Short: "Select an environment",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		name := args[0]
		path := filepath.Join("./environments", name)
		if _, err := os.Stat(path); os.IsNotExist(err) {
			fmt.Println("Environment does not exist:", name)
			return
		}
		err := os.WriteFile(".terrabutler_env", []byte(name), 0644)
		if err != nil {
			fmt.Println("Failed to select environment:", err)
			return
		}
		fmt.Println("Environment selected:", name)
	},
}

var envShowCmd = &cobra.Command{
	Use:   "show",
	Short: "Show the currently selected environment",
	Run: func(cmd *cobra.Command, args []string) {
		data, err := os.ReadFile(".terrabutler_env")
		if err != nil {
			fmt.Println("No environment selected or file missing")
			return
		}
		fmt.Println("Current environment:", string(data))
	},
}

func CreateEnv(name string) error {
	path := filepath.Join("./environments", name)
	if _, err := os.Stat(path); !os.IsNotExist(err) {
		return fmt.Errorf("the environment '%s' already exists", name)
	}
	return os.MkdirAll(path, os.ModePerm)
}

var envRenameCmd = &cobra.Command{
	Use:   "rename [old] [new]",
	Short: "Rename an existing environment",
	Args:  cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		oldName := args[0]
		newName := args[1]

		if !IsValidEnvName(newName) {
			fmt.Println("Invalid new environment name. Use only lowercase letters, numbers, hyphens or underscores.")
			return
		}

		envsPath := filepath.Join(".", "environments")
		oldPath := filepath.Join(envsPath, oldName)
		newPath := filepath.Join(envsPath, newName)

		if _, err := os.Stat(oldPath); os.IsNotExist(err) {
			fmt.Printf("Environment \"%s\" does not exist.\n", oldName)
			return
		}
		if _, err := os.Stat(newPath); err == nil {
			fmt.Printf("An environment with this name already exists \"%s\".\n", newName)
			return
		}
		if err := os.Rename(oldPath, newPath); err != nil {
			fmt.Println("Error while renaming:", err)
			return
		}
		fmt.Printf("Environment \"%s\" was successfully renamed to \"%s\".\n", oldName, newName)
		selected, err := GetSelectedEnv()
		if err == nil && selected == oldName {
			SetSelectedEnv(newName)
			fmt.Println("Selected environment was updated to the new name.")
		}
	},
}

var envSwitchCmd = &cobra.Command{
	Use:   "env switch [env]",
	Short: "Set default environment for Terraform",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		env := args[0]
		err := os.WriteFile(".terraform/environment", []byte(env), 0644)
		if err != nil {
			fmt.Printf("Error setting default environment: %v\n", err)
			os.Exit(1)
		}
		fmt.Printf("Current environment set to: %s\n", env)
	},
}

func GetSelectedEnv() (string, error) {
	data, err := os.ReadFile(".terrabutler_env")
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(data)), nil
}

func SetSelectedEnv(env string) error {
	return os.WriteFile(".terrabutler_env", []byte(env), 0644)
}

func IsValidEnvName(name string) bool {
	matched, _ := regexp.MatchString(`^[a-z0-9_-]+$`, name)
	return matched
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Show the version of terrabutler",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("terrabutler version 1.0.0")
	},
}

func EnvCmd() *cobra.Command {
	return envCmd
}

func EnvRenameCmd() *cobra.Command {
	return envRenameCmd
}

func Contains(slice []string, item string) bool {
	for _, v := range slice {
		if v == item {
			return true
		}
	}
	return false
}

func Remove(slice []string, item string) []string {
	result := []string{}
	for _, v := range slice {
		if v != item {
			result = append(result, v)
		}
	}
	return result
}

func init() {
	RootCmd.AddCommand(envCmd)
	envCmd.AddCommand(envNewCmd)
	envCmd.AddCommand(envListCmd)
	envCmd.AddCommand(envDeleteCmd)
	envCmd.AddCommand(envSelectCmd)
	envCmd.AddCommand(envShowCmd)
	envCmd.AddCommand(envRenameCmd)
	envCmd.AddCommand(envSwitchCmd)
	RootCmd.AddCommand(versionCmd)
}
