package middleware

import (
	"github.com/gorilla/csrf"
	"net/http"
)

func SkipCsrfCheckMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		next.ServeHTTP(w, csrf.UnsafeSkipCheck(r))
	})
}
