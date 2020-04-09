package csrf

import (
	"github.com/gorilla/csrf"
	"net/http"
)

const TokenHeaderName = "X-Csrf-Token"

func GetNewCsrfToken(w http.ResponseWriter, r *http.Request) {
	if r.Context().Value("UserID") == nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	w.Header().Set(TokenHeaderName, csrf.Token(r))
}
