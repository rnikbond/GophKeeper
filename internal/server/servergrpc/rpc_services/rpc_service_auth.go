package rpc_services

import (
	"GophKeeper/internal/server/app_services"
	"GophKeeper/pkg/md_ctx"
	pb "GophKeeper/pkg/proto/auth"
	"context"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var (
	InternalErrorText = "Internal server error"
)

type AuthServiceRPC struct {
	pb.AuthServiceServer

	auth   *app_services.AuthAppService
	logger *zap.Logger
}

// NewAuthServiceRPC - Создание эклемпляра gRPC сервиса авторизации и регистрации
func NewAuthServiceRPC(auth *app_services.AuthAppService) *AuthServiceRPC {
	serv := &AuthServiceRPC{
		auth:   auth,
		logger: zap.L(),
	}

	return serv
}

// Register - Регистрация нового пользователя.
func (serv *AuthServiceRPC) Register(ctx context.Context, in *pb.AuthRequest) (*pb.AuthResponse, error) {

	cred := app_services.Credential{
		Email:    in.Email,
		Password: in.Password,
	}

	tokenStr, err := serv.auth.Register(cred)
	if err != nil {

		switch err {
		case app_services.ErrAlreadyExists:
			return nil, status.Error(codes.AlreadyExists, err.Error())

		case app_services.ErrInvalidPassword:
			return nil, status.Error(codes.Unauthenticated, err.Error())
		}

		return nil, status.Error(codes.Internal, InternalErrorText)
	}

	return &pb.AuthResponse{
		Token: tokenStr,
	}, nil
}

// Login - Авторизация пользователя.
func (serv *AuthServiceRPC) Login(ctx context.Context, in *pb.AuthRequest) (*pb.AuthResponse, error) {

	cred := app_services.Credential{
		Email:    in.Email,
		Password: in.Password,
	}

	tkn, err := serv.auth.Login(cred)
	if err != nil {

		switch err {
		case app_services.ErrNotFound:
			return nil, status.Error(codes.NotFound, err.Error())

		case app_services.ErrInvalidPassword:
			return nil, status.Error(codes.Unauthenticated, err.Error())
		}

		return nil, status.Error(codes.Internal, InternalErrorText)
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
		// Ошибка Internal, т.к. Interceptor должен был положить email в ctx
		return nil, status.Error(codes.Internal, InternalErrorText)
	}

	if err := serv.auth.ChangePassword(email, in.Password); err != nil {
		serv.logger.Error("failed change password", zap.Error(err))
		return nil, status.Error(codes.Internal, InternalErrorText)
	}

	return &pb.Empty{}, nil
}
