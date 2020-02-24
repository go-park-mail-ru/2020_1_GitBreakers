package main

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
)

type Repository struct {
	Name          string            `json:"name"`
	Fullname      string            `json:"full_name"`
	Private       bool              `json:"private"`
	Owner         map[string]string `json:"owner"`
	Description   string            `json:"description"`
	Fork          bool              `json:"fork"`
	DefaultBranch string            `json:"default_branch"`
}

/*
Принимает json формата{
	name: "ffsfsfsf",
    full_name: "Deiklov/ffsfsfsf",
    private: false,
    "owner": {
    	login: "Deiklov",
    	id: 48653955,
    	site_admin: false
  		},
    description: "sfsfsf",
    fork: false,
    default_branch: "master"
}
*/
func (api *MyHandler) NewRepository(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close() // важный пункт!
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	reposData := &Repository{}
	json.Unmarshal(body, reposData)
	//save to db
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(`{"status":"200"}`))
}
