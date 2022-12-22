package storage

type UserStorage interface {
	Create(cred Credential) error
	Find(email string) (Credential, error)
	Delete(email string) error
	Update(email, password string) error
}

type Credential struct {
	Email    string
	Password string
}
