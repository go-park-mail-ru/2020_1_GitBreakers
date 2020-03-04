package main

import (
	"./form"
	"./repository"
	"encoding/json"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"golang.org/x/crypto/bcrypt"
	"net/http"
)

//TODO: передавать ключ через .env файл

func main() {
	//роутер только если к нам приходит json
	r := mux.NewRouter().Headers("Content-Type", "application/json").Subrouter()
	r.HandleFunc("/signup", notAuthMware(signupHandler)).Methods(http.MethodPost)
	r.HandleFunc("/login", notAuthMware(loginHandler)).Methods(http.MethodPost)
	r.HandleFunc("/logout", checkAuthMware(logoutHandler)).Methods(http.MethodPost)

	r.HandleFunc("/repo", checkAuthMware(repository.NewRepository)).Methods(http.MethodPost)
	r.HandleFunc("/repo", repository.GetRepository).Methods(http.MethodGet)
	http.Handle("/", r)

	headersOk := handlers.AllowedHeaders([]string{"*"})
	credentailsOK := handlers.AllowCredentials()
	methodsOk := handlers.AllowedMethods([]string{"GET", "HEAD", "POST", "PUT", "OPTIONS", "DELETE"})
	http.ListenAndServe(":8080", handlers.CORS(methodsOk, headersOk, credentailsOK)(r))
}

//При успех 200 код и json login email(string string)
//При ошибке 400 500 код и json error(string)
func signupHandler(w http.ResponseWriter, r *http.Request) {
	signupForm := &form.SignupForm{}
	err := json.NewDecoder(r.Body).Decode(signupForm) //в signupForm будут данные из json
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"error": "Invalid json must be field login, email, password"}`))
		return
	}

	login := signupForm.Login
	password := signupForm.Password
	email := signupForm.Email

	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	//добавили юзера в бд
	signupForm.Password = hash
	form.UserSlice[login] = *signupForm

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	session, _ := form.StoreSession.Get(r, "session_id")
	session.Values["login"] = login //привязали к сессии login юзера
	session.Save(r, w)              //поставит session_id
	//запишет json в ответ
	json.NewEncoder(w).Encode(map[string]string{"info": "You are registred", "login": login, "email": email})
}

//set cookie сделает
func loginHandler(w http.ResponseWriter, r *http.Request) {
	loginForm := &form.LoginForm{}
	err := json.NewDecoder(r.Body).Decode(loginForm)
	login := loginForm.Login
	password := loginForm.Password

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"error": "Invalid json must be field login, password"}`))
		return
	}

	hash := form.UserSlice[login].Password
	err = bcrypt.CompareHashAndPassword(hash, []byte(password)) //сравнили хеши паролей

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`{"error": "Incorrect login or password"}`))
		return
	}

	session, _ := form.StoreSession.Get(r, "session_id")
	session.Values["login"] = login //привязали к сессии login юзера
	session.Save(r, w)              //поставит session_id

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"info": "You are authorized"})

}

//сделает cookie просроченными
func logoutHandler(w http.ResponseWriter, r *http.Request) {
	session, _ := form.StoreSession.Get(r, "session_id")
	//делает куки просроченными
	session.Options.MaxAge = -1
	session.Save(r, w)
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"info": "You are unauthorized"})
}

//требует что юзер был авторизован
func checkAuthMware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		session, _ := form.StoreSession.Get(r, "session_id") //тащит сессию из реквеста
		if session.IsNew {
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(map[string]string{"error": "You must to be logged in"})
			return
		}
		//добавили нужные headers
		//ставим для cors
		if r.Header.Get("Origin") != "" {
			w.Header().Set("Access-Control-Allow-Origin", r.Header.Get("Origin"))
		}
		w.Header().Set("Content-Type", "application/json")
		next.ServeHTTP(w, r)
	}
}

//требует, чтобы юзер был не авторизован
func notAuthMware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		session, _ := form.StoreSession.Get(r, "session_id")
		if !session.IsNew {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]string{"error": "You must to log out before this action"})
			return
		}
		next.ServeHTTP(w, r)
	}
}
