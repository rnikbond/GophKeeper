//go:generate mockgen -source rpc_service_auth.go -destination mocks/rpc_service_auth_mock.go -package grpc_service_auth
package grpc_service_auth

import (
	"context"
	"errors"

	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"GophKeeper/internal/server/app_services/app_service_auth"
	"GophKeeper/internal/server/model/auth"
	"GophKeeper/pkg/errs"
	"GophKeeper/pkg/md_ctx"
	pb "GophKeeper/pkg/proto/auth"
)

type AuthApp interface {
	Login(in auth.Credential) (string, error)
	Register(in auth.Credential) (string, error)
	ChangePassword(email, password string) error
}

type AuthServiceRPC struct {
	pb.AuthServiceServer

	auth   AuthApp
	logger *zap.Logger
}

// NewAuthServiceRPC - Создание эклемпляра gRPC сервиса авторизации и регистрации
func NewAuthServiceRPC(auth AuthApp) *AuthServiceRPC {
	serv := &AuthServiceRPC{
		auth:   auth,
		logger: zap.L(),
	}

	return serv
}

// Register - Регистрация нового пользователя.
func (serv *AuthServiceRPC) Register(ctx context.Context, in *pb.AuthRequest) (*pb.AuthResponse, error) {

	cred := auth.Credential{
		Email:    in.Email,
		Password: in.Password,
	}

	tokenStr, err := serv.auth.Register(cred)
	if err != nil {

		if errors.Is(err, errs.ErrAlreadyExist) {
			return nil, status.Error(codes.AlreadyExists, err.Error())
		}

		if errors.Is(err, app_service_auth.ErrInvalidPassword) ||
			errors.Is(err, app_service_auth.ErrShortPassword) ||
			errors.Is(err, app_service_auth.ErrInvalidEmail) {
			return nil, status.Error(codes.InvalidArgument, err.Error())
		}

		return nil, status.Error(codes.Internal, errs.ErrInternal.Error())
	}

	return &pb.AuthResponse{
		Token: tokenStr,
	}, nil
}

// Login - Авторизация пользователя.
func (serv *AuthServiceRPC) Login(ctx context.Context, in *pb.AuthRequest) (*pb.AuthResponse, error) {

	cred := auth.Credential{
		Email:    in.Email,
		Password: in.Password,
	}

	tkn, err := serv.auth.Login(cred)
	if err != nil {

		if errors.Is(err, errs.ErrNotFound) {
			return nil, status.Error(codes.NotFound, err.Error())
		}

		if errors.Is(err, app_service_auth.ErrInvalidPassword) {
			return nil, status.Error(codes.InvalidArgument, err.Error())
		}

		return nil, status.Error(codes.Internal, errs.ErrInternal.Error())
	}

	return &pb.AuthResponse{
		Token: tkn,
	}, nil
}

// ChangePassword - Смена пароля пользователя.
func (serv *AuthServiceRPC) ChangePassword(ctx context.Context, in *pb.ChangePasswordRequest) (*pb.Empty, error) {

	email, ok := md_ctx.ValueFromContext(ctx, "email")

	if !ok {
		serv.logger.Error("failed found email in ctx metadata")
		// Internal, т.к. Interceptor должен был положить email в ctx
		return nil, status.Error(codes.Internal, errs.ErrInternal.Error())
	}

	if err := serv.auth.ChangePassword(email, string(in.Password)); err != nil {

		if errors.Is(err, app_service_auth.ErrShortPassword) {
			return nil, status.Error(codes.InvalidArgument, err.Error())
		}

		serv.logger.Error("failed change password", zap.Error(err))
		return nil, status.Error(codes.Internal, errs.ErrInternal.Error())
	}

	return &pb.Empty{}, nil
}
