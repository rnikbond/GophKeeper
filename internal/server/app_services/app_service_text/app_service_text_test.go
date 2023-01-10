package app_service_text

import (
	"GophKeeper/internal/storage/text_store"
	"testing"

	"github.com/stretchr/testify/require"

	"GophKeeper/internal/model/text"
	"GophKeeper/pkg/errs"
)

func TestTextAppService(t *testing.T) {

	store := text_store.NewMemoryStorage()
	serv := NewTextAppService(store)

	testDataOK := text.DataTextFull{
		MetaInfo: "note_private",
		Text:     "text text text",
	}

	testDataChange := text.DataTextFull{
		MetaInfo: "note_private",
		Text:     "qwerty123",
	}

	testDataGet := text.DataTextGet{
		MetaInfo: "note_private",
	}

	testDataFail := text.DataTextGet{
		MetaInfo: "note_private_1",
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
