package metrics

import (
	"github.com/dankru/Api_gateway_v2/internal/cache"
	"github.com/gofiber/fiber/v2"
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
		[]string{"status", "method", "path"},
	)
	CacheElementCount = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: "cache_element_count",
			Help: "Number of elements in the cache",
		})

	CacheSizeBytes = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: "cache_size_bytes",
			Help: "Size of cache in bytes",
		})
)

func InitMetrics(port string, cache *cache.CacheDecorator) {
	prometheus.MustRegister(HttpStatusMetric)
	prometheus.MustRegister(CacheElementCount)
	prometheus.MustRegister(CacheSizeBytes)

	startCacheMetricsCollector(cache, 1*time.Second)

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

func httpStatusMetricInc(statusCode int, method, path string) {
	HttpStatusMetric.WithLabelValues(http.StatusText(statusCode), method, path).Inc()
}

func PrometheusMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		err := c.Next()

		statusCode := c.Response().StatusCode()
		method := c.Method()
		path := c.Route().Path

		httpStatusMetricInc(statusCode, method, path)

		return err
	}
}

func startCacheMetricsCollector(cache *cache.CacheDecorator, interval time.Duration) {
	go func() {
		ticker := time.NewTicker(interval)
		defer ticker.Stop()

		for {
			select {
			case <-ticker.C:
				CacheElementCount.Set(float64(cache.ElementCount()))
				CacheSizeBytes.Set(float64(cache.SizeBytes()))
			}
		}
	}()
}
