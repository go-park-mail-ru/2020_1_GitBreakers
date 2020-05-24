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

func CreateAccessLogMiddleware(log logger.Logger) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			rand.Seed(time.Now().UnixNano())
			ctx := r.Context()
			nBig, err := crand.Int(crand.Reader, big.NewInt(99999999))
			if err != nil {
				log.HttpLogCallerError(ctx, ctx, err)
			}
			requestID := fmt.Sprintf("%016x", nBig.String())[:6]

			log.StartRequest(*r, requestID)

			start := time.Now()
			ctx = context.WithValue(
				ctx,
				log.GetRequestIdKey(),
				requestID,
			)

			next.ServeHTTP(w, r.WithContext(ctx))

			log.EndRequest(start, ctx)
		})
	}
}
