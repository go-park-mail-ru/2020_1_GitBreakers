package middleware

import (
	"github.com/go-park-mail-ru/2020_1_GitBreakers/pkg/logger"
	"net/http"
)

func CreateCheckAuthMiddleware(logger logger.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.Context().Value(UserIdContextValue) == nil {
				logger.HttpInfo(r.Context(), "user unauthorized", http.StatusUnauthorized)
				http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
			} else {
				next.ServeHTTP(w, r)
			}
		})
	}
}
