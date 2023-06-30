package middleware

import (
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"net/http"
)

// CustomValidator validates a structs exposed fields, and automatically validates nested structs
type CustomValidator struct {
	Validator *validator.Validate
}

// NewCustomValidator godoc
func NewCustomValidator(validator *validator.Validate) *CustomValidator {
	return &CustomValidator{
		Validator: validator,
	}
}

// Validate validates a structs
func (cv *CustomValidator) Validate(i interface{}) error {
	if err := cv.Validator.Struct(i); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	return nil
}
