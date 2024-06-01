package workspace

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"gitea.bridge.digital/bridgedigital/db-manager-client-cli-go/services"
	"gitea.bridge.digital/bridgedigital/db-manager-client-cli-go/services/response"
	"github.com/AlecAivazis/survey/v2"
)

func Workspace(token string) string {
	req, err := http.NewRequest("GET", services.WebServiceProfileUrl(), nil)
	if err != nil {
		fmt.Println("Something wrong with GET request data:", err)
		return ""
	}

	req.Header.Set("Accept", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)

	client := &http.Client{}

	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Invalid token:", err)
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

	type Data struct {
		Identifier string   `json:"identifier"`
		Workspaces []string `json:"workspaces"`
	}

	// Unmarshal the JSON data into the struct
	var workspaceData Data
	allWorkspaces := map[int]string{}

	wsErr := json.Unmarshal([]byte(data), &workspaceData)
	if wsErr != nil {
		response.WrongResponseObserver(data)
		return ""
	}

	if len(workspaceData.Workspaces) > 1 {
		for k, workspace := range workspaceData.Workspaces {
			allWorkspaces[k] = workspace
		}

		var selectedWorkspace int

		prompt := &survey.Select{
			Message: "Select workspace:",
			Options: workspaceData.Workspaces,
		}

		survey.AskOne(prompt, &selectedWorkspace)

		return allWorkspaces[selectedWorkspace]
	} else if len(workspaceData.Workspaces) == 1 {
		return workspaceData.Workspaces[0]
	} else {
		fmt.Println("You don't assigned to any workspace")
	}

	return ""
}
