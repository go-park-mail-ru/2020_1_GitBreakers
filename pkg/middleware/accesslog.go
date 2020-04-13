package middleware

import (
	"context"
	"fmt"
	"github.com/go-park-mail-ru/2020_1_GitBreakers/pkg/logger"
	"math/rand"
	"net/http"
	"time"
)

func CreateAccessLogMiddleware(requestIDKey int, log logger.Logger) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			rand.Seed(time.Now().UnixNano())
			requestID := fmt.Sprintf("%016x", rand.Int())[:6]
			log.StartRequest(*r, requestID)
			start := time.Now()
			ctx := r.Context()
			ctx = context.WithValue(ctx,
				requestIDKey,
				requestID,
			)
			next.ServeHTTP(w, r.WithContext(ctx))
			log.EndRequest(start, ctx)
		})
	}
}
