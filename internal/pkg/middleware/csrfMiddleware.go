package middleware

import (
	ownCsrf "github.com/go-park-mail-ru/2020_1_GitBreakers/internal/pkg/csrf"
	"github.com/gorilla/csrf"
	"net/http"
)

const CookieName string = "codehub_csrf"

func CreateCsrfMiddleware(secret []byte, trustedOrigins []string, secure bool, maxAge int64) (func(http.Handler) http.Handler) {
	return csrf.Protect(secret,
		csrf.RequestHeader(ownCsrf.TokenHeaderName),
		csrf.TrustedOrigins(trustedOrigins),
		csrf.Secure(secure),
		csrf.CookieName(CookieName),
		csrf.SameSite(csrf.SameSiteNoneMode),
		csrf.HttpOnly(true),
		csrf.MaxAge(int(maxAge)),
		csrf.Path("/"))
}
