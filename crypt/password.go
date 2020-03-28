package crypt

import (
	"bytes"
	"encoding/json"
	"github.com/yamamushi/kmud-2020/config"
	"github.com/yamamushi/kmud-2020/types"
	"io/ioutil"
	"log"
	"net/http"
)

func GetAuthToken(username string, password string, conf *config.Config) (types.AuthResponse, bool) {
	sha := Sha256Sum(password)
	hashedPass := string(sha)

	jsonData := map[string]string{"secret": conf.Crypt.AccountManagerSecret, "username": username, "hashedpass": hashedPass}
	jsonValue, _ := json.Marshal(jsonData)
	response, err := http.Post("http://"+conf.Cluster.AccountManagerHostname+"/auth", "application/json", bytes.NewBuffer(jsonValue))
	if err != nil {
		log.Println("Error: GetAuthToken auth request failed with error: " + err.Error())
		return types.AuthResponse{}, false
	} else {
		data, _ := ioutil.ReadAll(response.Body)
		output := types.AuthResponse{}
		err = json.Unmarshal(data, &output)
		if err != nil {
			log.Println("Error: GetAuthToken unmarshal failed with error: " + err.Error())
		}
		if output.AuthToken == "" {
			return output, false
		}
		return output, true
	}
	return types.AuthResponse{}, true
}
