package handlers

import (
	"../models"
	"encoding/json"
	"fmt"
	uuid "github.com/satori/go.uuid"
	"golang.org/x/crypto/sha3"
	"log"
	"net/http"
	"time"
)

type StoresContext models.StoresContext

func (context *StoresContext) TryProcessSessionId(r *http.Request) (*http.Cookie, models.CommonError) {
	sessionId, err := r.Cookie("session_id")
	if err == http.ErrNoCookie || sessionId == nil {
		return sessionId, models.NewModelError(``, http.StatusUnauthorized)
	}
	if sessionId.Value == "" {
		return sessionId, models.NewModelError(``, http.StatusUnprocessableEntity)
	}
	if in := context.SessionStore.HaveSession(sessionId.Value); in {
		return sessionId, nil
	}
	return sessionId, models.NewModelError(``, http.StatusUnauthorized)
}

func (context *StoresContext) TryProcessUser(w http.ResponseWriter, r *http.Request) (*models.User, error) {
	receivedUser := new(models.User)
	if err := json.NewDecoder(r.Body).Decode(receivedUser); err != nil {
		http.Error(w, `Invalid data`, http.StatusNotAcceptable)
		return receivedUser, err
	}

	sha3Hash := sha3.New256()
	_, err := sha3Hash.Write([]byte(receivedUser.Password))
	if err != nil {
		http.Error(w, `Cant hash password`, http.StatusInternalServerError)
		err = fmt.Errorf("[ERROR] cant create user hash from password=%s, err=%s",
			receivedUser.Password, err.Error())
		return receivedUser, err
	}

	receivedUser.Password = string(sha3Hash.Sum(nil))
	return receivedUser, nil
}

// Signup user
func (context *StoresContext) Signup(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", r.Header.Get("Origin"))
	w.Header().Set("Access-Control-Allow-Credentials", "true")
	w.Header().Set("Content-Type", "application/json")

	var err error
	var user *models.User

	if user, err = context.TryProcessUser(w, r); err != nil {
		log.Println(err.Error())
		return
	}

	sessionId, err := uuid.NewV4()
	if err != nil {
		http.Error(w, `Cant generate session UUID`, http.StatusInternalServerError)
		log.Printf("[ERROR] Cant generate session UUID, err=%s\n", err.Error())
		return
	}

	if err := context.UsersStore.StoreUser(*user); err != nil {
		_ = json.NewEncoder(w).Encode(&models.Result{
			Err: err.Error(),
		})
		return
	}
	_ = context.SessionStore.StoreSession(models.Session{
		ID:        sessionId.String(),
		UserLogin: user.Login,
	})

	expiration := time.Now().Add(10 * time.Hour)
	cookie := http.Cookie{
		Name:    "session_id",
		Value:   sessionId.String(),
		Expires: expiration,
		Path:    "/",
	}
	http.SetCookie(w, &cookie)

	if encErr := json.NewEncoder(w).Encode(&models.Result{Body: map[string]string{"status": "ok"}}); encErr != nil {
		http.Error(w, `cant encode user to json:`+encErr.Error(), http.StatusInternalServerError)
	}
}

func (context *StoresContext) Login(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", r.Header.Get("Origin"))
	w.Header().Set("Access-Control-Allow-Credentials", "true")
	w.Header().Set("Content-Type", "application/json")

	var err error
	var user *models.User

	if _, err = context.TryProcessSessionId(r); err == nil {
		http.Error(w, `Session already exist`, http.StatusBadRequest)
		return
	}

	if user, err = context.TryProcessUser(w, r); err != nil {
		log.Println(err.Error())
		return
	}

	dbUser, err := context.UsersStore.GetUser(user.Login)
	if err != nil || dbUser.Password != user.Password {
		_ = json.NewEncoder(w).Encode(&models.Result{
			Err: `incorrect login or password`,
		})
		return
	}

	sessionId, err := uuid.NewV4()
	if err != nil {
		http.Error(w, `Cant generate session UUID`, http.StatusInternalServerError)
		log.Printf("[ERROR] Cant generate session UUID, err=%s\n", err.Error())
		return
	}

	_ = context.SessionStore.StoreSession(models.Session{
		ID:        sessionId.String(),
		UserLogin: user.Login,
	})

	expiration := time.Now().Add(10 * time.Hour)
	cookie := http.Cookie{
		Name:    "session_id",
		Value:   sessionId.String(),
		Expires: expiration,
		Path:    "/",
	}
	http.SetCookie(w, &cookie)

	if encErr := json.NewEncoder(w).Encode(&models.Result{Body: map[string]string{"status": "ok"}}); encErr != nil {
		http.Error(w, `cant encode user to json:`+encErr.Error(), http.StatusInternalServerError)
	}
}

func (context *StoresContext) Logout(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "application/json")

	var err error
	var sessionId *http.Cookie

	if sessionId, err = context.TryProcessSessionId(r); err != nil {
		http.Error(w, `Session id not exist`, http.StatusBadRequest)
		return
	}

	sessionId.Expires = time.Now().AddDate(0, 0, -1)

	context.SessionStore.DeleteSession(sessionId.Value)
	http.SetCookie(w, sessionId)

	if encErr := json.NewEncoder(w).Encode(&models.Result{Body: map[string]string{"status": "ok"}}); encErr != nil {
		http.Error(w, `cant encode user to json:`+encErr.Error(), http.StatusInternalServerError)
	}
}
