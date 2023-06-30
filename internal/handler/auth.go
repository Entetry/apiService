// Package handler provides grpc user api
package handler

import (
	"net/http"

	"github.com/Entetry/authService/protocol/authService"

	"github.com/labstack/echo/v4"
	log "github.com/sirupsen/logrus"
)

// Auth handler struct
type Auth struct {
	authServiceClient authService.AuthGRPCServiceClient
}

// NewAuth creates new auth handler
func NewAuth(authServiceClient authService.AuthGRPCServiceClient) *Auth {
	return &Auth{authServiceClient: authServiceClient}
}

// SignIn godoc
func (a *Auth) SignIn(ctx echo.Context) error {
	request := new(signInRequest)
	err := ctx.Bind(request)
	if err != nil {
		log.Error(err)
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	err = ctx.Validate(request)
	if err != nil {
		log.Error(err)
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	protoResp, err := a.authServiceClient.SignIn(ctx.Request().Context(), &authService.SignInRequest{
		Username: request.Username,
		Password: request.Password,
	})
	if err != nil {
		return handleError(err)
	}

	return ctx.JSON(http.StatusOK, &tokenResponse{
		AccessToken:  protoResp.AccessToken,
		RefreshToken: protoResp.RefreshToken,
	})
}

// SignUp godoc
func (a *Auth) SignUp(ctx echo.Context) error {
	request := new(signUpRequest)
	err := ctx.Bind(request)
	if err != nil {
		log.Error(err)
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	err = ctx.Validate(request)
	if err != nil {
		log.Error(err)
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	_, err = a.authServiceClient.SignUp(ctx.Request().Context(), &authService.SignUpRequest{
		Username: request.Username,
		Password: request.Password,
		Email:    request.Email,
	})
	if err != nil {
		return handleError(err)
	}

	return ctx.String(http.StatusCreated, "Registration completed successfully")
}

// Refresh godoc
func (a *Auth) Refresh(ctx echo.Context) error {
	request := new(refreshTokenRequest)
	err := ctx.Bind(request)
	if err != nil {
		return handleError(err)
	}

	err = ctx.Validate(request)
	if err != nil {
		log.Error(err)
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	protoResp, err := a.authServiceClient.RefreshTokens(ctx.Request().Context(), &authService.RefreshTokensRequest{
		RefreshToken: request.RefreshToken,
		Username:     request.Username,
	})
	if err != nil {
		log.Error(err)
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return ctx.JSON(http.StatusOK, &tokenResponse{
		AccessToken:  protoResp.AccessToken,
		RefreshToken: protoResp.RefreshToken,
	})
}
