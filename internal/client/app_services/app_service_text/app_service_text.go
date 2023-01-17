package app_service_text

import (
	"crypto/rsa"
)

type TextOptions func(c *TextService)

type TextService struct {
	publicKey  *rsa.PublicKey
	privateKey *rsa.PrivateKey

	token string
}

// NewService - Создание экземпляра сервиса для текстовых данных.
func NewService(opts ...TextOptions) *TextService {
	serv := &TextService{}

	for _, opt := range opts {
		opt(serv)
	}

	return serv
}

func WithPublicKey(key *rsa.PublicKey) TextOptions {
	return func(serv *TextService) {
		serv.publicKey = key
	}
}

func WithPrivateKey(key *rsa.PrivateKey) TextOptions {
	return func(serv *TextService) {
		serv.privateKey = key
	}
}
