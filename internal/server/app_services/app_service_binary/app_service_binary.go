//go:generate mockgen -source app_service_binary.go -destination mocks/app_service_binary_mock.go -package app_services
package app_service_binary

import (
	"go.uber.org/zap"

	"GophKeeper/internal/model/binary"
	"GophKeeper/internal/storage/data_store/binary_store"
)

type BinaryApp interface {
	Create(in binary.DataFull) error
	Get(in binary.DataGet) (binary.DataFull, error)
	Delete(in binary.DataGet) error
	Change(in binary.DataFull) error
}

type BinaryAppService struct {
	store  binary_store.BinaryStorage
	logger *zap.Logger
}

func NewBinaryAppService(store binary_store.BinaryStorage) *BinaryAppService {
	return &BinaryAppService{
		store:  store,
		logger: zap.L(),
	}
}

func (serv BinaryAppService) Create(in binary.DataFull) error {
	return serv.store.Create(in)
}

func (serv BinaryAppService) Get(in binary.DataGet) (binary.DataFull, error) {
	return serv.store.Get(in)
}

func (serv BinaryAppService) Delete(in binary.DataGet) error {
	return serv.store.Delete(in)
}

func (serv BinaryAppService) Change(in binary.DataFull) error {
	return serv.store.Change(in)
}
