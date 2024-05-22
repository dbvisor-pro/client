package savekey

import (
	"fmt"
	"strings"

	"gitea.bridge.digital/bridgedigital/db-manager-client-cli-go/services"
	"github.com/AlecAivazis/survey/v2"
)

func Execute() string {
	var keyName, keyData string

	/* type KeyData struct {
		Name string
		Key  string
	} */

KEYNAME:
	fmt.Println("Enter public key name: ")
	fmt.Scanln(&keyName)

	keyName = strings.TrimSpace(keyName)

	if len(keyName) == 0 {
		fmt.Println("The key name cannot be empty")
		goto KEYNAME
	} else {
		if !services.IsKeyFileExist(keyName) {
			services.CreateKeyPubFile(keyName)
		} else {
			var (
				options        = []string{"Yes", "No"}
				selectedOption string
			)

			prompt := &survey.Select{
				Message: "Key file is already exists. Do you want to override existing file?",
				Options: options,
			}

			survey.AskOne(prompt, &selectedOption)

			if selectedOption == "No" {
				return keyName + ".pub"
			}
		}
	}

KEY:
	fmt.Println("Enter public key: ")
	fmt.Scanln(&keyData)

	keyData = strings.TrimSpace(keyData)

	if len(keyData) == 0 {
		fmt.Println("The key cannot be empty")
		goto KEY
	} else {
		return services.WriteKeyPubFile(keyData, keyName)
	}
}
