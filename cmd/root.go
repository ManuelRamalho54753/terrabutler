package cmd

import (
	"fmt"
	"os"
	"terrabutler/requirements"
	"terrabutler/settings"

	"github.com/spf13/cobra"
)

var RootCmd = &cobra.Command{
	Use:   "terrabutler",
	Short: "Terrabutler CLI",
	Long:  "Terrabutler is a CLI tool for managing environments and Terraform commands.",
}

// Executa o CLI e valida settings no início
func Execute() {
	validateSettings()
	requirements.CheckRequirements()

	if err := RootCmd.Execute(); err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}
}

// Verifica se a root está definida e valida o settings.yml
func validateSettings() {
	root := os.Getenv("TERRABUTLER_ROOT")
	if root == "" {
		fmt.Println("Variável TERRABUTLER_ROOT não está definida.")
		fmt.Println("Usa: export TERRABUTLER_ROOT=/caminho/para/repositorio")
		os.Exit(1)
	}

	// Validações de settings
	_, err := settings.LoadSettings(fmt.Sprintf("%s/configs/settings.yml", root))
	if err != nil {
		fmt.Println("Erro ao carregar ficheiro settings.yml:")
		fmt.Println("  ", err)
		os.Exit(1)
	}
}
