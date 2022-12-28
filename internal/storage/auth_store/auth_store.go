//go:generate mockgen -source auth_store.go -destination ../../../mocks/storage/auth_store_mock.go -package storage
package auth_store

// Credential - Учетные данные пользователя.
type Credential struct {
	// Email - Почтовый адрес.
	Email string
	// Password - Пароль.
	Password string
}

type AuthStorage interface {
	Create(cred Credential) error
	Find(email string) (Credential, error)
	Delete(email string) error
	Update(email, password string) error
}
