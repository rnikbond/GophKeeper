//go:generate mockgen -source app_service_text.go -destination mocks/app_service_text_mock.go -package app_services
package app_services

import (
	"go.uber.org/zap"

	"GophKeeper/internal/model/text"
	"GophKeeper/internal/storage/data_store/text_store"
)

type TextApp interface {
	Create(in text.DataTextFull) error
	Get(in text.DataTextGet) (text.DataTextFull, error)
	Delete(in text.DataTextGet) error
	Change(in text.DataTextFull) error
}

type TextAppService struct {
	store  text_store.TextStorage
	logger *zap.Logger
}

func NewTextAppService(store text_store.TextStorage) *TextAppService {
	return &TextAppService{
		store:  store,
		logger: zap.L(),
	}
}

func (serv TextAppService) Create(in text.DataTextFull) error {
	return serv.store.Create(in)
}

func (serv TextAppService) Get(in text.DataTextGet) (text.DataTextFull, error) {
	return serv.store.Get(in)
}

func (serv TextAppService) Delete(in text.DataTextGet) error {
	return serv.store.Delete(in)
}

func (serv TextAppService) Change(in text.DataTextFull) error {
	return serv.store.Change(in)
}
