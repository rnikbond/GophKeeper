package app_service_card

import "fmt"

type ErrCard struct {
	Value string
}

func NewErrCard(val string) ErrCard {
	return ErrCard{
		Value: val,
	}
}

func (es ErrCard) Error() string {
	return es.Value
}

var (
	ErrInvalidPeriod   = NewErrCard(fmt.Sprintf("invalid card period. Must been in format: %s", PeriodLayout))
	ErrInvalidNumber   = NewErrCard("invalid card number")
	ErrInvalidCVV      = NewErrCard("invalid CVV code card")
	ErrInvalidFullName = NewErrCard("invalid full name card holder")
)
