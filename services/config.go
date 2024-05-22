package services

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

const (
	//WebServiceUrl    string = "https://app.dbvisor.pro"
	WebServiceUrl    string = "https://db-manager.bridge2.digital"
	WebServiceApiUrl string = "api"
	AppName          string = "db-manager"
	ServiceToken     string = "SERVICE_TOKEN"
	EnvFileName      string = ".env.json"
)

type Config struct {
	ServiceToken string `json:"token"`
	Workspace    string `json:"workspace"`
	KeyFile      string `json:"key_file"`
}

type ConfigDataService interface {
	ConfigData()
}

// Config file operations
func EnvFilePath() string {
	return EnvFileName
}

func ConfigData(userData map[string]string) []Config {
	data := []Config{
		{
			ServiceToken: userData["token"],
			Workspace:    userData["workspace"],
			KeyFile:      userData["keyName"],
		},
	}

	return data
}

func IsEnvFileExist(msgSupress bool) bool {
	var result bool = true

	configDir, errDir := CurrentAppDir()
	if errDir != nil {
		fmt.Printf("Cannot get current APP directory: %W.\n", errDir)
		return false
	}

	_, err := os.ReadFile(configDir + "/" + EnvFileName)
	if err != nil {
		if !msgSupress {
			fmt.Printf("Env file not found. Please run: %s install.\n", AppName)
		}
		result = false
	}

	return result
}

func CreateEnvFile(config []Config) {
	configDir, errDir := CurrentAppDir()
	if errDir != nil {
		fmt.Printf("Cannot get current APP directory: %W.\n", errDir)
		return
	}

	file, err := os.Create(configDir + "/" + EnvFileName)
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

func WriteEnvFile(config []Config) {
	/* configDir, errDir := currentAppDir()
	if errDir != nil {
		fmt.Printf("Cannot get current APP directory: %W", errDir)
		return
	} */

	CreateEnvFile(config)
}

//End Config file operations

// Key file operations
func IsKeyFileExist(keyname string) bool {
	var result bool = true

	configDir, errDir := CurrentAppDir()
	if errDir != nil {
		fmt.Printf("Cannot get current APP directory: %W.\n", errDir)
		return false
	}

	_, err := os.ReadFile(configDir + "/" + keyname + ".pub")
	if err != nil {
		fmt.Printf("Key %s.pub file not found. A %s.pub key has been created.\n", keyname, keyname)
		result = false
	}

	return result
}

func CreateKeyPubFile(keyname string) string {
	configDir, errDir := CurrentAppDir()
	if errDir != nil {
		fmt.Printf("Cannot get current APP directory: %W.\n", errDir)
		return ""
	}

	keyFileName := keyname + ".pub"

	file, err := os.Create(configDir + "/" + keyFileName)
	if err != nil {
		fmt.Println("Cannot create key file:", err)
		return ""
	}

	defer file.Close()

	return keyFileName
}

func WriteKeyPubFile(keyData string, keyFileName string) string {
	configDir, errDir := CurrentAppDir()
	if errDir != nil {
		fmt.Printf("Cannot get current APP directory: %W.\n", errDir)
		return ""
	}

	data := []byte(keyData)

	keyFileName = keyFileName + ".pub"

	err := os.WriteFile(configDir+"/"+keyFileName, data, 0664)
	if err != nil {
		fmt.Println("Cannot write key file:", err)
	}

	return keyFileName
}

//End Key file operations

// API login_check
func WebServiceAuthUrl() string {
	return fmt.Sprintf("%s/%s/%s", WebServiceUrl, WebServiceApiUrl, "login_check")
}

// API profile
func WebServiceProfileUrl() string {
	return fmt.Sprintf("%s/%s/%s", WebServiceUrl, WebServiceApiUrl, "profile")
}

// API database list
func WebServiceDatabaseListUrl() string {
	return fmt.Sprintf("%s/%s/%s", WebServiceUrl, WebServiceApiUrl, "databases")
}

// API database dump
func WebServiceDatabaseDumpUrl() string {
	return fmt.Sprintf("%s/%s/%s", WebServiceUrl, WebServiceApiUrl, "database_dumps")
}

func CurrentAppDir() (string, error) {
	ex, err := os.Executable()

	dir := filepath.Dir(ex)
	if err != nil {
		//fmt.Errorf("Cannot get current APP directory: %W", err)
		return "", fmt.Errorf("can not get current app directory: %W", err)
	}

	return dir, err
}
