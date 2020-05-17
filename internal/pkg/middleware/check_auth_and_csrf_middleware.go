package middleware

import (
	"github.com/go-park-mail-ru/2020_1_GitBreakers/pkg/logger"
	"github.com/gorilla/csrf"
	"net/http"
)

func CreateCheckAuthAndCSRFMiddleware(secret []byte, trustedOrigins []string,
	secure bool, site csrf.SameSiteMode, maxAge int64, logger logger.Logger) func(http.Handler) http.Handler {

	checkAuthMiddleware := CreateCheckAuthMiddleware(logger)
	csrfMiddleware := CreateCSRFMiddleware(
		secret,
		trustedOrigins,
		secure,
		site,
		maxAge,
	)

	return func(next http.Handler) http.Handler {
		return checkAuthMiddleware(csrfMiddleware(next))
	}
}
