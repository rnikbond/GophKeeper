package text_store

import (
	"GophKeeper/internal/model/text"
	"GophKeeper/pkg/errs"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestTextStore_Memory(t *testing.T) {

	store := NewMemoryStorage()

	testDataOK := text.DataTextFull{
		MetaInfo: "www.ololo.com",
		Text:     "qwerty",
	}

	testDataChange := text.DataTextFull{
		MetaInfo: "www.ololo.com",
		Text:     "qwerty123",
	}

	testDataGet := text.DataTextGet{
		MetaInfo: "www.ololo.com",
	}

	testDataFail := text.DataTextGet{
		MetaInfo: "www.test.com",
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
