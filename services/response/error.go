package response

import (
	"encoding/json"
	"fmt"
)

type WrongData struct {
	Code int    `json:"code"`
	Msg  string `json:"message"`
}

func WrongResponseObserver(data []byte) {
	var wrongData WrongData

	dbErr := json.Unmarshal([]byte(data), &wrongData)
	if dbErr != nil {
		fmt.Println("Error decoding from json:", dbErr)
		return
	}

	if wrongData.Code == 401 && wrongData.Msg == "Invalid JWT Token" {
		fmt.Println("Your access token is expired, please, log in before continuing.")
		return
	} else {
		fmt.Printf("Code: %d. Message: %s \n", wrongData.Code, wrongData.Msg)
		return
	}
}
