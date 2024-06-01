package savekey

import (
	"fmt"
	"strings"

	"gitea.bridge.digital/bridgedigital/db-manager-client-cli-go/services"
	"gitea.bridge.digital/bridgedigital/db-manager-client-cli-go/services/envfile"
	"gitea.bridge.digital/bridgedigital/db-manager-client-cli-go/services/keypubfile"
	"github.com/AlecAivazis/survey/v2"
	"golang.org/x/exp/maps"
)

func Execute(isNew bool) string {
	if !isNew {
		reCreate()
		return ""
	}

	var (
		options        = []string{"Yes", "No"}
		selectedOption string
	)

	prompt := &survey.Select{
		Message: "Want to create a public Pem key for downloading?",
		Options: options,
	}

	survey.AskOne(prompt, &selectedOption)

	if selectedOption == "No" {
		return ""
	}

	var keyName, keyData string

	qKeyName := &survey.Question{
		Name:   "Key Name",
		Prompt: &survey.Input{Message: "Enter public key name:"},
		Validate: func(val interface{}) error {
			if str, _ := val.(string); len(strings.TrimSpace(str)) == 0 {
				return fmt.Errorf("the key name cannot be empty")
			}
			return nil
		},
	}

	survey.AskOne(qKeyName.Prompt, &keyName, survey.WithValidator(qKeyName.Validate))

	if !keypubfile.IsKeyFileExist(keyName) {
		keypubfile.CreateKeyPubFile(keyName)
	} else {
		prompt := &survey.Select{
			Message: "Key file is already exists. Do you want to override existing file?",
			Options: options,
		}

		survey.AskOne(prompt, &selectedOption)

		if selectedOption == "No" {
			return keyName + services.PubKeyExt
		}
	}

	qKeyData := &survey.Question{
		Name:   "Key Data",
		Prompt: &survey.Multiline{Message: "Enter public key:"},
		Validate: func(val interface{}) error {
			if str, _ := val.(string); len(strings.TrimSpace(str)) == 0 {
				return fmt.Errorf("the key cannot be empty")
			}
			return nil
		},
	}

	survey.AskOne(qKeyData.Prompt, &keyData, survey.WithValidator(qKeyData.Validate))

	return keypubfile.WriteKeyPubFile(keyData, keyName)
}

// Function for regenerating a key
func reCreate() {
	savedWorkspaces, err := envfile.ReadEnvFile()
	if err != nil {
		fmt.Println(err)
		return
	}

	var (
		selectedWorkspaceIndex int
		savedWorkspacesKeys    []string
		options                = []string{"Yes", "No"}
		selectedOption         string
	)

	savedWorkspacesKeys = maps.Keys(savedWorkspaces)

	promptW := &survey.Select{
		Message: "Select one of your saved workspaces:",
		Options: savedWorkspacesKeys,
	}

	survey.AskOne(promptW, &selectedWorkspaceIndex)

	savedConfigData := savedWorkspaces[savedWorkspacesKeys[selectedWorkspaceIndex]]
	selectedKey := savedConfigData.KeyFile

keyNameAsk:

	var keyName, keyData string

	qKeyNamePrompt := &survey.Input{
		Message: "Enter a new name or leave the field blank if you want to use the old one:",
		Help:    "The default key name is the previous key name",
		Default: selectedKey,
	}

	survey.AskOne(qKeyNamePrompt, &keyName)

	if len(keyName) == 0 {
		return
	}

	if !keypubfile.IsKeyFileExist(keyName) {
		if len(keyName) > 0 {
			keyName = keypubfile.CreateKeyPubFile(keyName)
		} else {
			return
		}
	} else {
		prompt := &survey.Select{
			Message: "Key file is already exists. Do you want to override existing file?",
			Options: options,
		}

		survey.AskOne(prompt, &selectedOption)

		if selectedOption == "No" {
			goto keyNameAsk
		}
	}

	qKeyData := &survey.Question{
		Name:   "Key Data",
		Prompt: &survey.Multiline{Message: "Enter public key:"},
		Validate: func(val interface{}) error {
			if str, _ := val.(string); len(strings.TrimSpace(str)) == 0 {
				return fmt.Errorf("the key cannot be empty")
			}
			return nil
		},
	}

	survey.AskOne(qKeyData.Prompt, &keyData, survey.WithValidator(qKeyData.Validate))

	if len(keyData) == 0 {
		return
	}

	keyName = keypubfile.WriteKeyPubFile(keyData, keyName)

	if envfile.IsEnvFileExist(false) {
		if len(keyName) > 0 {
			configData := map[string]string{
				"token":     savedConfigData.ServiceToken,
				"workspace": savedWorkspacesKeys[selectedWorkspaceIndex],
				"keyName":   keyName,
			}

			envfile.WriteEnvFile(envfile.ConfigData(configData))

			fmt.Println("The public key has been saved successfully")
		}
	}
}
