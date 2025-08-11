package middleware

import (
	"strconv"
	"time"

	"github.com/frenswifbenefits/myfren/internal/metrics"
	"github.com/gin-gonic/gin"
)

func PrometheusMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()

		// Process the request
		c.Next()

		// Record metrics
		duration := time.Since(start).Seconds()
		status := c.Writer.Status()
		method := c.Request.Method
		path := c.FullPath() // Use the defined route path, not the raw URL path

		if path == "" {
			path = "unknown" // Fallback if FullPath is not defined
		}

		metrics.DurationHistogram.WithLabelValues(method, path, strconv.Itoa(status)).Observe(duration)
	}
}
