package download

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"gitea.bridge.digital/bridgedigital/db-manager-client-cli-go/services"
	"gitea.bridge.digital/bridgedigital/db-manager-client-cli-go/services/response"
	"gitea.bridge.digital/bridgedigital/db-manager-client-cli-go/services/token"
	"gitea.bridge.digital/bridgedigital/db-manager-client-cli-go/services/workspace"
	"github.com/AlecAivazis/survey/v2"
)

const (
	DefaultDumpDBName string = "backup.sql"
	StatusReady       string = "ready"
)

func Execute(dbUid, dumpUid *string) {

	var defaultDbUid, defaultDumpUid string

	if dbUid == nil {
		defaultDbUid = getDbUid()
		dbUid = &defaultDbUid
	}

	if dumpUid == nil {
		defaultDumpUid = getDumpUid(*dbUid)
		dumpUid = &defaultDumpUid
	}

	if defaultDbUid == "" || defaultDumpUid == "" {
		return
	}

	fmt.Println(defaultDbUid)
	//fmt.Println(defasultDumpUid)
	fmt.Println(dbUid)

	fmt.Println("Downloading.....")
}

func getDbUid() string {
	req, err := http.NewRequest("GET", services.WebServiceDatabaseListUrl(), nil)
	if err != nil {
		fmt.Println("Something wrong with GET request data:", err)
		return ""
	}

	req.Header.Set("Accept", "application/json")
	req.Header.Set("Authorization", "Bearer "+token.Current())

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
		Name string `json:"name"`
		Uid  string `json:"uid"`
	}

	var (
		dbData        []Data
		allDbDataName []string
	)

	dbErr := json.Unmarshal([]byte(data), &dbData)
	if dbErr != nil {
		response.WrongResponseObserver(data)
		return ""
	}

	if len(dbData) > 0 {
		for _, uid := range dbData {
			allDbDataName = append(allDbDataName, uid.Name)
		}

		fmt.Println(dbData)
		fmt.Println(allDbDataName)

		var selectedDb int

		prompt := &survey.Select{
			Message: "Please select database to process with:",
			Options: allDbDataName,
		}

		survey.AskOne(prompt, &selectedDb)

		return dbData[selectedDb].Uid
	} else {
		fmt.Println("Not found active databases")
	}

	return ""
}

func getDumpUid(dbUid string) string {
	var (
		selectedWorkspace string = workspace.Workspace(token.Current())
		requestUrl        string = services.WebServiceDatabaseDumpUrl() + "?db.uid=" + dbUid + "&workspace=" + selectedWorkspace
	)

	req, err := http.NewRequest("GET", requestUrl, nil)
	if err != nil {
		fmt.Println("Something wrong with GET request data:", err)
		return ""
	}

	req.Header.Set("Accept", "application/json")
	req.Header.Set("Authorization", "Bearer "+token.Current())

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
		Uuid string `json:"uuid"`
	}

	var (
		dumps    []Data
		allDumps []string
	)

	dbErr := json.Unmarshal([]byte(data), &dumps)
	if dbErr != nil {
		response.WrongResponseObserver(data)
		return ""
	}

	if len(dumps) > 0 {
		for _, uid := range dumps {
			allDumps = append(allDumps, uid.Uuid)
		}

		//fmt.Println(dumps)
		//fmt.Println(allDumps)

		var selectedDb int

		prompt := &survey.Select{
			Message: "Please select dump to process with:",
			Options: allDumps,
		}

		survey.AskOne(prompt, &selectedDb)

		return dumps[selectedDb].Uuid
	} else {
		fmt.Println("Not found active dumps for selected DB")
	}

	return ""
}
