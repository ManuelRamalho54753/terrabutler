package main

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"terrabutler/cmd"
	"terrabutler/logger"
	"terrabutler/variables"

	"github.com/hashicorp/terraform-exec/tfexec"
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
		checkRequirements()
		site := args[0]
		execTerraformWithSDK(site, func(tf *tfexec.Terraform) error {
			return tf.Init(context.Background(), tfexec.Upgrade(true), tfexec.Reconfigure(true))
		})
	},
}

var tfApplyCmd = &cobra.Command{
	Use:   "apply [site]",
	Short: "Run terraform apply for a specific site",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		checkRequirements()
		site := args[0]
		execTerraformWithSDK(site, func(tf *tfexec.Terraform) error {
			return tf.Apply(context.Background())
		})
	},
}

var tfDestroyCmd = &cobra.Command{
	Use:   "destroy [site]",
	Short: "Run terraform destroy for a specific site",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		checkRequirements()
		site := args[0]
		execTerraformWithSDK(site, func(tf *tfexec.Terraform) error {
			return tf.Destroy(context.Background())
		})
	},
}

var tfOutputCmd = &cobra.Command{
	Use:   "output [site]",
	Short: "Show Terraform outputs for a specific site",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		checkRequirements()
		site := args[0]
		execTerraformWithSDK(site, func(tf *tfexec.Terraform) error {
			outputs, err := tf.Output(context.Background())
			if err != nil {
				return err
			}
			for k, v := range outputs {
				fmt.Printf("%s = %s\n", k, v.Value)
			}
			return nil
		})
	},
}

var tfShowCmd = &cobra.Command{
	Use:   "show [site]",
	Short: "Show Terraform state for a specific site",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		checkRequirements()
		site := args[0]
		execTerraformWithSDK(site, func(tf *tfexec.Terraform) error {
			output, err := tf.Show(context.Background())
			if err != nil {
				return err
			}
			fmt.Println(output)
			return nil
		})
	},
}

var tfRefreshCmd = &cobra.Command{
	Use:   "refresh [site]",
	Short: "Refresh Terraform state for a specific site",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		checkRequirements()
		site := args[0]
		execTerraformWithSDK(site, func(tf *tfexec.Terraform) error {
			return tf.Refresh(context.Background())
		})
	},
}

var tfGenVarsCmd = &cobra.Command{
	Use:   "generate-vars [env]",
	Short: "Generate tfvars for all sites for a given environment",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		env := args[0]
		logger.Log.Infof("Generating variable files for environment: %s", env)
		variables.GenerateVarFiles(env)
	},
}

func execTerraformWithSDK(site string, action func(*tfexec.Terraform) error) {
	root := os.Getenv("TERRABUTLER_ROOT")
	if root == "" {
		logger.Log.Error("TERRABUTLER_ROOT is not set")
		os.Exit(1)
	}

	sitePath := filepath.Join(root, fmt.Sprintf("site_%s", site))
	if _, err := os.Stat(sitePath); os.IsNotExist(err) {
		logger.Log.Errorf("Site directory does not exist: %s", sitePath)
		os.Exit(1)
	}

	tf, err := tfexec.NewTerraform(sitePath, "terraform")
	if err != nil {
		logger.Log.Errorf("Failed to initialize Terraform SDK: %v", err)
		os.Exit(1)
	}

	if err := action(tf); err != nil {
		logger.Log.Errorf("Error running Terraform command: %v", err)
		os.Exit(1)
	}

	logger.Log.Infof("Terraform command executed successfully for site: %s", site)
}

func init() {
	cmd.RootCmd.AddCommand(tfCmd)
	tfCmd.AddCommand(tfInitCmd)
	tfCmd.AddCommand(tfApplyCmd)
	tfCmd.AddCommand(tfDestroyCmd)
	tfCmd.AddCommand(tfOutputCmd)
	tfCmd.AddCommand(tfShowCmd)
	tfCmd.AddCommand(tfRefreshCmd)
	tfCmd.AddCommand(tfGenVarsCmd)
}
