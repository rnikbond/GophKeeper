//go:generate mockgen -source text_store.go -destination mocks/text_store_mock.go -package credential_store
package text_store

import (
	"GophKeeper/internal/server/model/text"
)

type TextStorage interface {
	Create(data text.DataTextFull) error
	Get(in text.DataTextGet) (text.DataTextFull, error)
	Delete(in text.DataTextGet) error
	Change(in text.DataTextFull) error
}
