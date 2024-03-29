package errs

type ErrKeeper struct {
	Value string
}

func NewErr(val string) ErrKeeper {
	return ErrKeeper{
		Value: val,
	}
}

func (es ErrKeeper) Error() string {
	return es.Value
}

var (
	ErrNotFound        = NewErr("not found")
	ErrAlreadyExist    = NewErr("already exist")
	ErrInvalidArgument = NewErr("invalid argument")
	ErrInternal        = NewErr("internal error")
	ErrCancel          = NewErr("operation canceled")
	ErrLargeData       = NewErr("large data")
)
