package app_services

import (
	"GophKeeper/internal/model/auth"
	storeMock "GophKeeper/internal/storage/auth_store/mocks"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	"GophKeeper/internal/storage/auth_store"
)

func TestAuthAppService_Login(t *testing.T) {

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	store := storeMock.NewMockAuthStorager(ctrl)

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
			waitErr:  ErrNotFound,
			storeErr: auth_store.ErrNotFound,
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
			storeErr: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			store.EXPECT().Find(tt.credServ.Email).Return(tt.credStore, tt.storeErr)

			authServ := NewAuthService(store)
			_, err := authServ.Login(tt.credServ)

			assert.Equal(t, err, tt.waitErr)
		})
	}
}

func TestAuthAppService_Register(t *testing.T) {

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	store := storeMock.NewMockAuthStorager(ctrl)

	tests := []struct {
		name      string
		cred      auth.Credential
		waitErr   error
		storeErr  error
		callStore bool
	}{
		{
			name: "Success",
			cred: auth.Credential{
				Email:    "test@email.com",
				Password: "testPassword",
			},
			waitErr:   nil,
			storeErr:  nil,
			callStore: true,
		},
		{
			name: "Invalid email",
			cred: auth.Credential{
				Email:    "test_email.com",
				Password: "testPassword",
			},
			waitErr:   ErrInvalidEmail,
			storeErr:  nil,
			callStore: false,
		},
		{
			name: "Short password",
			cred: auth.Credential{
				Email:    "test@email.com",
				Password: "pwd",
			},
			waitErr:   ErrShortPassword,
			storeErr:  nil,
			callStore: false,
		},
		{
			name: "User already exist",
			cred: auth.Credential{
				Email:    "test@email.com",
				Password: "passwordTest",
			},
			waitErr:   ErrAlreadyExists,
			storeErr:  auth_store.ErrAlreadyExists,
			callStore: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			if tt.callStore {
				store.EXPECT().Create(tt.cred).Return(tt.storeErr)
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

	store := storeMock.NewMockAuthStorager(ctrl)

	tests := []struct {
		name            string
		email           string
		password        string
		waitErr         error
		callStoreFind   bool
		callStoreUpdate bool
		errStore        error
	}{
		{
			name:            "Success",
			email:           "test@email.com",
			password:        "qwerty123",
			callStoreFind:   true,
			callStoreUpdate: true,
		},
		{
			name:            "Invalid password",
			email:           "test@email.com",
			password:        "",
			waitErr:         ErrShortPassword,
			callStoreFind:   false,
			callStoreUpdate: false,
		},
		{
			name:            "Unauthenticated",
			email:           "test@email.com",
			password:        "qwerty123",
			waitErr:         ErrUnauthenticated,
			callStoreFind:   true,
			callStoreUpdate: false,
			errStore:        auth_store.ErrNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			if tt.callStoreFind {
				store.EXPECT().Find(tt.email).Return(auth.Credential{}, tt.errStore)
			}

			if tt.callStoreUpdate {
				store.EXPECT().Update(tt.email, tt.password).Return(tt.errStore)
			}

			authServ := NewAuthService(store)
			err := authServ.ChangePassword(tt.email, tt.password)

			assert.Equal(t, tt.waitErr, err)
		})
	}
}
