package rpc_auth_service

import (
	"GophKeeper/internal/client/model/auth_model"
	"context"
	"fmt"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"GophKeeper/pkg/errs"
	pb "GophKeeper/pkg/proto/auth"
)

type AuthAppService interface {
	SignIn() auth_model.Credential
	SignUp() auth_model.Credential
}

type AuthService struct {
	rpc    pb.AuthServiceClient
	app    AuthAppService
	logger *zap.Logger

	Token string
}

// NewService - Создание экземпляра сервиса авторизации.
func NewService(conn *grpc.ClientConn, app AuthAppService) *AuthService {
	return &AuthService{
		rpc:    pb.NewAuthServiceClient(conn),
		app:    app,
		logger: zap.L(),
	}
}

// SignIn - Авторизация пользователя.
func (c *AuthService) SignIn() error {

	cred := c.app.SignIn()
	auth := &pb.AuthRequest{
		Email:    cred.Email,
		Password: cred.Password,
	}

	resp, err := c.rpc.Login(context.Background(), auth)
	if err != nil {
		if e, ok := status.FromError(err); ok {
			switch e.Code() {
			case codes.NotFound:
				return errs.ErrNotFound
			case codes.InvalidArgument:
				return errs.ErrInvalidArgument
			}
		}

		return fmt.Errorf("%s %w", err.Error(), errs.ErrInternal)
	}

	c.Token = resp.Token
	return nil
}

// SignUp - Регистрация пользователя.
func (c *AuthService) SignUp() error {

	cred := c.app.SignUp()
	auth := &pb.AuthRequest{
		Email:    cred.Email,
		Password: cred.Password,
	}

	resp, err := c.rpc.Register(context.Background(), auth)
	if err != nil {
		e, ok := status.FromError(err)
		if ok {
			switch e.Code() {
			case codes.AlreadyExists:
				return errs.ErrAlreadyExist
			case codes.InvalidArgument:
				return errs.ErrInvalidArgument
			}
		}

		fmt.Printf("%s. Grpc: %d - %s. %v\n", err.Error(), e.Code(), e.String(), errs.ErrInternal)
		return fmt.Errorf("%s. Grpc: %d - %s. %w", err.Error(), e.Code(), e.String(), errs.ErrInternal)
	}

	c.Token = resp.Token
	return nil
}
