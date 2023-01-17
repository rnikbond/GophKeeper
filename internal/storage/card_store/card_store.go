package card_store

import (
	"GophKeeper/internal/server/model/card"
)

type CardStorage interface {
	Create(data card.DataCardFull) error
	Get(in card.DataCardGet) (card.DataCardFull, error)
	Delete(in card.DataCardGet) error
	Change(in card.DataCardFull) error
}
