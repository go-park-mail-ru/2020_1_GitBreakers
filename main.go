package main

import (
	"./routes"
	"html/template"
	"net/http"
)

var templates *template.Template

func main() {
	templates = template.Must(template.ParseGlob("public/*.html"))
	r := routes.NewRouter()
	r.HandleFunc("/", indexHandler).Methods("GET")
	http.Handle("/", r)
	http.ListenAndServe(":8080", r)
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	templates.ExecuteTemplate(w, "index.html", nil)
}
