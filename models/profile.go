package models

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"io/ioutil"
	"net/http"
	"strings"
)

var profilMap = map[string]Profile{
	"antonelagin": Profile{
		Login:     "AntonElagin",
		Avatar:    "https://avatars.mds.yandex.net/get-pdb/1964870/cfaa9b42-d74b-40f7-93cc-1d777bb5654c/s1200",
		RegDate:   "12.02.1999",
		Followers: 1,
		Following: 2,
	},
	"keklol": Profile{
		Login:     "keeeekLooool",
		Avatar:    "https://avatars.mds.yandex.net/get-pdb/1964870/cfaa9b42-d74b-40f7-93cc-1d777bb5654c/s1200",
		RegDate:   "12.02.1999",
		Followers: 1,
		Following: 150,
	}}

type Result struct {
	Body interface{} `json:"body,omitempty"`
	Err  string      `json:"err,omitempty"`
}

type Profile struct {
	Login     string `json:"login"`
	Avatar    string `json:"avatar"`
	RegDate   string `json:"registrationDate"`
	Followers uint   `json:"followers"`
	Following uint   `json:"following"`
}

func UpdateProfile(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "application/json")
	prof := Profile{}
	body, errB := ioutil.ReadAll(r.Body)

	if errB != nil {
		json.NewEncoder(w).Encode(&Result{
			Err: "Server Error",
		})
		return
	}

	err := json.Unmarshal(body, &prof)
	if err != nil {
		json.NewEncoder(w).Encode(&Result{
			Err: "bad body",
		})
		return
	}

	json.NewEncoder(w).Encode(&Result{Body: prof})
}

func GetProfile(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "application/json")
	vars := mux.Vars(r)
	login := vars["login"]
	fmt.Println(login, strings.ToLower(login))
	profile, ok := profilMap[strings.ToLower(login)]
	if !ok {
		json.NewEncoder(w).Encode(&Result{
			Err: "Not found"})
		return
	}

	json.NewEncoder(w).Encode(&Result{Body: profile})
}
