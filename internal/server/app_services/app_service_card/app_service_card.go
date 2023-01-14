//go:generate mockgen -source app_service_card.go -destination mocks/app_service_card_mock.go -package app_services
package app_service_card

import (
	"go.uber.org/zap"

	"GophKeeper/internal/model/card"
	"GophKeeper/internal/storage/card_store"
)

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

	return serv.store.Create(in)
}

func (serv CardAppService) Get(in card.DataCardGet) (card.DataCardFull, error) {

	return serv.store.Get(in)
}

func (serv CardAppService) Delete(in card.DataCardGet) error {
	return serv.store.Delete(in)
}

func (serv CardAppService) Change(in card.DataCardFull) error {

	return serv.store.Change(in)
}
