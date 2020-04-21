package main

import (
	"fmt"
	"github.com/gorilla/mux"
	"google.golang.org/grpc"
	"log"
	"net/http"
)

func main() {
	grcpConn, err := grpc.Dial(
		"127.0.0.1:8080",
		grpc.WithInsecure(),
	)
	if err != nil {
		log.Fatal("cant connect to grpc")
	}
	defer grcpConn.Close()

	r := mux.NewRouter()
	r.HandleFunc("/session", signupHandler).Methods(http.MethodPost)
	r.HandleFunc("/session", loginHandler).Methods(http.MethodGet)
	r.HandleFunc("/session", logoutHandler).Methods(http.MethodDelete)
	fmt.Println("starting server at :8080")
	http.ListenAndServe(":8080", nil)
}
