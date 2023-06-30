package main

import (
	"context"
	"fmt"
	"github.com/Entetry/apiService/internal/config"
	"github.com/Entetry/apiService/internal/handler"
	"github.com/Entetry/apiService/internal/middleware"
	"github.com/Entetry/authService/protocol/authService"
	"github.com/Entetry/userService/protocol/userService"
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"net"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	log.SetFormatter(&log.TextFormatter{TimestampFormat: "2006-01-02 15:04:05", FullTimestamp: true})

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	cfg, err := config.New()

	if err != nil {
		log.Panicf("Main / Couldnt parse config \n %v", err)
	}

	userConn, err := grpc.Dial(cfg.UserEndpoint,
		grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Panicf("Couldn't connect to user service: %v", err)
	}
	userServiceClient := userService.NewUserServiceClient(userConn)
	defer func() {
		err = userConn.Close()
		if err != nil {
			log.Errorf("Main / userConn.Close() / \n %v", err)
			return
		}
	}()
	log.Info(cfg.UserEndpoint)

	authConn, err := grpc.Dial(cfg.AuthEndpoint,
		grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Panicf("Couldn't connect to auth service: %v", err)
	}
	authServiceClient := authService.NewAuthGRPCServiceClient(authConn)
	defer func() {
		err = userConn.Close()
		if err != nil {
			log.Errorf("Main / authConn.Close() / \n %v", err)
			return
		}
	}()

	e := echo.New()
	e.Validator = middleware.NewCustomValidator(validator.New())

	authHdr := handler.NewAuth(authServiceClient)
	authGroup := e.Group("/auth")
	authGroup.POST("/sign-up", authHdr.SignUp)
	authGroup.POST("/sign-in", authHdr.SignIn)
	authGroup.POST("/refresh-tokens", authHdr.Refresh)
	authMiddleware := middleware.NewJwt(authServiceClient)

	userHdr := handler.NewUser(userServiceClient)
	userGroup := e.Group("/users", authMiddleware.JWTMiddleware)
	userGroup.GET("/:id", userHdr.GetByID)
	userGroup.POST("", userHdr.Create)
	userGroup.DELETE("/:id", userHdr.Delete)

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGTERM)
	err = e.Start(fmt.Sprintf(":%d", cfg.Port))
	if err != nil {
		log.Fatal(err)
	}
	log.Info("Server started on ", cfg.Port)
	go func() {
		<-sigChan
		cancel()
		err = e.Shutdown(ctx)
		if err != nil {
			log.Errorf("can't stop server gracefully %v", err)
		}
	}()
	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", cfg.Port))
	if err != nil {
		log.Fatal(err)
	}
	err = e.Server.Serve(listener)
	if err != nil {
		log.Fatal(err)
	}
}
