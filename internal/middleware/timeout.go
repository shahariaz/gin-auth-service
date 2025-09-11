package middleware

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/shahariaz/gin-auth-service/internal/errs"
	"github.com/sirupsen/logrus"
)

func TimeoutMiddleware(timeout time.Duration, log *logrus.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(c.Request.Context(), timeout)
		defer cancel()

		c.Request = c.Request.WithContext(ctx)

		done := make(chan struct{})
		go func() {
			defer close(done)
			c.Next()
		}()

		select {
		case <-done:
		case <-ctx.Done():
			log.WithField("path", c.Request.URL.Path).Error("Request timed out")
			errs.HandleError(c, errs.NewAPIError(http.StatusGatewayTimeout, "Request timed out", ctx.Err()), log)
			c.Abort()
		}
	}
}
