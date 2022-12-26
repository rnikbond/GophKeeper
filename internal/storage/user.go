package storage

// Credential - Учетные данные пользователя.
type Credential struct {
	// Email - Почтовый адрес.
	Email string
	// Password - Пароль.
	Password string
}

type UserStorage interface {
	Create(cred Credential) error
	Find(email string) (Credential, error)
	Delete(email string) error
	Update(email, password string) error
}
