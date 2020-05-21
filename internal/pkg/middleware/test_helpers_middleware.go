package middleware

import (
	"context"
	"github.com/go-park-mail-ru/2020_1_GitBreakers/internal/pkg/models"
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
			ctx = context.WithValue(ctx, models.UserIDKey, int64(rand.Intn(max-min)+min))
		}
		next(w, r.WithContext(ctx))
	}
}
func SetMuxVars(next http.HandlerFunc, vars map[string]string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		req := mux.SetURLVars(r, vars)
		next(w, req)
	}
}
