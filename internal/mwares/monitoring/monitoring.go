package monitoring

import (
	"github.com/labstack/echo/v4"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type Monitoring struct {
	Hits     *prometheus.CounterVec
	Duration *prometheus.HistogramVec
}

func NewMonitoring(server *echo.Echo) *Monitoring {
	hits := prometheus.NewCounterVec(prometheus.CounterOpts{
		Name: "hits",
	}, []string{"status", "path", "method"})

	duration := prometheus.NewHistogramVec(prometheus.HistogramOpts{
		Name: "duration",
	}, []string{"status", "path", "method"})

	var monitoring = &Monitoring{
		Hits:     hits,
		Duration: duration,
	}

	prometheus.MustRegister(monitoring.Hits, monitoring.Duration)
	server.GET("/metrics", echo.WrapHandler(promhttp.Handler()))
	return monitoring
}
