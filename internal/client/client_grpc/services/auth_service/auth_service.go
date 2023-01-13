package auth_service

import (
	"GophKeeper/pkg/errs"
	pbAuth "GophKeeper/pkg/proto/auth"
	"context"
	"fmt"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"os"
)

type AuthService struct {
	rpc    pbAuth.AuthServiceClient
	logger *zap.Logger

	Token string
}

// NewService - Создание экземпляра сервиса авторизации.
func NewService(conn *grpc.ClientConn) *AuthService {
	return &AuthService{
		rpc:    pbAuth.NewAuthServiceClient(conn),
		logger: zap.L(),
	}
}

// SignIn - Авторизация пользователя.
func (c *AuthService) SignIn() error {

	auth := &pbAuth.AuthRequest{}

	fmt.Print("Email : ")
	fmt.Fscan(os.Stdin, &auth.Email)

	fmt.Print("Пароль: ")
	fmt.Fscan(os.Stdin, &auth.Password)

	resp, err := c.rpc.Login(context.Background(), auth)
	if err != nil {
		if e, ok := status.FromError(err); ok {
			switch e.Code() {
			case codes.NotFound:
				return errs.ErrNotFound
			case codes.InvalidArgument:
				return errs.ErrInvalidArgument
			case codes.Internal:
				return errs.ErrInternal
			}
		}

		return fmt.Errorf("%s %w", err.Error(), errs.ErrInternal)
	}

	c.Token = resp.Token
	return nil
}

// SignUp - Регистрация пользователя.
func (c *AuthService) SignUp() error {

	auth := &pbAuth.AuthRequest{}

	fmt.Print("Email: ")
	fmt.Fscan(os.Stdin, &auth.Email)

	fmt.Print("Password: ")
	fmt.Fscan(os.Stdin, &auth.Password)

	resp, err := c.rpc.Register(context.Background(), auth)
	if err != nil {
		if e, ok := status.FromError(err); ok {
			switch e.Code() {
			case codes.AlreadyExists:
				return errs.ErrAlreadyExist
			case codes.Internal:
				return errs.ErrInternal
			}
		}

		return fmt.Errorf("%s %w", err.Error(), errs.ErrInternal)
	}

	c.Token = resp.Token
	return nil
}
