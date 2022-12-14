package auth_service

import (
	"context"
	"log"

	"GophKeeper/internal/storage"
	pb "GophKeeper/pkg/proto/auth"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// Register Регистрация нового пользователя
func (serv *AuthService) Register(ctx context.Context, in *pb.AuthRequest) (*pb.AuthResponse, error) {

	u := storage.User{
		Email:    in.Email,
		Password: in.Password,
	}

	if err := serv.store.Create(u); err != nil {
		log.Printf("failed create user: %v", err)
		return nil, err
	}

	tokenStr, err := GenerateJWT(u.Email, serv.secretKey)
	if err != nil {
		log.Printf("%v: %v", ErrGenerateToken, err)
		return nil, status.Error(codes.Internal, ErrGenerateToken.Error())
	}

	return &pb.AuthResponse{
		Token: tokenStr,
	}, nil
}

// Login Авторизация пользоватедя
func (serv *AuthService) Login(ctx context.Context, in *pb.AuthRequest) (*pb.AuthResponse, error) {

	user, err := serv.store.Find(in.Email)
	if err != nil {
		return nil, err
	}

	if in.Password != user.Password {
		return nil, ErrInvalidAuthData
	}

	tokenStr, errJWT := GenerateJWT(user.Email, serv.secretKey)
	if errJWT != nil {
		log.Printf("%v: %v", ErrGenerateToken, errJWT)
		return nil, status.Error(codes.Internal, ErrGenerateToken.Error())
	}

	return &pb.AuthResponse{
		Token: tokenStr,
	}, nil
}

// ChangePassword Смена пароля пользователя
func (serv *AuthService) ChangePassword(ctx context.Context, in *pb.ChangePasswordRequest) (*pb.Empty, error) {

	tokenStr, errJWT := VerifyJWT(in.Token, serv.secretKey)
	if errJWT != nil {
		log.Printf("failed change password: %v : %v", ErrInvalidToken, errJWT)
		return nil, ErrInvalidToken
	}

	token := tokenStr.Claims.(*Token)

	if err := serv.store.Update(token.Email, in.Password); err != nil {
		log.Printf("failed change password: %v", err)
		return nil, err
	}

	return &pb.Empty{}, nil
}
