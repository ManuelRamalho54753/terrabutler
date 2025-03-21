package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

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

func createEnv(name string) error {
	path := fmt.Sprintf("./environments/%s", name)
	if _, err := os.Stat(path); !os.IsNotExist(err) {
		return fmt.Errorf("the environment '%s' already exists", name)
	}
	return os.MkdirAll(path, os.ModePerm)
}

func init() {
	envCmd.AddCommand(envNewCmd)
}
