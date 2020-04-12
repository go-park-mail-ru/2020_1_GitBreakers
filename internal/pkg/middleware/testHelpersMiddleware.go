package middleware

import (
	"context"
	"github.com/gorilla/mux"
	"math/rand"
	"net/http"
)

const (
	max = 100
	min = 20
)

func AuthMiddlewareMock(next http.HandlerFunc, isAuthorized bool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		if isAuthorized {
			ctx = context.WithValue(ctx, "UserID", rand.Intn(max-min)+min)
		}
		next(w, r.WithContext(ctx))
	}
}
func SetMuxVars(next http.HandlerFunc, key, value string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		req := mux.SetURLVars(r, map[string]string{key: value})
		next(w, req)
	}
}
