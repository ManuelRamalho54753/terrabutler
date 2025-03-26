package cmd

import (
	"fmt"
	"os"
	"path/filepath"

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
		if err := createEnv(name); err != nil {
			fmt.Println("Error creating environment:", err)
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

func createEnv(name string) error {
	path := filepath.Join("./environments", name)
	if _, err := os.Stat(path); !os.IsNotExist(err) {
		return fmt.Errorf("the environment '%s' already exists", name)
	}
	return os.MkdirAll(path, os.ModePerm)
}

func init() {
	RootCmd.AddCommand(envCmd)
	envCmd.AddCommand(envNewCmd)
	envCmd.AddCommand(envListCmd)
	envCmd.AddCommand(envDeleteCmd)
	envCmd.AddCommand(envSelectCmd)
	envCmd.AddCommand(envShowCmd)
}
