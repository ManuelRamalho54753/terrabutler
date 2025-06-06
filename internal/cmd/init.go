package cmd

import (
	"github.com/montblu/terrabutler/internal/inception"
	"github.com/spf13/cobra"
)

// Command `terrabutler init`
var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize the inception site",
	Run: func(cmd *cobra.Command, args []string) {
		inception.InitInception()
	},
}

func init() {
	RootCmd.AddCommand(initCmd)
}
