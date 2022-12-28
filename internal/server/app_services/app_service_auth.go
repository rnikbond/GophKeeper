//go:generate mockgen -source app_service_auth.go -destination ../../../mocks/server/app_services/app_service_auth_mock.go -package app_services
package app_services

import (
	"GophKeeper/internal/storage/auth_store"
	"GophKeeper/pkg/token"
	"go.uber.org/zap"
	"strings"
)

// TODO :: Вынести в config
var minPasswordLength = 6

type AuthApp interface {
	Login(in Credential) (string, error)
	Register(in Credential) (string, error)
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

// Credential - Учетные данные пользователя.
type Credential struct {
	// Email - Почтовый адрес.
	Email string
	// Password - Пароль.
	Password string
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
func (auth AuthAppService) Login(in Credential) (string, error) {

	userStore, err := auth.store.Find(in.Email)
	if err != nil {
		if err == auth_store.ErrNotFound {
			return ``, ErrNotFound
		}

		auth.logger.Error("failed find user", zap.Error(err))
		return ``, ErrInternal
	}

	if in.Password != userStore.Password {
		return ``, ErrInvalidPassword
	}

	tokenStr, errJWT := token.GenerateJWT(in.Email, auth.secretKey)
	if errJWT != nil {
		auth.logger.Error("failed generate JWT", zap.Error(errJWT))
		return ``, ErrInternal
	}

	return tokenStr, nil
}

// Register - Регистрация пользователя.
// При успешной регистрации возвращается JWT.
func (auth AuthAppService) Register(in Credential) (string, error) {

	cred := Credential{
		Email:    in.Email,
		Password: in.Password,
	}

	if errCred := checkCredential(cred); errCred != nil {
		return ``, errCred
	}

	storeCred := auth_store.Credential{
		Email:    cred.Email,
		Password: cred.Password,
	}

	if err := auth.store.Create(storeCred); err != nil {
		if err == auth_store.ErrAlreadyExists {
			return ``, ErrAlreadyExists
		}

		auth.logger.Error("failed create user", zap.Error(err))
		return ``, ErrInternal
	}

	tokenStr, errJWT := token.GenerateJWT(in.Email, auth.secretKey)
	if errJWT != nil {
		auth.logger.Error("failed generate JWT", zap.Error(errJWT))
		return ``, ErrInternal
	}

	return tokenStr, nil
}

// ChangePassword - Смена пароля пользователя.
func (auth AuthAppService) ChangePassword(email, password string) error {

	if len(password) < minPasswordLength {
		return ErrShortPassword
	}

	if _, err := auth.store.Find(email); err != nil {
		if err == auth_store.ErrNotFound {
			return ErrUnauthenticated
		}

		auth.logger.Error("failed find user", zap.Error(err), zap.String("email", email))
		return ErrInternal
	}

	if err := auth.store.Update(email, password); err != nil {
		auth.logger.Error("failed update user password", zap.Error(err))
		return ErrInternal
	}

	return nil
}

// checkCredential - Проверка корректности пароля и email.
func checkCredential(cred Credential) error {
	if len(cred.Password) < minPasswordLength {
		return ErrShortPassword
	}

	if !strings.Contains(cred.Email, "@") {
		return ErrInvalidEmail
	}

	return nil
}
