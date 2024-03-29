package app_service_auth

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

var (
	ErrInvalidEmail    = NewErr("invalid email")
	ErrInvalidPassword = NewErr("invalid password")
	ErrShortPassword   = NewErr("password must contain 6 or more characters")
	ErrUnauthenticated = NewErr("unauthenticated")
)
