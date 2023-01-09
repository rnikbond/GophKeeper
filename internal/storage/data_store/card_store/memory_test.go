package card_store

import (
	"GophKeeper/internal/model/card"
	"GophKeeper/pkg/errs"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func getTime(value string) time.Time {
	t, _ := time.Parse("01.2006", value)
	return t
}

func TestCardStore_Memory(t *testing.T) {

	store := NewMemoryStorage()

	testDataOK := card.DataCardFull{
		MetaInfo: "MirPay",
		Number:   "4648289760410976",
		Period:   getTime("10.2030"),
		CVV:      "111",
		FullName: "Test Test",
	}

	testDataChange := card.DataCardFull{
		MetaInfo: "MirPay",
		Number:   "4648289760410976",
		Period:   getTime("11.2030"),
		CVV:      "111",
		FullName: "Test Test",
	}

	testDataGet := card.DataCardGet{
		MetaInfo: "MirPay",
	}

	testDataFail := card.DataCardGet{
		MetaInfo: "GPay",
	}

	errCreate := store.Create(testDataOK)
	require.NoError(t, errCreate)

	data, errGet := store.Get(testDataGet)
	require.NoError(t, errGet)
	require.Equal(t, data, testDataOK)

	_, errGet = store.Get(testDataFail)
	require.Error(t, errGet, errs.ErrNotFound)

	errChange := store.Change(testDataChange)
	require.NoError(t, errChange)

	data, errGet = store.Get(testDataGet)
	require.NoError(t, errGet)
	require.Equal(t, data, testDataChange)

	errDel := store.Delete(testDataGet)
	require.NoError(t, errDel)

	_, errGet = store.Get(testDataGet)
	require.Error(t, errGet, errs.ErrNotFound)

	errChange = store.Change(testDataChange)
	require.Error(t, errGet, errs.ErrNotFound)

	errCreate = store.Create(testDataOK)
	require.NoError(t, errCreate)

	errCreate = store.Create(testDataOK)
	require.Error(t, errCreate, errs.ErrAlreadyExist)
}