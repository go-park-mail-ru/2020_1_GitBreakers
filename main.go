package main

import (
	"github.com/gorilla/mux"
	"net/http"
)

type User struct {
	ID       uint   `json:"id"`
	Username string `json:"username"`
	Password string `json:"password"`
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

func main() {
	r := mux.NewRouter()

	api := NewMyHandler()
	r.HandleFunc("/", api.Root)
	r.HandleFunc("/new/repository", api.NewRepository)

	http.ListenAndServe(":8080", r)
}
