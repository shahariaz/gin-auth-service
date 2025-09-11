package middleware

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/shahariaz/gin-auth-service/internal/errs"
	"github.com/sirupsen/logrus"
)

func ErrorMiddleware(log *logrus.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next() // Process the request first

		// Only handle if errors exist and response not written
		if len(c.Errors) > 0 && c.Writer.Status() == 0 {
			// Take the first error (or loop for all if needed)
			err := c.Errors[0]

			// Determine status and message based on error type
			status := http.StatusInternalServerError
			message := "Internal server error"

			// Categorize errors
			switch err.Type {
			case gin.ErrorTypeBind:
				status = http.StatusBadRequest
				message = "Invalid input: " + err.Error()
			case gin.ErrorTypePublic:
				status = http.StatusBadRequest
				message = err.Error()
			case gin.ErrorTypePrivate:
				// For internal errors, log details but return generic message
				log.WithFields(logrus.Fields{
					"path":   c.Request.URL.Path,
					"method": c.Request.Method,
					"ip":     c.ClientIP(),
					"user":   c.GetString("user"), // From JWT context if available
					"error":  err.Error(),
					"status": status,
				}).Error("Private error occurred")
				message = "Internal server error"
			default:
				if apiErr, ok := err.Err.(*errs.APIError); ok {
					status = apiErr.Status
					message = apiErr.Message
				}
			}

			// Log the error (structured)
			log.WithFields(logrus.Fields{
				"path":   c.Request.URL.Path,
				"method": c.Request.Method,
				"ip":     c.ClientIP(),
				"user":   c.GetString("user"),
				"status": status,
				"error":  err.Error(),
			}).Error("Request error handled")

			// Return consistent JSON error response
			c.JSON(status, gin.H{
				"error": message,
				// Optional: Add timestamp or request ID for tracing
				"timestamp": time.Now().Format(time.RFC3339),
			})
		}
	}
}
