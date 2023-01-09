//go:generate mockgen -source app_service_card.go -destination mocks/app_service_card_mock.go -package app_services
package app_services

import (
	"GophKeeper/internal/model/card"
	"GophKeeper/internal/storage/data_store/card_store"
	"github.com/EClaesson/go-luhn"
	"go.uber.org/zap"
	"strconv"
	"time"
)

var PeriodLayout = "01.2006"

type CardApp interface {
	Create(data card.DataCard) error
	Get(in card.DataCardGet) (card.DataCard, error)
	Delete(in card.DataCardGet) error
	Change(in card.DataCard) error
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

func (serv CardAppService) Create(in card.DataCard) error {

	data, err := convertCardData(in)
	if err != nil {
		return err
	}

	return serv.store.Create(data)
}

func (serv CardAppService) Get(in card.DataCardGet) (card.DataCard, error) {

	data, err := serv.store.Get(in)
	if err != nil {
		return card.DataCard{}, err
	}

	out := card.DataCard{
		MetaInfo: data.MetaInfo,
		Number:   data.Number,
		Period:   data.Period.Format(PeriodLayout),
		CVV:      data.CVV,
		FullName: data.FullName,
	}

	return out, nil
}

func (serv CardAppService) Delete(in card.DataCardGet) error {
	return serv.store.Delete(in)
}

func (serv CardAppService) Change(in card.DataCard) error {
	data, err := convertCardData(in)
	if err != nil {
		return err
	}

	return serv.store.Change(data)
}

func convertCardData(in card.DataCard) (card.DataCardFull, error) {
	t, errTime := time.Parse(PeriodLayout, in.Period)
	if errTime != nil {
		return card.DataCardFull{}, ErrInvalidPeriod
	}

	if ok, err := luhn.IsValid(in.Number); !ok || err != nil {
		return card.DataCardFull{}, ErrInvalidNumber
	}

	if len(in.CVV) != 3 {
		return card.DataCardFull{}, ErrInvalidCVV
	}

	if _, err := strconv.Atoi(in.CVV); err != nil {
		return card.DataCardFull{}, ErrInvalidCVV
	}

	if len(in.FullName) < 4 {
		return card.DataCardFull{}, ErrInvalidFullName
	}

	data := card.DataCardFull{
		MetaInfo: in.MetaInfo,
		Number:   in.Number,
		Period:   t,
		CVV:      in.CVV,
		FullName: in.FullName,
	}

	return data, nil
}