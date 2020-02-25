package routes

import (
	"../models"
	"github.com/gorilla/mux"
	"net/http"
)

func NewRouter() *mux.Router {
	r := mux.NewRouter()

	r.HandleFunc("/new/repository", models.NewRepository).Methods("POST")
	r.HandleFunc("/settings/profile", models.UpdateProfile).Methods("PUT")
	r.HandleFunc("/profile/{login}", models.GetProfile).Methods("GET")
	r.HandleFunc("/repository/{login}", models.GetRepositoryList).Methods("GET")

	fs := http.FileServer(http.Dir("./public"))
	r.PathPrefix("/public/").Handler(http.StripPrefix("/public", fs))

	return r
}
