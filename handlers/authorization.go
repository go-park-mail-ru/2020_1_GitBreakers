package handlers

import (
	"../models"
	"encoding/json"
	"fmt"
	uuid "github.com/satori/go.uuid"
	"golang.org/x/crypto/bcrypt"
	"log"
	"net/http"
	"strings"
)

type StoresContext models.StoresContext

type AuthError int

func (err AuthError) Error() string {
	return fmt.Sprintf("error code: %d", int(err))
}

func (context *StoresContext) TryProcessSessionId(r *http.Request) (string, AuthError) {
	sessionId, err := r.Cookie("session_id")
	if err == http.ErrNoCookie || sessionId == nil {
		return "", AuthError(http.StatusUnauthorized)
	}
	if sessionId.Value == "" {
		return "", AuthError(http.StatusUnprocessableEntity)
	}
	if in := context.SessionStore.HaveSession(sessionId.Value); in {
		return sessionId.Value, 0
	}
	return "", AuthError(http.StatusUnauthorized)
}

func (context *StoresContext) TryProcessUser(w http.ResponseWriter, r *http.Request) (*models.User, error) {
	receivedUser := new(models.User)
	if err := json.NewDecoder(r.Body).Decode(receivedUser); err != nil {
		http.Error(w, `Invalid data`, http.StatusNotAcceptable)
		return receivedUser, fmt.Errorf(`invalid data`)
	}

	if len(strings.Fields(receivedUser.Login)) != 1 {
		_ = json.NewEncoder(w).Encode(&models.Result{
			Err: "invalid login",
		})
		return receivedUser, fmt.Errorf(`invalid login`)
	}

	userPasswordHash, err := bcrypt.GenerateFromPassword([]byte(receivedUser.Password), bcrypt.DefaultCost)
	if err != nil {
		http.Error(w, `Cant hash password`, http.StatusInternalServerError)
		err = fmt.Errorf("[ERROR] cant create user hash from password=%s, err=%s",
			receivedUser.Password, err.Error())
		return receivedUser, err
	}
	receivedUser.Password = string(userPasswordHash)
	return receivedUser, nil
}

// Register user
func (context *StoresContext) Register(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
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

	_ = json.NewEncoder(w).Encode(&models.Result{
		Body: map[string]string{"status": "ok"},
	})
}

func (context *StoresContext) Login(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "application/json")

	var err error
	var user *models.User

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

	_ = json.NewEncoder(w).Encode(&models.Result{
		Body: map[string]string{"status": "ok"},
	})
}

func (context *StoresContext) Logout(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "application/json")

	var err error
	var sessionId string

	if sessionId, err = context.TryProcessSessionId(r); err != nil {
		http.Error(w, `Session id not exist`, http.StatusBadRequest)
		return
	}

	context.SessionStore.DeleteSession(sessionId)

	_ = json.NewEncoder(w).Encode(&models.Result{
		Body: map[string]string{"status": "ok"},
	})
}
