package cmd

import (
	"github.com/spf13/cobra"
)

// Comand to manage environments
var envCmd = &cobra.Command{
	Use:   "env",
	Short: "Manage environments",
}

func init() {
	RootCmd.AddCommand(envCmd)
}
