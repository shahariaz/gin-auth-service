package validation

import (
	"github.com/go-playground/validator/v10"
)

func NewValidator() *validator.Validate {
	v := validator.New()
	// Add custom validators if needed
	return v
}
