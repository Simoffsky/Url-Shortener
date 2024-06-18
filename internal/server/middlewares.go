package server

import (
	"net/http"
	"url-shorter/internal/metrics"
)

func WithMetrics(next http.Handler) http.Handler {
	return increaseRequestsMiddleware(next)
}
func increaseRequestsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		metrics.IncRequestsCounter(r.URL.Path)
		next.ServeHTTP(w, r)
	})
}
