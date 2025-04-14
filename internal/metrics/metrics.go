package metrics

import (
	"github.com/dankru/Api_gateway_v2/internal/cache"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/rs/zerolog/log"
	"net/http"
	"time"
)

var (
	HttpStatusMetric = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_responses_total",
			Help: "Count of http responses, labeled by status code and method",
		},
		[]string{"status", "method"},
	)
)

func InitMetrics(port string, cache *cache.CacheDecorator) {
	prometheus.MustRegister(HttpStatusMetric)
	http.Handle("/metrics", promhttp.Handler())
	go func() {
		log.Info().Msgf("starting metrics server on port: %s", port)
		server := &http.Server{
			Addr:         ":" + port,
			ReadTimeout:  5 * time.Second,
			WriteTimeout: 10 * time.Second,
		}
		if err := server.ListenAndServe(); err != nil {
			log.Fatal().Err(err).Msgf("Failed to start metrics server on port: %s", port)
		}
	}()
}
