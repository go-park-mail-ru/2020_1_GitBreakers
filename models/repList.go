package models

import (
	"encoding/json"
	// "fmt"
	"github.com/gorilla/mux"
	// "io/ioutil"
	"net/http"
	"strings"
)

var repMap = map[string][]RepositoryForList{
	"antonelagin": []RepositoryForList{
		RepositoryForList{
			Name:        "kekmda",
			Private:     false,
			Description: "lolkekcheburek",
			Stars:       50,
			LastUpdate:  "12.12.2020",
		},
		RepositoryForList{
			Name:        "rep2",
			Private:     false,
			Description: "lolkekcheburek",
			Stars:       50,
			LastUpdate:  "12.12.2020",
		},
		RepositoryForList{
			Name:        "rep12",
			Private:     false,
			Description: "lolkekcheburek",
			Stars:       50,
			LastUpdate:  "12.12.2020",
		},
		RepositoryForList{
			Name:        "rep3",
			Private:     false,
			Description: "lolkekcheburek",
			Stars:       50,
			LastUpdate:  "12.12.2020",
		},
	},
	"keklol": []RepositoryForList{
		RepositoryForList{
			Name:        "rep1",
			Private:     false,
			Description: "lolkekcheburek",
			Stars:       50,
			LastUpdate:  "12.12.2020",
		},
		RepositoryForList{
			Name:        "rep2",
			Private:     false,
			Description: "lolkekcheburek",
			Stars:       50,
			LastUpdate:  "12.12.2020",
		},
	},
}

type RepositoryForList struct {
	Name        string `json:"name"`
	Private     bool   `json:"private"`
	Description string `json:"description"`
	Stars       int    `json:"stars"`
	LastUpdate  string `json:"update"`
}

func GetRepositoryList(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "application/json")
	vars := mux.Vars(r)
	login := vars["login"]

	rep, ok := repMap[strings.ToLower(login)]
	if !ok {
		json.NewEncoder(w).Encode(&Result{
			Err: "Not found"})
		return
	}

	json.NewEncoder(w).Encode(&Result{Body: map[string][]RepositoryForList{"reps": rep}})
}
