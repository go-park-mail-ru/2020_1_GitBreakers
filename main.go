package main

import (
	"github.com/gorilla/mux"
	"html/template"
	"net/http"
)

var templates *template.Template

func main() {
	templates = template.Must(template.ParseGlob("public/*.html"))
	r := mux.NewRouter()

	r.HandleFunc("/new/repository", NewRepository).Methods("POST")
	r.HandleFunc("/", indexHandler).Methods("GET")
	fs := http.FileServer(http.Dir("./public"))
	r.PathPrefix("/public/").Handler(http.StripPrefix("/public", fs))
	http.Handle("/", r)
	http.ListenAndServe(":8080", r)
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	templates.ExecuteTemplate(w, "index.html", nil)
}
