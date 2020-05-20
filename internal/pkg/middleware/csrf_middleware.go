package middleware

import (
	"github.com/gorilla/csrf"
	"net/http"
)

const DefaultCSRFCookieName string = "codehub_csrf"

func CreateCSRFMiddleware(secret []byte, trustedOrigins []string,
	cookieName, requestHeaderName string,
	secure bool, site csrf.SameSiteMode, maxAge int64) func(http.Handler) http.Handler {
	return csrf.Protect(
		secret,
		csrf.RequestHeader(requestHeaderName),
		csrf.TrustedOrigins(trustedOrigins),
		csrf.Secure(secure),
		csrf.CookieName(cookieName),
		csrf.SameSite(site),
		csrf.HttpOnly(true),
		csrf.MaxAge(int(maxAge)),
		csrf.Path("/"),
	)
}
