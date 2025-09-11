package utils

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func HandleValidationError(c *gin.Context, err error, log *logrus.Logger) {
	log.WithError(err).Warn("Validation error")
	c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input: " + err.Error()})
}
