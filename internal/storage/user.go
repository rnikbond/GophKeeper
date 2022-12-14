package storage

type UserStorage interface {
	Create(user User) error
	Find(email string) (User, error)
	Delete(email string) error
	Update(email, password string) error
}

type User struct {
	Email    string
	Password string
}
