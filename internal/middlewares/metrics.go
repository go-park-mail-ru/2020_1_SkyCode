package middlewares

import (
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"strconv"
)

type MetricsController struct {
	Counter prometheus.Counter
	Hits *prometheus.CounterVec

}

func NewMetricsController(router *gin.Engine) *MetricsController {
	counter := prometheus.NewCounter(prometheus.CounterOpts{
		Name: "total_counter",
		Help: "Number of processed handlers",
	})

	hits := prometheus.NewCounterVec(prometheus.CounterOpts{
		Name: "hits",
	}, []string{"status", "path"})

	m := &MetricsController{
		Counter: counter,
		Hits: hits,
	}

	prometheus.MustRegister(counter, hits)

	router.GET("/metrics", gin.WrapH(promhttp.Handler()))

	router.Use(m.IncCounter())
	router.Use(m.Hit())

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
