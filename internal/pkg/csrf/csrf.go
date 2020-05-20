package csrf

import (
	"github.com/gorilla/csrf"
	"net/http"
)

const DefaultTokenHeaderName = "X-Csrf-Token"

type Handlers struct {
	TokenHeaderName string
}

func NewHandlers(csrfTokenHeaderName string) Handlers {
	return Handlers{TokenHeaderName: csrfTokenHeaderName}
}

func (h Handlers) GetNewCsrfTokenHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Expose-Headers", h.TokenHeaderName)
	w.Header().Set(h.TokenHeaderName, csrf.Token(r))
}
