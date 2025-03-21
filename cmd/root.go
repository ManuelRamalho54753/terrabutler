package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

// Comando principal do CLI
var RootCmd = &cobra.Command{
	Use:   "terrabutler",
	Short: "Terrabutler CLI",
	Long:  "Terrabutler is a CLI tool for managing environments and Terraform commands.",
}

// Função que inicia o CLI
func Execute() {
	if err := RootCmd.Execute(); err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}
}
