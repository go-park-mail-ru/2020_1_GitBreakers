package main

import (
	"./routes"
	"net/http"
)

func main() {
	r := routes.NewRouter()
	http.Handle("/", r)

	http.ListenAndServe(":8080", r)
}

