package main

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"io/ioutil"
	"net/http"
)

type User struct {
	ID       uint   `json:"id"`
	Username string `json:"username"`
	Password string `json:"password"`
}
type Repository struct {
	Name          string            `json:"name"`
	Fullname      string            `json:"full_name"`
	Private       bool              `json:"private"`
	Owner         map[string]string `json:"owner"`
	Description   string            `json:"description"`
	Fork          bool              `json:"fork"`
	DefaultBranch string            `json:"default_branch"`
}

type MyHandler struct {
	sessions map[string]uint
	users    map[string]*User
}

func NewMyHandler() *MyHandler {
	return &MyHandler{
		sessions: make(map[string]uint, 10),
		users: map[string]*User{
			"rvasily": {1, "rvasily", "love"},
		},
	}
}

//копирнул из репока с лекциями, оставил временно
func (api *MyHandler) Root(w http.ResponseWriter, r *http.Request) {
	authorized := false
	session, err := r.Cookie("session_id")
	if err == nil && session != nil {
		_, authorized = api.sessions[session.Value]
	}

	if authorized {
		w.Write([]byte("autrorized"))
	} else {
		w.Write([]byte("not autrorized"))
	}
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

func main() {
	r := mux.NewRouter()

	api := NewMyHandler()
	r.HandleFunc("/", api.Root)
	r.HandleFunc("/new/repository", api.NewRepository)

	http.ListenAndServe(":8080", r)
}
