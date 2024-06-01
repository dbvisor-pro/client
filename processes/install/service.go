package install

import (
	"fmt"
	"os"
	"path/filepath"

	"gitea.bridge.digital/bridgedigital/db-manager-client-cli-go/services"
	"gitea.bridge.digital/bridgedigital/db-manager-client-cli-go/services/envfile"
)

func Execute() {

	if envfile.IsEnvFileExist(true) {
		fmt.Println("Application has already installed")
		return
	} else {
		addToBash()
	}

	fmt.Println("The application has been installed successfully")
}

func addToBash() {
	configDir, errDir := services.CurrentAppDir()
	if errDir != nil {
		fmt.Printf("Cannot get current APP directory: %W.\n", errDir)
		return
	}

	command := fmt.Sprintf("export PATH=\"%s/bin:$PATH\" \n", configDir)
	bashProfileCandidates := []string{".bashrc", ".bash_profile"}
	homeDir := os.Getenv("HOME")

	for _, bashProfileCandidate := range bashProfileCandidates {
		candidateFilePath := filepath.Join(homeDir, bashProfileCandidate)

		if _, err := os.Stat(candidateFilePath); err == nil {
			file, err := os.OpenFile(candidateFilePath, os.O_APPEND|os.O_WRONLY, 0644)

			if err != nil {
				fmt.Printf("Error opening file %s: %v", candidateFilePath, err)
				continue
			}

			defer file.Close()

			if _, err := fmt.Fprintln(file, command); err != nil {
				fmt.Printf("Error writing to file %s: %v", candidateFilePath, err)
				return
			}
		}
	}
}
