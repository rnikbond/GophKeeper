package model

type ErrModel struct {
	Value string
}

func NewErr(val string) ErrModel {
	return ErrModel{
		Value: val,
	}
}

func (es ErrModel) Error() string {
	return es.Value
}

// Ошибки User
var (
	ErrNotFound        = NewErr("user not found")
	ErrAlreadyExists   = NewErr("user already exists")
	ErrInvalidEmail    = NewErr("invalid email")
	ErrInvalidPassword = NewErr("invalid password")
	ErrInvalidToken    = NewErr("invalid token")
	ErrShortPassword   = NewErr("password must contain 6 or more characters")
	ErrInternal        = NewErr("internal error")
)
