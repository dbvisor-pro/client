package services

import (
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
	PubKeyExt        string = ".pem"
)

// API login_check url
func WebServiceAuthUrl() string {
	return fmt.Sprintf("%s/%s/%s", WebServiceUrl, WebServiceApiUrl, "login_check")
}

// API profile url
func WebServiceProfileUrl() string {
	return fmt.Sprintf("%s/%s/%s", WebServiceUrl, WebServiceApiUrl, "profile")
}

// API database list url
func WebServiceDatabaseListUrl() string {
	return fmt.Sprintf("%s/%s/%s", WebServiceUrl, WebServiceApiUrl, "databases")
}

// API database dump url
func WebServiceDatabaseDumpUrl() string {
	return fmt.Sprintf("%s/%s/%s", WebServiceUrl, WebServiceApiUrl, "database_dumps")
}

// API database download link url
func WebServiceDownLoadLinkUrl() string {
	return fmt.Sprintf("%s/%s/%s", WebServiceUrl, WebServiceApiUrl, "get_download_link")
}

func CurrentAppDir() (string, error) {
	ex, err := os.Executable()

	dir := filepath.Dir(ex)
	if err != nil {
		return "", fmt.Errorf("can not get current app directory: %W", err)
	}

	return dir, err
}
