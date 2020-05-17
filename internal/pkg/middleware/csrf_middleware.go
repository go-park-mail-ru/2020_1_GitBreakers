package middleware

import (
	ownCsrf "github.com/go-park-mail-ru/2020_1_GitBreakers/internal/pkg/csrf"
	"github.com/gorilla/csrf"
	"net/http"
)

const CookieName string = "codehub_csrf"

func CreateCSRFMiddleware(secret []byte, trustedOrigins []string,
	secure bool, site csrf.SameSiteMode, maxAge int64) func(http.Handler) http.Handler {
	return csrf.Protect(
		secret,
		csrf.RequestHeader(ownCsrf.TokenHeaderName),
		csrf.TrustedOrigins(trustedOrigins),
		csrf.Secure(secure),
		csrf.CookieName(CookieName),
		csrf.SameSite(site),
		csrf.HttpOnly(true),
		csrf.MaxAge(int(maxAge)),
		csrf.Path("/"),
	)
}
