package app_services

import (
	"GophKeeper/internal/model/binary"
	"GophKeeper/internal/storage/data_store/binary_store"
	"GophKeeper/pkg/errs"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestBinaryAppService(t *testing.T) {

	store := binary_store.NewMemoryStorage()
	serv := NewBinaryAppService(store)

	testDataOK := binary.DataFull{
		Email:    "test@email.com",
		MetaInfo: "www.ololo.com",
		Bytes:    []byte("010101010101"),
	}

	testDataChange := binary.DataFull{
		Email:    "test@email.com",
		MetaInfo: "www.ololo.com",
		Bytes:    []byte("000000000000000000"),
	}

	testDataGet := binary.DataGet{
		Email:    "test@email.com",
		MetaInfo: "www.ololo.com",
	}

	testDataFail := binary.DataGet{
		Email:    "test@email.com",
		MetaInfo: "www.test.com",
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
