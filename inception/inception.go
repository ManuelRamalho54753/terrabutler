package inception

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"terrabutler/logger"
	"terrabutler/settings"
	"terrabutler/utils"
)

// InceptionInitNeeded checks if the inception site is initialized; exits if not
func InceptionInitNeeded() {
	if !inceptionInitCheck() {
		fmt.Println("\nInception is not initialized.")
		fmt.Println("Please run `terrabutler init` to initialize it.")
		os.Exit(1)
	}
}

// InitInception initializes the inception site with terraform init and sets default env
func InitInception() {
	root := os.Getenv("TERRABUTLER_ROOT")
	if root == "" {
		logger.Log.Fatal("TERRABUTLER_ROOT environment variable is not set")
	}

	inceptionPath := filepath.Join(root, "site_inception")
	terraformDir := filepath.Join(inceptionPath, ".terraform")
	environmentFile := filepath.Join(terraformDir, "environment")

	// If already initialized, exit early
	if utils.PathExists(environmentFile) {
		logger.Log.Infof("Inception site already initialized at %s", inceptionPath)
		return
	}

	s := settings.GetSettings()
	org := s.General.Organization
	env := s.Environments.Default.Name
	backendFile := filepath.Join(root, "backends", fmt.Sprintf("%s-%s-inception.tfvars", org, env))

	// Ensure .terraform folder exists before writing environment file
	if err := os.MkdirAll(terraformDir, 0755); err != nil {
		logger.Log.Fatalf("Error creating .terraform directory: %v", err)
	}

	logger.Log.Infof("Running `terraform init` in %s", inceptionPath)
	cmd := exec.Command("terraform", "init", "-backend-config", backendFile)
	cmd.Dir = inceptionPath
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		logger.Log.Fatalf("Error running terraform init: %v", err)
	}

	if err := os.WriteFile(environmentFile, []byte(env), 0644); err != nil {
		logger.Log.Fatalf("Error writing .terraform/environment: %v", err)
	}

	logger.Log.Infof("Inception initialized successfully in: %s", inceptionPath)
}

// inceptionInitCheck returns true if .terraform/environment exists
func inceptionInitCheck() bool {
	root := os.Getenv("TERRABUTLER_ROOT")
	if root == "" {
		logger.Log.Fatal("TERRABUTLER_ROOT environment variable is not set")
	}

	environmentFile := filepath.Join(root, "site_inception", ".terraform", "environment")
	return utils.PathExists(environmentFile)
}
