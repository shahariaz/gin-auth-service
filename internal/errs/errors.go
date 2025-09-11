package errs

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type APIError struct {
	Status  int
	Message string
	Err     error
}

func NewAPIError(status int, message string, err error) *APIError {
	return &APIError{Status: status, Message: message, Err: err}
}

func (e *APIError) Error() string {
	return e.Message
}

func HandleValidationError(c *gin.Context, err error, log *logrus.Logger) {
	log.WithError(err).Warn("Validation error")
	c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input: " + err.Error()})
}

func HandleError(c *gin.Context, apiErr *APIError, log *logrus.Logger) {
	log.WithFields(map[string]interface{}{
		"error":  apiErr.Err,
		"status": apiErr.Status,
	}).Error("API error")
	c.JSON(apiErr.Status, gin.H{"error": apiErr.Message})
}
