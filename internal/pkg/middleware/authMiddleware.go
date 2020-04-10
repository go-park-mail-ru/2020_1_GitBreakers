package middleware

import (
	"context"
	"github.com/go-park-mail-ru/2020_1_GitBreakers/internal/pkg/session"
	"github.com/go-park-mail-ru/2020_1_GitBreakers/internal/pkg/user"
	"net/http"
)

//содержит интерфейсы
type Middleware struct {
	SessDeliv session.SessDelivery
	UCUser    user.UCUser
}

func (Mdware *Middleware) AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		cookie, err := r.Cookie("session_id")
		if err != nil {
			next.ServeHTTP(w, r)
			return
		}
		sessModel, err := Mdware.SessDeliv.GetBySessID(cookie.Value)
		if err != nil {
			next.ServeHTTP(w, r)
			return
		}
		ctx = context.WithValue(ctx, "UserID", sessModel.UserID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
