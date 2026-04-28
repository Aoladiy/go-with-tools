package metrics

import (
	"fmt"
	"net/http"
	"time"

	"github.com/Aoladiy/go-with-tools/internal/config"
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var RequestsTotal = prometheus.NewCounterVec(
	prometheus.CounterOpts{
		Name: "http_requests_total",
		Help: "number of processed requests",
	},
	[]string{"method", "path", "status"},
)

var ErrorsTotal = prometheus.NewCounterVec(
	prometheus.CounterOpts{
		Name: "http_errors_total",
		Help: "number of errors",
	},
	[]string{"method", "path", "status"},
)

var RequestDurations = prometheus.NewHistogramVec(
	prometheus.HistogramOpts{
		Name:    "http_request_duration_seconds",
		Help:    "request durations",
		Buckets: []float64{0.001, 0.005, 0.01, 0.05, 0.1, 0.2},
	},
	[]string{"method", "path", "status"},
)

func register() {
	prometheus.MustRegister(
		RequestsTotal,
		ErrorsTotal,
		RequestDurations,
	)
}

func NewServer(c config.Config) *http.Server {
	register()

	r := gin.New()

	r.Use(gin.Recovery())
	r.GET("/metrics", gin.WrapH(promhttp.Handler()))

	server := &http.Server{
		Addr:         fmt.Sprintf("%s:%d", c.MetricsHost, c.MetricsPort),
		Handler:      r,
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	return server
}
