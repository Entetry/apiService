package middleware

import (
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
)

// CustomValidator validates a structs exposed fields, and automatically validates nested structs
type CustomValidator struct {
	Validator *validator.Validate
}

// NewCustomValidator godoc
func NewCustomValidator(v *validator.Validate) *CustomValidator {
	return &CustomValidator{
		Validator: v,
	}
}

// Validate validates a structs
func (cv *CustomValidator) Validate(i interface{}) error {
	if err := cv.Validator.Struct(i); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	return nil
}
