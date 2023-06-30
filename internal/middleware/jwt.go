// Package middleware provides grpc user api
package middleware

import (
	"net/http"
	"strings"

	"github.com/Entetry/authService/protocol/authService"
	"github.com/labstack/echo/v4"
)

// Jwt middleware auth struct
type Jwt struct {
	authServiceClient authService.AuthGRPCServiceClient
}

// NewJwt godoc
func NewJwt(authServiceClient authService.AuthGRPCServiceClient) *Jwt {
	return &Jwt{
		authServiceClient: authServiceClient,
	}
}

// JWTMiddleware handle tokens
func (j *Jwt) JWTMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		authHeader := c.Request().Header.Get("Authorization")
		if authHeader == "" {
			return echo.NewHTTPError(http.StatusUnauthorized, "Missing authorization token")
		}
		tokenString := strings.Split(authHeader, " ")[1]

		_, err := j.authServiceClient.ValidateTokens(c.Request().Context(), &authService.ValidateTokensRequest{AccessToken: tokenString})
		if err != nil {
			return echo.NewHTTPError(http.StatusUnauthorized, err.Error())
		}
		return next(c)
	}
}
