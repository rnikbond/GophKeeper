package auth_service

import (
	"context"

	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"GophKeeper/internal/server/model"
	"GophKeeper/internal/storage"
	pb "GophKeeper/pkg/proto/auth"
	"GophKeeper/pkg/token"
)

// Register - Регистрация нового пользователя.
func (serv *AuthService) Register(ctx context.Context, in *pb.AuthRequest) (*pb.AuthResponse, error) {

	cred := storage.Credential{
		Email:    in.Email,
		Password: in.Password,
	}

	tkn, err := serv.auth.Register(cred)
	if err != nil {

		switch err {
		case model.ErrAlreadyExists:
			return nil, status.Error(codes.AlreadyExists, err.Error())

		case model.ErrInvalidPassword:
			return nil, status.Error(codes.Unauthenticated, err.Error())
		}

		return nil, status.Error(codes.Internal, "Internal server error")
	}

	return &pb.AuthResponse{
		Token: tkn,
	}, nil
}

// Login - Авторизация пользователя.
func (serv *AuthService) Login(ctx context.Context, in *pb.AuthRequest) (*pb.AuthResponse, error) {

	cred := storage.Credential{
		Email:    in.Email,
		Password: in.Password,
	}

	tkn, err := serv.auth.Login(cred)
	if err != nil {

		switch err {
		case model.ErrNotFound:
			return nil, status.Error(codes.NotFound, err.Error())

		case model.ErrInvalidPassword:
			return nil, status.Error(codes.Unauthenticated, err.Error())
		}

		return nil, status.Error(codes.Internal, "Internal server error")
	}

	return &pb.AuthResponse{
		Token: tkn,
	}, nil
}

// ChangePassword - Смена пароля пользователя.
func (serv *AuthService) ChangePassword(ctx context.Context, in *pb.ChangePasswordRequest) (*pb.Empty, error) {

	logger := zap.L()

	tkn, errToken := token.ReadToken(ctx)
	if errToken != nil {
		return nil, status.Error(codes.PermissionDenied, "Invalid token")
	}

	if err := serv.auth.ChangePassword(tkn, in.Password); err != nil {
		if err == model.ErrInvalidToken {
			return nil, status.Error(codes.PermissionDenied, "Invalid token")
		}

		logger.Error("failed change password", zap.Error(err))
		return nil, status.Error(codes.Internal, "Internal server error")
	}

	return &pb.Empty{}, nil
}
