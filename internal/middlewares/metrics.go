package middlewares

import (
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"strconv"
	"time"
)

type MetricsController struct {
	Counter prometheus.Counter
	Hits *prometheus.CounterVec
	Duration *prometheus.HistogramVec

}

func NewMetricsController(router *gin.Engine) *MetricsController {
	counter := prometheus.NewCounter(prometheus.CounterOpts{
		Name: "total_counter",
		Help: "Number of processed handlers",
	})

	hits := prometheus.NewCounterVec(prometheus.CounterOpts{
		Name: "hits",
	}, []string{"status", "path"})

	duration := prometheus.NewHistogramVec(prometheus.HistogramOpts{
		Name:      "duration",
		Help:      "The latency of the HTTP requests.",
	}, []string{"path", "method", "code"})

	m := &MetricsController{
		Counter: counter,
		Hits: hits,
		Duration: duration,
	}

	prometheus.MustRegister(counter, hits, duration)

	router.GET("/metrics", gin.WrapH(promhttp.Handler()))

	router.Use(m.IncCounter())
	router.Use(m.Hit())
	router.Use(m.GetTime())

	return m
}

func (mC *MetricsController) IncCounter() gin.HandlerFunc {
	return func(c *gin.Context) {
		mC.Counter.Add(1)
		c.Next()
	}
}

func (mC *MetricsController) Hit() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()
		statusCode := strconv.Itoa(c.Writer.Status())
		mC.Hits.WithLabelValues(statusCode, c.Request.URL.String()).Inc()
	}
}

func (mC *MetricsController) GetTime() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		c.Next()
		end := time.Since(start).Seconds()
		statusCode := strconv.Itoa(c.Writer.Status())
		mC.Duration.WithLabelValues(c.Request.URL.String(), c.Request.Method, statusCode).Observe(end)
	}
}
