package binary_store

import (
	"GophKeeper/internal/model/binary"
	"GophKeeper/pkg/errs"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestBinaryStore_Memory(t *testing.T) {

	store := NewMemoryStorage()

	testDataOK := binary.DataFull{
		Email:    "test@email.com",
		MetaInfo: "www.ololo.com",
		Bytes:    []byte("00000000000000"),
	}

	testDataChange := binary.DataFull{
		Email:    "test@email.com",
		MetaInfo: "www.ololo.com",
		Bytes:    []byte("11111111111111"),
	}

	testDataGet := binary.DataGet{
		Email:    "test@email.com",
		MetaInfo: "www.ololo.com",
	}

	testDataFail := binary.DataGet{
		Email:    "test@email.com",
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
