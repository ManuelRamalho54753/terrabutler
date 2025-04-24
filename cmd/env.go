package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"

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
			fmt.Printf("Ambiente \"%s\" não existe.\n", oldName)
			return
		}
		if _, err := os.Stat(newPath); err == nil {
			fmt.Printf("Já existe um ambiente chamado \"%s\".\n", newName)
			return
		}
		if err := os.Rename(oldPath, newPath); err != nil {
			fmt.Println("Erro ao renomear:", err)
			return
		}
		fmt.Printf("Ambiente \"%s\" foi renomeado para \"%s\" com sucesso.\n", oldName, newName)
		selected, err := GetSelectedEnv()
		if err == nil && selected == oldName {
			SetSelectedEnv(newName)
			fmt.Println("Ambiente selecionado foi atualizado para o novo nome.")
		}
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

func RunEnvNew(name string) error {
	if !IsValidEnvName(name) {
		return fmt.Errorf("Invalid environment name. Use only lowercase letters, numbers, hyphens or underscores.")
	}
	return CreateEnv(name)
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

func init() {
	RootCmd.AddCommand(envCmd)
	envCmd.AddCommand(envNewCmd)
	envCmd.AddCommand(envListCmd)
	envCmd.AddCommand(envDeleteCmd)
	envCmd.AddCommand(envSelectCmd)
	envCmd.AddCommand(envShowCmd)
	envCmd.AddCommand(envRenameCmd)
	RootCmd.AddCommand(versionCmd)
}
