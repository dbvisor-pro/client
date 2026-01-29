/*
Copyright © 2024 Bridge Digital
*/
package config

import (
	"fmt"
	"os"
	"strings"

	"gitea.bridge.digital/bridgedigital/db-manager-client-cli-go/services"
	"gitea.bridge.digital/bridgedigital/db-manager-client-cli-go/services/envfile"
	"gitea.bridge.digital/bridgedigital/db-manager-client-cli-go/services/predefined"
)

func Execute(dumpPath string, serviceUrl string) {
	hasDumpPath := len(strings.TrimSpace(dumpPath)) > 0
	hasServiceUrl := len(strings.TrimSpace(serviceUrl)) > 0

	if !hasDumpPath && !hasServiceUrl {
		fmt.Println(predefined.BuildError("Please provide at least one option: --dump-path (-d) or --url (-u)"))
		return
	}

	if !envfile.IsEnvFileExist(false) {
		return
	}

	savedConfig, err := envfile.ReadEnvFile()
	if err != nil {
		fmt.Println(predefined.BuildError("Error:"), err)
		return
	}

	if hasDumpPath {
		errDumpPath := createDumpPathDir(dumpPath)
		if errDumpPath != nil {
			fmt.Println(errDumpPath)
			return
		}
		savedConfig.DownloadDumpPath = dumpPath
		fmt.Println(predefined.BuildOk("You have set the default path for database dumps"))
	}

	if hasServiceUrl {
		savedConfig.ServiceUrl = strings.TrimSpace(serviceUrl)
		services.ResetConfigCache()
		fmt.Println(predefined.BuildOk("You have set the service URL to: " + serviceUrl))
	}

	envfile.WriteEnvFile(savedConfig)
}

func createDumpPathDir(dumpPath string) error {
	if err := os.MkdirAll(dumpPath, os.ModePerm); err != nil {
		return fmt.Errorf(predefined.BuildError("can not get entered dump directory: %W"), err)
	}

	return nil
}
