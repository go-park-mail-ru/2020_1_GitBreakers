package csrf

import (
	"github.com/gorilla/csrf"
	"net/http"
)

const TokenHeaderName = "X-Csrf-Token"

func GetNewCsrfTokenHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Expose-Headers", TokenHeaderName)
	w.Header().Set(TokenHeaderName, csrf.Token(r))
}
