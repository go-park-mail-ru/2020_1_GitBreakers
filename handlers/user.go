package handlers

import (
	"../models"
	"encoding/json"
	"log"
	"net/http"
)

func (context *StoresContext) GetUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", r.Header.Get("Origin"))
	w.Header().Set("Access-Control-Allow-Credentials", "true")
	var err models.CommonError
	var sessionIdCookie *http.Cookie

	if sessionIdCookie, err = context.TryProcessSessionId(r); err != nil {
		http.Error(w, `session not exist`, err.GetCode())
		return
	}

	session, _ := context.SessionStore.GetSession(sessionIdCookie.Value)
	user, err := context.UsersStore.GetUser(session.UserLogin)
	if err != nil {
		log.Fatalf("[FATAL] %q\n", err)
	}

	user.Password = ""

	if encErr := json.NewEncoder(w).Encode(&user); encErr != nil {
		http.Error(w, `cant encode user to json:`+encErr.Error(), http.StatusInternalServerError)
	}
}
