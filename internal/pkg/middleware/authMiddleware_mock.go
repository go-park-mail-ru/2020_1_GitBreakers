package middleware

import (
	"context"
	"net/http"
)

func AuthMiddlewareMock(next http.HandlerFunc, userid int) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		ctx = context.WithValue(ctx, "UserID", userid)
		next(w, r.WithContext(ctx))
	}
}
