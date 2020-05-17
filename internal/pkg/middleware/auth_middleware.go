package middleware

import (
	"context"
	"github.com/go-park-mail-ru/2020_1_GitBreakers/internal/app/clients"
	"net/http"
)

const (
	UserIdContextValue    = "UserID"
	SessionIdContextValue = "session_id"
)

//содержит интерфейсы
type Middleware struct {
	SessDeliv *clients.SessClient
}

func (Mdware *Middleware) AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		cookie, err := r.Cookie(SessionIdContextValue)
		if err != nil {
			next.ServeHTTP(w, r)
			return
		}
		sessModel, err := Mdware.SessDeliv.GetSess(cookie.Value)
		if err != nil {
			next.ServeHTTP(w, r)
			return
		}
		ctx = context.WithValue(ctx, UserIdContextValue, sessModel.UserID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
