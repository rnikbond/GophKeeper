package app_service_auth

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"syscall"

	"github.com/fatih/color"
	"go.uber.org/zap"
	"golang.org/x/term"

	"GophKeeper/internal/client/model/auth_model"
	"GophKeeper/pkg/errs"
	"GophKeeper/pkg/secret"
)

type Sender interface {
	SignIn(auth_model.Credential) (string, error)
	SignUp(auth_model.Credential) (string, error)
}

type AuthOptions func(c *AuthService)

type AuthService struct {
	logger *zap.Logger
	salt   string
	Sender
}

// NewService - Создание экземпляра сервиса авторизации.
func NewService(s Sender, opts ...AuthOptions) *AuthService {
	serv := &AuthService{
		logger: zap.L(),
		Sender: s,
	}

	for _, opt := range opts {
		opt(serv)
	}

	return serv
}

func WithSalt(salt string) AuthOptions {
	return func(s *AuthService) {
		s.salt = salt
	}
}

func (serv AuthService) Token() (string, error) {

	stdin := bufio.NewReader(os.Stdin)

	for {

		fmt.Println("---------------")
		fmt.Println("[0] Завершить")
		fmt.Println("[1] Авторизация")
		fmt.Println("[2] Регистрация")
		fmt.Println("---------------")
		fmt.Print("-> ")

		var choice int

		_, err := fmt.Fscan(os.Stdin, &choice)
		stdin.ReadString('\n')

		if err != nil {
			continue
		}

		switch choice {
		case 0:
			return ``, errs.ErrCancel

		case 1:
			token, errToken := serv.signIn()
			if ok := serv.parseErr(errToken); ok {
				return token, nil
			}

		case 2:
			token, errToken := serv.signUp()
			if ok := serv.parseErr(errToken); ok {
				return token, nil
			}
		}
	}
}

// SignIn - Авторизация пользователя.
func (serv AuthService) signIn() (string, error) {

	//cred := serv.readCredential()

	cred := auth_model.Credential{
		Email:    "email@com",
		Password: "email123",
	}

	return serv.Sender.SignIn(cred)
}

// SignUp - Регистрация пользователя.
func (serv AuthService) signUp() (string, error) {

	//cred := serv.readCredential()

	cred := auth_model.Credential{
		Email:    "email@com",
		Password: "email123",
	}

	return serv.Sender.SignUp(cred)
}

func (serv AuthService) readCredential() auth_model.Credential {

	cred := auth_model.Credential{}
	reader := bufio.NewReader(os.Stdin)

	fmt.Print("Email: ")
	cred.Email, _ = reader.ReadString('\n')

	fmt.Print("Пароль: ")
	pwd, _ := term.ReadPassword(int(syscall.Stdin))

	cred.Password = secret.GeneratePasswordHash(string(pwd), serv.salt)

	return cred
}

func (serv AuthService) parseErr(err error) bool {

	if err == nil {
		return true
	}

	color.New(color.FgRed).Print("\tОшибка авторизации: ")

	switch {

	case errors.Is(err, errs.ErrAlreadyExist):
		fmt.Println("Такой Email же зарегистрирован")

	case errors.Is(err, errs.ErrNotFound):
		fmt.Println("Пользователь не найден")

	case errors.Is(err, errs.ErrInvalidArgument):
		fmt.Println("Неверный логин или пароль")

	default:
		fmt.Println("Внутренняя ошибка сервиса")
		serv.logger.Error("unknown error", zap.Error(err))
	}

	return false
}
