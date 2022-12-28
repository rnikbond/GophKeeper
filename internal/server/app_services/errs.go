package app_services

type ErrAppServices struct {
	Value string
}

func NewErr(val string) ErrAppServices {
	return ErrAppServices{
		Value: val,
	}
}

func (es ErrAppServices) Error() string {
	return es.Value
}

// Ошибки User
var (
	ErrNotFound        = NewErr("user not found")
	ErrAlreadyExists   = NewErr("user already exists")
	ErrInvalidEmail    = NewErr("invalid email")
	ErrInvalidPassword = NewErr("invalid password")
	ErrShortPassword   = NewErr("password must contain 6 or more characters")
	ErrUnauthenticated = NewErr("unauthenticated")
	ErrInternal        = NewErr("internal error")
)
