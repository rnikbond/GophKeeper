package auth_store

type ErrAuth struct {
	Value string
}

func NewErr(val string) ErrAuth {
	return ErrAuth{
		Value: val,
	}
}

func (es ErrAuth) Error() string {
	return es.Value
}

var (
	ErrInvalidPassword = NewErr("invalid password")
)
