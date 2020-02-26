package models

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
)

type Owner struct {
	Login     string `json:"login"`
	Id        int    `json:"id"`
	SiteAdmin bool   `json:"site_admin"`
}
type Repository struct {
	Name          string `json:"name"`
	Fullname      string `json:"full_name"`
	Private       bool   `json:"private"`
	Owner         Owner  `json:"owner"`
	Description   string `json:"description"`
	Fork          bool   `json:"fork"`
	DefaultBranch string `json:"default_branch"`
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
var repositorySlice = make([]Repository, 0)

func NewRepository(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	body, err := ioutil.ReadAll(r.Body)

	defer r.Body.Close() // важный пункт!
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	reposData := &Repository{}
	err = json.Unmarshal(body, reposData)
	if err != nil {
		w.Write([]byte(`{"status": "400", "error": "Invalid json"}`))
		return
	}
	repositorySlice = append(repositorySlice, *reposData) //save to db

	w.Write([]byte(`{"status":"200"}`))
	return
}
func GetRepository(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", r.Header.Get("Origin"))
	w.Header().Set("Access-Control-Allow-Credentials", "true")
	w.Header().Set("Content-Type", "application/json")
	name := r.FormValue("name")
	if name == "" {
		w.Write([]byte(`{"status": "400", "error": "Invalid get params"}`))
		return
	}
	for _, v := range repositorySlice {
		if v.Name == name {
			data, _ := json.Marshal(v)
			w.Write(data)
			return
		}
	}

	return
}
