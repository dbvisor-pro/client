/*
Copyright © 2024 Bridge Digital
*/
package cmd

import (
	"github.com/spf13/cobra"
)

// listCmd represents the list command
var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List commands",
	Long:  `Display a list of all commands.`,
	Run: func(cmd *cobra.Command, args []string) {
		rootCmd.Help()
	},
}

func init() {
	rootCmd.AddCommand(listCmd)
}
