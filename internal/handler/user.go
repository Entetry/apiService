package handler

import (
	"net/http"

	"github.com/Entetry/userService/protocol/userService"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// User handler struct
type User struct {
	userServiceClient userService.UserServiceClient
}

// NewUser creates new user handler
func NewUser(userServiceClient userService.UserServiceClient) *User {
	return &User{userServiceClient: userServiceClient}
}

// GetByID godoc
func (u *User) GetByID(ctx echo.Context) error {
	id, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest)
	}

	protoResp, err := u.userServiceClient.GetByID(ctx.Request().Context(), &userService.GetByIDRequest{
		Uuid: id.String(),
	})
	if err != nil {
		return handleError(err)
	}

	return ctx.JSON(http.StatusOK, &getByIDResponse{
		Username: protoResp.Name,
		Email:    protoResp.Email,
		ID:       protoResp.Uuid,
	})
}

// Create godoc
func (u *User) Create(ctx echo.Context) error {
	request := new(createRequest)
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

	protoResp, err := u.userServiceClient.Create(ctx.Request().Context(), &userService.CreateRequest{
		Username: request.Username,
		Email:    request.Email,
		Password: request.Password,
	})
	if err != nil {
		return handleError(err)
	}

	return ctx.JSON(http.StatusOK, &createResponse{
		ID: protoResp.Uuid,
	})
}

// Delete godoc
func (u *User) Delete(ctx echo.Context) error {
	id, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest)
	}

	_, err = u.userServiceClient.Delete(ctx.Request().Context(), &userService.DeleteRequest{
		Uuid: id.String(),
	})
	if err != nil {
		return handleError(err)
	}

	return ctx.JSON(http.StatusOK, struct{}{})
}

func handleError(err error) error {
	switch status.Code(err) {
	case codes.NotFound:
		return echo.NewHTTPError(http.StatusNotFound, err.Error())
	case codes.InvalidArgument:
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	default:
		log.Error(err)
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
}
