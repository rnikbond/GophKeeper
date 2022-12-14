package auth_service

import (
	"GophKeeper/internal/storage"
	pb "GophKeeper/pkg/proto/auth"
)

type (
	AuthService struct {
		pb.AuthServiceServer

		store     storage.UserStorage
		secretKey string
	}

	OptionsAuth func(*AuthService)
)

func NewAuthService(store storage.UserStorage, opts ...OptionsAuth) *AuthService {

	serv := &AuthService{
		store: store,
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
