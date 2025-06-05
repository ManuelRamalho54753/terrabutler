package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"terrabutler/utils"
	"terrabutler/variables"

	"github.com/spf13/cobra"
)

var variablesCmd = &cobra.Command{
	Use:   "variables",
	Short: "Generate tfvars files for the current environment",
	Run: func(cmd *cobra.Command, args []string) {
		root := os.Getenv("TERRABUTLER_ROOT")
		if root == "" {
			fmt.Println("TERRABUTLER_ROOT is not set")
			return
		}

		env, err := utils.GetCurrentEnv(filepath.Join(root, "site_inception"))
		if err != nil {
			fmt.Println("Error:", err)
			return
		}

		variables.GenerateVarFiles(env)
		fmt.Printf("✅ Variable files for environment '%s' generated successfully.\n", env)
	},
}

func init() {
	RootCmd.AddCommand(variablesCmd)
}
