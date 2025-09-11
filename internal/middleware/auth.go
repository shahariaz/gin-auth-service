package middleware

import (
	"errors"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/shahariaz/gin-auth-service/internal/errs"
	"github.com/shahariaz/gin-auth-service/internal/lib"
	"github.com/sirupsen/logrus"
)

func JWTAuthMiddleware(secret []byte, tokenStore lib.TokenStore, log *logrus.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		auth := c.GetHeader("Authorization")
		if auth == "" || !strings.HasPrefix(auth, "Bearer ") {
			log.Warn("Missing Bearer token")
			errs.HandleError(c, errs.NewAPIError(http.StatusUnauthorized, "Bearer token required", nil), log)
			c.Abort()
			return
		}

		tokenStr := auth[7:]
		isBlacklisted, err := tokenStore.IsBlacklisted(tokenStr)
		if err != nil {
			log.WithError(err).Error("Token store error")
			errs.HandleError(c, errs.NewAPIError(http.StatusInternalServerError, "Internal server error", err), log)
			c.Abort()
			return
		}
		if isBlacklisted {
			log.Warn("Blacklisted token used")
			errs.HandleError(c, errs.NewAPIError(http.StatusUnauthorized, "Token is blacklisted", nil), log)
			c.Abort()
			return
		}

		token, err := jwt.ParseWithClaims(tokenStr, &lib.TokenClaims{}, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, errors.New("unexpected signing method")
			}
			return secret, nil
		})

		if err != nil || !token.Valid {
			log.WithError(err).Warn("Invalid token")
			errs.HandleError(c, errs.NewAPIError(http.StatusUnauthorized, "Invalid or expired token", err), log)
			c.Abort()
			return
		}

		claims, ok := token.Claims.(*lib.TokenClaims)
		if !ok || claims.Username == "" || claims.Role == "" {
			log.Warn("Invalid claims")
			errs.HandleError(c, errs.NewAPIError(http.StatusUnauthorized, "Invalid token claims", nil), log)
			c.Abort()
			return
		}

		c.Set("user", claims.Username)
		c.Set("role", claims.Role)
		c.Set("user_id", claims.UserID)
		c.Next()
	}
}

func RequireRole(role string, log *logrus.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		userRole, exists := c.Get("role")
		if !exists || userRole != role {
			log.WithField("required_role", role).Warn("Role access denied")
			errs.HandleError(c, errs.NewAPIError(http.StatusForbidden, "Forbidden: "+role+" role required", nil), log)
			c.Abort()
			return
		}
		c.Next()
	}
}
