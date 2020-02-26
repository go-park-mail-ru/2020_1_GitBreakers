package routes

import (
	"../models"
	"github.com/gorilla/mux"
	"net/http"
)

func NewRouter() *mux.Router {
	r := mux.NewRouter()

	r.HandleFunc("/new/repository", models.NewRepository).Methods("POST")
	r.HandleFunc("/repository", models.GetRepository).Methods(http.MethodGet, http.MethodOptions)
	r.Use(mux.CORSMethodMiddleware(r))
	return r
}
