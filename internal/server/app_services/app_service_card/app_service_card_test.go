package app_service_card

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"GophKeeper/internal/server/model/card"
	"GophKeeper/internal/storage/card_store"
	"GophKeeper/pkg/errs"
)

func TestCardAppService_Chain(t *testing.T) {

	store := card_store.NewMemoryStorage()
	serv := NewCardAppService(store)

	testDataOK := card.DataCardFull{
		MetaInfo: "MirPay",
		Number:   "4648289760410976",
		Period:   "10.2030",
		CVV:      "111",
		FullName: "Test Test",
	}

	testDataChange := card.DataCardFull{
		MetaInfo: "MirPay",
		Number:   "4648289760410976",
		Period:   "11.2030",
		CVV:      "111",
		FullName: "Test Test",
	}

	testDataGet := card.DataCardGet{
		MetaInfo: "MirPay",
	}

	testDataFail := card.DataCardGet{
		MetaInfo: "GPay",
	}

	errCreate := serv.Create(testDataOK)
	require.NoError(t, errCreate)

	data, errGet := serv.Get(testDataGet)
	require.NoError(t, errGet)
	require.Equal(t, data, testDataOK)

	_, errGet = serv.Get(testDataFail)
	require.Error(t, errGet, errs.ErrNotFound)

	errChange := serv.Change(testDataChange)
	require.NoError(t, errChange)

	data, errGet = serv.Get(testDataGet)
	require.NoError(t, errGet)
	require.Equal(t, data, testDataChange)

	errDel := serv.Delete(testDataGet)
	require.NoError(t, errDel)

	_, errGet = serv.Get(testDataGet)
	require.Error(t, errGet, errs.ErrNotFound)

	errChange = serv.Change(testDataChange)
	require.Error(t, errGet, errs.ErrNotFound)

	errCreate = serv.Create(testDataOK)
	require.NoError(t, errCreate)

	errCreate = serv.Create(testDataOK)
	require.Error(t, errCreate, errs.ErrAlreadyExist)
}

func TestCardAppService_Create(t *testing.T) {

	store := card_store.NewMemoryStorage()
	serv := NewCardAppService(store)

	tests := []struct {
		name    string
		in      card.DataCardFull
		waitErr error
	}{
		{
			name: "Success",
			in: card.DataCardFull{
				MetaInfo: "MirPay",
				Number:   "4648289760410976",
				Period:   "10.2030",
				CVV:      "111",
				FullName: "Test Test",
			},
			waitErr: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			err := serv.Create(tt.in)
			assert.Equal(t, tt.waitErr, err)
		})
	}
}
