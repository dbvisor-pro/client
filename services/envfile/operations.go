package envfile

import (
	"encoding/json"
	"fmt"
	"os"

	"gitea.bridge.digital/bridgedigital/db-manager-client-cli-go/services"
	"golang.org/x/exp/maps"
)

type Workspace struct {
	ServiceToken string `json:"token"`
	KeyFile      string `json:"key_file"`
}

type ConfigDataService interface {
	ConfigData()
}

func ConfigData(userData map[string]string) map[string]Workspace {
	data := Workspace{
		ServiceToken: userData["token"],
		KeyFile:      userData["keyName"],
	}

	var configData = map[string]Workspace{}

	configData[userData["workspace"]] = data

	return configData
}

func IsEnvFileExist(msgSupress bool) bool {
	var result bool = true

	configDir, errDir := services.CurrentAppDir()
	if errDir != nil {
		fmt.Printf("Cannot get current APP directory: %W.\n", errDir)
		return false
	}

	_, err := os.Stat(configDir + "/" + services.EnvFileName)
	if err != nil {
		if !msgSupress {
			fmt.Printf("Env file not found. Please run: %s install.\n", services.AppName)
		}
		result = false
	}

	return result
}

func CreateEnvFile(config map[string]Workspace) {
	configDir, errDir := services.CurrentAppDir()
	if errDir != nil {
		fmt.Printf("Cannot get current APP directory: %W.\n", errDir)
		return
	}

	file, err := os.Create(configDir + "/" + services.EnvFileName)
	if err != nil {
		fmt.Println("Cannot create file:", err)
		return
	}

	defer file.Close()

	encoder := json.NewEncoder(file)
	err = encoder.Encode(config)
	if err != nil {
		fmt.Println("Cannot write config data to file:", err)
		return
	}
}

func WriteEnvFile(config map[string]Workspace) {
	keys := maps.Keys(config)

	workspace := keys[0]

	if !IsEnvFileExist(true) {
		CreateEnvFile(config)
	}

	configFromFile, err := ReadEnvFile()
	if err != nil {
		fmt.Printf("Error: %s", err)
		return
	}

	workspaceDataFromFile, v := configFromFile[workspace]
	if v {
		if len(config[workspace].KeyFile) > 0 {
			workspaceDataFromFile.KeyFile = config[workspace].KeyFile
		}

		workspaceDataFromFile.ServiceToken = config[workspace].ServiceToken

		configFromFile[workspace] = workspaceDataFromFile
	} else {
		configFromFile[workspace] = config[workspace]
	}

	data, errData := json.Marshal(configFromFile)
	if errData != nil {
		fmt.Println("Cannot encode config data: ", err)
		return
	}

	configDir, errDir := services.CurrentAppDir()
	if errDir != nil {
		fmt.Printf("Cannot get current APP directory: %W.\n", errDir)
		return
	}

	err = os.WriteFile(configDir+"/"+services.EnvFileName, data, 0644)
	if err != nil {
		fmt.Println("Cannot write to env file:", err)
		return
	}
}

func ReadEnvFile() (map[string]Workspace, error) {
	if !IsEnvFileExist(true) {
		return nil, fmt.Errorf("env file not found. Please use the login command to update it")
	}

	configDir, errDir := services.CurrentAppDir()
	if errDir != nil {
		return nil, fmt.Errorf("cannot get current APP directory: %W", errDir)
	}

	file, err := os.ReadFile(configDir + "/" + services.EnvFileName)
	if err != nil {
		return nil, fmt.Errorf("env file is not readable: %W", errDir)
	}

	config := make(map[string]Workspace)

	err = json.Unmarshal(file, &config)
	if err != nil {
		return nil, fmt.Errorf("the settings record is not readable: %W", errDir)
	}

	return config, nil
}
