//go:generate mockgen -source grpc_service_auth.go -destination mocks/grpc_service_auth_mock.go -package grpc_service_auth
package grpc_service_auth

import (
	"context"

	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"GophKeeper/internal/client/model/auth_model"
	"GophKeeper/pkg/errs"
	pb "GophKeeper/pkg/proto/auth"
)

type AuthService struct {
	rpc    pb.AuthServiceClient
	logger *zap.Logger
}

// NewService - Создание экземпляра gRPC сервиса авторизации.
func NewService(conn *grpc.ClientConn) *AuthService {
	return &AuthService{
		rpc:    pb.NewAuthServiceClient(conn),
		logger: zap.L(),
	}
}

// SignIn - Авторизация пользователя.
// При отсутствии ошибки возвращается тоекн авторизации.
func (c *AuthService) SignIn(cred auth_model.Credential) (string, error) {

	auth := &pb.AuthRequest{
		Email:    cred.Email,
		Password: cred.Password,
	}

	resp, err := c.rpc.Login(context.Background(), auth)
	if err != nil {
		if e, ok := status.FromError(err); ok {
			switch e.Code() {
			case codes.NotFound:
				return ``, errs.ErrNotFound
			case codes.InvalidArgument:
				return ``, errs.ErrInvalidArgument
			default:
				c.logger.Error("unknown gRPC error in SignUp",
					zap.Uint32("gRPC code", uint32(e.Code())),
					zap.String("gRPC text", e.String()))
			}
		}

		return ``, errs.ErrInternal
	}

	return resp.Token, nil
}

// SignUp - Регистрация пользователя.
// При отсутствии ошибки возвращается тоекн авторизации.
func (c *AuthService) SignUp(cred auth_model.Credential) (string, error) {

	auth := &pb.AuthRequest{
		Email:    cred.Email,
		Password: cred.Password,
	}

	resp, err := c.rpc.Register(context.Background(), auth)

	if err != nil {
		if e, ok := status.FromError(err); ok {
			switch e.Code() {
			case codes.AlreadyExists:
				return ``, errs.ErrAlreadyExist
			case codes.InvalidArgument:
				return ``, errs.ErrInvalidArgument
			default:
				c.logger.Error("unknown gRPC error in SignUp",
					zap.Uint32("gRPC code", uint32(e.Code())),
					zap.String("gRPC text", e.String()))
			}
		}

		return ``, errs.ErrInternal
	}

	return resp.Token, nil
}
