package auth_service

import (
	"GophKeeper/internal/server/model"
	pb "GophKeeper/pkg/proto/auth"
)

type AuthService struct {
	pb.AuthServiceServer

	auth      *model.AuthModel
	secretKey string
}

type OptionsAuth func(*AuthService)

func NewAuthService(auth *model.AuthModel, opts ...OptionsAuth) *AuthService {

	serv := &AuthService{
		auth: auth,
	}

	for _, opt := range opts {
		opt(serv)
	}

	return serv
}

func WithSecretKey(key string) OptionsAuth {
	return func(serv *AuthService) {
		serv.secretKey = key
	}
}
