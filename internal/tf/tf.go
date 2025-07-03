package tf

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/montblu/terrabutler/internal/logger"
	"github.com/montblu/terrabutler/internal/requirements"
	"github.com/montblu/terrabutler/internal/variables"

	"github.com/hashicorp/terraform-exec/tfexec"
	"github.com/spf13/cobra"
)

// TfCmd is the exported root command for Terraform-related operations
var TfCmd = &cobra.Command{
	Use:   "tf",
	Short: "Execute Terraform commands",
}

// ------------------------
// Subcommands
// ------------------------

var tfInitCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize Terraform for the selected environment",
	Args:  cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		requirements.CheckRequirements()
		envName := GetSelectedEnv()
		if envName == "" {
			cmd.PrintErrln("No environment selected. Please use `terrabutler env select [name]` first.")
			cmd.SilenceUsage = true
			cmd.SilenceErrors = true
			return fmt.Errorf("no environment selected")
		}
		err := execTerraformWithSDK(envName, func(tf *tfexec.Terraform) error {
			return tf.Init(context.Background(), tfexec.Upgrade(true), tfexec.Reconfigure(true))
		})
		if err != nil {
			cmd.PrintErrf("Terraform init failed: %v\n", err)
			cmd.SilenceUsage = true
			cmd.SilenceErrors = true
			return err
		}
		return nil
	},
}

var tfApplyCmd = &cobra.Command{
	Use:   "apply [site]",
	Short: "Run terraform apply for a specific site",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		requirements.CheckRequirements()
		site := args[0]
		err := execTerraformWithSDK(site, func(tf *tfexec.Terraform) error {
			return tf.Apply(context.Background())
		})
		if err != nil {
			fmt.Fprintf(os.Stderr, "Terraform apply failed: %v\n", err)
			os.Exit(1)
		}
	},
}

var tfDestroyCmd = &cobra.Command{
	Use:   "destroy [site]",
	Short: "Run terraform destroy for a specific site",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		requirements.CheckRequirements()
		site := args[0]
		err := execTerraformWithSDK(site, func(tf *tfexec.Terraform) error {
			return tf.Destroy(context.Background())
		})
		if err != nil {
			fmt.Fprintf(os.Stderr, "Terraform destroy failed: %v\n", err)
			os.Exit(1)
		}
	},
}

var tfOutputCmd = &cobra.Command{
	Use:   "output [site]",
	Short: "Show Terraform outputs for a specific site",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		requirements.CheckRequirements()
		site := args[0]
		err := execTerraformWithSDK(site, func(tf *tfexec.Terraform) error {
			outputs, err := tf.Output(context.Background())
			if err != nil {
				return err
			}
			for k, v := range outputs {
				fmt.Printf("%s = %s\n", k, v.Value)
			}
			return nil
		})
		if err != nil {
			fmt.Fprintf(os.Stderr, "Terraform output failed: %v\n", err)
			os.Exit(1)
		}
	},
}

var tfShowCmd = &cobra.Command{
	Use:   "show [site]",
	Short: "Show Terraform state for a specific site",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		requirements.CheckRequirements()
		site := args[0]
		err := execTerraformWithSDK(site, func(tf *tfexec.Terraform) error {
			output, err := tf.Show(context.Background())
			if err != nil {
				return err
			}
			fmt.Println(output)
			return nil
		})
		if err != nil {
			fmt.Fprintf(os.Stderr, "Terraform show failed: %v\n", err)
			os.Exit(1)
		}
	},
}

var tfRefreshCmd = &cobra.Command{
	Use:   "refresh [site]",
	Short: "Refresh Terraform state for a specific site",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		requirements.CheckRequirements()
		site := args[0]
		err := execTerraformWithSDK(site, func(tf *tfexec.Terraform) error {
			return tf.Refresh(context.Background())
		})
		if err != nil {
			fmt.Fprintf(os.Stderr, "Terraform refresh failed: %v\n", err)
			os.Exit(1)
		}
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

// ------------------------
// Helpers
// ------------------------

func execTerraformWithSDK(site string, action func(*tfexec.Terraform) error) error {
	root := os.Getenv("TERRABUTLER_ROOT")
	if root == "" {
		return fmt.Errorf("TERRABUTLER_ROOT is not set")
	}

	env := GetSelectedEnv()
	if env == "" {
		return fmt.Errorf("No environment selected. Please use `terrabutler env select [name]` first.")
	}

	sitePath := filepath.Join(root, "environments", env, site)
	if _, err := os.Stat(sitePath); os.IsNotExist(err) {
		return fmt.Errorf("Environment directory does not exist: %s", sitePath)
	}

	tf, err := tfexec.NewTerraform(sitePath, "terraform")
	if err != nil {
		return fmt.Errorf("Failed to initialize Terraform SDK: %v", err)
	}

	if err := action(tf); err != nil {
		return fmt.Errorf("Error running Terraform command: %v", err)
	}

	logger.Log.Infof("Terraform command executed successfully for site '%s' in environment '%s'", site, env)
	return nil
}

func GetSelectedEnv() string {
	data, err := os.ReadFile(".terrabutler_env")
	if err != nil {
		return ""
	}
	return strings.TrimSpace(string(data))
}

// ------------------------
// Init: register subcommands
// ------------------------

func init() {
	TfCmd.AddCommand(tfInitCmd)
	TfCmd.AddCommand(tfApplyCmd)
	TfCmd.AddCommand(tfDestroyCmd)
	TfCmd.AddCommand(tfOutputCmd)
	TfCmd.AddCommand(tfShowCmd)
	TfCmd.AddCommand(tfRefreshCmd)
	TfCmd.AddCommand(tfGenVarsCmd)
}
