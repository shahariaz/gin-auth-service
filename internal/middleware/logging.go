package middleware

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func LoggingMiddleware(log *logrus.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		method := c.Request.Method
		userAgent := c.Request.UserAgent()
		ip := c.ClientIP()

		c.Next()

		latency := time.Since(start)
		status := c.Writer.Status()

		log.WithFields(map[string]interface{}{
			"method":     method,
			"path":       path,
			"status":     status,
			"latency":    latency.String(),
			"ip":         ip,
			"user_agent": userAgent,
		}).Info("Request processed")
	}
}
