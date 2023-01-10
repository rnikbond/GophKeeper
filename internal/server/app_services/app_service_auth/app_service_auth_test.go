package app_service_auth

import (
	"GophKeeper/internal/storage/auth_store"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	"GophKeeper/internal/model/auth"
	storeMock "GophKeeper/internal/storage/auth_store/mocks"
	"GophKeeper/pkg/errs"
)

func TestAuthAppService_Login(t *testing.T) {

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	store := storeMock.NewMockAuthStorage(ctrl)

	tests := []struct {
		name      string
		credStore auth.Credential
		credServ  auth.Credential
		waitErr   error
		storeErr  error
	}{
		{
			name: "Success",
			credStore: auth.Credential{
				Email:    "test@emailcom",
				Password: "testPassword",
			},
			credServ: auth.Credential{
				Email:    "test@emailcom",
				Password: "testPassword",
			},
			waitErr:  nil,
			storeErr: nil,
		},
		{
			name:      "Unregistered",
			credStore: auth.Credential{},
			credServ: auth.Credential{
				Email:    "test@emailcom",
				Password: "testPassword",
			},
			waitErr:  errs.ErrNotFound,
			storeErr: errs.ErrNotFound,
		},
		{
			name: "Invalid password",
			credStore: auth.Credential{
				Email:    "test@emailcom",
				Password: "passwordTest",
			},
			credServ: auth.Credential{
				Email:    "test@emailcom",
				Password: "testPassword",
			},
			waitErr:  ErrInvalidPassword,
			storeErr: auth_store.ErrInvalidPassword,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			store.EXPECT().Find(tt.credServ).Return(tt.storeErr)

			authServ := NewAuthService(store)
			_, err := authServ.Login(tt.credServ)

			assert.Equal(t, err, tt.waitErr)
		})
	}
}

func TestAuthAppService_Register(t *testing.T) {

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	store := storeMock.NewMockAuthStorage(ctrl)

	tests := []struct {
		name      string
		cred      auth.Credential
		waitErr   error
		callStore bool
	}{
		{
			name: "Success",
			cred: auth.Credential{
				Email:    "test@email.com",
				Password: "testPassword",
			},
			waitErr:   nil,
			callStore: true,
		},
		{
			name: "Invalid email",
			cred: auth.Credential{
				Email:    "test_email.com",
				Password: "testPassword",
			},
			waitErr:   ErrInvalidEmail,
			callStore: false,
		},
		{
			name: "Short password",
			cred: auth.Credential{
				Email:    "test@email.com",
				Password: "pwd",
			},
			waitErr:   ErrShortPassword,
			callStore: false,
		},
		{
			name: "User already exist",
			cred: auth.Credential{
				Email:    "test@email.com",
				Password: "passwordTest",
			},
			waitErr:   errs.ErrAlreadyExist,
			callStore: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			if tt.callStore {
				store.EXPECT().Create(tt.cred).Return(tt.waitErr)
			}

			authServ := NewAuthService(store)
			_, err := authServ.Register(tt.cred)

			assert.Equal(t, err, tt.waitErr)
		})
	}
}

func TestAuthAppService_ChangePassword(t *testing.T) {

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	store := storeMock.NewMockAuthStorage(ctrl)

	tests := []struct {
		name            string
		email           string
		password        string
		waitErr         error
		callStoreUpdate bool
		errStore        error
	}{
		{
			name:            "Success",
			email:           "test@email.com",
			password:        "qwerty123",
			callStoreUpdate: true,
		},
		{
			name:            "Invalid password",
			email:           "test@email.com",
			password:        "",
			waitErr:         ErrShortPassword,
			callStoreUpdate: false,
		},
		{
			name:            "Unauthenticated",
			email:           "test@email.com",
			password:        "qwerty123",
			waitErr:         ErrUnauthenticated,
			callStoreUpdate: true,
			errStore:        errs.ErrNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			if tt.callStoreUpdate {
				store.EXPECT().Update(tt.email, tt.password).Return(tt.errStore)
			}

			authServ := NewAuthService(store)
			err := authServ.ChangePassword(tt.email, tt.password)

			assert.Equal(t, tt.waitErr, err)
		})
	}
}
