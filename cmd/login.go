/*
Copyright © 2024 Bridge Digital
*/
package cmd

import (
	"fmt"

	"github.com/dbvisor-pro/client/processes/login"
	"github.com/dbvisor-pro/client/services/predefined"
	"github.com/spf13/cobra"
)

// loginCmd represents the login command
var loginCmd = &cobra.Command{
	Use:   "login",
	Short: "Authorize to service",
	Long:  `Creating/updating a token and creating/editing a public key in the configuration file required for downloading database dumps.`,
	Run: func(cmd *cobra.Command, args []string) {
		var result string = login.Execute(cmd)
		fmt.Println(predefined.BuildOk(result))
	},
}

func init() {
	rootCmd.AddCommand(loginCmd)
}
