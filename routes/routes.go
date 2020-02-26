package routes

import (
	"../models"
	"github.com/gorilla/mux"
	"net/http"
)

func NewRouter() *mux.Router {
	r := mux.NewRouter()

	r.HandleFunc("/new/repository", models.NewRepository).Methods("POST")
	r.HandleFunc("/settings/profile", models.UpdateProfile).Methods("POST")
	r.HandleFunc("/profile/{login}", models.GetProfile).Methods("GET")
	r.HandleFunc("/repository/{login}", models.GetRepositoryList).Methods("GET")
	r.HandleFunc("/settings/avatar", models.UploadAvatar).Methods("POST")

	fs := http.FileServer(http.Dir("./public"))
	r.PathPrefix("/public/").Handler(http.StripPrefix("/public", fs))

	staticHandler := http.FileServer(http.Dir("./static"))
	r.PathPrefix("/static/").Handler(http.StripPrefix("/static", staticHandler))


	r.HandleFunc("/repository", models.GetRepository).Methods(http.MethodGet, http.MethodOptions)
	r.Use(mux.CORSMethodMiddleware(r))

	return r
}
