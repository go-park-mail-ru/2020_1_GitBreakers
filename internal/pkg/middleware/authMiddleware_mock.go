package middleware

import (
	"context"
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
