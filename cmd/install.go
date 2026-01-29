/*
Copyright © 2024 Bridge Digital
*/
package cmd

import (
	"github.com/dbvisor-pro/client/processes/install"
	"github.com/spf13/cobra"
)

// installCmd represents the install command
var installCmd = &cobra.Command{
	Use:   "install",
	Short: "Install application",
	Long:  `Installing a console application.`,
	Run: func(cmd *cobra.Command, args []string) {
		install.Execute()
	},
}

func init() {
	// Install command is disabled - uncomment to enable
	// rootCmd.AddCommand(installCmd)
}
