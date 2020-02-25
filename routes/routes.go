package routes

import (
	"github.com/gorilla/mux"
	"net/http"
	"../models"
)

func NewRouter() *mux.Router {
	r := mux.NewRouter()

	r.HandleFunc("/new/repository", models.NewRepository).Methods("POST")
	fs := http.FileServer(http.Dir("./public"))
	r.PathPrefix("/public/").Handler(http.StripPrefix("/public", fs))

	return r
}
