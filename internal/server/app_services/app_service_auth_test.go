package app_services

import (
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	"GophKeeper/internal/storage"
	mock "GophKeeper/mocks/interbal/storage"
)

func TestAuthAppService_Login(t *testing.T) {

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	store := mock.NewMockUserStorage(ctrl)

	tests := []struct {
		name      string
		credStore storage.Credential
		credServ  Credential
		waitErr   error
		storeErr  error
	}{
		{
			name: "Success",
			credStore: storage.Credential{
				Email:    "test@emailcom",
				Password: "testPassword",
			},
			credServ: Credential{
				Email:    "test@emailcom",
				Password: "testPassword",
			},
			waitErr:  nil,
			storeErr: nil,
		},
		{
			name:      "Unregistered",
			credStore: storage.Credential{},
			credServ: Credential{
				Email:    "test@emailcom",
				Password: "testPassword",
			},
			waitErr:  ErrNotFound,
			storeErr: storage.ErrNotFound,
		},
		{
			name: "Invalid password",
			credStore: storage.Credential{
				Email:    "test@emailcom",
				Password: "passwordTest",
			},
			credServ: Credential{
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

			auth := NewAuthService(store)
			_, err := auth.Login(tt.credServ)

			assert.Equal(t, err, tt.waitErr)
		})
	}
}

func TestAuthAppService_Register(t *testing.T) {

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	store := mock.NewMockUserStorage(ctrl)

	tests := []struct {
		name      string
		credStore storage.Credential
		credServ  Credential
		waitErr   error
		storeErr  error
		callStore bool
	}{
		{
			name: "Success",
			credStore: storage.Credential{
				Email:    "test@email.com",
				Password: "testPassword",
			},
			credServ: Credential{
				Email:    "test@email.com",
				Password: "testPassword",
			},
			waitErr:   nil,
			storeErr:  nil,
			callStore: true,
		},
		{
			name: "Invalid email",
			credServ: Credential{
				Email:    "test_email.com",
				Password: "testPassword",
			},
			waitErr:   ErrInvalidEmail,
			storeErr:  nil,
			callStore: false,
		},
		{
			name: "Short password",
			credServ: Credential{
				Email:    "test@email.com",
				Password: "pwd",
			},
			waitErr:   ErrShortPassword,
			storeErr:  nil,
			callStore: false,
		},
		{
			name: "User already exist",
			credStore: storage.Credential{
				Email:    "test@email.com",
				Password: "passwordTest",
			},
			credServ: Credential{
				Email:    "test@email.com",
				Password: "passwordTest",
			},
			waitErr:   ErrAlreadyExists,
			storeErr:  storage.ErrAlreadyExists,
			callStore: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			if tt.callStore {
				store.EXPECT().Create(tt.credStore).Return(tt.storeErr)
			}

			auth := NewAuthService(store)
			_, err := auth.Register(tt.credServ)

			assert.Equal(t, err, tt.waitErr)
		})
	}
}

func TestAuthAppService_ChangePassword(t *testing.T) {

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	store := mock.NewMockUserStorage(ctrl)

	tests := []struct {
		name     string
		email    string
		password string
	}{
		{},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			store.EXPECT().Find(tt.email).Return(storage.Credential{}, nil)
			store.EXPECT().Update(tt.email, tt.password).Return(nil)

			auth := NewAuthService(store)
			err := auth.ChangePassword(tt.email, tt.password)
			assert.NoError(t, err)
		})
	}
}
