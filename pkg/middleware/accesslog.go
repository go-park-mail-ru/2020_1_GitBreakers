package middleware

import (
	"context"
	crand "crypto/rand"
	"fmt"
	"github.com/go-park-mail-ru/2020_1_GitBreakers/pkg/logger"
	"math/big"
	"math/rand"
	"net/http"
	"time"
)

func CreateAccessLogMiddleware(requestIDKey int, log logger.Logger) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			rand.Seed(time.Now().UnixNano())
			nBig, _ := crand.Int(crand.Reader, big.NewInt(99999999))
			requestID := fmt.Sprintf("%016x", nBig.String())[:6]
			log.StartRequest(*r, requestID)
			start := time.Now()
			ctx := r.Context()

			reqIDkey := requestIDKey //for linter
			ctx = context.WithValue(ctx,
				reqIDkey,
				requestID,
			)
			next.ServeHTTP(w, r.WithContext(ctx))
			log.EndRequest(start, ctx)
		})
	}
}
