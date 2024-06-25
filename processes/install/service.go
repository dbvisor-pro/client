/*
Copyright © 2024 Bridge Digital
*/
package install

import (
	"fmt"
	"os"
	"os/exec"

	"gitea.bridge.digital/bridgedigital/db-manager-client-cli-go/services"
	"gitea.bridge.digital/bridgedigital/db-manager-client-cli-go/services/predefined"
)

const DestinationPath = "/usr/local/bin/"

func Execute() {
	_, err := os.Stat(DestinationPath + services.AppName)
	if err == nil {
		fmt.Println(predefined.BuildOk("Application has already installed"))
		return
	}

	configDir, errDir := services.CurrentAppDir()
	if errDir != nil {
		fmt.Printf(predefined.BuildError("Cannot get current APP directory: %W.\n"), errDir)
		return
	}

	var sourceAppPath string = configDir + "/bin/" + services.AppName

	_, errApp := os.Stat(sourceAppPath)
	if errApp != nil {
		fmt.Println(predefined.BuildError("Executable file missing"))
		return
	}

	errInst := createAppLink(sourceAppPath)
	if errInst != nil {
		fmt.Printf(predefined.BuildError("Error executing command: %v\n"), err)
		return
	}

	fmt.Println(predefined.BuildOk("The application has been installed successfully"))
}

func createAppLink(sourcePath string) error {
	cmd := exec.Command("sudo", "ln", "-s", sourcePath, DestinationPath)
	cmd.Env = os.Environ()

	_, err := cmd.CombinedOutput()
	if err != nil {
		return err
	}

	return nil
}
