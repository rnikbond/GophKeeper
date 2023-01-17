package binary_store

import (
	"testing"

	"github.com/stretchr/testify/require"

	"GophKeeper/internal/server/model/binary"
	"GophKeeper/pkg/errs"
)

func TestBinaryStore_Memory(t *testing.T) {

	store := NewMemoryStorage()

	testDataOK := binary.DataFull{
		MetaInfo: "prog.bin",
		Bytes:    []byte("00000000000000"),
	}

	testDataChange := binary.DataFull{
		MetaInfo: "prog.bin",
		Bytes:    []byte("11111111111111"),
	}

	testDataGet := binary.DataGet{
		MetaInfo: "prog.bin",
	}

	testDataFail := binary.DataGet{
		MetaInfo: "prog1.bin",
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
