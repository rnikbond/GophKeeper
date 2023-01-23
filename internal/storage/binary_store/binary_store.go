//go:generate mockgen -source binary_store.go -destination mocks/binary_store_mock.go -package binary_store
package binary_store

import (
	"GophKeeper/internal/server/model/binary"
)

type BinaryStorage interface {
	Create(in binary.DataFull) error
	Get(in binary.DataGet) (binary.DataFull, error)
	Delete(in binary.DataGet) error
	Change(in binary.DataFull) error
}
