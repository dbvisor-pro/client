/*
Copyright © 2024 Bridge Digital
*/
package cmd

import (
	"fmt"

	"github.com/dbvisor-pro/client/processes/savekey"
	"github.com/dbvisor-pro/client/services/predefined"
	"github.com/spf13/cobra"
)

// saveKeyCmd represents the saveKey command
var saveKeyCmd = &cobra.Command{
	Use:   "save-key",
	Short: "Add public key",
	Long:  `Creating/editing a PEM public key.`,
	Run: func(cmd *cobra.Command, args []string) {
		var result string = savekey.Execute(false, "")

		if result != "" {
			fmt.Println(predefined.BuildOk("The public key has been saved successfully"))
		}
	},
}

func init() {
	rootCmd.AddCommand(saveKeyCmd)
}
