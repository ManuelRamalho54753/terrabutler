package tests

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"text/template"

	"github.com/montblu/terrabutler/internal/logger"
	"github.com/montblu/terrabutler/internal/settings"
	"github.com/montblu/terrabutler/internal/utils"
)

// InceptionInitNeeded checks if the inception site is initialized; exits if not
func InceptionInitNeeded() {
	root := os.Getenv("TERRABUTLER_ROOT")
	if root == "" {
		logger.Log.Fatal("TERRABUTLER_ROOT environment variable is not set")
	}

	envFile := filepath.Join(root, "site_inception", ".terraform", "environment")
	if !utils.PathExists(envFile) {
		fmt.Println("Inception is not initialized.")
		fmt.Println("Please run `terrabutler init` to initialize it.")
		os.Exit(1)
	}
}

// InitInception initializes a site directory with tfvars and environment file
func InitInception(org, site string) {
	root := os.Getenv("TERRABUTLER_ROOT")
	if root == "" {
		logger.Log.Fatal("TERRABUTLER_ROOT environment variable is not set")
	}

	s := settings.GetSettingsFromRoot(root)
	org = s.General.Organization
	env := s.Environments.Default.Name
	region := s.Environments.Default.Region
	profile := s.Environments.Default.ProfileName

	// DEBUG: print the values to confirm they are loaded correctly
	fmt.Println("Org:", org)
	fmt.Println("Env:", env)
	fmt.Println("Region:", region)
	fmt.Println("ProfileName:", profile)

	// fallback protection in case the values are missing
	if org == "" || env == "" || region == "" || profile == "" {
		logger.Log.Fatalf("Missing values for Org: %q, Env: %q, Region: %q, ProfileName: %q", org, env, region, profile)
	}

	sitePath := filepath.Join(root, "environments", site)
	terraformDir := filepath.Join(sitePath, ".terraform")
	environmentFile := filepath.Join(terraformDir, "environment")
	tfvarsFile := filepath.Join(sitePath, "terraform.tfvars")
	envTplPath := filepath.Join(root, "internal", "configs", "templates", "env.tpl")
	defaultTFPath := filepath.Join(root, "internal", "configs", "default_tf_files")

	// Create directories (even if site exists, allow regenerating content)
	if err := os.MkdirAll(terraformDir, 0755); err != nil {
		logger.Log.Fatalf("Failed to create site directories: %v", err)
	}

	// Create .terraform/environment file
	if err := os.WriteFile(environmentFile, []byte(env), 0644); err != nil {
		logger.Log.Fatalf("Failed to write .terraform/environment: %v", err)
	}

	vars := map[string]string{
		"Org":         org,
		"Env":         env,
		"Region":      region,
		"ProfileName": profile,
	}

	// Generate terraform.tfvars from env.tpl
	envTpl, err := template.ParseFiles(envTplPath)
	if err != nil {
		logger.Log.Fatalf("Failed to read env.tpl template: %v", err)
	}
	tfFile, err := os.Create(tfvarsFile)
	if err != nil {
		logger.Log.Fatalf("Failed to create terraform.tfvars: %v", err)
	}
	defer tfFile.Close()
	if err := envTpl.Execute(tfFile, vars); err != nil {
		logger.Log.Fatalf("Failed to generate terraform.tfvars: %v", err)
	}

	// Copy default .tf files (overwrite always)
	entries, err := os.ReadDir(defaultTFPath)
	if err != nil {
		logger.Log.Fatalf("Failed to read default_tf_files: %v", err)
	}

	for _, entry := range entries {
		if !entry.IsDir() {
			src := filepath.Join(defaultTFPath, entry.Name())
			dst := filepath.Join(sitePath, entry.Name())
			data, err := os.ReadFile(src)
			if err != nil {
				logger.Log.Fatalf("Failed to read file %s: %v", src, err)
			}
			if err := os.WriteFile(dst, data, fs.ModePerm); err != nil {
				logger.Log.Fatalf("Failed to write file %s: %v", dst, err)
			}
		}
	}

	logger.Log.Infof("Site %s initialized successfully at: %s", site, sitePath)
}

// InitAllInceptionSites initializes all sites listed in settings.yml
func InitAllInceptionSites() {
	root := os.Getenv("TERRABUTLER_ROOT")
	if root == "" {
		logger.Log.Fatal("TERRABUTLER_ROOT environment variable is not set")
	}

	s := settings.GetSettingsFromRoot(root)
	org := s.General.Organization

	for _, site := range s.Sites.Ordered {
		InitInception(org, site)
	}
}
