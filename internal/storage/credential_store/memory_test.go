package credential_store

import (
	"testing"

	"github.com/stretchr/testify/require"

	"GophKeeper/internal/model/cred"
	"GophKeeper/pkg/errs"
)

func TestCredentialStore_Memory(t *testing.T) {

	store := NewMemoryStorage()

	testDataOK := cred.CredentialFull{
		Email:    "test@email.com",
		MetaInfo: "www.ololo.com",
		Password: "qwerty",
	}

	testDataChange := cred.CredentialFull{
		Email:    "test@email.com",
		MetaInfo: "www.ololo.com",
		Password: "qwerty123",
	}

	testDataGet := cred.CredentialGet{
		Email:    "test@email.com",
		MetaInfo: "www.ololo.com",
	}

	testDataFail := cred.CredentialGet{
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
