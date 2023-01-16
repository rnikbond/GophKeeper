package app_service_auth

import (
	"bufio"
	"fmt"
	"golang.org/x/term"
	"os"
	"syscall"

	"GophKeeper/internal/client/model/auth_model"
	"GophKeeper/pkg/secret"
)

type AuthOptions func(c *AuthService)

type AuthService struct {
	salt string
}

// NewService - Создание экземпляра сервиса авторизации.
func NewService(opts ...AuthOptions) *AuthService {
	serv := &AuthService{}

	for _, opt := range opts {
		opt(serv)
	}

	return serv
}

func WithSalt(salt string, opts ...AuthOptions) AuthOptions {
	return func(s *AuthService) {
		s.salt = salt
	}
}

// SignIn - Авторизация пользователя.
func (serv *AuthService) SignIn() auth_model.Credential {

	cred := auth_model.Credential{}
	reader := bufio.NewReader(os.Stdin)

	fmt.Print("Email: ")
	cred.Email, _ = reader.ReadString('\n')

	fmt.Print("Пароль: ")
	pwd, _ := term.ReadPassword(int(syscall.Stdin))

	cred.Password = secret.GeneratePasswordHash(string(pwd), serv.salt)

	return cred
}

// SignUp - Регистрация пользователя.
func (serv *AuthService) SignUp() auth_model.Credential {

	cred := auth_model.Credential{}
	reader := bufio.NewReader(os.Stdin)

	fmt.Print("Email: ")
	cred.Email, _ = reader.ReadString('\n')

	fmt.Print("Пароль: ")
	pwd, _ := term.ReadPassword(int(syscall.Stdin))

	cred.Password = secret.GeneratePasswordHash(string(pwd), serv.salt)

	return cred
}
