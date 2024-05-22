package login

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	saveKey "gitea.bridge.digital/bridgedigital/db-manager-client-cli-go/processes/savekey"
	"gitea.bridge.digital/bridgedigital/db-manager-client-cli-go/services"
	"gitea.bridge.digital/bridgedigital/db-manager-client-cli-go/services/response"
	workspacePac "gitea.bridge.digital/bridgedigital/db-manager-client-cli-go/services/workspace"
	"github.com/howeyc/gopass"
	"github.com/spf13/cobra"
)

func Execute(cmd *cobra.Command) string {

	credentials := map[string]string{
		"username": "",
		"password": "",
	}

	var token, username, workspace, keyFileName string

	//reader := bufio.NewReader(os.Stdin)

USERNAME:
	fmt.Println("Username: ")
	fmt.Scanln(&username)
	//username, _ := reader.ReadString('\n')

	if len(strings.TrimSpace(username)) == 0 {
		fmt.Println("The Username cannot be empty")
		goto USERNAME
	} else {
		credentials["username"] = username
	}

PASSWORD:
	fmt.Println("Password: ")

	password, _ := gopass.GetPasswdMasked()

	if len(strings.TrimSpace(string(password))) == 0 {
		fmt.Println("The Password cannot be empty")
		goto PASSWORD
	} else {
		credentials["password"] = string(password)
	}

	token = jwtToken(credentials)
	workspace = workspacePac.Workspace(token)

	//fmt.Println(workspace)

	if len(token) == 0 || len(workspace) == 0 {
		return ""
	}

	configData := map[string]string{
		"token":     token,
		"workspace": workspace,
		"keyName":   "",
	}

	if !services.IsEnvFileExist(false) {
		services.CreateEnvFile(services.ConfigData(configData))
	} else {
		keyFileName = saveKey.Execute()

		configData["keyName"] = keyFileName

		services.WriteEnvFile(services.ConfigData(configData))
	}

	return "You logged in successfully"
}

// Get token from server
func jwtToken(credentials map[string]string) string {
	credsInJson, err := json.Marshal(credentials)
	if err != nil {
		fmt.Println("Error encoding to json:", err)
		return ""
	}

	req, err := http.NewRequest("POST", services.WebServiceAuthUrl(), bytes.NewBuffer(credsInJson))
	if err != nil {
		fmt.Println("Something wrong with POST request data:", err)
		return ""
	}

	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}

	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Invalid credentials:", err)
		return ""
	}

	if resp == nil {
		return fmt.Sprint(http.StatusBadRequest)
	}

	defer resp.Body.Close()

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error:", err)
		return ""
	}

	var configData map[string]string

	configErr := json.Unmarshal(data, &configData)
	if configErr != nil {
		response.WrongResponseObserver(data)
		return ""
	}

	if len(configData["token"]) > 0 {
		return configData["token"]
	}

	return ""
}
