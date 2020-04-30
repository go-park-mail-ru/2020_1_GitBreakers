package middleware

import (
	"github.com/go-park-mail-ru/2020_1_GitBreakers/pkg/logger"
	"github.com/pkg/errors"
	"net/http"
)

func CreatePanicMiddleware(logger logger.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			defer func() {
				if e := recover(); e != nil {
					http.Error(w, "internal server error", http.StatusInternalServerError)
					if err := e.(error); err != nil {
						logger.LogError(errors.Cause(err), "panic catched in panic middleware")
					} else {
						logger.LogError(errors.Cause(err), "panic catched in panic middleware")
					}
				}
			}()
			next.ServeHTTP(w, r)
		})
	}
}
