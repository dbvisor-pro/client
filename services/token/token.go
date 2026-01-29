/*
Copyright © 2024 Bridge Digital
*/
package token

import (
	"encoding/json"
	"fmt"

	"github.com/dbvisor-pro/client/services"
	"github.com/dbvisor-pro/client/services/predefined"
	"github.com/dbvisor-pro/client/services/request"
	"github.com/dbvisor-pro/client/services/response"
)

// Get token from server
func JwtToken(credentials map[string]string) string {
	credsInJson, err := json.Marshal(credentials)
	if err != nil {
		fmt.Println(predefined.BuildError("Error encoding to json:"), err)
		return ""
	}

	data, err := request.CreatePostRequest(credsInJson, services.WebServiceAuthUrl(), nil)
	if err != nil {
		fmt.Println(err)
		return ""
	}

	var configData map[string]string

	configErr := json.Unmarshal(data, &configData)
	if configErr != nil {
		err := response.WrongResponseObserver(data)
		if err != nil {
			fmt.Println(err)
			return ""
		}
	}

	if len(configData["token"]) > 0 {
		return configData["token"]
	}

	return ""
}
