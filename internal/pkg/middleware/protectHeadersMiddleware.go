package middleware

import (
	"github.com/go-park-mail-ru/2020_1_GitBreakers/pkg/middleware"
	"net/http"
)

func ProtectHeadersMiddleware(next http.Handler) http.Handler {
	return middleware.CreateHeadersMiddleware(map[string]string{
		"X-XSS-Protection":       "1; mode=block",
		"X-Content-Type-Options": "nosniff;",
	})(next)
}
