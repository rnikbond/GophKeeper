package auth_service

import (
	"GophKeeper/pkg/errs"
	pbAuth "GophKeeper/pkg/proto/auth"
	"GophKeeper/pkg/secret"
	"bufio"
	"context"
	"fmt"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"os"
)

type AuthOptions func(c *AuthService)

type AuthService struct {
	rpc    pbAuth.AuthServiceClient
	logger *zap.Logger

	Token string
	salt  string
}

// NewService - Создание экземпляра сервиса авторизации.
func NewService(conn *grpc.ClientConn, opts ...AuthOptions) *AuthService {
	serv := &AuthService{
		rpc:    pbAuth.NewAuthServiceClient(conn),
		logger: zap.L(),
	}

	for _, opt := range opts {
		opt(serv)
	}

	return serv
}

func WithSalt(salt string) AuthOptions {
	return func(serv *AuthService) {
		serv.salt = salt
	}
}

// SignIn - Авторизация пользователя.
func (c *AuthService) SignIn() error {

	auth := &pbAuth.AuthRequest{}
	reader := bufio.NewReader(os.Stdin)

	fmt.Print("Email: ")
	auth.Email, _ = reader.ReadString('\n')

	fmt.Print("Пароль: ")
	auth.Password, _ = reader.ReadString('\n')

	//auth.Email = "test@mail.ru"
	//auth.Password = "test"

	auth.Password = secret.GeneratePasswordHash(auth.Password, c.salt)

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

	auth := &pbAuth.AuthRequest{}
	reader := bufio.NewReader(os.Stdin)

	fmt.Print("Email: ")
	auth.Email, _ = reader.ReadString('\n')

	fmt.Print("Пароль: ")
	auth.Password, _ = reader.ReadString('\n')

	//auth.Email = "test@mail.ru"
	//auth.Password = "test"

	auth.Password = secret.GeneratePasswordHash(auth.Password, c.salt)

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
