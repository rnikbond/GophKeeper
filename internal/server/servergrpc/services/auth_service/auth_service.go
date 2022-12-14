package auth_service

import (
	"GophKeeper/internal/storage"
	pb "GophKeeper/pkg/proto/auth"
	"context"
)

type AuthService struct {
	pb.AuthServiceServer
	store storage.UserStorage
}

func NewAuthService(store storage.UserStorage) *AuthService {
	return &AuthService{
		store: store,
	}
}

func (serv *AuthService) Register(ctx context.Context, in *pb.AuthRequest) (*pb.AuthResponse, error) {

	u := storage.User{
		Email:    in.Email,
		Password: in.Password,
	}

	if err := serv.store.Create(u); err != nil {
		return nil, err
	}

	return &pb.AuthResponse{
		Token: in.Email,
	}, nil
}

func (serv *AuthService) Login(ctx context.Context, in *pb.AuthRequest) (*pb.AuthResponse, error) {

	user, err := serv.store.Find(in.Email)
	if err != nil {
		return nil, err
	}

	return &pb.AuthResponse{
		Token: user.Email,
	}, nil
}

func (serv *AuthService) Logout(ctx context.Context, in *pb.LogoutRequest) (*pb.Empty, error) {

	if err := serv.store.Delete(in.Email); err != nil {
		return nil, err
	}

	return &pb.Empty{}, nil
}

func (serv *AuthService) ChangePassword(ctx context.Context, in *pb.ChangePasswordRequest) (*pb.Empty, error) {

	if len(in.Token) == 0 {
		return nil, ErrInvalidToken
	}
	if err := serv.store.Update(in.Email, in.Password); err != nil {
		return nil, err
	}

	return &pb.Empty{}, nil
}
