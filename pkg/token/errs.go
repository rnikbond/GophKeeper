package token

type ErrToken struct {
	Value string
}

func NewErr(val string) ErrToken {
	return ErrToken{
		Value: val,
	}
}

func (es ErrToken) Error() string {
	return es.Value
}

var (
	ErrTokenNotFound = NewErr("token not found")
)
