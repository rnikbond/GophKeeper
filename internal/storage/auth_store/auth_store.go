//go:generate mockgen -source auth_store.go -destination mocks/auth_store_mock.go -package storage
package auth_store

import (
	"GophKeeper/internal/server/model/auth"
)

type AuthStorage interface {
	Create(cred auth.Credential) error
	Find(cred auth.Credential) error
	Update(email, password string) error
	Delete(email string) error
}
