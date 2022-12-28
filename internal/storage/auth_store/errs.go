package auth_store

type ErrStorage struct {
	Value string
}

func NewErr(val string) ErrStorage {
	return ErrStorage{
		Value: val,
	}
}

func (es ErrStorage) Error() string {
	return es.Value
}

// Ошибки User
var (
	ErrNotFound      = NewErr("user not found")
	ErrAlreadyExists = NewErr("user already exists")
)
