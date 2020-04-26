package middleware

import (
	"github.com/go-park-mail-ru/2020_1_GitBreakers/pkg/middleware"
	"net/http"
)

func JsonContentTypeMiddleware(next http.Handler) http.Handler {
	return middleware.CreateHeadersMiddleware(map[string]string{
		"Content-Type": "application/json",
	})(next)
}
