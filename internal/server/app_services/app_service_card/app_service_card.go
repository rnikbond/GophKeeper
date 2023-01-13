//go:generate mockgen -source app_service_card.go -destination mocks/app_service_card_mock.go -package app_services
package app_service_card

import (
	"strconv"
	"time"

	"github.com/EClaesson/go-luhn"

	"go.uber.org/zap"

	"GophKeeper/internal/model/card"
	"GophKeeper/internal/storage/card_store"
)

var PeriodLayout = "01.2006"

type CardApp interface {
	Create(data card.DataCardFull) error
	Get(in card.DataCardGet) (card.DataCardFull, error)
	Delete(in card.DataCardGet) error
	Change(in card.DataCardFull) error
}

type CardAppService struct {
	store  card_store.CardStorage
	logger *zap.Logger
}

func NewCardAppService(store card_store.CardStorage) *CardAppService {
	return &CardAppService{
		store:  store,
		logger: zap.L(),
	}
}

func (serv CardAppService) Create(in card.DataCardFull) error {

	if err := checkCardData(in); err != nil {
		return err
	}

	return serv.store.Create(in)
}

func (serv CardAppService) Get(in card.DataCardGet) (card.DataCardFull, error) {

	return serv.store.Get(in)
}

func (serv CardAppService) Delete(in card.DataCardGet) error {
	return serv.store.Delete(in)
}

func (serv CardAppService) Change(in card.DataCardFull) error {

	if err := checkCardData(in); err != nil {
		return err
	}

	return serv.store.Change(in)
}

func checkCardData(in card.DataCardFull) error {

	if _, errTime := time.Parse(PeriodLayout, in.Period); errTime != nil {
		return ErrInvalidPeriod
	}

	if ok, err := luhn.IsValid(in.Number); !ok || err != nil {
		return ErrInvalidNumber
	}

	if len(in.CVV) != 3 {
		return ErrInvalidCVV
	}

	// Используется ParseUint - т.к. не должно быть отрицательного CVV. Например, "-12".
	if _, err := strconv.ParseUint(in.CVV, 10, 32); err != nil {
		return ErrInvalidCVV
	}

	if len(in.FullName) < 4 {
		return ErrInvalidFullName
	}

	return nil
}
