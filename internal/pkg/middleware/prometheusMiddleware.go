package middleware

import (
	"github.com/go-park-mail-ru/2020_1_GitBreakers/monitoring"
	customHttp "github.com/go-park-mail-ru/2020_1_GitBreakers/pkg/http"
	"github.com/prometheus/client_golang/prometheus"
	"net/http"
	"strconv"
)

func PrometheusMetricsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		newResponseWriter := &customHttp.ResponseWriter{ResponseWriter: w}
		path := r.URL.Path

		defer monitoring.Hits.WithLabelValues(strconv.Itoa(newResponseWriter.GetStatusCode()), path).Inc()

		timer := prometheus.NewTimer(monitoring.RequestDuration.With(
			prometheus.Labels{"path": path, "method": r.Method},
		))
		defer timer.ObserveDuration()

		next.ServeHTTP(newResponseWriter, r)
	})
}
