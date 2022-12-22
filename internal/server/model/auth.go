package model

import (
	"GophKeeper/internal/storage"
	"GophKeeper/pkg/token"
	"go.uber.org/zap"
	"strings"
)

type AuthModel struct {
	store     storage.UserStorage
	secretKey string
	logger    *zap.Logger
}

func NewAuthModel(store storage.UserStorage) *AuthModel {
	return &AuthModel{
		store:  store,
		logger: zap.L(),
	}
}

// Login - Авторизация пользователя.
func (auth AuthModel) Login(in storage.Credential) (string, error) {

	userStore, err := auth.store.Find(in.Email)
	if err != nil {
		if err == storage.ErrNotFound {
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
func (auth AuthModel) Register(in storage.Credential) (string, error) {

	cred := storage.Credential{
		Email:    in.Email,
		Password: in.Password,
	}

	if errCred := checkCredential(cred); errCred != nil {
		return ``, errCred
	}

	if err := auth.store.Create(cred); err != nil {
		if err == storage.ErrAlreadyExists {
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
func (auth AuthModel) ChangePassword(email, password string) error {

	if _, err := auth.store.Find(email); err != nil {
		auth.logger.Error("failed find user", zap.Error(err))
		return ErrInternal
	}

	if err := auth.store.Update(email, password); err != nil {
		auth.logger.Error("failed update user password", zap.Error(err))
		return ErrInternal
	}

	return nil
}

// checkCredential - Проверка корректности пароля и email.
func checkCredential(cred storage.Credential) error {
	if len(cred.Password) < 6 {
		return ErrShortPassword
	}

	if !strings.Contains(cred.Email, "@") {
		return ErrInvalidEmail
	}

	return nil
}
