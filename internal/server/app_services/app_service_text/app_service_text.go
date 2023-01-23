package app_service_text

import (
	"go.uber.org/zap"

	"GophKeeper/internal/server/model/text"
	"GophKeeper/internal/storage/text_store"
)

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
