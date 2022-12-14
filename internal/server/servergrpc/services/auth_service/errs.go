package auth_service

type ErrAuth struct {
	Value string
}

func NewErr(val string) ErrAuth {
	return ErrAuth{
		Value: val,
	}
}

func (err ErrAuth) Error() string {
	return err.Value
}

var (
	ErrInvalidToken    = NewErr("invalid token")
	ErrInvalidAuthData = NewErr("invalid login or password")
	ErrGenerateToken   = NewErr("error generate access token")
)
