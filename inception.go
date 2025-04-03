package main

import (
	"fmt"
	"os"
	"path/filepath"
	"terrabutler/logger"
	"terrabutler/utils"
)

// InitInception initializes the inception site directory
func InitInception() {
	root := os.Getenv("TERRABUTLER_ROOT")
	if root == "" {
		logger.Log.Fatal("TERRABUTLER_ROOT environment variable is not set")
	}

	inceptionPath := filepath.Join(root, "site_inception")
	if utils.PathExists(inceptionPath) {
		logger.Log.Infof("Inception site already initialized at %s", inceptionPath)
		return
	}

	err := os.MkdirAll(inceptionPath, 0755)
	if err != nil {
		logger.Log.Fatalf("Failed to create inception site directory: %v", err)
	}

	logger.Log.Infof("Inception site initialized at: %s", inceptionPath)
}

// InceptionInitNeeded checks if the inception site directory exists and prompts initialization if not
func InceptionInitNeeded() {
	root := os.Getenv("TERRABUTLER_ROOT")
	if root == "" {
		logger.Log.Fatal("TERRABUTLER_ROOT environment variable is not set")
	}

	inceptionPath := filepath.Join(root, "site_inception")
	if !utils.PathExists(inceptionPath) {
		fmt.Println("Inception site is not initialized. Run `terrabutler init` to initialize it.")
		os.Exit(1)
	}
}
