package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/spf13/cobra"
)

var tfCmd = &cobra.Command{
	Use:   "tf",
	Short: "Execute Terraform commands",
}

var tfInitCmd = &cobra.Command{
	Use:   "init [site]",
	Short: "Initialize Terraform for a specific site",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		site := args[0]
		runTerraform(site, "init", "-reconfigure")
	},
}

var tfApplyCmd = &cobra.Command{
	Use:   "apply [site]",
	Short: "Run terraform apply for a specific site",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		site := args[0]
		runTerraform(site, "apply", "-auto-approve")
	},
}

var tfDestroyCmd = &cobra.Command{
	Use:   "destroy [site]",
	Short: "Run terraform destroy for a specific site",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		site := args[0]
		runTerraform(site, "destroy", "-auto-approve")
	},
}

func runTerraform(site string, command string, arg string) {
	root := os.Getenv("TERRABUTLER_ROOT")
	if root == "" {
		fmt.Println("TERRABUTLER_ROOT is not set")
		os.Exit(1)
	}

	sitePath := filepath.Join(root, fmt.Sprintf("site_%s", site))
	if _, err := os.Stat(sitePath); os.IsNotExist(err) {
		fmt.Printf("Site directory does not exist: %s\n", sitePath)
		os.Exit(1)
	}

	cmd := exec.Command("terraform", command, arg)
	cmd.Dir = sitePath
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	fmt.Printf("Running 'terraform %s %s' in %s\n", command, arg, sitePath)
	if err := cmd.Run(); err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Terraform %s completed successfully for site %s\n", command, site)
}

func init() {
	RootCmd.AddCommand(tfCmd)
	tfCmd.AddCommand(tfInitCmd)
	tfCmd.AddCommand(tfApplyCmd)
	tfCmd.AddCommand(tfDestroyCmd)
}
