//go:generate mockgen -source app_service_auth.go -destination mocks/app_service_auth_mock.go -package app_services
package app_service_auth

import (
	"errors"
	"strings"

	"go.uber.org/zap"

	authModel "GophKeeper/internal/model/auth"
	"GophKeeper/internal/storage/auth_store"
	"GophKeeper/pkg/errs"
	"GophKeeper/pkg/token"
)

// TODO :: Вынести в config
var minPasswordLength = 6

type AuthApp interface {
	Login(in authModel.Credential) (string, error)
	Register(in authModel.Credential) (string, error)
	ChangePassword(email, password string) error
}

// AuthAppOption - определяет операцию сервиса авторизации.
type AuthAppOption func(serv *AuthAppService)

// AuthAppService отвеает за сервис авторизации и регистрации пользователя.
type AuthAppService struct {
	store     auth_store.AuthStorage
	logger    *zap.Logger
	secretKey string
}

// NewAuthService - Создание экземпляра сервиса авторизации.
func NewAuthService(store auth_store.AuthStorage, opts ...AuthAppOption) *AuthAppService {

	auth := &AuthAppService{
		store:  store,
		logger: zap.L(),
	}

	for _, opt := range opts {
		opt(auth)
	}

	return auth
}

// WithSecretKey - Инициализирует секретный ключ для генерации JWT.
func WithSecretKey(key string) AuthAppOption {
	return func(auth *AuthAppService) {
		auth.secretKey = key
	}
}

// Login - Авторизация пользователя.
// При успешной авторизации возвращается JWT.
func (auth AuthAppService) Login(in authModel.Credential) (string, error) {

	err := auth.store.Find(in)

	switch err {
	case nil:
		break

	case errs.ErrNotFound:
		return ``, errs.ErrNotFound

	case auth_store.ErrInvalidPassword:
		return ``, ErrInvalidPassword

	default:
		auth.logger.Error("failed find user", zap.Error(err))
		return ``, errs.ErrInternal
	}

	tokenStr, errJWT := token.GenerateJWT(in.Email, auth.secretKey)
	if errJWT != nil {
		auth.logger.Error("failed generate JWT", zap.Error(errJWT))
		return ``, errs.ErrInternal
	}

	return tokenStr, nil
}

// Register - Регистрация пользователя.
// При успешной регистрации возвращается JWT.
func (auth AuthAppService) Register(in authModel.Credential) (string, error) {

	cred := authModel.Credential{
		Email:    in.Email,
		Password: in.Password,
	}

	if errCred := checkCredential(cred); errCred != nil {
		return ``, errCred
	}

	if err := auth.store.Create(cred); err != nil {
		if err == errs.ErrAlreadyExist {
			return ``, err
		}

		auth.logger.Error("failed create user", zap.Error(err))
		return ``, errs.ErrInternal
	}

	tokenStr, errJWT := token.GenerateJWT(in.Email, auth.secretKey)
	if errJWT != nil {
		auth.logger.Error("failed generate JWT", zap.Error(errJWT))
		return ``, errs.ErrInternal
	}

	return tokenStr, nil
}

// ChangePassword - Смена пароля пользователя.
func (auth AuthAppService) ChangePassword(email, password string) error {

	if len(password) < minPasswordLength {
		return ErrShortPassword
	}

	if err := auth.store.Update(email, password); err != nil {

		if errors.Is(err, errs.ErrNotFound) {
			return ErrUnauthenticated
		}

		auth.logger.Error("failed update user password", zap.Error(err))
		return errs.ErrInternal
	}

	return nil
}

// checkCredential - Проверка корректности пароля и email.
func checkCredential(cred authModel.Credential) error {
	if len(cred.Password) < minPasswordLength {
		return ErrShortPassword
	}

	if !strings.Contains(cred.Email, "@") {
		return ErrInvalidEmail
	}

	return nil
}
