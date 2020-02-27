package models

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
)

type Repository struct {
	Name          string `json:"name"`
	Fullname      string `json:"full_name"`
	Private       bool   `json:"private"`
	Description   string `json:"description"`
	Fork          bool   `json:"fork"`
	OwnerId       int    `json:"ownerid"`
	DefaultBranch string `json:"default_branch"`
}

/*
Принимает json формата{
	name: "ffsfsfsf",
    full_name: "Deiklov/ffsfsfsf",
    private: false,
    "ownerid": 252525525,
    description: "sfsfsf",
    fork: false,
    default_branch: "master"
}
*/
var repositorySlice = []Repository{
	{
		Name:          "kekmda",
		Fullname:      "Deiklov/kekmda",
		Private:       false,
		Description:   "somerepository",
		Fork:          false,
		OwnerId:       141421414,
		DefaultBranch: "master",
	},
	{
		Name:          "somerepo",
		Fullname:      "Dmitry/somerepo",
		Private:       false,
		Description:   "some text",
		Fork:          false,
		OwnerId:       1342424,
		DefaultBranch: "dev",
	},
}

func NewRepository(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", r.Header.Get("Origin"))
	w.Header().Set("Access-Control-Allow-Credentials", "true")
	w.Header().Set("Content-Type", "application/json")
	body, err := ioutil.ReadAll(r.Body)

	defer r.Body.Close() // важный пункт!
	if err != nil {
		http.Error(w, `read body error`, http.StatusBadRequest)
		return
	}
	reposData := &Repository{}
	err = json.Unmarshal(body, reposData)
	if err != nil {
		http.Error(w, `invalid json`, http.StatusBadRequest)
		return
	}
	for _, v := range repositorySlice {
		if v.Fullname == reposData.Fullname {
			http.Error(w, `Already exsist this repo`, http.StatusBadRequest)
			return
		}
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
		http.Error(w, `invalid get params`, http.StatusBadRequest)
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
