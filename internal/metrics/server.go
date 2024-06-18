package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	HttpRequestsTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_requests_total",
			Help: "Количество HTTP-запросов.",
		},
		[]string{"path"},
	)
)

func IncRequestsCounter(path string) {
	HttpRequestsTotal.WithLabelValues(path).Inc()
}
